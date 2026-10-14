package saft

import (
	"bytes"
	"fmt"
	"io"

	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
)

// auditFileNamespace is the XML namespace of the AuditFile root element.
const auditFileNamespace = "urn:OECD:StandardAuditFile-Tax:PT_1.04_01"

// ClassTotals holds the values of a document class wrapper element that must be known
// before the documents of that class are streamed.
type ClassTotals struct {
	Entries  int
	Lines    int
	Debit    string
	Credit   string
	Quantity string
}

// StreamEncoder writes a SAF-T audit file incrementally to an output writer.
type StreamEncoder struct {
	w      io.Writer
	header *Header

	begun  bool
	ended  bool
	at     section
	opened map[section]bool

	masterFilesClosed bool
	classElem         string
}

// section identifies the part of an audit file the encoder is writing.
type section int

const (
	sectionRoot section = iota
	sectionMasterFiles
	sectionTaxTable
	sectionSourceDocuments
)

// NewStreamEncoder prepares a stream encoder. The issuer and period are expected to have
// been validated by the caller with ValidateIssuer and ValidatePeriod.
func NewStreamEncoder(w io.Writer, cfg config.Config, issuer *org.Party, per cal.Period) *StreamEncoder {
	return &StreamEncoder{
		w:      w,
		header: buildHeader(cfg, issuer, per),
		opened: make(map[section]bool),
	}
}

// enter moves the encoder into a section, opening it and everything it nests in, and
// closing whatever the previous section left open.
func (e *StreamEncoder) enter(target section) error {
	if !e.begun || e.ended {
		return fmt.Errorf("stream not active")
	}
	if target < e.at {
		return fmt.Errorf("sections must be written in schema order")
	}

	switch target {
	case sectionMasterFiles:
		if err := e.open(sectionMasterFiles); err != nil {
			return err
		}
	case sectionTaxTable:
		if err := e.open(sectionMasterFiles); err != nil {
			return err
		}
		if err := e.open(sectionTaxTable); err != nil {
			return err
		}
	case sectionSourceDocuments:
		if err := e.closeMasterFiles(); err != nil {
			return err
		}
		if err := e.open(sectionSourceDocuments); err != nil {
			return err
		}
		// A new class ends the previous one; the caller only asks to enter.
		if err := e.closeClass(); err != nil {
			return err
		}
	}

	e.at = target
	return nil
}

// open writes a section's opening tag the first time it is entered.
func (e *StreamEncoder) open(s section) error {
	if e.opened[s] {
		return nil
	}
	e.opened[s] = true

	switch s {
	case sectionMasterFiles:
		return e.writeRaw("\n  <MasterFiles>")
	case sectionTaxTable:
		return e.writeRaw("\n    <TaxTable>")
	case sectionSourceDocuments:
		return e.writeRaw("\n  <SourceDocuments>")
	}
	return nil
}

// closeMasterFiles closes the tax table and the master files section, emitting the empty
// form when nothing was ever written to them.
func (e *StreamEncoder) closeMasterFiles() error {
	if e.masterFilesClosed {
		return nil
	}
	e.masterFilesClosed = true

	if e.opened[sectionTaxTable] {
		if err := e.writeRaw("\n    </TaxTable>"); err != nil {
			return err
		}
	}
	if e.opened[sectionMasterFiles] {
		if err := e.writeRaw("\n  </MasterFiles>"); err != nil {
			return err
		}
	} else if err := e.writeRaw("\n  <MasterFiles></MasterFiles>"); err != nil {
		return err
	}
	return nil
}

// Begin writes the XML declaration, the audit file root element and its header.
func (e *StreamEncoder) Begin() error {
	e.begun = true

	if err := e.writeRaw(xmlHeader + "\n<AuditFile xmlns=\"" + auditFileNamespace + "\">\n"); err != nil {
		return err
	}
	return EncodeXML(e.header, &indentWriter{w: e.w, pad: []byte("  ")}, false)
}

// WriteMasterFile writes a stored customer or product XML fragment, which both sit
// directly inside the master files. The caller is responsible for the order the SAF-T
// schema fixes: customers, then products, then tax table entries.
func (e *StreamEncoder) WriteMasterFile(frag []byte) error {
	if err := e.enter(sectionMasterFiles); err != nil {
		return err
	}
	return e.writeFragment(frag, 2)
}

