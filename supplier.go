package saft

import "github.com/invopop/gobl/org"

// Supplier represents a supplier in a SAF-T document
type Supplier struct {
	SupplierTaxID string `xml:"SupplierTaxID"`
}

// BuildSupplier converts a GOBL party to a SAF-T supplier
func BuildSupplier(gsup *org.Party) *Supplier {
	if gsup == nil {
		return nil
	}

	return &Supplier{
		SupplierTaxID: supplierTaxID(gsup),
	}
}

func supplierTaxID(gpar *org.Party) string {
	if gpar == nil || gpar.TaxID == nil {
		return ""
	}

	return gpar.TaxID.Code.String()
}
