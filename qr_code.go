package saft

import (
	"fmt"
	"sort"
	"strings"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
)

// BuildQRCodeParams is the parameters for building a QR code
type BuildQRCodeParams struct {
	Supplier  *Supplier
	Customer  *Customer
	Document  SourceDocument
	ATCUD     string
	Signature *Signature
	SWCertNum int
}

// BuildQRCode builds a QR code for a given document
func BuildQRCode(qr *BuildQRCodeParams) string {
	parts := new(qrParts)

	parts.add("A", qr.Supplier.SupplierTaxID)
	parts.add("B", qr.Customer.CustomerTaxID)
	parts.add("C", qrCountry(qr.Customer))
	parts.add("D", qr.Document.GetDocumentType())
	parts.add("E", qr.Document.GetStatus())
	parts.add("F", strings.ReplaceAll(qr.Document.GetIssueDate(), "-", ""))
	parts.add("G", qr.Document.GetDocumentNo())
	parts.add("H", qr.ATCUD)

	vat := vatCategory(qr.Document.GetTaxTotal())

	if vat == nil {
		parts.add("I1", "0")
	} else {
		// Group VAT rates by region and tax rate
		groups := make(map[cbc.Code]map[cbc.Code]*tax.RateTotal)
		for _, r := range vat.Rates {
			tr := r.Ext.Get(addon.ExtKeyTaxRate)
			rg := r.Ext.Get(pt.ExtKeyRegion)

			if groups[rg] == nil {
				groups[rg] = make(map[cbc.Code]*tax.RateTotal)
			}

			groups[rg][tr] = addTaxTotals(groups[rg][tr], r)
		}

		// Set the QR keys from the groups
		keys := []string{"I", "J", "K"}
		for reg, rates := range groups {
			// Get the next key prefix
			if len(keys) == 0 {
				break
			}
			key := keys[0]
			keys = keys[1:]

			// Set the QR keys for the group
			parts.add(key+"1", reg)
			for code, rate := range rates {
				switch code {
				case addon.TaxRateExempt:
					parts.add(key+"2", amount(rate.Base))
				case addon.TaxRateReduced:
					parts.add(key+"3", amount(rate.Base))
					parts.add(key+"4", amount(rate.Amount))
				case addon.TaxRateIntermediate:
					parts.add(key+"5", amount(rate.Base))
					parts.add(key+"6", amount(rate.Amount))
				case addon.TaxRateNormal:
					parts.add(key+"7", amount(rate.Base))
					parts.add(key+"8", amount(rate.Amount))
				}
			}
		}
	}

	// TODO: L - Não sujeito / não tributável em IVA / outras situações
	// TODO: M - Imposto do Selo
	parts.add("N", qr.Document.GetTaxPayable())
	parts.add("O", qr.Document.GetGrossTotal())
	parts.add("Q", qr.Signature.HashCharacters())
	parts.add("R", qr.SWCertNum)

	return parts.encode()
}

type qrParts []string

func (qp *qrParts) add(k string, v any) {
	if v == nil || v == "" {
		return
	}
	*qp = append(*qp, fmt.Sprintf("%s:%v", k, v))
}

func (qp *qrParts) encode() string {
	sort.Strings(*qp)
	return strings.Join(*qp, "*")
}

func qrCountry(customer *Customer) string {
	if customer.BillingAddress.Country == ValueUnknown {
		return "PT" // QR code spec requires "PT" if unknown
	}
	return customer.BillingAddress.Country
}

func amount(a num.Amount) string {
	return a.Rescale(2).String() // as required by the QR code spec
}

func addTaxTotals(a, b *tax.RateTotal) *tax.RateTotal {
	if a == nil {
		a = new(tax.RateTotal)
	}

	if b == nil {
		b = new(tax.RateTotal)
	}

	return &tax.RateTotal{
		Base:   a.Base.MatchPrecision(b.Base).Add(b.Base),
		Amount: a.Amount.MatchPrecision(b.Amount).Add(b.Amount),
	}
}
