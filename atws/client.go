package atws

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/nbio/xml"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/gobl.pt.saft/atws/soap"
	"github.com/invopop/xmldsig"
)

// Connection provides a connection to the AT Web Services.
type Connection struct {
	conn *resty.Client

	env       Environment
	pkey      []byte
	swCertNum int
}

// Client provides a client for the AT Web Services.
type Client struct {
	*Connection

	username string
	password string

	client *Client // self-reference to be used in aliased types
}

// New creates a connection to the AT Web Services.
func New(conf *Config) (*Connection, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	cert, err := xmldsig.LoadCertificate(conf.CertificatePath, conf.CertificatePassword)
	if err != nil {
		return nil, fmt.Errorf("loading AT certificate: %w", err)
	}
	tlsConfig, err := cert.TLSAuthConfig()
	if err != nil {
		return nil, err
	}

	// For some reason, the xmldsig package sets the RootCAs to the client
	// certificate's RootCA, which is not necessary to validate the server
	// certificate. We set it back to nil here to use the system CA certificates
	// instead. This works with the AT servers certificates.
	//
	// Note: The AT periodically rotates their server certificates. Since they
	// are signed by standard root CAs in the system certificate store, no action
	// is typically required.
	tlsConfig.RootCAs = nil

	c := new(Connection)

	c.conn = resty.New()
	c.conn.SetTLSClientConfig(tlsConfig)
	c.conn.SetDebug(conf.Debug)
	c.conn.SetHeader("Content-Type", "application/xml")

	// Use github.com/nbio/xml for marshaling and unmarshaling XML
	c.conn.SetXMLMarshaler(xml.Marshal)
	c.conn.SetXMLUnmarshaler(xml.Unmarshal)

	c.pkey, err = os.ReadFile(conf.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading public key: %w", err)
	}
	c.swCertNum = conf.SWCertNum
	c.env = conf.Environment

	return c, nil
}

// NewClient creates a new authenticated client for the AT Web Services.
func (c *Connection) NewClient(username, password string) *Client {
	cl := &Client{
		Connection: c,

		username: username,
		password: password,
	}

	cl.client = cl

	return cl
}

// HTTPClient returns the underlying HTTP client (useful for mocking).
func (c *Connection) HTTPClient() *http.Client {
	return c.conn.GetClient()
}

// Series returns a client for the Series Web Service.
func (c *Client) Series() *SeriesClient {
	return (*SeriesClient)(c)
}

// Movements returns a client for the Movement Documents Web Service.
func (c *Client) Movements() *MovementsClient {
	return (*MovementsClient)(c)
}

// Invoices returns a client for the Invoices Web Service.
func (c *Client) Invoices() *InvoicesClient {
	return (*InvoicesClient)(c)
}

// Payments returns a client for the Payments Web Service.
func (c *Client) Payments() *PaymentsClient {
	return (*PaymentsClient)(c)
}

// WorkDocuments returns a client for the Work Documents Web Service.
func (c *Client) WorkDocuments() *WorkDocumentsClient {
	return (*WorkDocumentsClient)(c)
}

// call performs a SOAP request.
func (c *Client) call(action string, endpoints Endpoints, in, out interface{}) error {
	if in == nil {
		return ErrInvalidRequest
	}
	if endpoints == nil {
		panic("nil endpoints")
	}
	endpoint, ok := endpoints[c.env]
	if !ok {
		return errors.New("invalid environment")
	}

	auth, err := c.auth()
	if err != nil {
		return err
	}

	fault := new(soap.Fault)

	soapReq := soap.Envelop(in, auth)
	soapRes := soap.Envelop(out)
	soapErr := soap.Envelop(fault)

	res, err := c.conn.R().
		SetHeader("SOAPAction", action).
		SetBody(soapReq).
		SetResult(soapRes).
		SetError(soapErr).
		Post(endpoint)
	if err != nil {
		return err
	}
	if res.IsError() {
		if fault.Code != "" || fault.String != "" {
			return fault
		}
		// The fault could not be parsed, so include the raw body which is
		// otherwise lost and is usually the only clue as to what went wrong.
		return fmt.Errorf("HTTP %d: %s: %s", res.StatusCode(), res.Status(), truncate(res.String(), 512))
	}

	return nil
}

// truncate limits the length of a string so that unexpectedly large responses
// don't flood the logs.
func truncate(s string, limit int) string {
	s = strings.TrimSpace(s)
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "…"
}

// auth builds a new security header for the client. A new header needs to be
// generated right before each request as the time window for the header to be
// accepted is small.
func (c *Client) auth() (*soap.SecurityHeader, error) {
	return soap.NewSecurityHeader(c.username, c.password, c.pkey)
}
