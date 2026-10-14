package saft

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// WorkDocument represents a working document in a SAF-T report
type WorkDocument struct {
	DocumentNumber  string              `xml:"DocumentNumber"`
	ATCUD           string              `xml:"ATCUD"`
	DocumentStatus  *WorkStatus         `xml:"DocumentStatus"`
	Hash            string              `xml:"Hash,omitempty"`
	HashControl     string              `xml:"HashControl,omitempty"`
	Period          int                 `xml:"Period"`
	WorkDate        string              `xml:"WorkDate"`
	WorkType        string              `xml:"WorkType"`
	SourceID        string              `xml:"SourceID"`
	EACCode         string              `xml:"EACCode,omitempty"`
	SystemEntryDate string              `xml:"SystemEntryDate"`
	TransactionID   string              `xml:"TransactionID,omitempty"`
	CustomerID      string              `xml:"CustomerID"`
	Lines           []*Line             `xml:"Line"`
	DocumentTotals  *WorkDocumentTotals `xml:"DocumentTotals"`

	// store the original tax total for QR code generation
	GOBLTaxTotal *tax.Total `xml:"-"`
}

// WorkStatus represents the status of a working document
type WorkStatus struct {
	WorkStatus     string `xml:"WorkStatus"`
	WorkStatusDate string `xml:"WorkStatusDate"`
	Reason         string `xml:"Reason,omitempty"`
	SourceID       string `xml:"SourceID"`
	SourceBilling  string `xml:"SourceBilling"`
}

// WorkDocumentTotals represents the totals of a working document
type WorkDocumentTotals struct {
	TaxPayable string    `xml:"TaxPayable"`
	NetTotal   string    `xml:"NetTotal"`
	GrossTotal string    `xml:"GrossTotal"`
	Currency   *Currency `xml:"Currency,omitempty"`
}

// BuildWorkDocumentFromInvoice converts a GOBL invoice to a SAF-T work document
func BuildWorkDocumentFromInvoice(ginv *bill.Invoice, p *BuildParams) (*WorkDocument, error) {
	if ginv == nil {
		return nil, nil
	}

	// Run local validations not in the regime or the addon (e.g. for not supported fields)
	if err := validateInvoice(ginv); err != nil {
		return nil, err
	}
	if err := validateWorkLines(ginv.Lines); err != nil {
		return nil, err
	}

	if p == nil {
		p = new(BuildParams)
	}

	// Convert the invoice to EUR if not already in EUR. Keep the original invoice to build
	// the currency information.
	ginvOrig := ginv
	if ginv.Currency != currency.CodeEmpty && ginv.Currency != currency.EUR {
		var err error
		if ginv, err = ginv.ConvertInto(currency.EUR); err != nil {
			return nil, err
		}
	}

	// Build and return the SAF-T work document
	return &WorkDocument{
		DocumentNumber: DocNumber(ginv.Series, ginv.Code),
		ATCUD:          p.ATCUD,
		DocumentStatus: &WorkStatus{
			WorkStatus:     StatusNormal,
			WorkStatusDate: p.SystemEntryDate,
			SourceID:       p.SourceID,
			SourceBilling:  sourceBilling(ginv.Tax),
		},
		Hash:            p.Hash,
		HashControl:     hashControl(ginv.Tax, p),
		Period:          period(ginv.IssueDate),
		WorkDate:        date(&ginv.IssueDate),
		WorkType:        workDocType(ginv.Tax),
		SourceID:        fit(p.SourceID, 30),
		SystemEntryDate: p.SystemEntryDate,
		CustomerID:      customerID(p.Customer),
		Lines:           buildLines(ginv, p.Products),
		DocumentTotals:  buildWorkDocumentTotals(ginv, ginvOrig),

		GOBLTaxTotal: invoiceTaxTotal(ginv),
	}, nil
}

// BuildWorkDocumentFromOrder converts a GOBL order to a SAF-T work document
func BuildWorkDocumentFromOrder(gord *bill.Order, p *BuildParams) (*WorkDocument, error) {
	if gord == nil {
		return nil, nil
	}

	// Run local validations not in the regime or the addon (e.g. for not supported fields)
	if err := validateOrder(gord); err != nil {
		return nil, err
	}

	if p == nil {
		p = new(BuildParams)
	}

	// Convert the order to EUR if not already in EUR. Keep the original order to build
	// the currency information.
	gordOrig := gord
	if gord.Currency != currency.CodeEmpty && gord.Currency != currency.EUR {
		var err error
		if gord, err = gord.ConvertInto(currency.EUR); err != nil {
			return nil, err
		}
	}

	// Build and return the SAF-T work document
	return &WorkDocument{
		DocumentNumber: DocNumber(gord.Series, gord.Code),
		ATCUD:          p.ATCUD,
		DocumentStatus: &WorkStatus{
			WorkStatus:     StatusNormal,
			WorkStatusDate: p.SystemEntryDate,
			SourceID:       p.SourceID,
			SourceBilling:  sourceBilling(gord.Tax),
		},
		Hash:            p.Hash,
		HashControl:     hashControl(gord.Tax, p),
		Period:          period(gord.IssueDate),
		WorkDate:        date(&gord.IssueDate),
		WorkType:        workDocType(gord.Tax),
		SourceID:        fit(p.SourceID, 30),
		SystemEntryDate: p.SystemEntryDate,
		CustomerID:      customerID(p.Customer),
		Lines:           buildLinesFromOrder(gord, p.Products),
		DocumentTotals:  buildWorkDocumentTotalsFromOrder(gord, gordOrig),

		GOBLTaxTotal: orderTaxTotal(gord),
	}, nil
}

