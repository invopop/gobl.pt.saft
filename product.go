package saft

import (
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
)

// deductionDefaultDesc is the product description of advance payment deductions without
// a reason.
const deductionDefaultDesc = "Adiantamento"

// Product represents a product or service in the SAF-T file
type Product struct {
	ProductType       string          `xml:"ProductType"`
	ProductCode       string          `xml:"ProductCode"`
	ProductGroup      string          `xml:"ProductGroup,omitempty"`
	ProductDesc       string          `xml:"ProductDescription"`
	ProductNumberCode string          `xml:"ProductNumberCode"`
	CustomsDetails    *CustomsDetails `xml:"CustomsDetails,omitempty"`
}

// CustomsDetails represents the customs details of a product
type CustomsDetails struct {
	CNCode   string `xml:"CNCode,omitempty"`
	UNNumber string `xml:"UNNumber,omitempty"`
}

// BuildProducts builds the list of products from GOBL lines
func BuildProducts(gls []*bill.Line) []*Product {
	if gls == nil {
		return nil
	}

	prods := make([]*Product, len(gls))
	for i, gl := range gls {
		prods[i] = BuildProduct(gl)
	}
	return prods
}

// BuildProduct builds a product from a GOBL invoice line
func BuildProduct(gline *bill.Line) *Product {
	if gline == nil || gline.Item == nil {
		return nil
	}

	p := &Product{
		ProductType: gline.Item.Ext.Get(addon.ExtKeyProductType).String(),
		ProductDesc: fit(gline.Item.Name, 200),
	}
	for _, id := range gline.Item.Identities {
		if id.Key == org.IdentityKeyEAN {
			p.ProductNumberCode = fit(id.Code.String(), 60)
		}
	}

	p.ProductCode = hashID(p)
	if p.ProductNumberCode == "" {
		p.ProductNumberCode = p.ProductCode
	}
	return p
}

// BuildDeductionProducts builds the list of products for the advance payment deductions of
// a GOBL invoice, in the same order as their SAF-T lines.
func BuildDeductionProducts(ginv *bill.Invoice) []*Product {
	ds := deductions(ginv)
	if len(ds) == 0 {
		return nil
	}

	prods := make([]*Product, len(ds))
	for i, d := range ds {
		prods[i] = BuildDeductionProduct(d)
	}
	return prods
}

// BuildDeductionProduct builds a product from an advance payment deduction, using the
// "other" product type that SAF-T reserves for received advances.
func BuildDeductionProduct(d *bill.Discount) *Product {
	if d == nil {
		return nil
	}

	desc := d.Reason
	if desc == "" {
		desc = deductionDefaultDesc
	}

	p := &Product{
		ProductType: addon.ProductTypeOther.String(),
		ProductDesc: fit(desc, 200),
	}
	p.ProductCode = hashID(p)
	p.ProductNumberCode = p.ProductCode
	return p
}
