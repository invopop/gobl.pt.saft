package saft

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// transportDocTypes are the SAF-T document types that can appear as preceding
// references on a stock movement or as order references on an invoice.
var transportDocTypes = []cbc.Code{
	addon.MovementTypeDeliveryNote,
	addon.MovementTypeWaybill,
	addon.MovementTypeFixedAssets,
	addon.MovementTypeConsignment,
	addon.MovementTypeReturn,
}

// correctableTypes are the document types that can be corrected
var correctableTypes = []cbc.Code{
	addon.InvoiceTypeStandard,
	addon.InvoiceTypeSimplified,
	addon.InvoiceTypeInvoiceReceipt,
}

// Invoice represents an invoice in a SAF-T report
type Invoice struct {
	InvoiceNo        string            `xml:"InvoiceNo"`
	ATCUD            string            `xml:"ATCUD"`
	DocumentStatus   *InvoiceStatus    `xml:"DocumentStatus"`
	Hash             string            `xml:"Hash"`
	HashControl      string            `xml:"HashControl"`
	Period           int               `xml:"Period"`
	InvoiceDate      string            `xml:"InvoiceDate"`
	InvoiceType      string            `xml:"InvoiceType"`
	SpecialRegimes   *SpecialRegimes   `xml:"SpecialRegimes"`
	SourceID         string            `xml:"SourceID"`
	EACCode          string            `xml:"EACCode,omitempty"`
	SystemEntryDate  string            `xml:"SystemEntryDate"`
	CustomerID       string            `xml:"CustomerID"`
	Lines            []*Line           `xml:"Line"`
	DocumentTotals   *DocumentTotals   `xml:"DocumentTotals"`
	WithholdingTaxes []*WithholdingTax `xml:"WithholdingTax,omitempty"`

	// store the original tax total for QR code generation
	GOBLTaxTotal *tax.Total `xml:"-"`
}

// InvoiceStatus represents the status of an invoice
type InvoiceStatus struct {
	InvoiceStatus     string `xml:"InvoiceStatus"`
	InvoiceStatusDate string `xml:"InvoiceStatusDate"`
	Reason            string `xml:"Reason,omitempty"`
	SourceID          string `xml:"SourceID"`
	SourceBilling     string `xml:"SourceBilling"`
}

// DocumentTotals represents the totals of an invoice
type DocumentTotals struct {
	TaxPayable  string           `xml:"TaxPayable"`
	NetTotal    string           `xml:"NetTotal"`
	GrossTotal  string           `xml:"GrossTotal"`
	Currency    *Currency        `xml:"Currency,omitempty"`
	Settlements []*Settlement    `xml:"Settlement,omitempty"`
	Payments    []*PaymentMethod `xml:"Payment,omitempty"`
}

// SpecialRegimes represents the special regimes of an invoice
type SpecialRegimes struct {
	SelfBillingIndicator         int `xml:"SelfBillingIndicator"`
	CashVATSchemeIndicator       int `xml:"CashVATSchemeIndicator"`
	ThirdPartiesBillingIndicator int `xml:"ThirdPartiesBillingIndicator"`
}

// Line represents a line in an invoice
type Line struct {
	LineNumber          int                 `xml:"LineNumber"`
	OrderReferences     []*OrderReference   `xml:"OrderReferences,omitempty"`
	ProductCode         string              `xml:"ProductCode"`
	ProductDesc         string              `xml:"ProductDescription"`
	Quantity            string              `xml:"Quantity"`
	UnitOfMeasure       string              `xml:"UnitOfMeasure"`
	UnitPrice           string              `xml:"UnitPrice"`
	TaxBase             string              `xml:"TaxBase,omitempty"`
	TaxPointDate        string              `xml:"TaxPointDate,omitempty"`
	References          []*References       `xml:"References,omitempty"`
	Description         string              `xml:"Description"`
	DebitAmount         string              `xml:"DebitAmount,omitempty"`
	CreditAmount        string              `xml:"CreditAmount,omitempty"`
	Tax                 *Tax                `xml:"Tax"`
	TaxExemptionReason  string              `xml:"TaxExemptionReason,omitempty"`
	TaxExemptionCode    string              `xml:"TaxExemptionCode,omitempty"`
	SettlementAmount    string              `xml:"SettlementAmount,omitempty"`
	ProductSerialNumber string              `xml:"ProductSerialNumber,omitempty"`
	CustomsInformation  *CustomsInformation `xml:"CustomsInformation,omitempty"`
}

