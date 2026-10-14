package saft

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// Signature represents an invoice signature
type Signature struct {
	IssueDate       string
	SystemEntryDate string
	DocumentNo      string
	GrossTotal      string
	PreviousHash    string

	Hash string
}

// BuildSignature builds a signature for a given document.
func BuildSignature(doc SourceDocument, prevSig *Signature, pkey []byte) (*Signature, error) {
	sig := &Signature{
		IssueDate:       doc.GetIssueDate(),
		SystemEntryDate: Timestamp(),
		DocumentNo:      doc.GetDocumentNo(),
		GrossTotal:      doc.GetGrossTotal(),
		PreviousHash:    previousHash(prevSig),
	}

	if err := sig.signRSA(pkey); err != nil {
		return nil, err
	}

	return sig, nil
}

// HashCharacters extracts a 4-char snippet of the signature hash
func (s *Signature) HashCharacters() string {
	return HashCharacters(s.Hash)
}

// HashCharacters extracts a 4-char snippet from a full hash string
func HashCharacters(hash string) string {
	if len(hash) < 31 {
		return ""
	}
	return string([]byte{hash[0], hash[10], hash[20], hash[30]})
}

// Msg returns the message to be signed
func (s *Signature) Msg() string {
	return strings.Join([]string{
		s.IssueDate,
		s.SystemEntryDate,
		s.DocumentNo,
		s.GrossTotal,
		s.PreviousHash,
	}, ";")
}

// SeqNo extracts the sequence number from the document number
func (s *Signature) SeqNo() int {
	_, _, seq := DocNumberParts(s.DocumentNo)
	return seq
}

// Series extracts the series from the document number
func (s *Signature) Series() string {
	_, ser, _ := DocNumberParts(s.DocumentNo)
	return ser.String()
}

func (s *Signature) signRSA(pkey []byte) error {
	// Decode the private key
	block, _ := pem.Decode(pkey)
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("decoding: invalid private key")
	}

	// Parse the private key
	rsaPriv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parsing: %w", err)
	}

	// Sign the message
	hash := sha1.New()
	hash.Write([]byte(s.Msg()))
	hashed := hash.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaPriv.(*rsa.PrivateKey), crypto.SHA1, hashed)
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	// Encode the signature
	s.Hash = base64.StdEncoding.EncodeToString(sig)

	return nil
}

func previousHash(prevSig *Signature) string {
	if prevSig == nil {
		return ""
	}
	return prevSig.Hash
}
