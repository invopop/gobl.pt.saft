package addon_test

import (
	"testing"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SAF-T (PT) reports monetary values with the currency's precision, so a
// converter will normally run the document through `RoundToCurrency` before
// serialising it. These tests pin down that the addon's normalizers and rules
// survive that second calculation pass, and that amounts the AT receives fit
// in two decimal places.

// awkwardInvoice builds an invoice whose line amounts do not fit in the
// currency's precision under the default `precise` rounding rule.
func awkwardInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := validInvoice()
	inv.Lines = []*bill.Line{
		{
			Quantity: num.MakeAmount(3, 0),
			Item: &org.Item{
				Name:  "Awkward Item",
				Price: num.NewAmount(10005, 3), // 10.005
				Unit:  org.UnitHour,
			},
			Taxes: tax.Set{
				{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
			},
		},
		{
			Quantity: num.MakeAmount(7, 0),
			Item: &org.Item{
				Name:  "Another Awkward Item",
				Price: num.NewAmount(3333, 3), // 3.333
				Unit:  org.UnitKilogram,
			},
			Taxes: tax.Set{
				{Category: tax.CategoryVAT, Rate: tax.RateIntermediate},
			},
		},
	}
	return inv
}

// assertCurrencyPrecision checks that no amount carries more decimals than
// the currency supports, which is what the AT will accept.
func assertCurrencyPrecision(t *testing.T, cur currency.Code, amounts ...*num.Amount) {
	t.Helper()
	def := cur.Def()
	require.NotNil(t, def)
	for _, a := range amounts {
		if a == nil {
			continue
		}
		assert.LessOrEqual(t, a.Exp(), def.Subunits,
			"amount %s exceeds the precision of %s", a.String(), cur)
	}
}

func TestRoundingPreciseCalculates(t *testing.T) {
	inv := awkwardInvoice(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))

	// The default `precise` rule keeps the extra precision in the line
	// totals, which is fine for GOBL but not for the AT.
	assert.Equal(t, "30.015", inv.Lines[0].Sum.String())
	assert.Equal(t, "23.331", inv.Lines[1].Sum.String())
}

func TestRoundingToCurrency(t *testing.T) {
	inv := awkwardInvoice(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.RoundToCurrency())
	require.NoError(t, rules.Validate(inv))

	for _, l := range inv.Lines {
		assertCurrencyPrecision(t, currency.EUR, l.Sum, l.Total)
	}
	assertCurrencyPrecision(t, currency.EUR,
		&inv.Totals.Sum, &inv.Totals.Total, &inv.Totals.TotalWithTax, &inv.Totals.Payable,
	)
	for _, cat := range inv.Totals.Taxes.Categories {
		for _, r := range cat.Rates {
			assertCurrencyPrecision(t, currency.EUR, &r.Base, &r.Amount)
		}
	}

	// Unit prices keep their precision: the AT accepts more decimals on
	// `UnitPrice` than on the amounts derived from it.
	assert.Equal(t, "10.005", inv.Lines[0].Item.Price.String())
}

func TestRoundingToCurrencyKeepsExtensions(t *testing.T) {
	inv := awkwardInvoice(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.RoundToCurrency())
	require.NoError(t, rules.Validate(inv))

	// The second calculation pass must leave the SAF-T codes in place.
	vat := inv.Lines[0].Taxes.Get(tax.CategoryVAT)
	require.NotNil(t, vat)
	assert.Equal(t, cbc.Code("NOR"), vat.Ext.Get(addon.ExtKeyTaxRate))
	assert.Equal(t, cbc.Code("PT"), vat.Ext.Get("pt-region"))
	assert.Equal(t, org.UnitHour, inv.Lines[0].Item.Unit)
	assert.Equal(t, cbc.Code("S"), inv.Lines[0].Item.Ext.Get(addon.ExtKeyProductType))
}

func TestRoundingToCurrencyKeepsExemptionNotes(t *testing.T) {
	inv := awkwardInvoice(t)
	inv.Lines = inv.Lines[:1]
	inv.Lines[0].Taxes = tax.Set{
		{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyExemption: "M07",
			}),
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	require.Len(t, inv.Lines[0].Notes, 1)

	// Normalization runs again inside RoundToCurrency, and must not append a
	// second exemption note nor drop the one already there.
	require.NoError(t, inv.RoundToCurrency())
	require.NoError(t, rules.Validate(inv))
	require.Len(t, inv.Lines[0].Notes, 1)
	assert.Equal(t, cbc.Code("M07"), inv.Lines[0].Notes[0].Code)
}

func TestRoundingCurrencyRuleWithPricesIncluded(t *testing.T) {
	inv := awkwardInvoice(t)
	inv.Tax.PricesInclude = tax.CategoryVAT
	inv.Tax.Rounding = tax.RoundingRuleCurrency
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))

	// Bases, tax amounts and the document totals must all add up.
	vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
	require.NotNil(t, vat)
	sum := num.MakeAmount(0, 2)
	for _, r := range vat.Rates {
		assertCurrencyPrecision(t, currency.EUR, &r.Base, &r.Amount)
		sum = sum.Add(r.Base)
	}
	assert.Equal(t, inv.Totals.Total.String(), sum.String())
	assert.Equal(t,
		inv.Totals.Total.Add(vat.Amount).String(),
		inv.Totals.TotalWithTax.String(),
	)
}
