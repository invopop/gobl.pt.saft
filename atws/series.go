package atws

// Endpoints for the Series Web Service
var seriesEndpoints = Endpoints{
	EnvTesting:    "https://servicos.portaldasfinancas.gov.pt:722/SeriesWSService",
	EnvProduction: "https://servicos.portaldasfinancas.gov.pt:422/SeriesWSService",
}

// Enumerations
const (
	// Series types
	SeriesTypeNormal   = "N" // Normal
	SeriesTypeTraining = "F" // Formação
	SeriesTypeRecovery = "R" // Recuperação

	// Document classes
	DocClassInvoices   = "SI" // Faturas e documentos retificativos
	DocClassTransport  = "MG" // Documentos de Transporte
	DocClassConference = "WD" // Documentos de Conferência
	DocClassPayments   = "PY" // Recibos

	// Document types
	DocTypeInvoice            = "FT" // Fatura
	DocTypeSimplified         = "FS" // Fatura Simplificada
	DocTypeInvoiceReceipt     = "FR" // Fatura-Recibo
	DocTypeDebitNote          = "ND" // Nota de Débito
	DocTypeCreditNote         = "NC" // Nota de Crédito
	DocTypeDelivery           = "GR" // Guia de Remessa
	DocTypeTransport          = "GT" // Guia de Transporte
	DocTypeOwnAssets          = "GA" // Guia de Movimentação de Ativos Fixos Próprios
	DocTypeConsignment        = "GC" // Guia de Consignação
	DocTypeReturn             = "GD" // Guia ou Nota de Devolução
	DocTypeConsult            = "CM" // Consultas de Mesa
	DocTypeConsignmentCredit  = "CC" // Crédito de Consignação
	DocTypeConsignmentInvoice = "FC" // Fatura de Consignação
	DocTypeWorksheets         = "FO" // Folhas de Obra
	DocTypeOrderForm          = "NE" // Nota de Encomenda
	DocTypeOthers             = "OU" // Outros
	DocTypeBudgets            = "OR" // Orçamentos
	DocTypeProForma           = "PF" // Pró-Forma
	DocTypePremium            = "RP" // Prémio ou Recibo de Prémio
	DocTypeChargeback         = "RE" // Estorno ou Recibo de Estorno
	DocTypeCoInsurers         = "CS" // Imputação a Cosseguradoras
	DocTypeLeadCoInsurer      = "LD" // Imputação a Cosseguradora Líder
	DocTypeReinsurance        = "RA" // Resseguro Aceite
	DocTypeReceiptCash        = "RC" // Recibo emitido no âmbito do regime de IVA de Caixa
	DocTypeReceiptOther       = "RG" // Outros recibos emitidos

	// Processing methods
	ProcessingMethodSoftware = "PI" // Programa Informático de Faturação
	ProcessingMethodPortal   = "PF" // Portal das Finanças
	ProcessingMethodOther    = "OM" // Outros Meios Eletrónicos
)

// DocTypeToClass maps document types to document classes
var DocTypeToClass = map[string]string{
	DocTypeInvoice:            DocClassInvoices,
	DocTypeSimplified:         DocClassInvoices,
	DocTypeInvoiceReceipt:     DocClassInvoices,
	DocTypeDebitNote:          DocClassInvoices,
	DocTypeCreditNote:         DocClassInvoices,
	DocTypeDelivery:           DocClassTransport,
	DocTypeTransport:          DocClassTransport,
	DocTypeOwnAssets:          DocClassTransport,
	DocTypeConsignment:        DocClassTransport,
	DocTypeReturn:             DocClassTransport,
	DocTypeConsult:            DocClassConference,
	DocTypeConsignmentCredit:  DocClassConference,
	DocTypeConsignmentInvoice: DocClassConference,
	DocTypeWorksheets:         DocClassConference,
	DocTypeOrderForm:          DocClassConference,
	DocTypeOthers:             DocClassConference,
	DocTypeBudgets:            DocClassConference,
	DocTypeProForma:           DocClassConference,
	DocTypePremium:            DocClassConference,
	DocTypeChargeback:         DocClassConference,
	DocTypeCoInsurers:         DocClassConference,
	DocTypeLeadCoInsurer:      DocClassConference,
	DocTypeReinsurance:        DocClassConference,
	DocTypeReceiptCash:        DocClassPayments,
	DocTypeReceiptOther:       DocClassPayments,
}

