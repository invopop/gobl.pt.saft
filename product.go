package saft

import (
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
)

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
