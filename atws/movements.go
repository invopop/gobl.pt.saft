package atws

// Endpoints for the Movement Documents Web Service
var movementsEndpoints = Endpoints{
	EnvTesting:    "https://servicos.portaldasfinancas.gov.pt:701/sgdtws/documentosTransporte",
	EnvProduction: "https://servicos.portaldasfinancas.gov.pt:401/sgdtws/documentosTransporte",
}

// Enumerations
const (
	// Movement types
	MovementTypeDelivery    = "GR" // Guia de Remessa
	MovementTypeTransport   = "GT" // Guia de Transporte
	MovementTypeOwnAssets   = "GA" // Guia de Movimentação de Ativos Fixos Próprios
	MovementTypeConsignment = "GC" // Guia de Consignação
	MovementTypeReturn      = "GD" // Guia ou Nota de Devolução

	// Movement status
	MovementStatusNormal     = "N" // Normal
	MovementStatusThirdParty = "T" // Por conta de terceiros
	MovementStatusCancelled  = "A" // Documento anulado
)

// Result codes
// nolint:misspell // Texts in Portuguese should not be flagged as misspelled.
const (
	// Success
	MovementReturnCodeSuccess        = 0 // Documento registado com sucesso.
	MovementReturnCodeSuccessNoATCUD = 0 // O documento não apresenta um ATCUD válido.

	// Authentication
	MovementReturnCodeUserNotFilled             = 1  // Utilizador não preenchido.
	MovementReturnCodeUserSizeIncorrect         = 2  // Tamanho do utilizador incorreto.
	MovementReturnCodeInvalidNIF                = 3  // NIF inválido.
	MovementReturnCodeInvalidUserFormat         = 4  // Utilizador com formato inválido.
	MovementReturnCodeInvalidSubUserFormat      = 5  // Subutilizador com formato inválido.
	MovementReturnCodePasswordNotFilled         = 6  // Senha não preenchida.
	MovementReturnCodeInvalidBase64             = 7  // Codificação Base64 inválida.
	MovementReturnCodeInvalidPublicKey          = 8  // Cifra da chave pública inválida.
	MovementReturnCodeTimestampNotFilled        = 9  // Timestamp não preenchido.
	MovementReturnCodeInvalidTimestampFormat    = 10 // Formato do timestamp inválido.
	MovementReturnCodeCredentialExpired         = 11 // Validade da credencial expirada.
	MovementReturnCodeSymmetricKeyNotFilled     = 12 // Chave simétrica não preenchida.
	MovementReturnCodeSymmetricKeyRepeated      = 13 // Chave simétrica repetida.
	MovementReturnCodePasswordDigestNotFilled   = 14 // Digest da senha não preenchido.
	MovementReturnCodeDigestMismatch            = 15 // O Digest não corresponde ao esperado.
	MovementReturnCodeInvalidSessionKeyCreated  = 16 // Chave de sessão inválida. Não foi possível decifrar o campo Created.
	MovementReturnCodeInvalidSessionKeyPassword = 17 // Chave de sessão inválida. Não foi possível decifrar o campo Password.
	MovementReturnCodeInvalidSessionKeyDigest   = 18 // Chave de sessão inválida. Não foi possível decifrar o campo Digest.
	MovementReturnCodeRequestDateNotFilled      = 19 // Data de criação do pedido não preenchida.
	MovementReturnCodeRequestKeyNotFilled       = 20 // Chave do pedido não preenchida.
	MovementReturnCodeInvalidSOAPRequest        = 33 // Pedido SOAP inválido.
	MovementReturnCodeHeaderMissing             = 50 // Header inexistente ou vazio.
	MovementReturnCodeHeaderNIFNotFilled        = 51 // O NIF não está preenchido no Header.
	MovementReturnCodePermissionCheckFailed     = 52 // Não foi possível verificar se o utilizador tem permissões para aceder a esta operação.
	MovementReturnCodeNoPermission              = 54 // Não tem permissões para aceder a esta operação.
	MovementReturnCodePasswordValidation        = 99 // Erro na validação da senha (Senha errada, acesso suspenso, etc.).

	// Alerts
	MovementReturnCodeStartDateInPastAlert = -100 // A data início de transporte é inferior à data atual.

	// Errors
	MovementReturnCodeInvalidParameter              = -1     // Parâmetro de entrada inválido.
	MovementReturnCodeDocumentNotFound              = -2     // O número do Documento Global não corresponde a nenhum Documento de Transporte existente.
	MovementReturnCodeDocumentAlreadyExists         = -3     // Já foi inserido um Documento de Transporte com o número fornecido.
	MovementReturnCodeDocumentAlreadyCancelled      = -4     // Já foi anulado o Documento de Transporte com o número fornecido.
	MovementReturnCodeInvalidDocumentStatus         = -5     // Estado de Documento de Transporte inválido.
	MovementReturnCodeStartDateInPast               = -6     // A Data de início de transporte não pode ser anterior à data atual.
	MovementReturnCodeSenderNIFMismatch             = -7     // O NIF do Remetente não corresponde ao NIF do Header do pedido.
	MovementReturnCodeATCodeWithoutDocument         = -8     // O Código AT está preenchido mas não existe Documento de Transporte.
	MovementReturnCodeDocumentNumberMismatch        = -9     // O nº de Documento de Transporte é diferente do anteriormente fornecido.
	MovementReturnCodeSenderNoActivity              = -10    // O Remetente não tem atividade registada.
	MovementReturnCodeAcquirerNIFMismatch           = -11    // O NIF do Adquirente não corresponde ao anteriormente fornecido.
	MovementReturnCodeCannotVerifySenderActivity    = -12    // Não foi possível verificar se o Remetente tem atividade aberta.
	MovementReturnCodeCannotChangeAfterStartDate    = -13    // Não pode ser alterado um Documento de Transporte quando a Data de Início já decorreu.
	MovementReturnCodeInvalidGlobalDocReferences    = -14    // As referências de Documento Global para Documentos com tipo de emissão 4, têm que ser do tipo 3.
	MovementReturnCodeCannotAddReferences           = -15    // Não pode adicionar referências para outros documentos.
	MovementReturnCodeCannotChangeIssueDate         = -16    // Não pode alterar a data de emissão do documento.
	MovementReturnCodeCommunicationTooLate          = -17    // A comunicação de documentos só é permitida até 3 meses após o início do transporte.
	MovementReturnCodeCannotChangeDocumentType      = -18    // O tipo de documento não pode ser alterado.
	MovementReturnCodeSenderNIFMismatchPrevious     = -19    // O NIF do Remetente não corresponde ao anteriormente fornecido.
	MovementReturnCodeInvalidRecipientNIF           = -20    // O NIF do destinatário não é válido.
	MovementReturnCodeCommunicationTooEarly         = -21    // O documento deve ser comunicado com um máximo de 3 meses de antecedência.
	MovementReturnCodeInvalidATCUDFormat            = -22    // Código ATCUD inválido. ATCUD não preenchido ou formato inválido.
	MovementReturnCodeInvalidATCUDSeries            = -23    // Código ATCUD inválido. ATCUD não corresponde a série válida.
	MovementReturnCodeATCUDAlreadyAssociated        = -24    // Código ATCUD inválido. ATCUD já se encontra associado a um documento de transporte.
	MovementReturnCodeCannotVerifyATCUD             = -25    // Código ATCUD inválido. Não foi possível verificar se o ATCUD corresponde a uma série válida.
	MovementReturnCodeATCUDNotBelongsToDocument     = -26    // Código ATCUD inválido. ATCUD indicado não pertence ao documento.
	MovementReturnCodeATCUDInternalSeries           = -27    // Código ATCUD inválido. Código pertence a uma série interna e reservada da aplicação SGDT.
	MovementReturnCodeInternal                      = -99    // Erro interno.
	MovementReturnCodeInsufficientData              = -20001 // Dados recebidos insuficientes (todos os parâmetros são obrigatórios).
	MovementReturnCodeInvalidDataFormat             = -20002 // Dados recebidos com um formato inválido (número documento e/ou atcud).
	MovementReturnCodeInvalidNIFValidationPair      = -20005 // Par NIF/Código de Validação inválido para Classe/Tipo de Documento.
	MovementReturnCodeSeriesNotValidOnDate          = -20006 // Conjunto NIF/Código de Validação/Classe/Tipo de Documento existente.
	MovementReturnCodeInvalidSeriesValidationPair   = -20007 // Par Série/Código de Validação inválido.
	MovementReturnCodeDocumentATCUDSequenceMismatch = -20008 // Sequência do nº de Documento e do ATCUD diferentes.
	MovementReturnCodeInvalidDocumentSequence       = -20009 // Sequência do nº de Documento inválida para o intervalo comunicado.
	MovementReturnCodeUnexpected                    = -20000 // Erro não esperado inespecífico.
)

