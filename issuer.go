package saft

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// ValidateIssuer validates a GOBL party as a SAF-T issuer
func ValidateIssuer(val any) error {
	giss, _ := val.(*org.Party)
	if giss == nil {
		return nil
	}

	return validation.ValidateStruct(giss,
		validation.Field(&giss.Name, validation.Required),
		validation.Field(&giss.TaxID,
			validation.Required,
			validation.By(validateIssuerTaxID),
		),
		validation.Field(&giss.Addresses,
			validation.Required,
			validation.Each(validation.By(validateIssuerAddress)),
		),
	)
}

func validateIssuerTaxID(val any) error {
	gt, _ := val.(*tax.Identity)
	if gt == nil {
		return nil
	}

	return validation.ValidateStruct(gt,
		validation.Field(&gt.Country,
			validation.Required,
			validation.In(l10n.PT.Tax()).Error("must be PT"),
		),
		validation.Field(&gt.Code, validation.Required),
	)
}

func validateIssuerAddress(val any) error {
	ga, _ := val.(*org.Address)
	if ga == nil {
		return nil
	}

	return validation.ValidateStruct(ga,
		validation.Field(&ga.Street, validation.Required),
		validation.Field(&ga.Locality, validation.Required),
		validation.Field(&ga.Code, validation.Required),
		validation.Field(&ga.Country, validation.Required),
	)
}
