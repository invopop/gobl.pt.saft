package saft

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/validation"
)

// Header constants
const (
	AuditFileVersion   = "1.04_01"
	CurrencyCode       = "EUR"
	TaxAccountingBasis = "F" // F - Faturação (Invoicing)
	TaxEntity          = "Global"
)

// AuditFile represents the root element of a SAF-T file
type AuditFile struct {
	XMLName struct{} `xml:"urn:OECD:StandardAuditFile-Tax:PT_1.04_01 AuditFile"`

	Header          *Header          `xml:"Header"`
	MasterFiles     *MasterFiles     `xml:"MasterFiles"`
	SourceDocuments *SourceDocuments `xml:"SourceDocuments,omitempty"`
}

// Header represents the header of a SAF-T file
type Header struct {
	AuditFileVersion          string   `xml:"AuditFileVersion"`
	CompanyID                 string   `xml:"CompanyID"`
	TaxRegistrationNumber     string   `xml:"TaxRegistrationNumber"`
	TaxAccountingBasis        string   `xml:"TaxAccountingBasis"`
	CompanyName               string   `xml:"CompanyName"`
	BusinessName              string   `xml:"BusinessName,omitempty"`
	CompanyAddress            *Address `xml:"CompanyAddress"`
	FiscalYear                string   `xml:"FiscalYear"`
	StartDate                 string   `xml:"StartDate"`
	EndDate                   string   `xml:"EndDate"`
	CurrencyCode              string   `xml:"CurrencyCode"`
	DateCreated               string   `xml:"DateCreated"`
	TaxEntity                 string   `xml:"TaxEntity"`
	ProductCompanyTaxID       string   `xml:"ProductCompanyTaxID"`
	SoftwareCertificateNumber int      `xml:"SoftwareCertificateNumber"`
	ProductID                 string   `xml:"ProductID"`
	ProductVersion            string   `xml:"ProductVersion"`
	HeaderComment             string   `xml:"HeaderComment,omitempty"`
	Telephone                 string   `xml:"Telephone,omitempty"`
	Fax                       string   `xml:"Fax,omitempty"`
	Email                     string   `xml:"Email,omitempty"`
	Website                   string   `xml:"Website,omitempty"`
}

// MasterFiles represents the master files of a SAF-T file
type MasterFiles struct {
	Customers []*Customer `xml:"Customer,omitempty"`
	Products  []*Product  `xml:"Product,omitempty"`
	TaxTable  *TaxTable   `xml:"TaxTable,omitempty"`
}

// SourceDocuments represents the source documents of a SAF-T file
type SourceDocuments struct {
	SalesInvoices    *SalesInvoices    `xml:"SalesInvoices,omitempty"`
	MovementOfGoods  *MovementOfGoods  `xml:"MovementOfGoods,omitempty"`
	WorkingDocuments *WorkingDocuments `xml:"WorkingDocuments,omitempty"`
	Payments         *Payments         `xml:"Payments,omitempty"`
}

// SalesInvoices represents the sales invoices of a SAF-T file
type SalesInvoices struct {
	NumberOfEntries int        `xml:"NumberOfEntries"`
	TotalDebit      string     `xml:"TotalDebit"`
	TotalCredit     string     `xml:"TotalCredit"`
	Invoices        []*Invoice `xml:"Invoice,omitempty"`
}

// Payments represents the payments of a SAF-T file
type Payments struct {
	NumberOfEntries int        `xml:"NumberOfEntries"`
	TotalDebit      string     `xml:"TotalDebit"`
	TotalCredit     string     `xml:"TotalCredit"`
	Payments        []*Payment `xml:"Payment,omitempty"`
}

// WorkingDocuments represents the working documents of a SAF-T file
type WorkingDocuments struct {
	NumberOfEntries int             `xml:"NumberOfEntries"`
	TotalDebit      string          `xml:"TotalDebit"`
	TotalCredit     string          `xml:"TotalCredit"`
	WorkDocuments   []*WorkDocument `xml:"WorkDocument,omitempty"`
}

// AuditFileParams provides the parameters to build a SAF-T audit file
type AuditFileParams struct {
	config.Config

	Issuer *org.Party
	Period cal.Period

	Products        []*Product
	Customers       []*Customer
	TaxTableEntries []*TaxTableEntry

	SourceDocuments []SourceDocument
}

// BuildAuditFile builds a SAF-T audit file from the given input parameters
func BuildAuditFile(p *AuditFileParams) (*AuditFile, error) {
	if p == nil {
		p = new(AuditFileParams)
	}

	if err := ValidateIssuer(p.Issuer); err != nil {
		return nil, fmt.Errorf("invalid supplier: %w", err)
	}

	if err := ValidatePeriod(p.Period); err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	sis := make([]*Invoice, 0)
	pys := make([]*Payment, 0)
	wds := make([]*WorkDocument, 0)
	mgd := make([]*StockMovement, 0)

	for _, sd := range p.SourceDocuments {
		switch doc := sd.(type) {
		case *Invoice:
			sis = append(sis, doc)
		case *Payment:
			pys = append(pys, doc)
		case *WorkDocument:
			wds = append(wds, doc)
		case *StockMovement:
			mgd = append(mgd, doc)
		default:
			return nil, fmt.Errorf("unknown source document type: %T", doc)
		}
	}

	return &AuditFile{
		Header: buildHeader(p.Config, p.Issuer, p.Period),
		MasterFiles: &MasterFiles{
			Products:  p.Products,
			Customers: p.Customers,
			TaxTable:  buildTaxTable(p.TaxTableEntries),
		},
		SourceDocuments: &SourceDocuments{
			SalesInvoices:    buildSalesInvoices(sis),
			MovementOfGoods:  buildMovementOfGoods(mgd),
			WorkingDocuments: buildWorkingDocuments(wds),
			Payments:         buildPayments(pys),
		},
	}, nil
}