// MovementReturnSuccessCodes lists all success codes for movement documents
var MovementReturnSuccessCodes = []int{
	MovementReturnCodeSuccess,
	MovementReturnCodeSuccessNoATCUD,
	MovementReturnCodeStartDateInPastAlert,
}

// MovementsClient provides a client for the Movement Documents Web Service
type MovementsClient Client

// SendStockMovement represents a request to send a movement document
type SendStockMovement struct {
	XMLName               struct{} `xml:"https://servicos.portaldasfinancas.gov.pt/sgdtws/documentosTransporte/ dt:envioDocumentoTransporteRequestElem"`
	TaxRegistrationNumber string   `xml:"TaxRegistrationNumber"`
	CompanyName           string   `xml:"CompanyName"`
	CompanyAddress        *Address `xml:"CompanyAddress"`
	DocumentNumber        string   `xml:"DocumentNumber"`
	ATCUD                 string   `xml:"ATCUD,omitempty"`
	ATDocCodeID           string   `xml:"ATDocCodeID,omitempty"`
	MovementStatus        string   `xml:"MovementStatus"`
	MovementDate          string   `xml:"MovementDate"`
	MovementType          string   `xml:"MovementType"`
	CustomerTaxID         string   `xml:"CustomerTaxID,omitempty"`
	SupplierTaxID         string   `xml:"SupplierTaxID,omitempty"`
	CustomerAddress       *Address `xml:"CustomerAddress,omitempty"`
	CustomerName          string   `xml:"CustomerName,omitempty"`
	AddressTo             *Address `xml:"AddressTo,omitempty"`
	AddressFrom           *Address `xml:"AddressFrom"`
	MovementEndTime       string   `xml:"MovementEndTime,omitempty"`
	MovementStartTime     string   `xml:"MovementStartTime"`
	VehicleID             string   `xml:"VehicleID,omitempty"`
	Lines                 []*Line  `xml:"Line"`
}

