package realtime

import (
	"errors"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/atws"
)

// BuildSendStockMovement builds a SendStockMovement request from a SAF-T AuditFile.
func BuildSendStockMovement(af *saft.AuditFile) (*atws.SendStockMovement, error) {
	if af == nil {
		return nil, errors.New("audit file must not be nil")
	}

	if af.Header == nil {
		return nil, errors.New("audit file header must not be nil")
	}

	if af.SourceDocuments == nil || af.SourceDocuments.MovementOfGoods == nil {
		return nil, errors.New("audit file must contain movement of goods")
	}

	if len(af.SourceDocuments.MovementOfGoods.StockMovements) == 0 {
		return nil, errors.New("audit file must contain at least one stock movement")
	}

	if af.MasterFiles == nil || len(af.MasterFiles.Customers) == 0 {
		return nil, errors.New("audit file must contain at least one customer")
	}

	// Build the SendStockMovement from SAF-T structures
	return buildSendStockMovement(af), nil
}

func buildSendStockMovement(af *saft.AuditFile) *atws.SendStockMovement {
	hd := af.Header
	sm := af.SourceDocuments.MovementOfGoods.StockMovements[0]
	cus := af.MasterFiles.Customers[0]

	req := &atws.SendStockMovement{
		TaxRegistrationNumber: hd.TaxRegistrationNumber,
		CompanyName:           hd.CompanyName,
		CompanyAddress:        buildAddress(hd.CompanyAddress),
		DocumentNumber:        sm.DocumentNumber,
		ATCUD:                 sm.ATCUD,
		ATDocCodeID:           sm.ATDocCodeID,
		MovementStatus:        movementStatus(sm),
		MovementDate:          sm.MovementDate,
		MovementType:          sm.MovementType,
		CustomerTaxID:         cus.CustomerTaxID,
		CustomerName:          cus.CompanyName,
		AddressTo:             addressTo(sm),
		AddressFrom:           addressFrom(sm),
		MovementEndTime:       sm.MovementEndTime,
		MovementStartTime:     sm.MovementStartTime,
		Lines:                 buildLines(sm.Lines),
	}

	// The CustomerAddress is optional in the realtime transport API. Also, the
	// "Desconhecido" placeholders used in SAF-T don't fit the stricter SOAP schema of
	// the delivery documents web services. So we only map the address if the customer
	// tax ID is different from the generic "Consumidor final" one.
	if cus.CustomerTaxID != saft.GenericCustomerTaxID {
		req.CustomerAddress = buildAddress(cus.BillingAddress)
	}

	return req
}

func buildLines(sls []*saft.Line) []*atws.Line {
	if sls == nil {
		return nil
	}

	ls := make([]*atws.Line, len(sls))
	for i, sl := range sls {
		ls[i] = &atws.Line{
			OrderReferences:    buildOrderReferences(sl),
			ProductDescription: sl.ProductDesc,
			Quantity:           sl.Quantity,
			UnitOfMeasure:      sl.UnitOfMeasure,
			UnitPrice:          sl.UnitPrice,
		}
	}
	return ls
}

func buildAddress(addr *saft.Address) *atws.Address {
	if addr == nil {
		return nil
	}

	return &atws.Address{
		Addressdetail: addr.AddressDetail,
		City:          addr.City,
		PostalCode:    addr.PostalCode,
		Country:       addr.Country,
	}
}

func buildOrderReferences(sl *saft.Line) []*atws.OrderReferences {
	if sl == nil || sl.OrderReferences == nil {
		return nil
	}

	ors := make([]*atws.OrderReferences, len(sl.OrderReferences))
	for i, or := range sl.OrderReferences {
		ors[i] = &atws.OrderReferences{
			OriginatingON: or.OriginatingON,
			OrderDate:     "", // Field doesn't exist for stock movements in the WSDL. Leave empty to omit from XML.
		}
	}
	return ors
}

func movementStatus(sm *saft.StockMovement) string {
	if sm.DocumentStatus == nil {
		return ""
	}
	return sm.DocumentStatus.MovementStatus
}

func addressTo(sm *saft.StockMovement) *atws.Address {
	if sm.ShipTo == nil {
		return nil
	}
	return buildAddress(sm.ShipTo.Address)
}

func addressFrom(sm *saft.StockMovement) *atws.Address {
	if sm.ShipFrom == nil {
		return nil
	}
	return buildAddress(sm.ShipFrom.Address)
}
