package realtime

import (
	"errors"
	"fmt"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/atws"
)

// BuildRegisterPayment builds a RegisterPaymentRequest from a SAF-T AuditFile.
func BuildRegisterPayment(af *saft.AuditFile) (*atws.RegisterPayment, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.Payments == nil {
		return nil, errors.New("audit file must contain payments")
	}

	if len(af.SourceDocuments.Payments.Payments) == 0 {
		return nil, errors.New("audit file must contain at least one payment")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	pmt := af.SourceDocuments.Payments.Payments[0]
	if pmt.PaymentType != "RC" {
		return nil, fmt.Errorf("payment type %s is not supported for real-time reporting", pmt.PaymentType)
	}

	// Build the RegisterPaymentRequest from SAF-T structures
	return buildRegisterPaymentRequest(af), nil
}

func buildRegisterPaymentRequest(af *saft.AuditFile) *atws.RegisterPayment {
	hd := af.Header
	pmt := af.SourceDocuments.Payments.Payments[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.RegisterPayment{
		EFaturaMDVersion:          atws.EFaturaMDVersion,
		AuditFileVersion:          hd.AuditFileVersion,
		TaxRegistrationNumber:     hd.TaxRegistrationNumber,
		TaxEntity:                 hd.TaxEntity,
		SoftwareCertificateNumber: hd.SoftwareCertificateNumber,
		PaymentData:               buildPaymentData(pmt, cus),
	}
}

func buildPaymentData(pmt *saft.Payment, cus *saft.Customer) *atws.PaymentData {
	return &atws.PaymentData{
		PaymentRefNo:         pmt.PaymentRefNo,
		ATCUD:                pmt.ATCUD,
		TransactionDate:      pmt.TransactionDate,
		PaymentType:          pmt.PaymentType,
		CustomerTaxID:        cus.CustomerTaxID,
		CustomerTaxIDCountry: customerTaxIDCountry(cus),
		DocumentStatus:       buildPaymentStatus(pmt),
		SystemEntryDate:      pmt.SystemEntryDate,
		LineSummaries:        buildPaymentLineSummaries(pmt.Lines),
		DocumentTotals:       buildDocumentTotals(pmt.DocumentTotals),
		WithholdingTaxes:     buildWithholdingTaxes(pmt.WithholdingTaxes),
	}
}

func buildPaymentStatus(pmt *saft.Payment) *atws.PaymentStatus {
	if pmt.DocumentStatus == nil {
		return nil
	}
	return &atws.PaymentStatus{
		PaymentStatus:     pmt.DocumentStatus.PaymentStatus,
		PaymentStatusDate: pmt.DocumentStatus.PaymentStatusDate,
	}
}

func buildPaymentLineSummaries(lines []*saft.PaymentLine) []*atws.PaymentLineSummary {
	if lines == nil {
		return nil
	}

	ls := make([]*atws.PaymentLineSummary, len(lines))
	for i, line := range lines {
		ls[i] = &atws.PaymentLineSummary{
			SourceDocumentIDs:    buildSourceDocumentIDs(line),
			SettlementAmount:     line.SettlementAmount,
			DebitCreditIndicator: debitCreditIndicatorForPayment(line),
			Amount:               amountForPayment(line),
			Tax:                  buildTaxForPayment(line),
			TaxExemptionCode:     line.TaxExemptionCode,
		}
	}
	return compactPaymentLineSummaries(ls)
}

func buildSourceDocumentIDs(line *saft.PaymentLine) []*atws.SourceDocumentID {
	if line == nil || line.SourceDocumentIDs == nil {
		return nil
	}

	ids := make([]*atws.SourceDocumentID, len(line.SourceDocumentIDs))
	for i, sd := range line.SourceDocumentIDs {
		ids[i] = &atws.SourceDocumentID{
			OriginatingON: sd.OriginatingON,
			InvoiceDate:   sd.InvoiceDate,
		}
	}
	return ids
}

func debitCreditIndicatorForPayment(line *saft.PaymentLine) string {
	if line == nil {
		return ""
	}
	if line.DebitAmount != "" {
		return "D"
	}
	return "C"
}

func amountForPayment(line *saft.PaymentLine) string {
	if line == nil {
		return ""
	}
	if line.DebitAmount != "" {
		return line.DebitAmount
	}
	return line.CreditAmount
}

func buildTaxForPayment(line *saft.PaymentLine) *atws.Tax {
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

func buildWithholdingTaxes(taxes []*saft.WithholdingTax) []*atws.WithholdingTax {
	if taxes == nil {
		return nil
	}

	wts := make([]*atws.WithholdingTax, len(taxes))
	for i, wt := range taxes {
		wts[i] = &atws.WithholdingTax{
			WithholdingTaxType:   wt.WithholdingTaxType,
			WithholdingTaxAmount: wt.WithholdingTaxAmount,
		}
	}
	return wts
}

func compactPaymentLineSummaries(lines []*atws.PaymentLineSummary) []*atws.PaymentLineSummary {
	var clines []*atws.PaymentLineSummary

	for _, line := range lines {
		cline := findMatchingPaymentLineSummary(line, clines)
		if cline != nil {
			cline.Amount = sumAmounts(cline.Amount, line.Amount)
		} else {
			clines = append(clines, line)
		}
	}
	return clines
}

func findMatchingPaymentLineSummary(line *atws.PaymentLineSummary, clines []*atws.PaymentLineSummary) *atws.PaymentLineSummary {
	for _, cline := range clines {
		if matchingPaymentLineSummaries(line, cline) {
			return cline
		}
	}
	return nil
}

func matchingPaymentLineSummaries(l1, l2 *atws.PaymentLineSummary) bool {
	if len(l1.SourceDocumentIDs) != len(l2.SourceDocumentIDs) {
		return false
	}
	for i, sd1 := range l1.SourceDocumentIDs {
		sd2 := l2.SourceDocumentIDs[i]
		if sd1.OriginatingON != sd2.OriginatingON {
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

// BuildChangePaymentStatus builds a ChangePaymentStatusRequest from a SAF-T AuditFile.
func BuildChangePaymentStatus(af *saft.AuditFile) (*atws.ChangePaymentStatus, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.Payments == nil {
		return nil, errors.New("audit file must contain payments")
	}

	if len(af.SourceDocuments.Payments.Payments) == 0 {
		return nil, errors.New("audit file must contain at least one payment")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	pmt := af.SourceDocuments.Payments.Payments[0]
	if pmt.PaymentType != "RC" {
		return nil, fmt.Errorf("payment type %s is not supported for real-time reporting", pmt.PaymentType)
	}

	// Build the ChangePaymentStatusRequest from SAF-T structures
	return buildChangePaymentStatusRequest(af), nil
}

func buildChangePaymentStatusRequest(af *saft.AuditFile) *atws.ChangePaymentStatus {
	hd := af.Header
	pmt := af.SourceDocuments.Payments.Payments[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.ChangePaymentStatus{
		EFaturaMDVersion:      atws.EFaturaMDVersion,
		TaxRegistrationNumber: hd.TaxRegistrationNumber,
		PaymentHeader:         buildPaymentHeader(pmt, cus),
		PaymentStatus:         buildNewPaymentStatus(pmt),
	}
}

func buildPaymentHeader(pmt *saft.Payment, cus *saft.Customer) *atws.PaymentHeader {
	return &atws.PaymentHeader{
		PaymentRefNo:         pmt.PaymentRefNo,
		ATCUD:                pmt.ATCUD,
		TransactionDate:      pmt.TransactionDate,
		PaymentType:          pmt.PaymentType,
		CustomerTaxID:        cus.CustomerTaxID,
		CustomerTaxIDCountry: customerTaxIDCountry(cus),
	}
}

func buildNewPaymentStatus(pmt *saft.Payment) *atws.NewPaymentStatus {
	if pmt.DocumentStatus == nil {
		return nil
	}
	return &atws.NewPaymentStatus{
		PaymentStatus:     pmt.DocumentStatus.PaymentStatus,
		PaymentStatusDate: pmt.DocumentStatus.PaymentStatusDate,
	}
}
