package addon_test

import (
	"testing"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// deductionInvoice returns an invoice with two advance payment invoices as preceding
// documents and one deduction for each of them.
func deductionInvoice() *bill.Invoice {
	inv := validInvoice()
	inv.Preceding = []*org.DocumentRef{
		{Series: "FT SERIES-A", Code: "1", IssueDate: cal.NewDate(2022, 12, 1)},
		{Series: "FT SERIES-A", Code: "2", IssueDate: cal.NewDate(2022, 12, 2)},
	}
	inv.Discounts = []*bill.Discount{
		deduction(1000, "FT SERIES-A/1"),
		deduction(2000, "FT SERIES-A/2"),
	}
	return inv
}

// deduction returns an advance payment deduction for the given amount in cents and reference.
func deduction(cents int64, ref cbc.Code) *bill.Discount {
	d := &bill.Discount{
		Key:    addon.DiscountKeyDeduction,
		Reason: "Deposit",
		Amount: num.MakeAmount(cents, 2),
		Taxes:  tax.Set{{Category: tax.CategoryVAT, Rate: "general"}},
	}
	if ref != "" {
		d.Ext = tax.ExtensionsOf(cbc.CodeMap{addon.ExtKeyReference: ref})
	}
	return d
}

func TestDeductionNormalization(t *testing.T) {
	t.Run("defaults reference to the only preceding document", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{{Series: "FT SERIES-A", Code: "1"}}
		inv.Discounts = []*bill.Discount{deduction(1000, "")}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FT SERIES-A/1", inv.Discounts[0].Ext.Get(addon.ExtKeyReference).String())
	})

	t.Run("keeps an explicit reference", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{{Series: "FT SERIES-A", Code: "1"}}
		inv.Discounts = []*bill.Discount{deduction(1000, "FT SERIES-A/9")}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "FT SERIES-A/9", inv.Discounts[0].Ext.Get(addon.ExtKeyReference).String())
	})

	t.Run("no default with several preceding documents", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts[0].Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		assert.False(t, inv.Discounts[0].Ext.Has(addon.ExtKeyReference))
	})

	t.Run("other discounts untouched", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{{Series: "FT SERIES-A", Code: "1"}}
		inv.Discounts = []*bill.Discount{{Reason: "Promo", Amount: num.MakeAmount(1000, 2)}}
		require.NoError(t, inv.Calculate())
		assert.False(t, inv.Discounts[0].Ext.Has(addon.ExtKeyReference))
	})
}

func TestDeductionValidation(t *testing.T) {
	t.Run("valid deductions on several preceding documents", func(t *testing.T) {
		inv := deductionInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid deduction without reference and one preceding", func(t *testing.T) {
		inv := validInvoice()
		inv.Preceding = []*org.DocumentRef{{Series: "FT SERIES-A", Code: "1"}}
		inv.Discounts = []*bill.Discount{deduction(1000, "")}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid deductions in simplified and invoice-receipt", func(t *testing.T) {
		for _, typ := range []cbc.Code{addon.InvoiceTypeSimplified, addon.InvoiceTypeInvoiceReceipt} {
			inv := deductionInvoice()
			inv.Series = cbc.Code(typ.String() + " SERIES-A")
			require.NoError(t, inv.Calculate())
			inv.Tax.Ext = inv.Tax.Ext.Set(addon.ExtKeyInvoiceType, typ)
			if typ == addon.InvoiceTypeInvoiceReceipt {
				inv.Totals.Due = num.NewAmount(0, 2)
			}
			assert.NoError(t, rules.Validate(inv), typ)
		}
	})

	t.Run("several preceding documents without deductions", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-11]")
	})

	t.Run("deductions in credit notes", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Series = "NC SERIES-A"
		inv.Preceding = inv.Preceding[:1]
		inv.Preceding[0].Reason = "Refund"
		inv.Discounts = inv.Discounts[:1]
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(addon.ExtKeyInvoiceType, addon.InvoiceTypeCreditNote)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-18] discount deductions are only supported in invoices (FT, FS, FR)")
	})

	t.Run("deductions in debit notes", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Series = "ND SERIES-A"
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(addon.ExtKeyInvoiceType, addon.InvoiceTypeDebitNote)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-18] discount deductions are only supported in invoices (FT, FS, FR)")
	})

	t.Run("deductions in work documents", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Series = "PF SERIES-A"
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Delete(addon.ExtKeyInvoiceType).Set(addon.ExtKeyWorkType, addon.WorkTypeProforma)
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-18]")
	})

	t.Run("missing reference with several preceding documents", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts[1].Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-DISCOUNT-04] ($.discounts[1].ext) discount deductions require 'pt-saft-reference' extension")
	})

	t.Run("missing reference with a single preceding document is normalized", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Preceding = inv.Preceding[:1]
		inv.Discounts = inv.Discounts[:1]
		inv.Discounts[0].Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("reference not in preceding documents", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts[1] = deduction(2000, "FT SERIES-A/3")
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-21] discount deduction 'pt-saft-reference' must match a preceding document")
	})

	t.Run("reference to a document that is not an invoice", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Preceding[1] = &org.DocumentRef{Series: "GT SERIES-A", Code: "1"}
		inv.Discounts[1] = deduction(2000, "GT SERIES-A/1")
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-DISCOUNT-03] ($.discounts[1].ext) invoice deducted preceding document must be an invoice (FT, FS, FR)")
	})

	t.Run("deduction without preceding documents", func(t *testing.T) {
		inv := validInvoice()
		inv.Discounts = []*bill.Discount{deduction(1000, "")}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-INVOICE-20] ($.preceding) discount deductions require a preceding document")
	})

	t.Run("deduction without VAT", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts[0].Taxes = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-DISCOUNT-01] ($.discounts[0].taxes) invoice discount deduction taxes must include VAT category")
	})

	t.Run("deduction with zero amount", func(t *testing.T) {
		inv := deductionInvoice()
		inv.Discounts[0].Amount = num.MakeAmount(0, 2)
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-PT-SAFT-BILL-DISCOUNT-02] ($.discounts[0].amount) invoice discount amount must be greater than 0")
	})
}
