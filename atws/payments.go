package atws

// Endpoints for the Payments Web Service
var paymentsEndpoints = Endpoints{
	EnvTesting:    "https://servicos.portaldasfinancas.gov.pt:723/fatcorews/ws/",
	EnvProduction: "https://servicos.portaldasfinancas.gov.pt:423/fatcorews/ws/",
}

// Result codes
// nolint:misspell // Texts in Portuguese should not be flagged as misspelled.
const (
	// Success
	PaymentResponseCodeSuccess = 0 // Operação efetuada com sucesso.

	// Authentication
	PaymentResponseCodeUserNotFilled             = 1  // Utilizador não preenchido.
	PaymentResponseCodeUserSizeIncorrect         = 2  // Tamanho do utilizador incorreto.
	PaymentResponseCodeInvalidUserFormat         = 4  // Utilizador com formato inválido.
	PaymentResponseCodeInvalidSubUserFormat      = 5  // Subutilizador com formato inválido.
	PaymentResponseCodePasswordNotFilled         = 6  // Senha não preenchida.
	PaymentResponseCodeInvalidBase64             = 7  // Codificação Base64 inválida.
	PaymentResponseCodeInvalidPublicKey          = 8  // Cifra da chave pública inválida.
	PaymentResponseCodeInvalidCreatedFormat      = 9  // Formato do campo Created inválido.
	PaymentResponseCodeCredentialExpired         = 10 // Validade da credencial expirada.
	PaymentResponseCodeInvalidSymmetricKey       = 11 // Chave simétrica inválida.
	PaymentResponseCodeSymmetricKeyRepeated      = 12 // Chave simétrica repetida.
	PaymentResponseCodeInvalidPasswordStructure  = 13 // Estrutura da senha inválida.
	PaymentResponseCodeInvalidSessionKeyCreated  = 16 // Chave de sessão inválida. Não foi possível decifrar o campo Created.
	PaymentResponseCodeInvalidSessionKeyPassword = 17 // Chave de sessão inválida. Não foi possível decifrar o campo Password.
	PaymentResponseCodeInvalidSessionKeyDigest   = 18 // Chave de sessão inválida. Não foi possível decifrar o campo Digest.
	PaymentResponseCodeRequestDateNotFilled      = 19 // Data de criação do pedido não preenchida.
	PaymentResponseCodeRequestKeyNotFilled       = 20 // Chave do pedido não preenchida.
	PaymentResponseCodeInvalidSOAPRequest        = 33 // Pedido SOAP inválido.
	PaymentResponseCodePasswordValidationError   = 99 // Erro na validação da senha (Senha errada, acesso suspenso, etc.).

	// Receipt errors
	PaymentResponseCodeNoPermissionForIssuerNIF  = -16 // Utilizador não tem permissões para registar documentos com o NIF de emitente indicado.
	PaymentResponseCodeInvalidNIF                = -17 // NIF inválido.
	PaymentResponseCodeInvalidAnomalousValues    = -23 // Documento inválido por valores anómalos.
	PaymentResponseCodeInvalidIssueDate          = -24 // Data de emissão inválida.
	PaymentResponseCodeDocumentAlreadyRegistered = -25 // O documento já foi registado pelo emitente.
	PaymentResponseCodeInvalidInputValues        = -30 // Valores de entrada inválidos.
	PaymentResponseCodeSystemError               = -99 // Erro de sistema. Por favor volte a tentar mais tarde.
)

// PaymentsClient provides a client for the Payments Web Service
type PaymentsClient Client

// RegisterPayment represents a request to register a payment
type RegisterPayment struct {
	XMLName                   struct{}     `xml:"http://factemi.at.min_financas.pt/documents RegisterPaymentRequest"`
	EFaturaMDVersion          string       `xml:"eFaturaMDVersion"`
	AuditFileVersion          string       `xml:"AuditFileVersion"`
	TaxRegistrationNumber     string       `xml:"TaxRegistrationNumber"`
	TaxEntity                 string       `xml:"TaxEntity"`
	SoftwareCertificateNumber int          `xml:"SoftwareCertificateNumber"`
	PaymentData               *PaymentData `xml:"PaymentData"`
	RegistrationChannel       *Channel     `xml:"CanalRegisto,omitempty"`
}

