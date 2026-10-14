package saft

import (
	"bytes"
	"encoding/xml"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const xmlHeader = `<?xml version="1.0" encoding="Windows-1252"?>`

// MarshalXML encodes the given value to XML according to the SAF-T specification and
// returns the resulting byte slice. If the root parameter is true, the value is
// considered a root document and the XML header is prepended to the output.
func MarshalXML(v any, root bool) ([]byte, error) {
	var buf bytes.Buffer
	if err := EncodeXML(v, &buf, root); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalXML decodes the given XML data into the given value according to the SAF-T
// specification. The input byte slice is assumed to be encoded in Windows-1252.
func UnmarshalXML(data []byte, v any) error {
	return DecodeXML(bytes.NewReader(data), v)
}

// EncodeXML encodes the given value to XML in Windows-1252 as required by SAF-T.
// If root is true, the XML header is prepended. Characters unsupported by
// Windows-1252 are normalized or escaped as HTML entities.
func EncodeXML(v any, w io.Writer, root bool) error {
	// Prepend the XML header if root document
	if root {
		if _, err := w.Write([]byte(xmlHeader + "\n")); err != nil {
			return err
		}
	}

	// First, encode to XML in UTF-8 to a buffer
	var buf bytes.Buffer
	xmlEnc := xml.NewEncoder(&buf)
	xmlEnc.Indent("", "  ")
	if err := xmlEnc.Encode(v); err != nil {
		return err
	}
	if err := xmlEnc.Flush(); err != nil {
		return err
	}

	// Second, write the buffer encoded to Windows-1252
	writer := transform.NewWriter(w, transform.Chain(
		norm.NFC, // Normalize Unicode to precomposed form to preserve accented characters
		encoding.HTMLEscapeUnsupported(charmap.Windows1252.NewEncoder()), // Escape unsupported characters as HTML entities
	))
	if _, err := writer.Write(buf.Bytes()); err != nil {
		return err
	}

	return writer.Close()
}

// DecodeXML decodes the given XML data into the given value according to the SAF-T
// specification. The input reader is assumed to be encoded in Windows-1252.
func DecodeXML(r io.Reader, v any) error {
	// Wrap the reader in a Windows-1252 decoder
	decoder := charmap.Windows1252.NewDecoder()
	reader := transform.NewReader(r, decoder)

	// Decode the XML data
	xmlDec := xml.NewDecoder(reader)
	// A whole audit file declares Windows-1252 in its XML header, and the decoder
	// refuses any declared encoding other than UTF-8 unless a CharsetReader is set. The
	// reader above has already converted the bytes, so this one passes them through
	// untouched. Stored fragments carry no header and never reach it.
	xmlDec.CharsetReader = func(_ string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	return xmlDec.Decode(v)
}
