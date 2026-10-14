// Package soap provides types and functions to work with SOAP messages.
package soap

import "fmt"

// Envelope represents a SOAP envelope
type Envelope struct {
	XMLName struct{} `xml:"http://schemas.xmlsoap.org/soap/envelope/ soap:Envelope"`

	Header struct {
		Headers []interface{} `xml:",any"`
	} `xml:"soap:Header"`

	Body struct {
		Content interface{} `xml:",any"`
	} `xml:"soap:Body"`
}

// Fault represents a SOAP fault
type Fault struct { // nolint:errname
	XMLName struct{} `xml:"soap:Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
}

// Error implements the error interface.
func (f *Fault) Error() string {
	return fmt.Sprintf("SOAP Fault %s: %s", f.Code, f.String)
}

// Envelop builds a new SOAP envelope for the given body and headers
func Envelop(body interface{}, headers ...interface{}) *Envelope {
	e := new(Envelope)
	e.Header.Headers = headers
	e.Body.Content = body
	return e
}
