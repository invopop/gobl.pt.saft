package atws

// Endpoints for the Work Documents Web Service
var workDocumentsEndpoints = Endpoints{
	EnvTesting:    "https://servicos.portaldasfinancas.gov.pt:723/fatcorews/ws/",
	EnvProduction: "https://servicos.portaldasfinancas.gov.pt:423/fatcorews/ws/",
}

// Result codes
// nolint:misspell // Texts in Portuguese should not be flagged as misspelled.
const (
	// Success
	WorkResponseCodeSuccess = 0 // Operação efetuada com sucesso.

	// Authentication
	WorkResponseCodeUserNotFilled             = 1  // Utilizador não preenchido.
	WorkResponseCodeUserSizeIncorrect         = 2  // Tamanho do utilizador incorreto.
	WorkResponseCodeInvalidNIFAuth            = 3  // NIF inválido.
	WorkResponseCodeInvalidUserFormat         = 4  // Utilizador com formato inválido.
	WorkResponseCodeInvalidSubUserFormat      = 5  // Subutilizador com formato inválido.
	WorkResponseCodePasswordNotFilled         = 6  // Senha não preenchida.
	WorkResponseCodeInvalidBase64             = 7  // Codificação Base64 inválida.
	WorkResponseCodeInvalidPublicKey          = 8  // Cifra da chave pública inválida.
	WorkResponseCodeInvalidCreatedFormat      = 9  // Formato do campo Created inválido.
	WorkResponseCodeCredentialExpired         = 10 // Validade da credencial expirada.
	WorkResponseCodeInvalidSymmetricKey       = 11 // Chave simétrica inválida.
	WorkResponseCodeSymmetricKeyRepeated      = 12 // Chave simétrica repetida.
	WorkResponseCodeInvalidPasswordStructure  = 13 // Estrutura da senha inválida.
	WorkResponseCodeInvalidSessionKeyCreated  = 16 // Chave de sessão inválida. Não foi possível decifrar o campo Created.
	WorkResponseCodeInvalidSessionKeyPassword = 17 // Chave de sessão inválida. Não foi possível decifrar o campo Password.
	WorkResponseCodeInvalidSessionKeyDigest   = 18 // Chave de sessão inválida. Não foi possível decifrar o campo Digest.
	WorkResponseCodeRequestDateNotFilled      = 19 // Data de criação do pedido não preenchida.
	WorkResponseCodeRequestKeyNotFilled       = 20 // Chave do pedido não preenchida.
	WorkResponseCodeInvalidSOAPRequest        = 33 // Pedido SOAP inválido.
	WorkResponseCodePasswordValidationError   = 99 // Erro na validação da senha (Senha errada, acesso suspenso, etc.).

	// Document conference errors
	WorkResponseCodeNoPermissionForIssuerNIF   = -16 // Utilizador não tem permissões para registar documentos com o NIF de emitente indicado.
	WorkResponseCodeInvalidNIF                 = -17 // NIF inválido.
	WorkResponseCodeInvalidAnomalousValues     = -20 // Documento inválido por valores anómalos.
	WorkResponseCodeInvalidIssueDate           = -21 // Data de emissão inválida.
	WorkResponseCodeDocumentAlreadyRegistered  = -22 // O documento já foi registado pelo emitente.
	WorkResponseCodeInvalidInputValues         = -29 // Valores de entrada inválidos.
	WorkResponseCodeInvalidStatusDate          = -45 // A data de estado do documento não pode ser posterior à data atual nem anterior à última data de estado já comunicada.
	WorkResponseCodeMissingMandatoryParameters = -51 // Parâmetros de entrada obrigatórios em falta.
	WorkResponseCodeSystemError                = -99 // Erro de sistema. Por favor volte a tentar mais tarde.
)

// WorkDocumentsClient provides a client for the Work Documents Web Service
type WorkDocumentsClient Client

