package saft

import (
	"errors"

	"github.com/invopop/gobl.pt.saft/goblwrap"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

// DocumentClass represents the class of source document
type DocumentClass string

// Document class constants
const (
	DocumentClassInvoice   DocumentClass = "SI" // Faturas e documentos retificativos
	DocumentClassTransport DocumentClass = "MG" // Documentos de Transporte
	DocumentClassWork      DocumentClass = "WD" // Documentos de Conferência
	DocumentClassPayment   DocumentClass = "PY" // Recibos
)

// SourceDocument is the interface providing common accessors for SAFT source documents
type SourceDocument interface {
	GetDocumentClass() DocumentClass
	GetDocumentType() string
	GetDocumentNo() string
	GetIssueDate() string
	GetCustomerID() string
	GetGrossTotal() string
	GetStatus() string
	GetStatusDate() string
	GetStatusReason() string
	SetStatus(status, reason, sourceID string)
	GetTaxTotal() *tax.Total
	GetTaxPayable() string
	GetDocumentTotals() *DocumentTotals
}

// BuildParams are the parameters for building a source document
type BuildParams struct {
	ATCUD             string
	Customer          *Customer
	Products          []*Product
	SourceID          string
	Hash              string
	SigningKeyVersion string
	SystemEntryDate   string
}

// NewSourceDocument builds a new source document for a given document class.
func NewSourceDocument(docClass DocumentClass) SourceDocument {
	switch docClass {
	case DocumentClassInvoice:
		return new(Invoice)
	case DocumentClassPayment:
		return new(Payment)
	case DocumentClassWork:
		return new(WorkDocument)
	case DocumentClassTransport:
		return new(StockMovement)
	default:
		return nil
	}
}

// BuildSourceDocument builds a document from a given document and parameters.
func BuildSourceDocument(doc any, params *BuildParams) (SourceDocument, error) {
	doc = goblwrap.Unwrap(doc) // make sure we have an actual GOBL type

	switch doc := doc.(type) {
	case *bill.Invoice:
		if IsWorkDocument(doc) {
			return BuildWorkDocumentFromInvoice(doc, params)
		}
		return BuildInvoice(doc, params)
	case *bill.Payment:
		return BuildPayment(doc, params)
	case *bill.Delivery:
		return BuildStockMovement(doc, params)
	case *bill.Order:
		return BuildWorkDocumentFromOrder(doc, params)
	}
	return nil, errors.New("unknown document type")
}
