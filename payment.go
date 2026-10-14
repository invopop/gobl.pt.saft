package saft

import (
	"errors"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Payment represents a payment in a SAF-T file
type Payment struct {
	PaymentRefNo     string            `xml:"PaymentRefNo"`
	ATCUD            string            `xml:"ATCUD"`
	Period           int               `xml:"Period,omitempty"`
	TransactionID    string            `xml:"TransactionID,omitempty"`
	TransactionDate  string            `xml:"TransactionDate"`
	PaymentType      string            `xml:"PaymentType"`
	Description      string            `xml:"Description,omitempty"`
	SystemID         string            `xml:"SystemID,omitempty"`
	DocumentStatus   *PaymentStatus    `xml:"DocumentStatus"`
	PaymentMethods   []*PaymentMethod  `xml:"PaymentMethod,omitempty"`
	SourceID         string            `xml:"SourceID"`
	SystemEntryDate  string            `xml:"SystemEntryDate"`
	CustomerID       string            `xml:"CustomerID"`
	Lines            []*PaymentLine    `xml:"Line"`
	DocumentTotals   *DocumentTotals   `xml:"DocumentTotals"`
	WithholdingTaxes []*WithholdingTax `xml:"WithholdingTax,omitempty"`

	// store the original tax total for QR code generation
	GOBLTaxTotal *tax.Total `xml:"-"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus struct {
	PaymentStatus     string `xml:"PaymentStatus"`
	PaymentStatusDate string `xml:"PaymentStatusDate"`
	Reason            string `xml:"Reason,omitempty"`
	SourceID          string `xml:"SourceID"`
	SourcePayment     string `xml:"SourcePayment"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	PaymentMechanism string `xml:"PaymentMechanism,omitempty"`
	PaymentAmount    string `xml:"PaymentAmount"`
	PaymentDate      string `xml:"PaymentDate"`
}

// PaymentLine represents a payment line
type PaymentLine struct {
	LineNumber         int                 `xml:"LineNumber"`
	SourceDocumentIDs  []*SourceDocumentID `xml:"SourceDocumentID"`
	SettlementAmount   string              `xml:"SettlementAmount,omitempty"`
	DebitAmount        string              `xml:"DebitAmount,omitempty"`
	CreditAmount       string              `xml:"CreditAmount,omitempty"`
	Tax                *PaymentTax         `xml:"Tax,omitempty"`
	TaxExemptionReason string              `xml:"TaxExemptionReason,omitempty"`
	TaxExemptionCode   string              `xml:"TaxExemptionCode,omitempty"`
}

// SourceDocumentID represents a source document ID of a payment line
type SourceDocumentID struct {
	OriginatingON string `xml:"OriginatingON"`
	InvoiceDate   string `xml:"InvoiceDate"`
	Description   string `xml:"Description,omitempty"`
}

// PaymentTax represents the tax info of a payment line
type PaymentTax struct {
	TaxType          string `xml:"TaxType"`
	TaxCountryRegion string `xml:"TaxCountryRegion"`
	TaxCode          string `xml:"TaxCode"`
	TaxPercentage    string `xml:"TaxPercentage,omitempty"`
	TaxAmount        string `xml:"TaxAmount,omitempty"`
}

// WithholdingTax represents a withholding tax of a payment
type WithholdingTax struct {
	WithholdingTaxType        string `xml:"WithholdingTaxType,omitempty"`
	WithholdingTaxDescription string `xml:"WithholdingTaxDescription,omitempty"`
	WithholdingTaxAmount      string `xml:"WithholdingTaxAmount"`
}

// BuildPayment builds a payment from a GOBL payment
func BuildPayment(gpmt *bill.Payment, p *BuildParams) (*Payment, error) {
	if gpmt == nil {
		return nil, nil
	}

	// Run local validations not in the regime or the addon (e.g. for not supported fields)
	if err := validatePayment(gpmt); err != nil {
		return nil, err
	}

	if p == nil {
		p = new(BuildParams)
	}

	gtt := mergeTaxTotals(gpmt)

	return &Payment{
		PaymentRefNo: DocNumber(gpmt.Series, gpmt.Code),
		ATCUD:        p.ATCUD,
		Period:       period(gpmt.IssueDate),
		// TODO: TransactionID
		TransactionDate: date(&gpmt.IssueDate),
		PaymentType:     gpmt.Ext.Get(addon.ExtKeyPaymentType).String(),
		Description:     description(gpmt), // TODO: review this mapping
		SystemID:        gpmt.GetUUID().String(),
		DocumentStatus: &PaymentStatus{
			PaymentStatus:     StatusNormal,
			PaymentStatusDate: p.SystemEntryDate,
			SourceID:          fit(p.SourceID, 30),
			SourcePayment:     sourcePayment(gpmt),
		},
		PaymentMethods:   buildPaymentMethods(gpmt),
		SourceID:         fit(p.SourceID, 30),
		SystemEntryDate:  p.SystemEntryDate,
		CustomerID:       customerID(p.Customer),
		Lines:            buildPaymentLines(gpmt),
		DocumentTotals:   buildPaymentDocumentTotals(gpmt, gtt),
		WithholdingTaxes: buildWithholdingTaxes(gtt, nil),

		GOBLTaxTotal: gtt,
	}, nil
}

// ---- SourceDocument interface implementation ---- //

// GetDocumentClass returns the document class of the payment
func (p *Payment) GetDocumentClass() DocumentClass {
	return DocumentClassPayment
}

// GetStatus returns the status of the payment
func (p *Payment) GetStatus() string {
	if p.DocumentStatus == nil {
		return ""
	}
	return p.DocumentStatus.PaymentStatus
}

// GetStatusDate returns the status date of the payment
func (p *Payment) GetStatusDate() string {
	if p.DocumentStatus == nil {
		return ""
	}
	return p.DocumentStatus.PaymentStatusDate
}

// GetStatusReason returns the status reason of the payment
func (p *Payment) GetStatusReason() string {
	if p.DocumentStatus == nil {
		return ""
	}
	return p.DocumentStatus.Reason
}

// SetStatus sets the status of the payment
func (p *Payment) SetStatus(status string, reason string, sourceID string) {
	p.DocumentStatus = &PaymentStatus{
		PaymentStatus:     status,
		Reason:            fit(reason, 50),
		PaymentStatusDate: Timestamp(),
		SourceID:          fit(sourceID, 30),
		SourcePayment:     addon.SourceBillingProduced.String(),
	}
}

// GetDocumentNo returns the document number of the payment
func (p *Payment) GetDocumentNo() string {
	return p.PaymentRefNo
}

// GetGrossTotal returns the gross total of the payment
func (p *Payment) GetGrossTotal() string {
	if p.DocumentTotals == nil {
		return ""
	}

	return p.DocumentTotals.GrossTotal
}

// GetIssueDate returns the issue date of the payment
func (p *Payment) GetIssueDate() string {
	return p.TransactionDate
}

// GetCustomerID returns the customer ID of the payment
func (p *Payment) GetCustomerID() string {
	return p.CustomerID
}

// GetDocumentType returns the document type of the payment
func (p *Payment) GetDocumentType() string {
	return p.PaymentType
}

// GetTaxPayable returns the tax payable of the payment
func (p *Payment) GetTaxPayable() string {
	if p.DocumentTotals == nil {
		return ""
	}

	return p.DocumentTotals.TaxPayable
}

// GetTaxTotal returns the tax total of the payment
func (p *Payment) GetTaxTotal() *tax.Total {
	return p.GOBLTaxTotal
}

// GetDocumentTotals returns the document totals of the payment
func (p *Payment) GetDocumentTotals() *DocumentTotals {
	return p.DocumentTotals
}

// ---- //

func validatePayment(gpmt *bill.Payment) error {
	return validation.ValidateStruct(gpmt,
		validation.Field(&gpmt.Supplier,
			validation.By(ValidateIssuer),
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&gpmt.Currency,
			validation.In(currency.EUR).Error("must be EUR"),
			validation.Skip,
		),
		validation.Field(&gpmt.Lines,
			validation.Each(
				validation.By(validatePaymentLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validatePaymentLine(val any) error {
	line, _ := val.(*bill.PaymentLine)
	if line == nil {
		return nil
	}

	// TODO: These should probably be moved to the addon, maybe as normalizations.
	// They ensure the provided calculations are consistent.
	if line.Tax == nil {
		return errors.New("tax not given")
	}

	vat := line.Tax.Category(tax.CategoryVAT)
	if vat == nil {
		return errors.New("category VAT not given")
	}

	if len(vat.Rates) == 0 {
		return errors.New("no category VAT rate given")
	}

	rate := vat.Rates[0]
	if !rate.Amount.Equals(vat.Amount) {
		return errors.New("rate amount does not match the category amount")
	}

	if !line.Tax.Sum.Equals(vat.Amount) {
		return errors.New("tax sum does not match the category amount")
	}

	percent := rate.Percent
	if percent == nil {
		percent = &num.PercentageZero
	}

	if !amountsWithin(percent.Of(rate.Base), rate.Amount, num.MakeAmount(1, 2)) {
		return errors.New("rate amount is not the percentage of the base")
	}

	if !amountsWithin(line.Amount.Divide(percent.Factor()), rate.Base, num.MakeAmount(1, 2)) {
		return errors.New("line amount is not consistent with the tax base and percentage")
	}

	return nil
}

func buildPaymentMethods(gpmt *bill.Payment) []*PaymentMethod {
	if gpmt == nil || len(gpmt.Methods) == 0 {
		return nil
	}

	pms := make([]*PaymentMethod, 0, len(gpmt.Methods))
	for _, m := range gpmt.Methods {
		if m == nil {
			continue
		}
		pms = append(pms, buildPaymentMethod(m, &gpmt.IssueDate))
	}

	return pms
}

func buildPaymentMethod(m *pay.Record, issueDate *cal.Date) *PaymentMethod {
	d := issueDate
	if m.Date != nil {
		d = m.Date
	}

	return &PaymentMethod{
		PaymentMechanism: m.Ext.Get(addon.ExtKeyPaymentMeans).String(),
		PaymentAmount:    m.Amount.String(),
		PaymentDate:      date(d),
	}
}

func buildPaymentLines(gpmt *bill.Payment) []*PaymentLine {
	if gpmt == nil {
		return nil
	}

	var lines []*PaymentLine
	for _, line := range gpmt.Lines {
		lines = append(lines, buildPaymentLine(line))
	}

	return lines
}

func buildPaymentLine(line *bill.PaymentLine) *PaymentLine {
	if line == nil {
		return nil
	}

	rate := paymentVATRate(line)

	return &PaymentLine{
		LineNumber:         line.Index,
		SourceDocumentIDs:  buildSourceDocumentIDs(line.Document),
		CreditAmount:       paymentCreditAmount(line, rate),
		DebitAmount:        paymentDebitAmount(line, rate),
		Tax:                buildPaymentTax(rate),
		TaxExemptionReason: paymentTaxExemptionReason(line, rate),
		TaxExemptionCode:   paymentTaxExemptionCode(rate),
	}
}

func paymentVATRate(line *bill.PaymentLine) *tax.RateTotal {
	if line == nil {
		return nil
	}

	if vat := vatCategory(line.Tax); vat != nil && len(vat.Rates) > 0 {
		// According to point 4.4.4.14.6. of Portaria nª 302/2016, each different
		// tax rate (even for the same document) must appear in different payment
		// line. The saft addon enforces only one rate is present in each line.
		return vat.Rates[0]
	}

	return nil
}

func paymentTaxExemptionCode(rate *tax.RateTotal) string {
	if rate == nil {
		return ""
	}

	return rate.Ext.Get(addon.ExtKeyExemption).String()
}

func paymentTaxExemptionReason(line *bill.PaymentLine, rate *tax.RateTotal) string {
	code := paymentTaxExemptionCode(rate)
	return taxExemptionReason(code, line.Notes)
}

func buildPaymentTax(rate *tax.RateTotal) *PaymentTax {
	if rate == nil {
		return nil
	}

	return &PaymentTax{
		TaxType:          TaxTypeVAT,
		TaxCountryRegion: rate.Ext.Get(pt.ExtKeyRegion).String(),
		TaxCode:          rate.Ext.Get(addon.ExtKeyTaxRate).String(),
		TaxPercentage:    taxPercentage(rate.Percent),
	}
}

func buildSourceDocumentIDs(doc *org.DocumentRef) []*SourceDocumentID {
	if doc == nil {
		return nil
	}

	return []*SourceDocumentID{
		{
			OriginatingON: DocNumber(doc.Series, doc.Code),
			InvoiceDate:   date(doc.IssueDate),
			Description:   fit(doc.Description, 200),
		},
	}
}

func buildPaymentDocumentTotals(gpmt *bill.Payment, gtt *tax.Total) *DocumentTotals {
	if gpmt == nil {
		return nil
	}

	return &DocumentTotals{
		TaxPayable: totalTaxPayable(gtt),
		NetTotal:   totalNetTotal(gpmt, gtt),
		GrossTotal: gpmt.Total.RescaleUp(2).String(), // ensure 2 decimal places (required for sealing)
	}
}

func description(gpmt *bill.Payment) string {
	if gpmt == nil {
		return ""
	}

	for _, n := range gpmt.Notes {
		if n != nil && n.Key == org.NoteKeyGeneral && n.Text != "" {
			return fit(n.Text, 200)
		}
	}

	// Fall back to the first method's description for backwards
	// compatibility with documents that used the legacy singular
	// `method.detail` field.
	if len(gpmt.Methods) > 0 && gpmt.Methods[0] != nil {
		return fit(gpmt.Methods[0].Description, 200)
	}

	return ""
}

func paymentCreditAmount(line *bill.PaymentLine, rate *tax.RateTotal) string {
	if line == nil || rate == nil || line.Refund {
		return ""
	}

	return rate.Base.String()
}

func paymentDebitAmount(line *bill.PaymentLine, rate *tax.RateTotal) string {
	if line == nil || rate == nil || !line.Refund {
		return ""
	}

	return rate.Base.String()
}

func totalTaxPayable(gtt *tax.Total) string {
	if gtt == nil {
		return ""
	}

	return vatAmount(gtt).String()
}

func totalNetTotal(gpmt *bill.Payment, gtt *tax.Total) string {
	if gpmt == nil || gtt == nil {
		return ""
	}

	vat := vatAmount(gtt)
	return gpmt.Total.Subtract(vat).String()
}

func vatAmount(gtt *tax.Total) num.Amount {
	vat := vatCategory(gtt)
	if vat == nil {
		return num.AmountZero
	}

	return vat.Amount
}

func vatCategory(gtt *tax.Total) *tax.CategoryTotal {
	if gtt == nil {
		return nil
	}

	return gtt.Category(tax.CategoryVAT)
}

func mergeTaxTotals(gpmt *bill.Payment) *tax.Total {
	var tt *tax.Total
	for _, pl := range gpmt.Lines {
		if pl.Tax == nil {
			continue
		}
		if tt == nil {
			tt = pl.Tax
		} else {
			tt = tt.Merge(pl.Tax)
		}
	}
	return tt
}

func sourcePayment(gpmt *bill.Payment) string {
	if gpmt == nil {
		return ""
	}

	return gpmt.Ext.Get(addon.ExtKeySource).String()
}

// amountsWithin checks if the difference between two amounts is within a certain
// precision. Used to compare amounts that may have been rounded differently but are still
// considered equal for the purposes of validation.
func amountsWithin(a1, a2, precision num.Amount) bool {
	diff := a1.Subtract(a2).Abs()
	return diff.Compare(precision) <= 0
}