// RegisterWork represents a request to register a work document
type RegisterWork struct {
	XMLName                   struct{}  `xml:"http://factemi.at.min_financas.pt/documents RegisterWorkRequest"`
	EFaturaMDVersion          string    `xml:"eFaturaMDVersion"`
	AuditFileVersion          string    `xml:"AuditFileVersion"`
	TaxRegistrationNumber     string    `xml:"TaxRegistrationNumber"`
	TaxEntity                 string    `xml:"TaxEntity"`
	SoftwareCertificateNumber int       `xml:"SoftwareCertificateNumber"`
	WorkData                  *WorkData `xml:"WorkData"`
	RegistrationChannel       *Channel  `xml:"CanalRegisto,omitempty"`
}

// RegisterWorkResponse represents a response to a RegisterWorkRequest
type RegisterWorkResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents RegisterWorkResponse"`
	Response *Response `xml:"Response"`
}

// WorkData represents the work document data structure
type WorkData struct {
	DocumentNumber       string          `xml:"DocumentNumber"`
	ATCUD                string          `xml:"ATCUD"`
	WorkDate             string          `xml:"WorkDate"`
	WorkType             string          `xml:"WorkType"`
	CustomerTaxID        string          `xml:"CustomerTaxID"`
	CustomerTaxIDCountry string          `xml:"CustomerTaxIDCountry"`
	DocumentStatus       *WorkStatus     `xml:"DocumentStatus"`
	HashCharacters       string          `xml:"HashCharacters"`
	EACCode              string          `xml:"EACCode,omitempty"`
	SystemEntryDate      string          `xml:"SystemEntryDate"`
	LineSummaries        []*LineSummary  `xml:"LineSummary"`
	DocumentTotals       *DocumentTotals `xml:"DocumentTotals"`
}

// WorkStatus represents the status of a work document
type WorkStatus struct {
	WorkStatus     string `xml:"WorkStatus"`
	WorkStatusDate string `xml:"WorkStatusDate"`
}

// WorkHeader represents the header information for a work document
type WorkHeader struct {
	DocumentNumber       string `xml:"DocumentNumber"`
	ATCUD                string `xml:"ATCUD"`
	WorkDate             string `xml:"WorkDate"`
	WorkType             string `xml:"WorkType"`
	CustomerTaxID        string `xml:"CustomerTaxID"`
	CustomerTaxIDCountry string `xml:"CustomerTaxIDCountry"`
}

// NewWorkStatus represents the new status for a work document
type NewWorkStatus struct {
	WorkStatus     string `xml:"WorkStatus"`
	WorkStatusDate string `xml:"WorkStatusDate"`
}

// ChangeWorkStatus represents a request to change a work document status
type ChangeWorkStatus struct {
	XMLName               struct{}       `xml:"http://factemi.at.min_financas.pt/documents ChangeWorkStatusRequest"`
	EFaturaMDVersion      string         `xml:"eFaturaMDVersion"`
	TaxRegistrationNumber string         `xml:"TaxRegistrationNumber"`
	WorkHeader            *WorkHeader    `xml:"WorkHeader"`
	WorkStatus            *NewWorkStatus `xml:"WorkStatus"`
	RegistrationChannel   *Channel       `xml:"CanalRegisto,omitempty"`
}

// ChangeWorkStatusResponse represents a response to a ChangeWorkStatusRequest
type ChangeWorkStatusResponse struct {
	XMLName  struct{}  `xml:"http://factemi.at.min_financas.pt/documents ChangeWorkStatusResponse"`
	Response *Response `xml:"Response"`
}

// Register registers a work document
func (c *WorkDocumentsClient) Register(req *RegisterWork) (*RegisterWorkResponse, error) {
	res := new(RegisterWorkResponse)
	return res, c.call("", req, res)
}

// ChangeStatus changes the status of a work document
func (c *WorkDocumentsClient) ChangeStatus(req *ChangeWorkStatus) (*ChangeWorkStatusResponse, error) {
	res := new(ChangeWorkStatusResponse)
	return res, c.call("", req, res)
}

func (c *WorkDocumentsClient) call(action string, in, out any) error {
	return c.client.call(action, workDocumentsEndpoints, in, out)
}
