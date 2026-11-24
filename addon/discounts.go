package addon

import (
	"fmt"
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// DiscountKeyDeduction identifies a document discount that deducts an advance payment.
const DiscountKeyDeduction cbc.Key = "deduction"

// deductibleTypes are the document types that can deduct advance payments.
var deductibleTypes = []cbc.Code{
	InvoiceTypeStandard,
	InvoiceTypeSimplified,
	InvoiceTypeInvoiceReceipt,
}

// normalizeInvoiceDiscounts sets the reference of every deduction without one to the
// invoice's only preceding document.
func normalizeInvoiceDiscounts(inv *bill.Invoice) {
	if len(inv.Preceding) == 1 && inv.Preceding[0] != nil {
		ref := precedingDocNumber(inv.Preceding[0])
		for _, d := range inv.Discounts {
			if isDeduction(d) && !d.Ext.Has(ExtKeyReference) {
				d.Ext = d.Ext.Set(ExtKeyReference, ref)
			}
		}
	}
}

// billDiscountRules validates the document discounts.
func billDiscountRules() *rules.Set {
	return rules.For(new(bill.Discount),
		rules.When(is.Func("is deduction", discountIsDeduction),
			rules.Field("taxes",
				rules.Assert("01", "invoice discount deduction taxes must include VAT category", tax.SetHasCategory(tax.CategoryVAT)),
			),
			rules.Field("amount",
				rules.Assert("02", "invoice discount amount must be greater than 0", num.Positive),
			),
			rules.Field("ext",
				rules.Assert("04", fmt.Sprintf("discount deductions require '%s' extension", ExtKeyReference),
					tax.ExtensionsRequire(ExtKeyReference),
				),
				rules.Assert("03", "invoice deducted preceding document must be an invoice (FT, FS, FR)",
					is.Func("reference is deductible", extReferenceDeductible),
				),
			),
		),
	)
}

// invoiceHasDeductions returns true when the invoice has at least one deduction discount.
func invoiceHasDeductions(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && slices.ContainsFunc(inv.Discounts, isDeduction)
}

// invoiceIsDeductible returns true when the invoice's document type can deduct advance payments.
func invoiceIsDeductible(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && invoiceDocType(inv).In(deductibleTypes...)
}

// invoiceDeductionReferencesFound checks that every deduction reference matches one of
// the invoice's preceding documents.
func invoiceDeductionReferencesFound(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, d := range inv.Discounts {
		if !isDeduction(d) || !d.Ext.Has(ExtKeyReference) {
			continue
		}
		ref := d.Ext.Get(ExtKeyReference)
		if !slices.ContainsFunc(inv.Preceding, func(p *org.DocumentRef) bool {
			return p != nil && precedingDocNumber(p) == ref
		}) {
			return false
		}
	}
	return true
}

// extReferenceDeductible checks that the reference extension, when set, is the document
// number of an invoice that can be deducted.
func extReferenceDeductible(val any) bool {
	ext, ok := tax.ExtensionsFromValue(val)
	if !ok || !ext.Has(ExtKeyReference) {
		return true
	}
	return docNumberType(ext.Get(ExtKeyReference)).In(deductibleTypes...)
}

// discountIsDeduction returns true when the discount deducts an advance payment.
func discountIsDeduction(val any) bool {
	d, ok := val.(*bill.Discount)
	return ok && isDeduction(d)
}

// isDeduction returns true when the discount is an advance payment deduction.
func isDeduction(d *bill.Discount) bool {
	return d != nil && d.Key == DiscountKeyDeduction
}

// precedingDocNumber returns the SAF-T document number of a preceding document
// (e.g. "FT SERIES-A/1").
func precedingDocNumber(ref *org.DocumentRef) cbc.Code {
	return ref.Series.JoinWith("/", ref.Code)
}

// docNumberType returns the document type prefix of a SAF-T document number
// (e.g. "FT" for "FT SERIES-A/1").
func docNumberType(docNo cbc.Code) cbc.Code {
	typ, _, _ := strings.Cut(docNo.String(), " ")
	return cbc.Code(typ)
}
