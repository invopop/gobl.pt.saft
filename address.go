package saft

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Address represents an address in the SAF-T file
type Address struct {
	BuildingNumber string `xml:"BuildingNumber,omitempty"`
	StreetName     string `xml:"StreetName,omitempty"`
	AddressDetail  string `xml:"AddressDetail"`
	City           string `xml:"City"`
	PostalCode     string `xml:"PostalCode"`
	Region         string `xml:"Region,omitempty"`
	Country        string `xml:"Country"`
}

// BuildAddress builds a SAF-T address from a GOBL address
func BuildAddress(gparty *org.Party) *Address {
	if gparty == nil {
		return nil
	}

	var gaddr *org.Address
	if len(gparty.Addresses) > 0 {
		gaddr = gparty.Addresses[0]
	}
	if gaddr == nil {
		gaddr = new(org.Address)
	}

	return &Address{
		BuildingNumber: fit(gaddr.Number, 10),
		StreetName:     fit(gaddr.Street, 200),
		AddressDetail:  addressDetail(gaddr),
		City:           city(gaddr),
		PostalCode:     postalCode(gaddr),
		Region:         fit(gaddr.Region, 50),
		Country:        country(gaddr, gparty.TaxID),
	}
}

func addressDetail(gaddr *org.Address) string {
	if gaddr == nil || gaddr.LineOne() == "" {
		return ValueUnknown
	}

	return fit(gaddr.LineOne(), 210)
}

func city(gaddr *org.Address) string {
	if gaddr == nil || gaddr.Locality == "" {
		return ValueUnknown
	}

	return fit(gaddr.Locality, 50)
}

func postalCode(gaddr *org.Address) string {
	if gaddr == nil || gaddr.Code == "" {
		return ValueUnknown
	}

	return fit(gaddr.Code.String(), 20)
}

func country(gaddr *org.Address, gtid *tax.Identity) string {
	// Use the address country if set
	if gaddr != nil && gaddr.Country != "" {
		return gaddr.Country.String()
	}

	// Use the tax ID country if set
	if gtid != nil {
		if c := isoCountryCode(gtid.Country); c != "" {
			return c
		}
	}

	// No country is set, return unknown value
	return ValueUnknown
}

func isoCountryCode(tc l10n.TaxCountryCode) string {
	// Return the code if it is a valid ISO code
	c := tc.Code()
	if c.ISO().Validate() == nil {
		return c.String()
	}

	// Try the alternative code (if any)
	def := l10n.Countries().Code(c)
	if def != nil && def.AltCode != l10n.CodeEmpty && def.AltCode.ISO().Validate() == nil {
		return def.AltCode.String()
	}

	// Tax country code doesn't have a related ISO code
	return ""
}
