package atws

// Payload constants
const (
	EFaturaMDVersion = "0.0.1"
)

// Response codes
// nolint:misspell // Texts in Portuguese should not be flagged as misspelled.
const (
	// Success
	InvoiceResponseCodeSuccess = 0 // Operação efetuada com sucesso.

	// Authentication
	InvoiceResponseCodeUserNotFilled             = 1  // Utilizador não preenchido.
	InvoiceResponseCodeUserSizeIncorrect         = 2  // Tamanho do utilizador incorreto.
	InvoiceResponseCodeInvalidNIFAuth            = 3  // NIF inválido.
	InvoiceResponseCodeInvalidUserFormat         = 4  // Utilizador com formato inválido.
	InvoiceResponseCodeInvalidSubUserFormat      = 5  // Subutilizador com formato inválido.
	InvoiceResponseCodePasswordNotFilled         = 6  // Senha não preenchida.
	InvoiceResponseCodeInvalidBase64             = 7  // Codificação Base64 inválida.
	InvoiceResponseCodeInvalidPublicKey          = 8  // Cifra da chave pública inválida.
	InvoiceResponseCodeInvalidCreatedFormat      = 9  // Formato do campo Created inválido.
	InvoiceResponseCodeCredentialExpired         = 10 // Validade da credencial expirada.
	InvoiceResponseCodeInvalidSymmetricKey       = 11 // Chave simétrica inválida.
	InvoiceResponseCodeSymmetricKeyRepeated      = 12 // Chave simétrica repetida.
	InvoiceResponseCodeInvalidPasswordStructure  = 13 // Estrutura da senha inválida.
	InvoiceResponseCodeInvalidSessionKeyCreated  = 16 // Chave de sessão inválida. Não foi possível decifrar o campo Created.
	InvoiceResponseCodeInvalidSessionKeyPassword = 17 // Chave de sessão inválida. Não foi possível decifrar o campo Password.
	InvoiceResponseCodeInvalidSessionKeyDigest   = 18 // Chave de sessão inválida. Não foi possível decifrar o campo Digest.
	InvoiceResponseCodeRequestDateNotFilled      = 19 // Data de criação do pedido não preenchida.
	InvoiceResponseCodeRequestKeyNotFilled       = 20 // Chave do pedido não preenchida.
	InvoiceResponseCodeInvalidSOAPRequest        = 33 // Pedido SOAP inválido.
	InvoiceResponseCodePasswordValidationError   = 99 // Erro na validação da senha (Senha errada, acesso suspenso, etc.).

	// Document errors
	InvoiceResponseCodeInvalidDocument              = -1  // Documento inválido.
	InvoiceResponseCodeCancelledRegisteredDocument  = -2  // Anulação de documento registado.
	InvoiceResponseCodeDocumentRegisteredDiffValues = -3  // Documento registado com valores diferentes.
	InvoiceResponseCodeInvalidIssueDate             = -4  // Data de emissão inválida.
	InvoiceResponseCodeInvalidDocumentType          = -5  // Tipo de documento inválido.
	InvoiceResponseCodeInvalidDocumentStatus        = -6  // Estado de documento inválido.
	InvoiceResponseCodeInvalidAnomalousValues       = -7  // Documento inválido por valores anómalos.
	InvoiceResponseCodeInvalidExemptionReasonCode   = -8  // Código de motivo de isenção inválido.
	InvoiceResponseCodeIntegratedDocument           = -9  // Documento integrado (Rejeição retificada).
	InvoiceResponseCodeDocumentAlreadyRegistered    = -10 // O documento já foi registado pelo emitente.
	InvoiceResponseCodeNoPermissionForIssuerNIF     = -16 // Utilizador não tem permissões para registar documentos com o NIF de emitente indicado.
	InvoiceResponseCodeInvalidNIF                   = -17 // NIF inválido.
	InvoiceResponseCodeInvalidInputValues           = -28 // Valores de entrada inválidos.
	InvoiceResponseCodeInvalidCertificateNumber     = -39 // Número de certificado inválido.
	InvoiceResponseCodeMissingMandatoryParameters   = -50 // Parâmetros de entrada obrigatórios em falta.
	InvoiceResponseCodeInternalError                = -97 // Erro interno. Por favor tente mais tarde.
	InvoiceResponseCodeSystemError                  = -99 // Erro de sistema. Por favor volte a tentar mais tarde.
)

// Endpoints for the Invoices Web Service
var invoicesEndpoints = Endpoints{
	EnvTesting:    "https://servicos.portaldasfinancas.gov.pt:723/fatcorews/ws/",
	EnvProduction: "https://servicos.portaldasfinancas.gov.pt:423/fatcorews/ws/",
}

// InvoicesClient provides a client for the Invoices Web Service
type InvoicesClient Client

// RegisterInvoice represents a request to register an invoice
type RegisterInvoice struct {
	XMLName                   struct{}     `xml:"http://factemi.at.min_financas.pt/documents RegisterInvoiceRequest"`
	EFaturaMDVersion          string       `xml:"eFaturaMDVersion"`
	AuditFileVersion          string       `xml:"AuditFileVersion"`
	TaxRegistrationNumber     string       `xml:"TaxRegistrationNumber"`
	TaxEntity                 string       `xml:"TaxEntity"`
	SoftwareCertificateNumber int          `xml:"SoftwareCertificateNumber"`
	InvoiceData               *InvoiceData `xml:"InvoiceData"`
	RegistrationChannel       *Channel     `xml:"CanalRegisto,omitempty"`
}

// RegisterInvoiceResponse represents a response to a RegisterInvoiceRequest
type RegisterInvoiceResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents RegisterInvoiceResponse"`
	Response *Response `xml:"Response"`
}