// OrderReference represents an order reference in a line
type OrderReference struct {
	OriginatingON string `xml:"OriginatingON,omitempty"`
	OrderDate     string `xml:"OrderDate,omitempty"`
}

// References represents a reference of a line
type References struct {
	Reference string `xml:"Reference,omitempty"`
	Reason    string `xml:"Reason,omitempty"`
}

// Tax represents the tax information of a line
type Tax struct {
	TaxType          string `xml:"TaxType"`
	TaxCountryRegion string `xml:"TaxCountryRegion,omitempty"`
	TaxCode          string `xml:"TaxCode,omitempty"`
	TaxPercentage    string `xml:"TaxPercentage"`
}

// Currency represents the currency information of an invoice
type Currency struct {
	CurrencyCode   string `xml:"CurrencyCode"`
	CurrencyAmount string `xml:"CurrencyAmount"`
	ExchangeRate   string `xml:"ExchangeRate"`
}

// CustomsInformation represents customs information for a line
type CustomsInformation struct {
	ARCNo     []string `xml:"ARCNo,omitempty"`
	IECAmount string   `xml:"IECAmount,omitempty"`
}

// Settlement represents a settlement of an invoice
type Settlement struct {
	SettlementDiscount string `xml:"SettlementDiscount,omitempty"`
	SettlementAmount   string `xml:"SettlementAmount,omitempty"`
	SettlementDate     string `xml:"SettlementDate,omitempty"`
	PaymentTerms       string `xml:"PaymentTerms,omitempty"`
}

