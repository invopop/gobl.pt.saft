package realtime

import (
	"errors"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/atws"
	"github.com/invopop/gobl/num"
)

// BuildRegisterInvoice builds a RegisterInvoiceRequest from a SAF-T AuditFile.
func BuildRegisterInvoice(af *saft.AuditFile) (*atws.RegisterInvoice, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.SalesInvoices == nil {
		return nil, errors.New("audit file must contain sales invoices")
	}

	if len(af.SourceDocuments.SalesInvoices.Invoices) == 0 {
		return nil, errors.New("audit file must contain at least one invoice")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	// Build the RegisterInvoiceRequest from SAF-T structures
	return buildRegisterInvoiceRequest(af), nil
}

func buildRegisterInvoiceRequest(af *saft.AuditFile) *atws.RegisterInvoice {
	hd := af.Header
	inv := af.SourceDocuments.SalesInvoices.Invoices[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.RegisterInvoice{
		EFaturaMDVersion:          atws.EFaturaMDVersion,
		AuditFileVersion:          hd.AuditFileVersion,
		TaxRegistrationNumber:     hd.TaxRegistrationNumber,
		TaxEntity:                 hd.TaxEntity,
		SoftwareCertificateNumber: hd.SoftwareCertificateNumber,
		InvoiceData:               buildInvoiceData(inv, cus),
	}
}

func buildInvoiceData(inv *saft.Invoice, cus *saft.Customer) *atws.InvoiceData {
	return &atws.InvoiceData{
		InvoiceNo:              inv.InvoiceNo,
		ATCUD:                  inv.ATCUD,
		InvoiceDate:            inv.InvoiceDate,
		InvoiceType:            inv.InvoiceType,
		SelfBillingIndicator:   selfBillingIndicator(inv),
		CustomerTaxID:          cus.CustomerTaxID,
		CustomerTaxIDCountry:   customerTaxIDCountry(cus),
		DocumentStatus:         buildInvoiceStatus(inv),
		HashCharacters:         saft.HashCharacters(inv.Hash),
		CashVATSchemeIndicator: cashVATSchemeIndicator(inv),
		PaperLessIndicator:     paperLessIndicator(inv),
		EACCode:                inv.EACCode,
		SystemEntryDate:        inv.SystemEntryDate,
		LineSummaries:          buildLineSummaries(inv.Lines),
		DocumentTotals:         buildDocumentTotals(inv.DocumentTotals),
		WithholdingTaxes:       buildWithholdingTaxes(inv.WithholdingTaxes),
	}
}

func buildInvoiceStatus(inv *saft.Invoice) *atws.InvoiceStatus {
	if inv.DocumentStatus == nil {
		return nil
	}
	return &atws.InvoiceStatus{
		InvoiceStatus:     inv.DocumentStatus.InvoiceStatus,
		InvoiceStatusDate: inv.DocumentStatus.InvoiceStatusDate,
	}
}

func buildLineSummaries(lines []*saft.Line) []*atws.LineSummary {
	if lines == nil {
		return nil
	}

	ls := make([]*atws.LineSummary, len(lines))
	for i, line := range lines {
		ls[i] = &atws.LineSummary{
			OrderReferences:      buildOrderReferences(line),
			TaxPointDate:         line.TaxPointDate,
			References:           buildReferences(line),
			DebitCreditIndicator: debitCreditIndicator(line),
			Amount:               amount(line),
			Tax:                  buildTax(line),
			TaxExemptionCode:     line.TaxExemptionCode,
		}
	}
	return compactLineSummaries(ls)
}

func buildReferences(line *saft.Line) []string {
	if line == nil || line.References == nil {
		return nil
	}

	refs := make([]string, len(line.References))
	for i, ref := range line.References {
		refs[i] = ref.Reference
	}
	return refs
}

func buildTax(line *saft.Line) *atws.Tax {
	if line == nil || line.Tax == nil {
		return nil
	}

	return &atws.Tax{
		TaxType:          line.Tax.TaxType,
		TaxCountryRegion: line.Tax.TaxCountryRegion,
		TaxCode:          line.Tax.TaxCode,
		TaxPercentage:    line.Tax.TaxPercentage,
	}
}

func buildDocumentTotals(dt *saft.DocumentTotals) *atws.DocumentTotals {
	if dt == nil {
		return nil
	}

	return &atws.DocumentTotals{
		TaxPayable: dt.TaxPayable,
		NetTotal:   dt.NetTotal,
		GrossTotal: dt.GrossTotal,
	}
}

func customerTaxIDCountry(cus *saft.Customer) string {
	if cus == nil || cus.BillingAddress == nil {
		return ""
	}
	return cus.BillingAddress.Country
}

func selfBillingIndicator(inv *saft.Invoice) int {
	if inv == nil || inv.SpecialRegimes == nil {
		return 0
	}
	return inv.SpecialRegimes.SelfBillingIndicator
}

func cashVATSchemeIndicator(inv *saft.Invoice) int {
	if inv == nil || inv.SpecialRegimes == nil {
		return 0
	}
	return inv.SpecialRegimes.CashVATSchemeIndicator
}

func paperLessIndicator(_ *saft.Invoice) int {
	// TODO: support other values
	return 0
}

func debitCreditIndicator(line *saft.Line) string {
	if line == nil {
		return ""
	}
	if line.DebitAmount != "" {
		return "D"
	}
	return "C"
}

func amount(line *saft.Line) string {
	if line == nil {
		return ""
	}
	if line.DebitAmount != "" {
		return line.DebitAmount
	}
	return line.CreditAmount
}

// compactLineSummaries compacts an array of line invoice summaries according to the specs
func compactLineSummaries(lines []*atws.LineSummary) []*atws.LineSummary {
	var clines []*atws.LineSummary

	for _, line := range lines {
		cline := findMatchingLineSummary(line, clines)
		if cline != nil {
			cline.Amount = sumAmounts(cline.Amount, line.Amount)
		} else {
			clines = append(clines, line)
		}
	}

	return clines
}

func findMatchingLineSummary(line *atws.LineSummary, clines []*atws.LineSummary) *atws.LineSummary {
	for _, cline := range clines {
		if matchingLineSummaries(line, cline) {
			return cline
		}
	}
	return nil
}

// matchingLineSummaries determines whether two lines must be compacted together
func matchingLineSummaries(l1, l2 *atws.LineSummary) bool {
	if len(l1.OrderReferences) != len(l2.OrderReferences) {
		return false
	}
	for i, or1 := range l1.OrderReferences {
		or2 := l2.OrderReferences[i]
		if or1.OriginatingON != or2.OriginatingON {
			return false
		}
	}
	if l1.TaxPointDate != l2.TaxPointDate {
		return false
	}
	if len(l1.References) != len(l2.References) {
		return false
	}
	for i, r1 := range l1.References {
		r2 := l2.References[i]
		if r1 != r2 {
			return false
		}
	}
	if l1.DebitCreditIndicator != l2.DebitCreditIndicator {
		return false
	}
	if l1.Tax.TaxType != l2.Tax.TaxType {
		return false
	}
	if l1.Tax.TaxCountryRegion != l2.Tax.TaxCountryRegion {
		return false
	}
	if l1.Tax.TaxCode != l2.Tax.TaxCode {
		return false
	}
	if l1.TaxExemptionCode != l2.TaxExemptionCode {
		return false
	}

	return true
}

func sumAmounts(a1, a2 string) string {
	na1, _ := num.AmountFromString(a1)
	na2, _ := num.AmountFromString(a2)
	nasum := na1.MatchPrecision(na2).Add(na2)
	return nasum.String()
}

// BuildChangeInvoiceStatus builds a ChangeInvoiceStatusRequest from a SAF-T AuditFile.
func BuildChangeInvoiceStatus(af *saft.AuditFile) (*atws.ChangeInvoiceStatus, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.SalesInvoices == nil {
		return nil, errors.New("audit file must contain sales invoices")
	}

	if len(af.SourceDocuments.SalesInvoices.Invoices) == 0 {
		return nil, errors.New("audit file must contain at least one invoice")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	// Build the ChangeInvoiceStatusRequest from SAF-T structures
	return buildChangeInvoiceStatusRequest(af), nil
}

func buildChangeInvoiceStatusRequest(af *saft.AuditFile) *atws.ChangeInvoiceStatus {
	hd := af.Header
	inv := af.SourceDocuments.SalesInvoices.Invoices[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.ChangeInvoiceStatus{
		EFaturaMDVersion:      atws.EFaturaMDVersion,
		TaxRegistrationNumber: hd.TaxRegistrationNumber,
		InvoiceHeader:         buildInvoiceHeader(inv, cus),
		InvoiceStatus:         buildNewInvoiceStatus(inv),
	}
}

func buildInvoiceHeader(inv *saft.Invoice, cus *saft.Customer) *atws.InvoiceHeader {
	return &atws.InvoiceHeader{
		InvoiceNo:            inv.InvoiceNo,
		ATCUD:                inv.ATCUD,
		InvoiceDate:          inv.InvoiceDate,
		InvoiceType:          inv.InvoiceType,
		SelfBillingIndicator: selfBillingIndicator(inv),
		CustomerTaxID:        cus.CustomerTaxID,
		CustomerTaxIDCountry: customerTaxIDCountry(cus),
	}
}

func buildNewInvoiceStatus(inv *saft.Invoice) *atws.NewInvoiceStatus {
	if inv.DocumentStatus == nil {
		return nil
	}
	return &atws.NewInvoiceStatus{
		InvoiceStatus:     inv.DocumentStatus.InvoiceStatus,
		InvoiceStatusDate: inv.DocumentStatus.InvoiceStatusDate,
	}
}
