package saft

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/goblwrap"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// sealable defines the methods a GOBL document must implement to be sealed.
type sealable interface {
	GetSupplier() *org.Party
	GetCustomer() *org.Party
	GetCode() cbc.Code
	SetCode(code cbc.Code)
	GetSeries() cbc.Code
	GetIssueDate() cal.Date
	GetRegimeDef() *tax.RegimeDef
	GetExt() tax.Extensions
}

// Seal generates the signature, ATCUD, and QR code for a GOBL invoice envelope and
// appends them as stamps to the envelope.
func Seal(env *gobl.Envelope, prevSig *Signature, valCode string, cfg *config.Config) (*Signature, error) {
	// Validate envelope for sealing
	if err := validateEnvelope(env); err != nil {
		return nil, err
	}

	// Extract the document from the envelope
	gdoc := goblwrap.Extract(env).(sealable)
	if gdoc == nil {
		return nil, errors.New("envelope must contain a sealable document")
	}

	// Validate document for sealing
	if err := validateDocument(gdoc); err != nil {
		return nil, err
	}

	// Assign the code to the document
	assignCode(gdoc, prevSig)

	// Generate the seal components
	sr, err := seal(gdoc, prevSig, valCode, cfg)
	if err != nil {
		return nil, err
	}

	// Add the stamps to the envelope
	setStamps(env, sr, cfg)

	return sr.sig, nil
}

// validateEnvelope validates the envelope can be sealed. These validations specific
// to the sealing process, not to the regime, the addon or the saft conversion.
func validateEnvelope(env *gobl.Envelope) error {
	if env.Signed() {
		return errors.New("envelope must not be signed")
	}

	if env.Head != nil && env.Head.Stamp(pt.StampProviderATHash) != nil {
		return errors.New("envelope already sealed")
	}

	return nil
}

// validateDocument validates the document can be sealed. These validations specific
// to the sealing process, not to the regime, the addon or the saft conversion.
func validateDocument(doc sealable) error {
	errs := make(validation.Errors)

	if doc.GetCode() != cbc.CodeEmpty {
		errs["code"] = errors.New("must not be set")
	}

	if !isExternal(doc) { // Issue date of external documents cannot be validated because it is not controlled by the application
		loc := doc.GetRegimeDef().TimeLocation()
		today := config.TodayIn(loc)
		if d := doc.GetIssueDate(); d != today {
			errs["issue_date"] = fmt.Errorf(
				"must be today in %s (%s)",
				loc.String(),
				today.Time().Format("2006-01-02"),
			)
		}

		if len(errs) > 0 {
			return errs
		}
	}

	return nil
}

func assignCode(doc sealable, prevSig *Signature) {
	// Get the last sequence number
	var seqNo int
	if prevSig != nil {
		seqNo = prevSig.SeqNo()
	}

	// Increment the sequence number
	seqNo++

	// Convert to code and set it
	doc.SetCode(cbc.Code(strconv.Itoa(seqNo)))
}

type sealResult struct {
	sig   *Signature
	atcud string
	qr    string
}

func seal(gdoc sealable, prevSig *Signature, valCode string, cfg *config.Config) (*sealResult, error) {
	// Build the SAFT entities
	sup := BuildSupplier(gdoc.GetSupplier())
	cus := BuildCustomer(gdoc.GetCustomer())
	doc, err := BuildSourceDocument(gdoc, nil)
	if err != nil {
		return nil, fmt.Errorf("building document: %w", err)
	}

	// Generate the signature
	key, err := cfg.SigningKey()
	if err != nil {
		return nil, fmt.Errorf("reading signing key: %w", err)
	}

	sig, err := BuildSignature(doc, prevSig, key)
	if err != nil {
		return nil, fmt.Errorf("building signature: %w", err)
	}

	// Generate the ATCUD
	atcud := BuildATCUD(doc, valCode)

	// Generate the QR code
	qr := BuildQRCode(&BuildQRCodeParams{
		Supplier:  sup,
		Customer:  cus,
		Document:  doc,
		ATCUD:     atcud,
		Signature: sig,
		SWCertNum: cfg.CertificateNumber,
	})

	return &sealResult{
		sig:   sig,
		atcud: atcud,
		qr:    qr,
	}, nil
}

func setStamps(env *gobl.Envelope, sr *sealResult, cfg *config.Config) {
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATATCUD,
		Value:    sr.atcud,
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATHash,
		Value:    sr.sig.HashCharacters(),
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATHashFull,
		Value:    sr.sig.Hash,
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATQR,
		Value:    sr.qr,
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATAppID,
		Value:    strconv.Itoa(cfg.CertificateNumber),
	})
	env.Head.AddStamp(&head.Stamp{
		Provider: pt.StampProviderATTimestamp,
		Value:    sr.sig.SystemEntryDate,
	})
}

func isExternal(doc sealable) bool {
	return doc.GetExt().Has(addon.ExtKeySourceRef)
}