// Result codes
const (
	// Common
	ResultTechnicalError = 5000 // Erro técnico não especificado.

	// Register
	ResultRegisterSuccess          = 2001 // Série registada com sucesso.
	ResultRegisterAlreadyExists    = 4001 // Não foi possível registar a Série. A Série indicada para o tipo de documento já foi registada.
	ResultRegisterMissingFields    = 4002 // Não foi possível registar a Série. Por favor preencha os campos assinalados.
	ResultRegisterInvalidFields    = 4003 // Não foi possível registar a Série. Por favor corrija os erros nos campos assinalados.
	ResultRegisterNoActivity       = 4043 // O NIF do sujeito passivo para o qual está a registar a Série não possui atividade aberta.
	ResultRegisterInvalidDocClass  = 4045 // O valor indicado na Classe do Documento deve corresponder a um valor predefinido
	ResultRegisterInvalidDocType   = 4046 // O valor indicado no Tipo de Documento deve corresponder a um valor válido de acordo com a Classe do Documento indicada.
	ResultRegisterInvalidStartDate = 4049 // O valor indicado no início previsto de utilização deve ser superior ou igual à data atual do sistema.

	// Fetch
	ResultFetchSuccess      = 2002 // Pesquisa efetuada com sucesso.
	ResultFetchInvalidDates = 4048 // A data que indicou no campo A (data fim) não pode ser inferior à data que indicou no campo DE (data início).

	// Finish
	ResultFinishSuccess        = 2004 // Série finalizada com sucesso.
	ResultFinishError          = 4005 // Não foi possível finalizar a Série.
	ResultFinishInvalidLastDoc = 4047 // O valor indicado no último documento emitido deve ser superior ao indicado no início da sequência.

	// Cancel
	ResultCancelSuccess       = 2003 // Série anulada com sucesso.
	ResultCancelError         = 4004 // Não foi possível anular a Série.
	ResultCancelInvalidReason = 4051 // O valor indicado no motivo de anulação deve corresponder a um valor pré-definido.
)

// SeriesClient provides a client for the Series Web Service
type SeriesClient Client

// RegisterSeries represents a request to register a series
type RegisterSeries struct {
	XMLName          struct{} `xml:"http://at.gov.pt/ at:registarSerie"`
	Code             string   `xml:"serie,omitempty"`
	Type             string   `xml:"tipoSerie,omitempty"`
	DocClass         string   `xml:"classeDoc,omitempty"`
	DocType          string   `xml:"tipoDoc,omitempty"`
	FirstSeqNum      int      `xml:"numInicialSeq,omitempty"`
	StartDate        string   `xml:"dataInicioPrevUtiliz,omitempty"`
	SWCertNum        int      `xml:"numCertSWFatur"`
	ProcessingMethod string   `xml:"meioProcessamento,omitempty"`
}

// RegisterSeriesResponse represents a response to a RegisterSeries request
type RegisterSeriesResponse struct {
	XMLName    struct{}    `xml:"registarSerieResponse"`
	SeriesInfo *SeriesInfo `xml:"registarSerieResp>infoSerie,omitempty"`
	ResultInfo *ResultInfo `xml:"registarSerieResp>infoResultOper,omitempty"`
}

// FetchSeries represents a request to fetch series
type FetchSeries struct {
	XMLName          struct{} `xml:"http://at.gov.pt/ at:consultarSeries"`
	Code             string   `xml:"serie,omitempty"`
	Type             string   `xml:"tipoSerie,omitempty"`
	DocClass         string   `xml:"classeDoc,omitempty"`
	DocType          string   `xml:"tipoDoc,omitempty"`
	ValidationCode   string   `xml:"codValidacaoSerie,omitempty"`
	StartDate        string   `xml:"dataRegistoDe,omitempty"`
	EndDate          string   `xml:"dataRegistoAte,omitempty"`
	Status           string   `xml:"estado,omitempty"`
	ProcessingMethod string   `xml:"meioProcessamento,omitempty"`
}

// FetchSeriesResponse represents a response to a FetchSeries request
type FetchSeriesResponse struct {
	XMLName    struct{}      `xml:"consultarSeriesResponse"`
	SeriesInfo []*SeriesInfo `xml:"consultarSeriesResp>infoSerie,omitempty"`
	ResultInfo *ResultInfo   `xml:"consultarSeriesResp>infoResultOper,omitempty"`
}

// FinishSeries represents a request to finish a series
type FinishSeries struct {
	XMLName          struct{} `xml:"http://at.gov.pt/ at:finalizarSerie"`
	Code             string   `xml:"serie,omitempty"`
	DocClass         string   `xml:"classeDoc,omitempty"`
	DocType          string   `xml:"tipoDoc,omitempty"`
	ValidationCode   string   `xml:"codValidacaoSerie,omitempty"`
	LastIssuedSeqNum int      `xml:"seqUltimoDocEmitido,omitempty"`
	Reason           string   `xml:"justificacao,omitempty"`
}