// InvoiceData represents the invoice data structure
type InvoiceData struct {
	InvoiceNo              string            `xml:"InvoiceNo"`
	ATCUD                  string            `xml:"ATCUD"`
	InvoiceDate            string            `xml:"InvoiceDate"`
	InvoiceType            string            `xml:"InvoiceType"`
	SelfBillingIndicator   int               `xml:"SelfBillingIndicator"`
	CustomerTaxID          string            `xml:"CustomerTaxID"`
	CustomerTaxIDCountry   string            `xml:"CustomerTaxIDCountry"`
	DocumentStatus         *InvoiceStatus    `xml:"DocumentStatus"`
	HashCharacters         string            `xml:"HashCharacters"`
	CashVATSchemeIndicator int               `xml:"CashVATSchemeIndicator"`
	PaperLessIndicator     int               `xml:"PaperLessIndicator"`
	EACCode                string            `xml:"EACCode,omitempty"`
	SystemEntryDate        string            `xml:"SystemEntryDate"`
	LineSummaries          []*LineSummary    `xml:"LineSummary"`
	DocumentTotals         *DocumentTotals   `xml:"DocumentTotals"`
	WithholdingTaxes       []*WithholdingTax `xml:"WithholdingTax,omitempty"`
}

// InvoiceStatus represents the status of an invoice
type InvoiceStatus struct {
	InvoiceStatus     string `xml:"InvoiceStatus"`
	InvoiceStatusDate string `xml:"InvoiceStatusDate"`
}

// LineSummary represents a line summary in an invoice
type LineSummary struct {
	OrderReferences      []*OrderReferences `xml:"OrderReferences,omitempty"` //*
	TaxPointDate         string             `xml:"TaxPointDate"`              //*
	References           []string           `xml:"Reference,omitempty"`       //*
	DebitCreditIndicator string             `xml:"DebitCreditIndicator"`
	TotalTaxBase         string             `xml:"TotalTaxBase,omitempty"`
	Amount               string             `xml:"Amount,omitempty"`
	Tax                  *Tax               `xml:"Tax"`                        // *
	TaxExemptionCode     string             `xml:"TaxExemptionCode,omitempty"` //*
}

// Tax represents tax information
type Tax struct {
	TaxType          string `xml:"TaxType"`          // *
	TaxCountryRegion string `xml:"TaxCountryRegion"` // *
	TaxCode          string `xml:"TaxCode"`          // *
	TaxPercentage    string `xml:"TaxPercentage,omitempty"`
	TotalTaxAmount   string `xml:"TotalTaxAmount,omitempty"`
}

// DocumentTotals represents document totals
type DocumentTotals struct {
	TaxPayable string `xml:"TaxPayable"`
	NetTotal   string `xml:"NetTotal"`
	GrossTotal string `xml:"GrossTotal"`
}

// WithholdingTax represents withholding tax information
type WithholdingTax struct {
	WithholdingTaxType   string `xml:"WithholdingTaxType,omitempty"`
	WithholdingTaxAmount string `xml:"WithholdingTaxAmount"`
}

// Channel represents channel information
type Channel struct {
	System  string `xml:"Sistema"`
	Version string `xml:"Versao,omitempty"`
}

// Response represents a response from the webservice
type Response struct {
	ResponseCode  int    `xml:"CodigoResposta"`
	Message       string `xml:"Mensagem"`
	OperationDate string `xml:"DataOperacao"`
}

// InvoiceHeader represents the header information for an invoice
type InvoiceHeader struct {
	InvoiceNo            string `xml:"InvoiceNo"`
	ATCUD                string `xml:"ATCUD"`
	InvoiceDate          string `xml:"InvoiceDate"`
	InvoiceType          string `xml:"InvoiceType"`
	SelfBillingIndicator int    `xml:"SelfBillingIndicator"`
	CustomerTaxID        string `xml:"CustomerTaxID"`
	CustomerTaxIDCountry string `xml:"CustomerTaxIDCountry"`
}

// NewInvoiceStatus represents the new status for an invoice
type NewInvoiceStatus struct {
	InvoiceStatus     string `xml:"InvoiceStatus"`
	InvoiceStatusDate string `xml:"InvoiceStatusDate"`
}

// ChangeInvoiceStatus represents a request to change an invoice status
type ChangeInvoiceStatus struct {
	XMLName               struct{}          `xml:"http://factemi.at.min_financas.pt/documents ChangeInvoiceStatusRequest"`
	EFaturaMDVersion      string            `xml:"eFaturaMDVersion"`
	TaxRegistrationNumber string            `xml:"TaxRegistrationNumber"`
	InvoiceHeader         *InvoiceHeader    `xml:"InvoiceHeader"`
	InvoiceStatus         *NewInvoiceStatus `xml:"InvoiceStatus"`
	RegistrationChannel   *Channel          `xml:"CanalRegisto,omitempty"`
}

// ChangeInvoiceStatusResponse represents a response to a ChangeInvoiceStatusRequest
type ChangeInvoiceStatusResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents ChangeInvoiceStatusResponse"`
	Response *Response `xml:"Response"`
}

// Register registers an invoice
func (c *InvoicesClient) Register(req *RegisterInvoice) (*RegisterInvoiceResponse, error) {
	res := new(RegisterInvoiceResponse)
	return res, c.call("", req, res)
}

// ChangeStatus changes the status of an invoice
func (c *InvoicesClient) ChangeStatus(req *ChangeInvoiceStatus) (*ChangeInvoiceStatusResponse, error) {
	res := new(ChangeInvoiceStatusResponse)
	return res, c.call("", req, res)
}

func (c *InvoicesClient) call(action string, in, out any) error {
	return c.client.call(action, invoicesEndpoints, in, out)
}
