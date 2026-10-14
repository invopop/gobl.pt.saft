// Package saft implements the conversion of GOBL invoices, payments, orders and
// delivery documents into Portuguese SAF-T (PT) structures, as well as the sealing
// of those documents: the AT signature (hash), ATCUD and QR code.
//
// The GOBL addon itself lives in the addon subpackage; the client for the AT web
// services in atws, and the builders of its real-time reporting requests in realtime.
package saft

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/goblwrap"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
)

// Status codes
const (
	StatusNormal      = "N" // Normal
	StatusSelfInvoice = "S" // Autofaturação
	StatusCanceled    = "A" // Anulado
	StatusSummary     = "R" // Resumo
	StatusBilled      = "F" // Faturado
)

const (
	// TimeLocation used as reference for SAF-T documents
	TimeLocation = "Europe/Lisbon"

	// TaxTypeVAT is the TaxType value for VAT
	TaxTypeVAT = "IVA"
)

// convertible is the interface that GOBL documents must implement to be converted to
// SAF-T structures.
type convertible interface {
	GetSupplier() *org.Party
	GetCustomer() *org.Party
	GetBillLines() []*bill.Line
	GetBillTotals() *bill.Totals
}

// ConvertResult provides the SAF-T structures resulting from a conversion
type ConvertResult struct {
	Document        SourceDocument
	Supplier        *Supplier
	Customer        *Customer
	Products        []*Product
	TaxTableEntries []*TaxTableEntry
}

// Convert maps a GOBL envelope to SAF-T structures
func Convert(env *gobl.Envelope, sourceID string, cfg *config.Config) (*ConvertResult, error) {
	// Validate inputs
	atcud := env.Head.GetStamp(pt.StampProviderATATCUD)
	if atcud == nil {
		return nil, errors.New("ATCUD stamp must be present")
	}

	hash := env.Head.GetStamp(pt.StampProviderATHashFull)
	if hash == nil {
		return nil, errors.New("full hash stamp must be present")
	}

	ts := env.Head.GetStamp(pt.StampProviderATTimestamp)
	if ts == nil {
		return nil, errors.New("timestamp stamp must be present")
	}

	if len(sourceID) > 30 {
		return nil, errors.New("sourceID must be less than 30 characters")
	}

	// Extract the document from the envelope
	gdoc := goblwrap.Extract(env).(convertible)

	// Build the SAF-T master files
	sup := BuildSupplier(gdoc.GetSupplier())
	cus := BuildCustomer(gdoc.GetCustomer())
	pds := BuildProducts(gdoc.GetBillLines())
	tes := BuildTaxTableEntries(gdoc.GetBillTotals())

	// Build the SAF-T source document
	bp := &BuildParams{
		ATCUD:             atcud.Value,
		SourceID:          sourceID,
		Hash:              hash.Value,
		SigningKeyVersion: cfg.SigningKeyVersion,
		SystemEntryDate:   ts.Value,
		Customer:          cus,
		Products:          pds,
	}

	doc, err := BuildSourceDocument(gdoc, bp)
	if err != nil {
		return nil, err
	}

	// Return the result
	return &ConvertResult{
		Document:        doc,
		Supplier:        sup,
		Customer:        cus,
		Products:        uniqueBy(pds, func(p *Product) any { return p.ProductCode }),
		TaxTableEntries: tes,
	}, nil
}

// hashID generates a unique ID for an entity from its hash. The function will produce the
// same ID for any entity with the same content. To be called after the entity has been
// populated with data.
func hashID(e any) string {
	data, _ := MarshalXML(e, false)
	hash := sha256.Sum256(data)
	return fit(hex.EncodeToString(hash[:]), 30)
}

func uniqueBy[T any](items []*T, comparer func(*T) any) []*T {
	unique := make([]*T, 0, len(items))
	seen := make(map[interface{}]bool)
	for _, i := range items {
		key := comparer(i)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, i)
		}
	}
	return unique
}

// DocNumber generates a SAF-T document number from a GOBL series and a code
func DocNumber(series, code cbc.Code) string {
	return series.JoinWith("/", code).String()
}

// DocNumberParts splits a SAF-T document number into its components: document type,
// series code and sequence number
func DocNumberParts(c string) (cbc.Code, cbc.Code, int) {
	parts := strings.Split(c, " ")
	if len(parts) < 2 {
		return cbc.CodeEmpty, cbc.CodeEmpty, 0
	}
	typ := cbc.Code(parts[0])

	parts = strings.Split(parts[1], "/")
	ser := cbc.Code(parts[0])
	if len(parts) < 2 {
		return typ, ser, 0
	}

	seq, err := strconv.Atoi(parts[1])
	if err != nil {
		return typ, ser, 0
	}
	return typ, ser, seq
}

// IsWorkDocument returns whether a GOBL invoice represents a SAF-T work document
func IsWorkDocument(ginv *bill.Invoice) bool {
	return ginv.Tax != nil && ginv.Tax.Ext.Has(addon.ExtKeyWorkType)
}

// Timestamp returns the current timestamp in the SAF-T format
func Timestamp() string {
	return config.Now().In(pt.New().TimeLocation()).Format("2006-01-02T15:04:05")
}

func period(date cal.Date) int {
	return int(date.Time().Month())
}

func date(date *cal.Date) string {
	if date == nil {
		return ""
	}
	return date.Time().Format("2006-01-02")
}

// fit truncates a string to a maximum length
func fit(s string, m int) string {
	if len(s) > m {
		return s[:m]
	}
	return s
}