// BuildInvoice converts a GOBL invoice to a SAF-T invoice
func BuildInvoice(ginv *bill.Invoice, p *BuildParams) (*Invoice, error) {
	if ginv == nil {
		return nil, nil
	}

	// Run local validations not in the regime or the addon (e.g. for not supported fields)
	if err := validateInvoice(ginv); err != nil {
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

	gtt := invoiceTaxTotal(ginv)

	// Build and return the SAF-T invoice
	return &Invoice{
		InvoiceNo: DocNumber(ginv.Series, ginv.Code),
		ATCUD:     p.ATCUD,
		DocumentStatus: &InvoiceStatus{
			InvoiceStatus:     StatusNormal,
			InvoiceStatusDate: p.SystemEntryDate,
			SourceID:          fit(p.SourceID, 30),
			SourceBilling:     sourceBilling(ginv.Tax),
		},
		Hash:        p.Hash,
		HashControl: hashControl(ginv.Tax, p),
		Period:      period(ginv.IssueDate),
		InvoiceDate: date(&ginv.IssueDate),
		InvoiceType: invoiceType(ginv),
		SpecialRegimes: &SpecialRegimes{
			CashVATSchemeIndicator: cashVATScheme(ginv),
		},
		SourceID:         fit(p.SourceID, 30),
		SystemEntryDate:  p.SystemEntryDate,
		CustomerID:       customerID(p.Customer),
		Lines:            buildLines(ginv, p.Products),
		DocumentTotals:   buildDocumentTotals(ginv, ginvOrig),
		WithholdingTaxes: buildWithholdingTaxes(gtt, ginv.Tax),

		GOBLTaxTotal: gtt,
	}, nil
}

// ---- SourceDocument interface implementation ---- //

// GetDocumentClass returns the document class of the invoice
func (i *Invoice) GetDocumentClass() DocumentClass {
	return DocumentClassInvoice
}

// GetStatus returns the status of the invoice
func (i *Invoice) GetStatus() string {
	if i.DocumentStatus == nil {
		return ""
	}
	return i.DocumentStatus.InvoiceStatus
}

// SetStatus sets the status of the invoice
func (i *Invoice) SetStatus(status, reason, sourceID string) {
	i.DocumentStatus = &InvoiceStatus{
		InvoiceStatus:     status,
		InvoiceStatusDate: Timestamp(),
		Reason:            fit(reason, 50),
		SourceID:          fit(sourceID, 30),
		SourceBilling:     addon.SourceBillingProduced.String(),
	}
}

// GetStatusDate returns the status date of the invoice
func (i *Invoice) GetStatusDate() string {
	if i.DocumentStatus == nil {
		return ""
	}
	return i.DocumentStatus.InvoiceStatusDate
}

// GetStatusReason returns the status reason of the invoice
func (i *Invoice) GetStatusReason() string {
	if i.DocumentStatus == nil {
		return ""
	}
	return i.DocumentStatus.Reason
}

// GetDocumentNo returns the document number of the invoice
func (i *Invoice) GetDocumentNo() string {
	return i.InvoiceNo
}

// GetGrossTotal returns the gross total of the invoice
func (i *Invoice) GetGrossTotal() string {
	if i.DocumentTotals == nil {
		return ""
	}
	return i.DocumentTotals.GrossTotal
}

// GetIssueDate returns the issue date of the invoice
func (i *Invoice) GetIssueDate() string {
	return i.InvoiceDate
}

// GetCustomerID returns the customer ID of the invoice
func (i *Invoice) GetCustomerID() string {
	return i.CustomerID
}

// GetDocumentType returns the document type of the invoice
func (i *Invoice) GetDocumentType() string {
	return i.InvoiceType
}

// GetTaxPayable returns the tax payable of the invoice
func (i *Invoice) GetTaxPayable() string {
	if i.DocumentTotals == nil {
		return ""
	}
	return i.DocumentTotals.TaxPayable
}

// GetTaxTotal returns the tax total of the invoice
func (i *Invoice) GetTaxTotal() *tax.Total {
	return i.GOBLTaxTotal
}

// GetDocumentTotals returns the document totals of the invoice
func (i *Invoice) GetDocumentTotals() *DocumentTotals {
	return i.DocumentTotals
}

// ---- //

// TODO: some of the validations here belong in the regime or the addon, but to avoid
// breaking the existing integrations with IX, we'll keep them here for now.
func validateInvoice(ginv *bill.Invoice) error {
	if ginv == nil {
		return nil
	}

	return validation.ValidateStruct(ginv,
		validation.Field(&ginv.Tax,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&ginv.Supplier,
			validation.By(validateSupplier),
			validation.By(ValidateIssuer),
			validation.Skip,
		),
		validation.Field(&ginv.Preceding,
			validation.When(!isCorrection(ginv),
				// The addon already ensures there's at most one
				validation.Each(
					validation.By(validateTransportPreceding),
					validation.Skip,
				),
				validation.Skip,
			),
			validation.When(isCorrection(ginv),
				validation.Required,
				validation.Length(1, 1).Error("only one preceding document is allowed"),
				validation.Each(
					validation.By(validateCorrectionPreceding),
					validation.Skip,
				),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&ginv.Discounts,
			validation.Empty.Error("not supported, use line discounts instead"),
			validation.Skip,
		),
		validation.Field(&ginv.Charges,
			validation.Empty.Error("not supported"),
			validation.Skip,
		),
		validation.Field(&ginv.Lines,
			validation.Each(
				validation.By(validateLine),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&ginv.Totals,
			validation.By(validateTotals(ginv)),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(val any) error {
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

func validateSupplier(val any) error {
	supplier, _ := val.(*org.Party)
	if supplier == nil {
		return nil
	}

	return validation.ValidateStruct(supplier,
		// Only commercial companies must show office and capital
		validation.Field(&supplier.Registration,
			validation.When(isCommercialCompany(supplier),
				validation.Required,
				validation.By(validateSupplierRegistration),
			),
			validation.Skip,
		),
	)
}

func validateSupplierRegistration(val any) error {
	reg, _ := val.(*org.Registration)
	if reg == nil {
		return nil
	}

	return validation.ValidateStruct(reg,
		validation.Field(&reg.Office, validation.Required),
		validation.Field(&reg.Capital, validation.Required),
	)
}

func validateCorrectionPreceding(val any) error {
	ref, _ := val.(*org.DocumentRef)
	if ref == nil {
		return nil
	}

	// Validate the document type is correctable
	docNo := DocNumber(ref.Series, ref.Code)
	docType, _, _ := DocNumberParts(docNo)
	if !docType.In(correctableTypes...) {
		return fmt.Errorf("document type (%s) is not supported", docType)
	}

	return validation.ValidateStruct(ref,
		validation.Field(&ref.Reason,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateTransportPreceding(val any) error {
	ref, _ := val.(*org.DocumentRef)
	if ref == nil {
		return nil
	}

	// Validate the document type is a transport document
	docNo := DocNumber(ref.Series, ref.Code)
	docType, _, _ := DocNumberParts(docNo)
	if !docType.In(transportDocTypes...) {
		return fmt.Errorf("document type (%s) is not supported", docType)
	}

	return nil
}

func validateLine(val any) error {
	gline, _ := val.(*bill.Line)
	if gline == nil {
		return nil
	}

	return validation.ValidateStruct(gline,
		validation.Field(&gline.Charges,
			validation.Empty.Error("not supported"),
			validation.Skip,
		),
	)
}

func validateTotals(ginv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		totals, _ := val.(*bill.Totals)
		if totals == nil {
			return nil
		}

		return validation.ValidateStruct(totals,
			validation.Field(&totals.Payable,
				validation.By(validateTotalsPayable(ginv)),
				validation.Skip,
			),
		)
	}
}

func validateTotalsPayable(ginv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		payable, ok := val.(num.Amount)
		if !ok {
			return nil
		}

		if ginv.Tax == nil {
			return nil
		}

		if ginv.Tax.Ext.Get(addon.ExtKeyInvoiceType) != addon.InvoiceTypeSimplified {
			// Only simplified invoices are affected by the following rules.
			return nil
		}

		// Simplified invoices under 1000 EUR are allowed.
		// Source: alínea b) do n.º 1 do artigo 40.º do CIVA
		if payable.Compare(num.MakeAmount(1000, 0)) != 1 {
			return nil
		}

		return fmt.Errorf("too high according to Artigo 40.º do CIVA")
	}
}

func buildDocumentTotals(ginv *bill.Invoice, ginvOrig *bill.Invoice) *DocumentTotals {
	if ginv == nil || ginv.Totals == nil {
		return nil
	}

	return &DocumentTotals{
		TaxPayable:  taxPayable(ginv),
		NetTotal:    ginv.Totals.Total.String(),
		GrossTotal:  ginv.Totals.TotalWithTax.RescaleUp(2).String(), // ensure 2 decimal places (required for sealing)
		Currency:    buildCurrency(ginvOrig),
		Payments:    buildInvoicePayments(ginv),
		Settlements: buildInvoiceSettlements(ginv),
	}
}

func buildCurrency(ginv *bill.Invoice) *Currency {
	if ginv == nil || ginv.Currency == currency.CodeEmpty || ginv.Currency == currency.EUR {
		return nil
	}

	rate := currency.MatchExchangeRate(ginv.ExchangeRates, ginv.Currency, currency.EUR)
	if rate == nil {
		return nil
	}

	return &Currency{
		CurrencyCode:   ginv.Currency.String(),
		CurrencyAmount: ginv.Totals.TotalWithTax.String(),
		ExchangeRate:   rate.Amount.String(),
	}
}

func buildLines(ginv *bill.Invoice, prods []*Product) []*Line {
	if ginv == nil || len(ginv.Lines) != len(prods) {
		return nil
	}

	lines := make([]*Line, len(ginv.Lines))

	for i, gl := range ginv.Lines {
		prod := prods[i]
		lines[i] = buildLine(ginv, gl, prod)
	}

	return lines
}

func buildLine(ginv *bill.Invoice, gline *bill.Line, prod *Product) *Line {
	if gline == nil || gline.Item == nil || prod == nil {
		return nil
	}

	return &Line{
		LineNumber:         gline.Index,
		OrderReferences:    buildOrderReferences(ginv.Preceding, gline),
		ProductCode:        fit(prod.ProductCode, 60),
		ProductDesc:        fit(prod.ProductDesc, 200),
		Quantity:           gline.Quantity.String(),
		UnitOfMeasure:      unitOfMeasure(gline),
		UnitPrice:          unitPrice(gline),
		TaxPointDate:       taxPointDate(ginv),
		References:         buildReferences(ginv),
		Description:        fit(prod.ProductDesc, 200),
		DebitAmount:        debitAmount(ginv, gline),
		CreditAmount:       creditAmount(ginv, gline),
		Tax:                buildTax(gline),
		TaxExemptionReason: lineTaxExemptionReason(gline),
		TaxExemptionCode:   lineTaxExemptionCode(gline),
		SettlementAmount:   settlementAmount(gline),
	}
}

func buildOrderReferences(gprec []*org.DocumentRef, gline *bill.Line) []*OrderReference {
	var refs []*OrderReference
	refs = appendOrderReferencesFromLine(refs, gline)
	refs = appendOrderReferencesFromPreceding(refs, gprec)
	return refs
}

// appendOrderReferencesFromLine build and appends the per-line order reference
func appendOrderReferencesFromLine(refs []*OrderReference, gline *bill.Line) []*OrderReference {
	if gline == nil || gline.Order == "" {
		return refs
	}

	return append(refs, &OrderReference{
		OriginatingON: gline.Order.String(),
	})
}

// appendOrderReferencesFromPreceding builds and appends OrderReference from preceding information.
func appendOrderReferencesFromPreceding(refs []*OrderReference, gprec []*org.DocumentRef) []*OrderReference {
	for _, ref := range gprec {
		if ref == nil {
			continue
		}
		docNo := DocNumber(ref.Series, ref.Code)
		docType, _, _ := DocNumberParts(docNo)
		if !docType.In(transportDocTypes...) {
			continue
		}
		refs = append(refs, &OrderReference{
			OriginatingON: docNo,
			OrderDate:     date(ref.IssueDate),
		})
	}
	return refs
}

func buildReferences(ginv *bill.Invoice) []*References {
	if ginv == nil || len(ginv.Preceding) == 0 {
		return nil
	}

	var refs []*References
	for _, ref := range ginv.Preceding {
		docNo := DocNumber(ref.Series, ref.Code)
		docType, _, _ := DocNumberParts(docNo)
		if docType.In(transportDocTypes...) {
			continue // transport docs go to OrderReferences, not References
		}
		refs = append(refs, &References{
			Reference: docNo,
			Reason:    fit(ref.Reason, 50),
		})
	}

	return refs
}

func buildTax(gline *bill.Line) *Tax {
	if gline == nil {
		return nil
	}

	vat := gline.Taxes.Get(tax.CategoryVAT)
	if vat == nil {
		return nil
	}

	return &Tax{
		TaxType:          TaxTypeVAT,
		TaxCountryRegion: vat.Ext.Get(pt.ExtKeyRegion).String(),
		TaxCode:          vat.Ext.Get(addon.ExtKeyTaxRate).String(),
		TaxPercentage:    taxPercentage(vat.Percent),
	}
}

func buildWithholdingTaxes(tt *tax.Total, bt *bill.Tax) []*WithholdingTax {
	if tt == nil {
		return nil
	}

	// Map each retained category code to its tax.Note text (the description).
	descs := make(map[string]string)
	if bt != nil {
		for _, n := range bt.Notes {
			if n != nil {
				descs[n.Category.String()] = n.Text
			}
		}
	}

	var wts []*WithholdingTax
	for _, ct := range tt.Categories {
		if !ct.Retained {
			continue
		}
		wts = append(wts, &WithholdingTax{
			WithholdingTaxType:        ct.Code.String(),
			WithholdingTaxDescription: fit(descs[ct.Code.String()], 60),
			WithholdingTaxAmount:      ct.Amount.RescaleUp(2).String(), // ensure 2 decimal places
		})
	}
	return wts
}

func buildInvoicePayments(ginv *bill.Invoice) []*PaymentMethod {
	if ginv == nil || ginv.Payment == nil {
		return nil
	}

	pms := make([]*PaymentMethod, len(ginv.Payment.Advances))
	for i, rec := range ginv.Payment.Advances {
		pms[i] = buildInvoicePayment(rec)
	}

	return pms
}

func buildInvoicePayment(rec *pay.Record) *PaymentMethod {
	if rec == nil {
		return nil
	}

	return &PaymentMethod{
		PaymentMechanism: rec.Ext.Get(addon.ExtKeyPaymentMeans).String(),
		PaymentAmount:    rec.Amount.String(),
		PaymentDate:      date(rec.Date),
	}
}

func buildInvoiceSettlements(ginv *bill.Invoice) []*Settlement {
	if ginv == nil || ginv.Payment == nil || ginv.Payment.Terms == nil {
		return nil
	}

	settlements := make([]*Settlement, len(ginv.Payment.Terms.DueDates))
	for i, d := range ginv.Payment.Terms.DueDates {
		settlements[i] = &Settlement{
			PaymentTerms: paymentTerms(d),
		}
	}

	return settlements
}

func invoiceType(ginv *bill.Invoice) string {
	if ginv == nil || ginv.Tax == nil {
		return ""
	}

	return ginv.Tax.Ext.Get(addon.ExtKeyInvoiceType).String()
}

func cashVATScheme(ginv *bill.Invoice) int {
	if ginv == nil || ginv.Tax == nil {
		return 0
	}

	code := ginv.Tax.Ext.Get(addon.ExtKeyCashVAT).String()
	val, err := strconv.Atoi(code)
	if err != nil {
		return 0
	}
	return val
}

func customerID(c *Customer) string {
	if c == nil {
		return ""
	}

	return fit(c.CustomerID, 30)
}

func unitOfMeasure(gline *bill.Line) string {
	if gline == nil || gline.Item == nil {
		return ""
	}

	return fit(string(gline.Item.Unit), 20)
}

func unitPrice(gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	// The SAF-T PT spec requires the UnitPrice to be the price per unit with
	// discounts applied but not taxes. GOBL doesn't provide this value directly, so
	// we calculate it.
	if len(gline.Discounts) == 0 {
		return amountOrZero(gline.Item.Price).String()
	}
	if gline.Total == nil {
		return ""
	}
	return gline.Total.Upscale(2).Divide(gline.Quantity).String()
}

func taxPointDate(ginv *bill.Invoice) string {
	if ginv == nil {
		return ""
	}

	return date(ginv.ValueDate)
}

func debitAmount(ginv *bill.Invoice, gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	if !isCredited(ginv) {
		return amountOrZero(gline.Total).String()
	}

	return ""
}

func creditAmount(ginv *bill.Invoice, gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	if isCredited(ginv) {
		return amountOrZero(gline.Total).String()
	}

	return ""
}

func isCredited(ginv *bill.Invoice) bool {
	if ginv == nil || ginv.Tax == nil {
		return false
	}

	// All invoice types except credit notes and consignment credits are credited to the customer
	return ginv.Tax.Ext.Get(addon.ExtKeyInvoiceType) != addon.InvoiceTypeCreditNote &&
		ginv.Tax.Ext.Get(addon.ExtKeyWorkType) != addon.WorkTypeConsignmentCredit
}

func taxPercentage(gper *num.Percentage) string {
	if gper == nil {
		return "0"
	}

	return gper.StringWithoutSymbol()
}

func lineTaxExemptionCode(gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	vat := gline.Taxes.Get(tax.CategoryVAT)
	if vat == nil {
		return ""
	}

	return taxExemptionCode(vat.Ext)
}

func lineTaxExemptionReason(gline *bill.Line) string {
	code := lineTaxExemptionCode(gline)
	return taxExemptionReason(code, gline.Notes)
}

func taxExemptionCode(ext tax.Extensions) string {
	return ext.Get(addon.ExtKeyExemption).String()
}

func taxExemptionReason(code string, notes []*org.Note) string {
	ers := taxExemptionReasons(notes)

	r := ers[code]

	if len(r) > 60 {
		// Truncate to max length allowed by the SAF-T PT spec
		r = r[:60]
	}

	return fit(r, 60)
}

func taxExemptionReasons(notes []*org.Note) map[string]string {
	rs := make(map[string]string)
	for _, note := range notes {
		if note.Key == org.NoteKeyLegal && note.Src == addon.ExtKeyExemption {
			rs[note.Code.String()] = note.Text
		}
	}
	return rs
}

func settlementAmount(gline *bill.Line) string {
	if gline == nil || len(gline.Discounts) == 0 {
		return ""
	}

	total := num.AmountZero
	for _, d := range gline.Discounts {
		total = total.MatchPrecision(d.Amount).Add(d.Amount)
	}

	return total.String()
}

func taxPayable(ginv *bill.Invoice) string {
	if ginv == nil || ginv.Totals == nil || ginv.Totals.Taxes == nil {
		return ""
	}

	vat := ginv.Totals.Taxes.Category(tax.CategoryVAT)
	if vat == nil {
		return ""
	}

	return vat.Amount.String()
}

func amountOrZero(ga *num.Amount) num.Amount {
	if ga == nil {
		return num.AmountZero
	}

	return *ga
}

func invoiceTaxTotal(ginv *bill.Invoice) *tax.Total {
	if ginv == nil || ginv.Totals == nil || ginv.Totals.Taxes == nil {
		return nil
	}

	return ginv.Totals.Taxes
}

func paymentTerms(d *pay.DueDate) string {
	if d == nil {
		return ""
	}

	terms := d.Notes
	if terms == "" {
		terms = fmt.Sprintf("%s EUR por pagar com vencimento em %s", d.Amount.String(), d.Date)
	}

	return fit(terms, 100)
}

// sourceBilling returns the document's SAF-T SourceBilling. The addon requires
// and defaults the extension, so an empty result means the document never went
// through normalization and the XSD will reject it.
func sourceBilling(gtax *bill.Tax) string {
	if gtax == nil {
		return ""
	}

	return gtax.Ext.Get(addon.ExtKeySource).String()
}

// hashControl appends the source reference of a recovered or manually issued
// document to the signing key version, as the SAFPTHashControl pattern requires.
func hashControl(gtax *bill.Tax, p *BuildParams) string {
	hc := p.SigningKeyVersion
	if gtax != nil && gtax.Ext.Has(addon.ExtKeySourceRef) {
		hc = fmt.Sprintf("%s-%s", hc, gtax.Ext.Get(addon.ExtKeySourceRef))
	}
	return hc
}

func isCorrection(ginv *bill.Invoice) bool {
	return ginv != nil && ginv.Tax != nil && ginv.Tax.Ext.Get(addon.ExtKeyInvoiceType).In(
		addon.InvoiceTypeCreditNote,
		addon.InvoiceTypeDebitNote,
	)
}

// isCommercialCompany reports whether the supplier is a PT pessoa coletiva (NIF prefix 5)
func isCommercialCompany(supplier *org.Party) bool {
	return supplier != nil &&
		supplier.TaxID != nil &&
		supplier.TaxID.Country == l10n.PT.Tax() &&
		strings.HasPrefix(supplier.TaxID.Code.String(), "5")
}