// ValidatePeriod checks that a reporting period is complete and already elapsed.
func ValidatePeriod(per cal.Period) error {
	if err := rules.Validate(&per); err != nil {
		return err
	}
	today := config.TodayIn(pt.New().TimeLocation())
	return validation.ValidateStruct(&per,
		validation.Field(&per.Start, cal.DateBefore(today)),
		validation.Field(&per.End, cal.DateBefore(today)),
	)
}

func buildTaxTable(tes []*TaxTableEntry) *TaxTable {
	if len(tes) == 0 {
		return nil
	}

	return &TaxTable{TaxTableEntries: tes}
}

func buildSalesInvoices(sis []*Invoice) *SalesInvoices {
	if len(sis) == 0 {
		return nil
	}

	var tot DocTotals
	for _, inv := range sis {
		tot = tot.Merge(reportableDocTotals(inv))
	}

	return &SalesInvoices{
		NumberOfEntries: len(sis),
		TotalDebit:      tot.Debit.String(),
		TotalCredit:     tot.Credit.String(),
		Invoices:        sis,
	}
}

func buildPayments(pys []*Payment) *Payments {
	if len(pys) == 0 {
		return nil
	}

	var tot DocTotals
	for _, py := range pys {
		tot = tot.Merge(reportableDocTotals(py))
	}
	return &Payments{
		NumberOfEntries: len(pys),
		TotalDebit:      tot.Debit.String(),
		TotalCredit:     tot.Credit.String(),
		Payments:        pys,
	}
}

func buildWorkingDocuments(wds []*WorkDocument) *WorkingDocuments {
	if len(wds) == 0 {
		return nil
	}

	var tot DocTotals
	for _, wd := range wds {
		tot = tot.Merge(reportableDocTotals(wd))
	}

	return &WorkingDocuments{
		NumberOfEntries: len(wds),
		TotalDebit:      tot.Debit.String(),
		TotalCredit:     tot.Credit.String(),
		WorkDocuments:   wds,
	}
}

func buildMovementOfGoods(sms []*StockMovement) *MovementOfGoods {
	if len(sms) == 0 {
		return nil
	}

	var tot DocTotals
	lines := 0
	for _, sm := range sms {
		tot = tot.Merge(reportableDocTotals(sm))
		lines += len(sm.Lines)
	}

	return &MovementOfGoods{
		NumberOfMovementLines: lines,
		TotalQuantityIssued:   tot.Quantity.String(),
		StockMovements:        sms,
	}
}

// reportableDocTotals returns the document's totals, or zero totals when the document
// is excluded from the report sums because it is canceled or billed.
func reportableDocTotals(doc SourceDocument) DocTotals {
	if st := doc.GetStatus(); st == StatusCanceled || st == StatusBilled {
		return DocTotals{}
	}
	return ComputeDocTotals(doc)
}

func buildHeader(cfg config.Config, issuer *org.Party, per cal.Period) *Header {
	if issuer == nil {
		return nil
	}

	return &Header{
		AuditFileVersion:          AuditFileVersion,
		CompanyID:                 issuer.TaxID.Code.String(),
		TaxRegistrationNumber:     issuer.TaxID.Code.String(),
		TaxAccountingBasis:        TaxAccountingBasis,
		CompanyName:               fit(issuer.Name, 100),
		BusinessName:              fit(issuer.Alias, 60),
		CompanyAddress:            address(issuer),
		FiscalYear:                fiscalYear(per),
		StartDate:                 per.Start.String(),
		EndDate:                   per.End.String(),
		CurrencyCode:              CurrencyCode,
		DateCreated:               config.Now().Format("2006-01-02"),
		TaxEntity:                 TaxEntity,
		ProductCompanyTaxID:       fit(cfg.ProducerTaxID, 30),
		SoftwareCertificateNumber: cfg.CertificateNumber,
		ProductID:                 fit(cfg.ProductID(), 255),
		ProductVersion:            fit(cfg.AppVersion, 30),
		Telephone:                 telephone(issuer),
		Email:                     email(issuer),
		Website:                   website(issuer),
	}
}

func fiscalYear(per cal.Period) string {
	return strconv.Itoa(per.Start.Year)
}

func telephone(gisr *org.Party) string {
	if len(gisr.Telephones) == 0 {
		return ""
	}
	return fit(gisr.Telephones[0].Number, 20)
}

func email(gisr *org.Party) string {
	if len(gisr.Emails) == 0 {
		return ""
	}
	return fit(gisr.Emails[0].Address, 254)
}

func website(gisr *org.Party) string {
	if len(gisr.Websites) == 0 {
		return ""
	}
	return fit(gisr.Websites[0].URL, 60)
}

func address(gisr *org.Party) *Address {
	return BuildAddress(gisr)
}