// RegisterPaymentResponse represents a response to a RegisterPaymentRequest
type RegisterPaymentResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents RegisterPaymentResponse"`
	Response *Response `xml:"Response"`
}

// PaymentData represents the payment data structure
type PaymentData struct {
	PaymentRefNo         string                `xml:"PaymentRefNo"`
	ATCUD                string                `xml:"ATCUD"`
	TransactionDate      string                `xml:"TransactionDate"`
	PaymentType          string                `xml:"PaymentType"`
	CustomerTaxID        string                `xml:"CustomerTaxID"`
	CustomerTaxIDCountry string                `xml:"CustomerTaxIDCountry"`
	DocumentStatus       *PaymentStatus        `xml:"DocumentStatus"`
	SystemEntryDate      string                `xml:"SystemEntryDate"`
	LineSummaries        []*PaymentLineSummary `xml:"LineSummary"`
	DocumentTotals       *DocumentTotals       `xml:"DocumentTotals"`
	WithholdingTaxes     []*WithholdingTax     `xml:"WithholdingTax,omitempty"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus struct {
	PaymentStatus     string `xml:"PaymentStatus"`
	PaymentStatusDate string `xml:"PaymentStatusDate"`
}

// PaymentLineSummary represents a line summary in a payment
type PaymentLineSummary struct {
	SourceDocumentIDs    []*SourceDocumentID `xml:"SourceDocumentID"`
	SettlementAmount     string              `xml:"SettlementAmount,omitempty"`
	DebitCreditIndicator string              `xml:"DebitCreditIndicator"`
	Amount               string              `xml:"Amount"`
	Tax                  *Tax                `xml:"Tax,omitempty"`
	TaxExemptionCode     string              `xml:"TaxExemptionCode,omitempty"`
}

// SourceDocumentID represents a source document ID
type SourceDocumentID struct {
	OriginatingON string `xml:"OriginatingON"`
	InvoiceDate   string `xml:"InvoiceDate"`
}

// PaymentHeader represents the header information for a payment
type PaymentHeader struct {
	PaymentRefNo         string `xml:"PaymentRefNo"`
	ATCUD                string `xml:"ATCUD"`
	TransactionDate      string `xml:"TransactionDate"`
	PaymentType          string `xml:"PaymentType"`
	CustomerTaxID        string `xml:"CustomerTaxID"`
	CustomerTaxIDCountry string `xml:"CustomerTaxIDCountry"`
}

// NewPaymentStatus represents the new status for a payment
type NewPaymentStatus struct {
	PaymentStatus     string `xml:"PaymentStatus"`
	PaymentStatusDate string `xml:"PaymentStatusDate"`
}

// ChangePaymentStatus represents a request to change a payment status
type ChangePaymentStatus struct {
	XMLName               struct{}          `xml:"http://factemi.at.min_financas.pt/documents ChangePaymentStatusRequest"`
	EFaturaMDVersion      string            `xml:"eFaturaMDVersion"`
	TaxRegistrationNumber string            `xml:"TaxRegistrationNumber"`
	PaymentHeader         *PaymentHeader    `xml:"PaymentHeader"`
	PaymentStatus         *NewPaymentStatus `xml:"PaymentStatus"`
	RegistrationChannel   *Channel          `xml:"CanalRegisto,omitempty"`
}

// ChangePaymentStatusResponse represents a response to a ChangePaymentStatusRequest
type ChangePaymentStatusResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents ChangePaymentStatusResponse"`
	Response *Response `xml:"Response"`
}

// Register registers a payment
func (c *PaymentsClient) Register(req *RegisterPayment) (*RegisterPaymentResponse, error) {
	res := new(RegisterPaymentResponse)
	return res, c.call("", req, res)
}

// ChangeStatus changes the status of a payment
func (c *PaymentsClient) ChangeStatus(req *ChangePaymentStatus) (*ChangePaymentStatusResponse, error) {
	res := new(ChangePaymentStatusResponse)
	return res, c.call("", req, res)
}

func (c *PaymentsClient) call(action string, in, out any) error {
	return c.client.call(action, paymentsEndpoints, in, out)
}