// WriteTaxTableEntry writes a stored tax table entry XML fragment, one level deeper than
// the other master files as it is wrapped in a tax table of its own.
func (e *StreamEncoder) WriteTaxTableEntry(frag []byte) error {
	if err := e.enter(sectionTaxTable); err != nil {
		return err
	}
	return e.writeFragment(frag, 3)
}

// BeginClass opens the wrapper element of a document class, writing its entry count and
// total amounts. Classes must be written in the SAF-T schema order: SI, MG, WD, PY.
func (e *StreamEncoder) BeginClass(dc DocumentClass, t ClassTotals) error {
	if err := e.enter(sectionSourceDocuments); err != nil {
		return err
	}

	var body string
	switch dc {
	case DocumentClassInvoice:
		e.classElem, body = "SalesInvoices", entriesAndTotals(t)
	case DocumentClassTransport:
		e.classElem, body = "MovementOfGoods", movementTotals(t)
	case DocumentClassWork:
		e.classElem, body = "WorkingDocuments", entriesAndTotals(t)
	case DocumentClassPayment:
		e.classElem, body = "Payments", entriesAndTotals(t)
	default:
		return fmt.Errorf("unknown document class: %s", dc)
	}

	return e.writeRaw("\n    <" + e.classElem + ">" + body)
}

// WriteDocument writes a stored source document XML fragment into the currently open
// document class.
func (e *StreamEncoder) WriteDocument(frag []byte) error {
	if e.classElem == "" {
		return fmt.Errorf("no document class open")
	}
	return e.writeFragment(frag, 3)
}

// End closes all open sections and the audit file root element.
func (e *StreamEncoder) End() error {
	if !e.begun || e.ended {
		return fmt.Errorf("stream not active")
	}

	if err := e.closeMasterFiles(); err != nil {
		return err
	}
	if err := e.closeClass(); err != nil {
		return err
	}
	e.ended = true

	if e.opened[sectionSourceDocuments] {
		if err := e.writeRaw("\n  </SourceDocuments>"); err != nil {
			return err
		}
	} else if err := e.writeRaw("\n  <SourceDocuments></SourceDocuments>"); err != nil {
		return err
	}

	return e.writeRaw("\n</AuditFile>")
}

func (e *StreamEncoder) closeClass() error {
	if e.classElem == "" {
		return nil
	}
	elem := e.classElem
	e.classElem = ""
	return e.writeRaw("\n    </" + elem + ">")
}

// writeFragment writes a stored Windows-1252 XML fragment on a new line, indenting
// every line of the fragment to the given element depth.
func (e *StreamEncoder) writeFragment(frag []byte, depth int) error {
	if err := e.writeRaw("\n"); err != nil {
		return err
	}
	iw := &indentWriter{w: e.w, pad: bytes.Repeat([]byte("  "), depth)}
	_, err := iw.Write(frag)
	return err
}

func (e *StreamEncoder) writeRaw(s string) error {
	_, err := io.WriteString(e.w, s)
	return err
}

func movementTotals(t ClassTotals) string {
	return fmt.Sprintf("\n      <NumberOfMovementLines>%d</NumberOfMovementLines>\n      <TotalQuantityIssued>%s</TotalQuantityIssued>", t.Lines, t.Quantity)
}

func entriesAndTotals(t ClassTotals) string {
	return fmt.Sprintf("\n      <NumberOfEntries>%d</NumberOfEntries>\n      <TotalDebit>%s</TotalDebit>\n      <TotalCredit>%s</TotalCredit>", t.Entries, t.Debit, t.Credit)
}

// indentWriter prefixes every line written through it with a fixed padding. It is safe
// on Windows-1252 XML output because the encoding is single-byte and literal newlines
// in character data are always escaped, so any newline byte is structural.
type indentWriter struct {
	w       io.Writer
	pad     []byte
	midLine bool
}

func (iw *indentWriter) Write(p []byte) (int, error) {
	var total int
	for len(p) > 0 {
		if !iw.midLine {
			if _, err := iw.w.Write(iw.pad); err != nil {
				return total, err
			}
			iw.midLine = true
		}

		i := bytes.IndexByte(p, '\n')
		if i < 0 {
			n, err := iw.w.Write(p)
			return total + n, err
		}

		n, err := iw.w.Write(p[:i+1])
		total += n
		if err != nil {
			return total, err
		}
		iw.midLine = false
		p = p[i+1:]
	}
	return total, nil
}