// ---- SourceDocument interface implementation ---- //

// GetDocumentClass returns the document class of the work document
func (w *WorkDocument) GetDocumentClass() DocumentClass {
	return DocumentClassWork
}

// GetStatus returns the status of the work document
func (w *WorkDocument) GetStatus() string {
	if w.DocumentStatus == nil {
		return ""
	}
	return w.DocumentStatus.WorkStatus
}

// GetStatusDate returns the status date of the work document
func (w *WorkDocument) GetStatusDate() string {
	if w.DocumentStatus == nil {
		return ""
	}
	return w.DocumentStatus.WorkStatusDate
}

// GetStatusReason returns the status reason of the work document
func (w *WorkDocument) GetStatusReason() string {
	if w.DocumentStatus == nil {
		return ""
	}
	return w.DocumentStatus.Reason
}

// SetStatus sets the status of the work document
func (w *WorkDocument) SetStatus(status, reason, sourceID string) {
	w.DocumentStatus = &WorkStatus{
		WorkStatus:     status,
		WorkStatusDate: Timestamp(),
		Reason:         fit(reason, 50),
		SourceID:       sourceID,
		SourceBilling:  w.sourceBilling(),
	}
}

// sourceBilling returns the document's current source billing, defaulting to
// "produced" when the document has no status yet.
func (w *WorkDocument) sourceBilling() string {
	if w.DocumentStatus != nil && w.DocumentStatus.SourceBilling != "" {
		return w.DocumentStatus.SourceBilling
	}
	return addon.SourceBillingProduced.String()
}

// GetDocumentNo returns the document number of the work document
func (w *WorkDocument) GetDocumentNo() string {
	return w.DocumentNumber
}

// GetGrossTotal returns the gross total of the work document
func (w *WorkDocument) GetGrossTotal() string {
	if w.DocumentTotals == nil {
		return ""
	}
	return w.DocumentTotals.GrossTotal
}

// GetIssueDate returns the issue date of the work document
func (w *WorkDocument) GetIssueDate() string {
	return w.WorkDate
}

// GetCustomerID returns the customer ID of the work document
func (w *WorkDocument) GetCustomerID() string {
	return w.CustomerID
}

// GetDocumentType returns the document type of the work document
func (w *WorkDocument) GetDocumentType() string {
	return w.WorkType
}

// GetTaxPayable returns the tax payable of the work document
func (w *WorkDocument) GetTaxPayable() string {
	if w.DocumentTotals == nil {
		return ""
	}
	return w.DocumentTotals.TaxPayable
}

// GetTaxTotal returns the tax total of the work document
func (w *WorkDocument) GetTaxTotal() *tax.Total {
	return w.GOBLTaxTotal
}

// GetDocumentTotals returns the document totals of the work document
func (w *WorkDocument) GetDocumentTotals() *DocumentTotals {
	return nil
}

// ---- //

func buildWorkDocumentTotals(ginv *bill.Invoice, ginvOrig *bill.Invoice) *WorkDocumentTotals {
	if ginv == nil || ginv.Totals == nil {
		return nil
	}

	return &WorkDocumentTotals{
		TaxPayable: taxPayable(ginv),
		NetTotal:   ginv.Totals.Total.String(),
		GrossTotal: ginv.Totals.TotalWithTax.RescaleUp(2).String(), // ensure 2 decimal places (required for sealing)
		Currency:   buildCurrency(ginvOrig),
	}
}

// workDocType extracts the work document type from a GOBL invoice
func workDocType(gtax *bill.Tax) string {
	if gtax == nil {
		return ""
	}

	return gtax.Ext.Get(addon.ExtKeyWorkType).String()
}

func buildWorkDocumentTotalsFromOrder(gord *bill.Order, gordOrig *bill.Order) *WorkDocumentTotals {
	if gord == nil || gord.Totals == nil {
		return nil
	}

	return &WorkDocumentTotals{
		TaxPayable: taxPayableFromOrder(gord),
		NetTotal:   gord.Totals.Total.String(),
		GrossTotal: gord.Totals.TotalWithTax.RescaleUp(2).String(), // ensure 2 decimal places (required for sealing)
		Currency:   buildCurrencyFromOrder(gordOrig),
	}
}