// FinishSeriesResponse represents a response to a FinishSeries request
type FinishSeriesResponse struct {
	XMLName    struct{}    `xml:"finalizarSerieResponse"`
	SeriesInfo *SeriesInfo `xml:"finalizarSerieResp>infoSerie,omitempty"`
	ResultInfo *ResultInfo `xml:"finalizarSerieResp>infoResultOper,omitempty"`
}

// CancelSeries represents a request to cancel a series
type CancelSeries struct {
	XMLName        struct{} `xml:"http://at.gov.pt/ at:anularSerie"`
	Code           string   `xml:"serie,omitempty"`
	DocClass       string   `xml:"classeDoc,omitempty"`
	DocType        string   `xml:"tipoDoc,omitempty"`
	ValidationCode string   `xml:"codValidacaoSerie,omitempty"`
	Reason         string   `xml:"motivo,omitempty"`
	NotUsedAck     bool     `xml:"naoUtilizadaAck,omitempty"`
}

// CancelSeriesResponse represents a response to a CancelSeries request
type CancelSeriesResponse struct {
	XMLName    struct{}    `xml:"anularSerieResponse"`
	SeriesInfo *SeriesInfo `xml:"anularSerieResp>infoSerie,omitempty"`
	ResultInfo *ResultInfo `xml:"anularSerieResp>infoResultOper,omitempty"`
}

// SeriesInfo represents information about a series
type SeriesInfo struct {
	Code             string `xml:"serie,omitempty"`
	Type             string `xml:"tipoSerie,omitempty"`
	DocClass         string `xml:"classeDoc,omitempty"`
	DocType          string `xml:"tipoDoc,omitempty"`
	FirstSeqNum      int    `xml:"numInicialSeq,omitempty"`
	LastSeqNum       int    `xml:"numFinalSeq,omitempty"`
	StartDate        string `xml:"dataInicioPrevUtiliz,omitempty"`
	LastIssuedSeqNum int    `xml:"seqUltimoDocEmitido,omitempty"`
	ProcessingMethod string `xml:"meioProcessamento,omitempty"`
	SWCertNum        int    `xml:"numCertSWFatur,omitempty"`
	ValidationCode   string `xml:"codValidacaoSerie,omitempty"`
	RegDate          string `xml:"dataRegisto,omitempty"`
	Status           string `xml:"estado,omitempty"`
	StatusReason     string `xml:"motivoEstado,omitempty"`
	Justification    string `xml:"justificacao,omitempty"`
	StatusDate       string `xml:"dataEstado,omitempty"`
	RegistrarNIF     string `xml:"nifComunicou,omitempty"`
}

// ResultInfo represents information about the result of an operation
type ResultInfo struct {
	Code ResultCode `xml:"codResultOper,omitempty"`
	Msg  string     `xml:"msgResultOper,omitempty"`
}

// ResultCode represents a result code
type ResultCode int

// Success returns wether the result code is a success code
func (c ResultCode) Success() bool {
	return c >= 2000 && c < 3000
}

// ClientError returns wether the result code is a client error code
func (c ResultCode) ClientError() bool {
	return c >= 4000 && c < 5000
}

// ServerError returns wether the result code is a server error code
func (c ResultCode) ServerError() bool {
	return c >= 5000
}

// Register registers a new series
func (c *SeriesClient) Register(req *RegisterSeries) (*RegisterSeriesResponse, error) {
	if req != nil {
		req.SWCertNum = c.swCertNum
	}
	res := new(RegisterSeriesResponse)
	return res, c.call("RegistrarSerie", req, res)
}

// Fetch fetches existing series
func (c *SeriesClient) Fetch(req *FetchSeries) (*FetchSeriesResponse, error) {
	res := new(FetchSeriesResponse)
	return res, c.call("ConsultarSeries", req, res)
}

// Finish finishes an existing series
func (c *SeriesClient) Finish(req *FinishSeries) (*FinishSeriesResponse, error) {
	res := new(FinishSeriesResponse)
	return res, c.call("FinalizarSerie", req, res)
}

// Cancel cancels an existing series
func (c *SeriesClient) Cancel(req *CancelSeries) (*CancelSeriesResponse, error) {
	res := new(CancelSeriesResponse)
	return res, c.call("AnularSerie", req, res)
}

func (c *SeriesClient) call(action string, in, out any) error {
	return c.client.call(action, seriesEndpoints, in, out)
}
