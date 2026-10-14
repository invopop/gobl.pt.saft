package saft

import (
	"fmt"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
)

// TaxTable represents the tax table in the SAF-T file
type TaxTable struct {
	TaxTableEntries []*TaxTableEntry `xml:"TaxTableEntry"`
}

// TaxTableEntry represents an entry of that SAF-T tax table
type TaxTableEntry struct {
	TaxType          string `xml:"TaxType"`
	TaxCountryRegion string `xml:"TaxCountryRegion"`
	TaxCode          string `xml:"TaxCode"`
	Description      string `xml:"Description"`
	TaxPercentage    string `xml:"TaxPercentage,omitempty"`
	TaxAmount        string `xml:"TaxAmount,omitempty"`

	TaxTableEntryID string `xml:"-"` // Auxiliary field for unique identification
}

// BuildTaxTableEntries builds the tax table entries from GOBL tax totals
func BuildTaxTableEntries(gts *bill.Totals) []*TaxTableEntry {
	if gts == nil || gts.Taxes == nil {
		return nil
	}

	vat := gts.Taxes.Category(tax.CategoryVAT)
	if vat == nil {
		return nil
	}

	ttes := make([]*TaxTableEntry, len(vat.Rates))
	for i, rt := range vat.Rates {
		ttes[i] = BuildTaxTableEntry(rt)
	}

	return ttes
}

// BuildTaxTableEntry builds a tax table entry from a GOBL tax rate
func BuildTaxTableEntry(rt *tax.RateTotal) *TaxTableEntry {
	if rt == nil {
		return nil
	}

	tte := &TaxTableEntry{
		TaxType:          TaxTypeVAT,
		TaxCountryRegion: rt.Ext.Get(pt.ExtKeyRegion).String(),
		TaxCode:          rt.Ext.Get(addon.ExtKeyTaxRate).String(),
		TaxPercentage:    taxPercentage(rt.Percent),
	}
	tte.Description = fmt.Sprintf("%s-%s-%s", tte.TaxCode, tte.TaxCountryRegion, tte.TaxPercentage)
	tte.TaxTableEntryID = hashID(tte)
	return tte
}
