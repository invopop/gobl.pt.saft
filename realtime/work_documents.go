package realtime

import (
	"errors"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/atws"
)

// BuildRegisterWork builds a RegisterWorkRequest from a SAF-T AuditFile.
func BuildRegisterWork(af *saft.AuditFile) (*atws.RegisterWork, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.WorkingDocuments == nil {
		return nil, errors.New("audit file must contain working documents")
	}

	if len(af.SourceDocuments.WorkingDocuments.WorkDocuments) == 0 {
		return nil, errors.New("audit file must contain at least one work document")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	// Build the RegisterWorkRequest from SAF-T structures
	return buildRegisterWorkRequest(af), nil
}

func buildRegisterWorkRequest(af *saft.AuditFile) *atws.RegisterWork {
	hd := af.Header
	wd := af.SourceDocuments.WorkingDocuments.WorkDocuments[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.RegisterWork{
		EFaturaMDVersion:          atws.EFaturaMDVersion,
		AuditFileVersion:          hd.AuditFileVersion,
		TaxRegistrationNumber:     hd.TaxRegistrationNumber,
		TaxEntity:                 hd.TaxEntity,
		SoftwareCertificateNumber: hd.SoftwareCertificateNumber,
		WorkData:                  buildWorkData(wd, cus),
	}
}

func buildWorkData(wd *saft.WorkDocument, cus *saft.Customer) *atws.WorkData {
	return &atws.WorkData{
		DocumentNumber:       wd.DocumentNumber,
		ATCUD:                wd.ATCUD,
		WorkDate:             wd.WorkDate,
		WorkType:             wd.WorkType,
		CustomerTaxID:        cus.CustomerTaxID,
		CustomerTaxIDCountry: customerTaxIDCountry(cus),
		DocumentStatus:       buildWorkStatus(wd),
		HashCharacters:       saft.HashCharacters(wd.Hash),
		EACCode:              wd.EACCode,
		SystemEntryDate:      wd.SystemEntryDate,
		LineSummaries:        buildLineSummaries(wd.Lines),
		DocumentTotals:       buildDocumentTotalsFromWork(wd.DocumentTotals),
	}
}

func buildWorkStatus(wd *saft.WorkDocument) *atws.WorkStatus {
	if wd.DocumentStatus == nil {
		return nil
	}
	return &atws.WorkStatus{
		WorkStatus:     wd.DocumentStatus.WorkStatus,
		WorkStatusDate: wd.DocumentStatus.WorkStatusDate,
	}
}

func buildDocumentTotalsFromWork(wdt *saft.WorkDocumentTotals) *atws.DocumentTotals {
	if wdt == nil {
		return nil
	}

	return &atws.DocumentTotals{
		TaxPayable: wdt.TaxPayable,
		NetTotal:   wdt.NetTotal,
		GrossTotal: wdt.GrossTotal,
	}
}

// BuildChangeWorkStatus builds a ChangeWorkStatusRequest from a SAF-T AuditFile.
func BuildChangeWorkStatus(af *saft.AuditFile) (*atws.ChangeWorkStatus, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.WorkingDocuments == nil {
		return nil, errors.New("audit file must contain working documents")
	}

	if len(af.SourceDocuments.WorkingDocuments.WorkDocuments) == 0 {
		return nil, errors.New("audit file must contain at least one work document")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	// Build the ChangeWorkStatusRequest from SAF-T structures
	return buildChangeWorkStatusRequest(af), nil
}

func buildChangeWorkStatusRequest(af *saft.AuditFile) *atws.ChangeWorkStatus {
	hd := af.Header
	wd := af.SourceDocuments.WorkingDocuments.WorkDocuments[0]
	cus := af.MasterFiles.Customers[0]

	return &atws.ChangeWorkStatus{
		EFaturaMDVersion:      atws.EFaturaMDVersion,
		TaxRegistrationNumber: hd.TaxRegistrationNumber,
		WorkHeader:            buildWorkHeader(wd, cus),
		WorkStatus:            buildNewWorkStatus(wd),
	}
}

func buildWorkHeader(wd *saft.WorkDocument, cus *saft.Customer) *atws.WorkHeader {
	return &atws.WorkHeader{
		DocumentNumber:       wd.DocumentNumber,
		ATCUD:                wd.ATCUD,
		WorkDate:             wd.WorkDate,
		WorkType:             wd.WorkType,
		CustomerTaxID:        cus.CustomerTaxID,
		CustomerTaxIDCountry: customerTaxIDCountry(cus),
	}
}

func buildNewWorkStatus(wd *saft.WorkDocument) *atws.NewWorkStatus {
	if wd.DocumentStatus == nil {
		return nil
	}
	return &atws.NewWorkStatus{
		WorkStatus:     wd.DocumentStatus.WorkStatus,
		WorkStatusDate: wd.DocumentStatus.WorkStatusDate,
	}
}
