package saft

import (
	"errors"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// MovementOfGoods represents the movement of goods in a SAF-T file
type MovementOfGoods struct {
	NumberOfMovementLines int              `xml:"NumberOfMovementLines"`
	TotalQuantityIssued   string           `xml:"TotalQuantityIssued"`
	StockMovements        []*StockMovement `xml:"StockMovement,omitempty"`
}

// StockMovement represents a stock movement in a SAF-T file
type StockMovement struct {
	DocumentNumber    string          `xml:"DocumentNumber"`
	ATCUD             string          `xml:"ATCUD"`
	DocumentStatus    *MovementStatus `xml:"DocumentStatus"`
	Hash              string          `xml:"Hash"`
	HashControl       string          `xml:"HashControl"`
	Period            int             `xml:"Period,omitempty"`
	MovementDate      string          `xml:"MovementDate"`
	MovementType      string          `xml:"MovementType"`
	SystemEntryDate   string          `xml:"SystemEntryDate"`
	TransactionID     string          `xml:"TransactionID,omitempty"`
	CustomerID        string          `xml:"CustomerID,omitempty"`
	SupplierID        string          `xml:"SupplierID,omitempty"`
	SourceID          string          `xml:"SourceID"`
	EACCode           string          `xml:"EACCode,omitempty"`
	MovementComments  string          `xml:"MovementComments,omitempty"`
	ShipTo            *ShippingPoint  `xml:"ShipTo,omitempty"`
	ShipFrom          *ShippingPoint  `xml:"ShipFrom,omitempty"`
	MovementEndTime   string          `xml:"MovementEndTime,omitempty"`
	MovementStartTime string          `xml:"MovementStartTime"`
	ATDocCodeID       string          `xml:"ATDocCodeID,omitempty"`
	Lines             []*Line         `xml:"Line"`
	DocumentTotals    *DocumentTotals `xml:"DocumentTotals"`

	// store the original tax total for QR code generation
	GOBLTaxTotal *tax.Total `xml:"-"`
}

// MovementStatus represents the status of a stock movement
type MovementStatus struct {
	MovementStatus     string `xml:"MovementStatus"`
	MovementStatusDate string `xml:"MovementStatusDate"`
	Reason             string `xml:"Reason,omitempty"`
	SourceID           string `xml:"SourceID"`
	SourceBilling      string `xml:"SourceBilling"`
}

// ShippingPoint represents a shipping point in a SAF-T file
type ShippingPoint struct {
	DeliveryID   []string `xml:"DeliveryID,omitempty"`
	DeliveryDate string   `xml:"DeliveryDate,omitempty"`
	WarehouseID  string   `xml:"WarehouseID,omitempty"`
	LocationID   string   `xml:"LocationID,omitempty"`
	Address      *Address `xml:"Address,omitempty"`
}

// BuildStockMovement converts a GOBL delivery to a SAF-T stock movement
func BuildStockMovement(gdlv *bill.Delivery, p *BuildParams) (*StockMovement, error) {
	if gdlv == nil {
		return nil, nil
	}

	// Run local validations not in the regime or the addon (e.g. for not supported fields)
	if err := validateDelivery(gdlv); err != nil {
		return nil, err
	}

	if p == nil {
		p = new(BuildParams)
	}

	// Build and return the SAF-T stock movement
	return &StockMovement{
		DocumentNumber:    DocNumber(gdlv.Series, gdlv.Code),
		ATCUD:             p.ATCUD,
		DocumentStatus:    buildMovementStatus(StatusNormal, "", p.SystemEntryDate, p.SourceID),
		Hash:              p.Hash,
		HashControl:       p.SigningKeyVersion,
		Period:            period(gdlv.IssueDate),
		MovementDate:      date(&gdlv.IssueDate),
		MovementType:      movementType(gdlv),
		SystemEntryDate:   p.SystemEntryDate,
		CustomerID:        customerID(p.Customer),
		SourceID:          fit(p.SourceID, 30),
		MovementComments:  movementComments(gdlv),
		ShipTo:            buildShippingPoint(gdlv.Receiver, gdlv.ReceiveDate),
		ShipFrom:          buildShippingPoint(gdlv.Despatcher, gdlv.DespatchDate),
		MovementEndTime:   movementEndTime(gdlv),
		MovementStartTime: movementStartTime(gdlv),
		Lines:             buildMovementLines(gdlv, p.Products),
		DocumentTotals:    buildMovementDocumentTotals(gdlv),

		GOBLTaxTotal: deliveryTaxTotal(gdlv),
	}, nil
}

// ---- SourceDocument interface implementation ---- //

// GetDocumentClass returns the document class of the stock movement
func (sm *StockMovement) GetDocumentClass() DocumentClass {
	return DocumentClassTransport
}

// GetStatus returns the status of the stock movement
func (sm *StockMovement) GetStatus() string {
	if sm.DocumentStatus == nil {
		return ""
	}
	return sm.DocumentStatus.MovementStatus
}

// GetStatusDate returns the status date of the stock movement
func (sm *StockMovement) GetStatusDate() string {
	if sm.DocumentStatus == nil {
		return ""
	}
	return sm.DocumentStatus.MovementStatusDate
}

// GetStatusReason returns the status reason of the stock movement
func (sm *StockMovement) GetStatusReason() string {
	if sm.DocumentStatus == nil {
		return ""
	}
	return sm.DocumentStatus.Reason
}

// SetStatus sets the status of the stock movement
func (sm *StockMovement) SetStatus(status, reason, sourceID string) {
	sm.DocumentStatus = buildMovementStatus(status, reason, Timestamp(), sourceID)
}

// GetDocumentNo returns the document number of the stock movement
func (sm *StockMovement) GetDocumentNo() string {
	return sm.DocumentNumber
}

// GetGrossTotal returns the gross total of the stock movement
func (sm *StockMovement) GetGrossTotal() string {
	if sm.DocumentTotals == nil {
		return ""
	}
	return sm.DocumentTotals.GrossTotal
}

// GetIssueDate returns the issue date of the stock movement
func (sm *StockMovement) GetIssueDate() string {
	return sm.MovementDate
}

// GetCustomerID returns the customer ID of the stock movement
func (sm *StockMovement) GetCustomerID() string {
	return sm.CustomerID
}

// GetDocumentType returns the document type of the stock movement
func (sm *StockMovement) GetDocumentType() string {
	return sm.MovementType
}

// GetTaxPayable returns the tax payable of the stock movement
func (sm *StockMovement) GetTaxPayable() string {
	if sm.DocumentTotals == nil {
		return ""
	}
	return sm.DocumentTotals.TaxPayable
}

// GetTaxTotal returns the tax total of the stock movement
func (sm *StockMovement) GetTaxTotal() *tax.Total {
	return sm.GOBLTaxTotal
}

// GetDocumentTotals returns the document totals of the stock movement
func (sm *StockMovement) GetDocumentTotals() *DocumentTotals {
	return sm.DocumentTotals
}

// ---- //

func validateDelivery(gdlv *bill.Delivery) error {
	return validation.ValidateStruct(gdlv,
		validation.Field(&gdlv.Supplier,
			validation.By(ValidateIssuer),
			validation.Skip,
		),
		validation.Field(&gdlv.Tax,
			validation.By(validateDeliveryTax),
			validation.Skip,
		),
		validation.Field(&gdlv.Currency,
			validation.In(currency.EUR).Error("must be EUR"),
			validation.Skip,
		),
		validation.Field(&gdlv.Discounts,
			validation.Empty.Error("not supported, use line discounts instead"),
			validation.Skip,
		),
		validation.Field(&gdlv.Charges,
			validation.Empty.Error("not supported"),
			validation.Skip,
		),
		validation.Field(&gdlv.Lines,
			validation.Each(
				validation.By(validateLine),
				validation.Skip,
			),
			validation.Skip,
		),

		// The validations below are not required by the SAF-T specification, but
		// they are by the real-time web service (TODO: move them to the addon)
		validation.Field(&gdlv.DespatchDate,
			validation.By(validateDespatchDate(gdlv)),
			validation.Skip,
		),
		validation.Field(&gdlv.ReceiveDate,
			cal.DateAfter(dateOrZero(gdlv.DespatchDate)),
			validation.Skip,
		),
		validation.Field(&gdlv.Despatcher,
			validation.Required,
			validation.By(validateAddressPresence),
			validation.Skip,
		),
		validation.Field(&gdlv.Customer,
			validation.By(validateAddressPresence),
			validation.Skip,
		),
		validation.Field(&gdlv.Receiver,
			validation.By(validateAddressPresence),
			validation.Skip,
		),
		validation.Field(&gdlv.Preceding,
			validation.Each(
				validation.By(validateTransportPreceding),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateDespatchDate(gdlv *bill.Delivery) validation.RuleFunc {
	return func(val any) error {
		date, _ := val.(*cal.Date)
		if date == nil {
			return nil
		}

		if date.Before(gdlv.IssueDate.AddMonths(-3).AddDays(1)) {
			return errors.New("must not be more than 3 months earlier than the issue date")
		}

		if date.After(gdlv.IssueDate.AddMonths(3)) {
			return errors.New("must not be more than 3 months after the issue date")
		}

		return nil
	}
}

func validateAddressPresence(val any) error {
	gpar, _ := val.(*org.Party)
	if gpar == nil {
		return nil
	}

	return validation.ValidateStruct(gpar,
		validation.Field(&gpar.Addresses,
			validation.Required,
			validation.Each(
				validation.By(validateAddress),
				validation.Skip,
			),
		),
	)
}

func validateAddress(val any) error {
	addr, _ := val.(*org.Address)
	if addr == nil {
		return nil
	}

	return validation.ValidateStruct(addr,
		validation.Field(&addr.Code, validation.Required, validation.Skip),
		validation.Field(&addr.Country, validation.Required, validation.Skip),
	)
}

func validateDeliveryTax(val any) error {
	gtax, _ := val.(*bill.Tax)
	if gtax == nil {
		return nil
	}

	return validation.ValidateStruct(gtax,
		validation.Field(&gtax.Ext,
			tax.ExtensionsExclude(addon.ExtKeySource),    // Not supported
			tax.ExtensionsExclude(addon.ExtKeySourceRef), // Not supported
			validation.Skip,
		),
	)
}

func buildMovementStatus(status, reason, date, sourceID string) *MovementStatus {
	return &MovementStatus{
		MovementStatus:     status,
		Reason:             fit(reason, 50),
		MovementStatusDate: date,
		SourceID:           sourceID,
		SourceBilling:      addon.SourceBillingProduced.String(),
	}
}

func buildShippingPoint(party *org.Party, deliveryDate *cal.Date) *ShippingPoint {
	if party == nil {
		return nil
	}

	return &ShippingPoint{
		DeliveryDate: date(deliveryDate),
		Address:      BuildAddress(party),
	}
}

func buildMovementLines(gdlv *bill.Delivery, prods []*Product) []*Line {
	if gdlv == nil || len(gdlv.Lines) != len(prods) {
		return nil
	}

	lines := make([]*Line, len(gdlv.Lines))

	for i, gl := range gdlv.Lines {
		prod := prods[i]
		lines[i] = buildMovementLine(gdlv, gl, prod)
	}

	return lines
}

func buildMovementLine(gdlv *bill.Delivery, gline *bill.Line, prod *Product) *Line {
	if gdlv == nil || gline == nil || gline.Item == nil || prod == nil {
		return nil
	}

	return &Line{
		LineNumber:         gline.Index,
		OrderReferences:    buildOrderReferences(gdlv.Preceding, gline),
		ProductCode:        fit(prod.ProductCode, 60),
		ProductDesc:        fit(prod.ProductDesc, 200),
		Quantity:           gline.Quantity.String(),
		UnitOfMeasure:      unitOfMeasure(gline),
		UnitPrice:          unitPrice(gline),
		Description:        fit(prod.ProductDesc, 200),
		DebitAmount:        movementDebitAmount(gdlv, gline),
		CreditAmount:       movementCreditAmount(gdlv, gline),
		Tax:                buildTax(gline),
		TaxExemptionReason: lineTaxExemptionReason(gline),
		TaxExemptionCode:   lineTaxExemptionCode(gline),
		SettlementAmount:   settlementAmount(gline),
	}
}

func buildMovementDocumentTotals(gdlv *bill.Delivery) *DocumentTotals {
	if gdlv == nil || gdlv.Totals == nil {
		return &DocumentTotals{
			TaxPayable: "0.00",
			NetTotal:   "0.00",
			GrossTotal: "0.00",
		}
	}

	return &DocumentTotals{
		TaxPayable: deliveryTaxPayable(gdlv),
		NetTotal:   gdlv.Totals.Total.String(),
		GrossTotal: gdlv.Totals.TotalWithTax.RescaleUp(2).String(), // ensure 2 decimal places (required for sealing)
	}
}

func movementType(gdlv *bill.Delivery) string {
	if gdlv == nil || gdlv.Tax == nil {
		return ""
	}

	return gdlv.Tax.Ext.Get(addon.ExtKeyMovementType).String()
}

func movementComments(gdlv *bill.Delivery) string {
	if gdlv == nil || len(gdlv.Notes) == 0 {
		return ""
	}

	// Use the first note as the movement comments
	return fit(gdlv.Notes[0].Text, 60)
}

func movementEndTime(gdlv *bill.Delivery) string {
	if gdlv == nil || gdlv.ReceiveDate == nil {
		return ""
	}

	// GOBL only provides the despatch date without time. We default the time to the
	// end of the day (23:59:59) consistently with movementStartTime.
	return gdlv.ReceiveDate.String() + "T23:59:59"
}

func movementStartTime(gdlv *bill.Delivery) string {
	if gdlv == nil || gdlv.DespatchDate == nil {
		return ""
	}

	// GOBL only provides the despatch date without time and the AT requires the
	// document to be issued before the despatch. To avoid invalid chronologies on
	// same-day deliveries, we default the time to the end of the day (23:59:59) to
	// ensure it is always after the current time.
	return gdlv.DespatchDate.String() + "T23:59:59"
}

func movementDebitAmount(gdlv *bill.Delivery, gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	if !isOutbound(gdlv) {
		return amountOrZero(gline.Total).String()
	}

	return ""
}

func movementCreditAmount(gdlv *bill.Delivery, gline *bill.Line) string {
	if gline == nil {
		return ""
	}

	if !isOutbound(gdlv) {
		return ""
	}

	return amountOrZero(gline.Total).String()
}

func isOutbound(gdlv *bill.Delivery) bool {
	if gdlv == nil || gdlv.Tax == nil {
		return false
	}

	return gdlv.Tax.Ext.Get(addon.ExtKeyMovementType) != addon.MovementTypeReturn
}

func deliveryTaxPayable(gdlv *bill.Delivery) string {
	if gdlv == nil || gdlv.Totals == nil || gdlv.Totals.Taxes == nil {
		return ""
	}

	vat := gdlv.Totals.Taxes.Category(tax.CategoryVAT)
	if vat == nil {
		return ""
	}

	return vat.Amount.String()
}

func deliveryTaxTotal(gdlv *bill.Delivery) *tax.Total {
	if gdlv == nil || gdlv.Totals == nil || gdlv.Totals.Taxes == nil {
		return nil
	}

	return gdlv.Totals.Taxes
}

func dateOrZero(date *cal.Date) cal.Date {
	if date == nil {
		return cal.Date{}
	}
	return *date
}