func buildCurrencyFromOrder(gord *bill.Order) *Currency {
	if gord == nil || gord.Currency == currency.CodeEmpty || gord.Currency == currency.EUR {
		return nil
	}

	rate := currency.MatchExchangeRate(gord.ExchangeRates, gord.Currency, currency.EUR)
	if rate == nil {
		return nil
	}

	return &Currency{
		CurrencyCode:   gord.Currency.String(),
		CurrencyAmount: gord.Totals.TotalWithTax.String(),
		ExchangeRate:   rate.Amount.String(),
	}
}

func buildLinesFromOrder(gord *bill.Order, prods []*Product) []*Line {
	if gord == nil || len(gord.Lines) != len(prods) {
		return nil
	}

	lines := make([]*Line, len(gord.Lines))

	for i, gl := range gord.Lines {
		prod := prods[i]
		lines[i] = buildLineFromOrder(gord, gl, prod)
	}

	return lines
}

func buildLineFromOrder(gord *bill.Order, gline *bill.Line, prod *Product) *Line {
	if gord == nil || gline == nil || gline.Item == nil || prod == nil {
		return nil
	}

	return &Line{
		LineNumber:         gline.Index,
		OrderReferences:    buildOrderReferences(gord.Preceding, gline),
		ProductCode:        fit(prod.ProductCode, 60),
		ProductDesc:        fit(prod.ProductDesc, 200),
		Quantity:           gline.Quantity.String(),
		UnitOfMeasure:      unitOfMeasure(gline),
		UnitPrice:          unitPrice(gline),
		TaxPointDate:       orderTaxPointDate(gord),
		Description:        fit(prod.ProductDesc, 200),
		CreditAmount:       orderCreditAmount(gline),
		Tax:                buildTax(gline),
		TaxExemptionReason: lineTaxExemptionReason(gline),
		TaxExemptionCode:   lineTaxExemptionCode(gline),
		SettlementAmount:   settlementAmount(gline),
	}
}

func taxPayableFromOrder(gord *bill.Order) string {
	if gord == nil || gord.Totals == nil || gord.Totals.Taxes == nil {
		return ""
	}

	vat := gord.Totals.Taxes.Category(tax.CategoryVAT)
	if vat == nil {
		return ""
	}

	return vat.Amount.String()
}

// TODO: some of the validations here belong in the regime or the addon, but to avoid
// breaking the existing integrations with IX, we'll keep them here for now.
func validateWorkLines(lines []*bill.Line) error {
	return validation.Validate(lines,
		validation.Each(
			validation.By(validateWorkLine),
			validation.Skip,
		),
	)
}

// validateWorkLine ensures a work document line carries the VAT details SAF-T
// requires. The addon enforces the same rules upstream, but callers may not
// re-run them, and unlike the SalesInvoices branch the WorkingDocuments branch
// of the XSD leaves the whole tax block optional and its exemption assertions
// commented out, so an incomplete line would be filed silently.
func validateWorkLine(val any) error {
	gline, _ := val.(*bill.Line)
	if gline == nil {
		return nil
	}

	if err := validateLine(gline); err != nil {
		return err
	}

	vat := gline.Taxes.Get(tax.CategoryVAT)
	if vat == nil {
		return errors.New("must include the VAT category")
	}

	if vat.Ext.Get(addon.ExtKeyTaxRate) != addon.TaxRateExempt {
		return nil
	}

	if lineTaxExemptionCode(gline) == "" {
		return fmt.Errorf("exempt VAT requires the '%s' extension", addon.ExtKeyExemption)
	}
	if lineTaxExemptionReason(gline) == "" {
		return errors.New("exempt VAT requires a note with the exemption reason")
	}

	return nil
}

func validateOrder(gord *bill.Order) error {
	return validation.ValidateStruct(gord,
		validation.Field(&gord.Supplier,
			validation.By(ValidateIssuer),
			validation.Skip,
		),
		validation.Field(&gord.Tax,
			validation.By(validateOrderTax),
			validation.Skip,
		),
		validation.Field(&gord.Discounts,
			validation.Empty.Error("not supported, use line discounts instead"),
			validation.Skip,
		),
		validation.Field(&gord.Charges,
			validation.Empty.Error("not supported"),
			validation.Skip,
		),
		validation.Field(&gord.Lines,
			validation.Each(
				validation.By(validateWorkLine),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateOrderTax(val any) error {
	gtax, _ := val.(*bill.Tax)
	if gtax == nil {
		return nil
	}

	return validation.ValidateStruct(gtax,
		validation.Field(&gtax.Ext,
			tax.ExtensionsExcludeCodes(addon.ExtKeySource, addon.SourceBillingIntegrated), // Not supported
			validation.Skip,
		),
	)
}

func orderTaxPointDate(gord *bill.Order) string {
	if gord == nil {
		return ""
	}

	return date(gord.ValueDate)
}

func orderCreditAmount(gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	return amountOrZero(gline.Total).String()
}

func orderTaxTotal(gord *bill.Order) *tax.Total {
	if gord == nil || gord.Totals == nil || gord.Totals.Taxes == nil {
		return nil
	}

	return gord.Totals.Taxes
}