// SendStockMovementResponse represents a response to a SendStockMovement request
type SendStockMovementResponse struct {
	XMLName        struct{}          `xml:"envioDocumentoTransporteResponseElem"`
	ResponseStatus []*ResponseStatus `xml:"ResponseStatus"`
	DocumentNumber string            `xml:"DocumentNumber,omitempty"`
	ATCUD          string            `xml:"ATCUD,omitempty"`
	ATDocCodeID    string            `xml:"ATDocCodeID,omitempty"`
}

// Line represents a line item in a movement document
type Line struct {
	OrderReferences    []*OrderReferences `xml:"OrderReferences,omitempty"`
	ProductDescription string             `xml:"ProductDescription"`
	Quantity           string             `xml:"Quantity"`
	UnitOfMeasure      string             `xml:"UnitOfMeasure"`
	UnitPrice          string             `xml:"UnitPrice"`
}

// OrderReferences represents references to an originating document
type OrderReferences struct {
	OriginatingON string `xml:"OriginatingON,omitempty"`
	OrderDate     string `xml:"OrderDate,omitempty"`
}

// Address represents a Portuguese address structure
type Address struct {
	Addressdetail string `xml:"Addressdetail,omitempty"`
	City          string `xml:"City,omitempty"`
	PostalCode    string `xml:"PostalCode,omitempty"`
	Country       string `xml:"Country,omitempty"`
}

// ResponseStatus represents the status of a response
type ResponseStatus struct {
	ReturnCode    int    `xml:"ReturnCode"`
	ReturnMessage string `xml:"ReturnMessage,omitempty"`
}

// Send sends a movement document
func (c *MovementsClient) Send(req *SendStockMovement) (*SendStockMovementResponse, error) {
	res := new(SendStockMovementResponse)
	return res, c.call("https://servicos.portaldasfinancas.gov.pt/sgdtws/documentosTransporte/", req, res)
}

func (c *MovementsClient) call(action string, in, out any) error {
	return c.client.call(action, movementsEndpoints, in, out)
}
