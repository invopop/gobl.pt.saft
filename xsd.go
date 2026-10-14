package saft

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/xml"
	"fmt"
	"sync"

	"github.com/lestrrat-go/helium"
	"github.com/lestrrat-go/helium/xsd"
	"golang.org/x/text/unicode/norm"
)

// XMLNamespace is the target namespace of the SAF-T PT XSD schema.
const XMLNamespace = "urn:OECD:StandardAuditFile-Tax:PT_1.04_01"

//go:embed schemas/saftpt1.04_01.xsd
var schemaData []byte

// loadSchema parses and compiles the embedded SAF-T XSD schema only once.
var loadSchema = sync.OnceValues(func() (*xsd.Schema, error) {
	ctx := context.Background()
	doc, err := helium.NewParser().Parse(ctx, schemaData)
	if err != nil {
		return nil, fmt.Errorf("parsing schema: %w", err)
	}
	return xsd.NewCompiler().Label("saftpt1.04_01.xsd").Compile(ctx, doc)
})

// ValidateXSD validates the XML form of the audit file against the official SAF-T PT
// XSD schema. Returns the full list of validation errors or nil if none.
func (af *AuditFile) ValidateXSD(ctx context.Context) []error {
	return validateFragment(ctx, "AuditFile", af)
}

// ValidateXSD validates the XML against the official SAF-T PT XSD schema.
// Returns the full list of validations error or nil if none.
func (r *ConvertResult) ValidateXSD(ctx context.Context) []error {
	var errs []error

	sds, err := wrapSourceDocument(r.Document)
	if err != nil {
		errs = append(errs, err)
	} else {
		errs = append(errs, validateFragment(ctx, "SourceDocuments", sds)...)
	}

	if r.Customer != nil {
		errs = append(errs, validateFragment(ctx, "Customer", r.Customer)...)
	}
	for _, p := range r.Products {
		errs = append(errs, validateFragment(ctx, "Product", p)...)
	}
	for _, te := range r.TaxTableEntries {
		errs = append(errs, validateFragment(ctx, "TaxTableEntry", te)...)
	}

	return errs
}

// wrapSourceDocument wraps a single source document in a SourceDocuments element.
func wrapSourceDocument(doc SourceDocument) (*SourceDocuments, error) {
	sds := new(SourceDocuments)
	switch d := doc.(type) {
	case *Invoice:
		sds.SalesInvoices = buildSalesInvoices([]*Invoice{d})
	case *Payment:
		sds.Payments = buildPayments([]*Payment{d})
	case *WorkDocument:
		sds.WorkingDocuments = buildWorkingDocuments([]*WorkDocument{d})
	case *StockMovement:
		sds.MovementOfGoods = buildMovementOfGoods([]*StockMovement{d})
	default:
		return nil, fmt.Errorf("unknown source document type: %T", doc)
	}
	return sds, nil
}

// validateFragment validates a single SAF-T structure as a standalone XML document,
// wrapping it in a root element of the given name qualified with the SAF-T namespace.
func validateFragment(ctx context.Context, name string, v any) []error {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	start := xml.StartElement{Name: xml.Name{Space: XMLNamespace, Local: name}}
	if err := enc.EncodeElement(v, start); err != nil {
		return []error{fmt.Errorf("encoding %s: %w", name, err)}
	}
	if err := enc.Close(); err != nil {
		return []error{fmt.Errorf("encoding %s: %w", name, err)}
	}

	return validateXML(ctx, name, norm.NFC.Bytes(buf.Bytes()))
}

// validateXML parses the given XML document and validates it against the compiled
// SAF-T XSD schema, accumulating every violation reported by the validator.
func validateXML(ctx context.Context, label string, data []byte) []error {
	schema, err := loadSchema()
	if err != nil {
		return []error{fmt.Errorf("loading SAF-T schema: %w", err)}
	}

	doc, err := helium.NewParser().Parse(ctx, data)
	if err != nil {
		return []error{fmt.Errorf("parsing XML: %w", err)}
	}

	ec := new(errorCollector)
	err = xsd.NewValidator(schema).Label(label).ErrorHandler(ec).Validate(ctx, doc)
	if err == nil {
		return nil
	}
	if len(ec.errs) == 0 {
		return []error{err}
	}
	return ec.errs
}

// errorCollector accumulates the individual validation errors delivered by the
// helium XSD validator.
type errorCollector struct {
	errs []error
}

// Handle implements the helium.ErrorHandler interface.
func (c *errorCollector) Handle(_ context.Context, err error) {
	c.errs = append(c.errs, err)
}
