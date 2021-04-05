package precoprazo

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseUrl       = "http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx/CalcPrecoPrazo?"
	urlComplement = "StrRetorno=xml"
)

type request struct {
	enterpriseContract    string
	enterprisePassword    string
	serviceType           service
	originPostalCode      validPostalCode
	destinationPostalCode validPostalCode
	weight                validInt
	shape                 validInt
	length                validFloat
	height                validFloat
	width                 validFloat
	diameter              validFloat
	deliverInHands        customBool
	declaredValue         validFloat
	notifyReceiving       customBool
}

type Response struct {
	Codigo                int         `xml:"Servicos>cServico>Codigo"`
	Valor                 CustomFloat `xml:"Servicos>cServico>Valor"`
	PrazoEntrega          int         `xml:"Servicos>cServico>PrazoEntrega"`
	ValorMaoPropria       CustomFloat `xml:"Servicos>cServico>ValorMaoPropria"`
	ValorAvisoRecebimento CustomFloat `xml:"Servicos>cServico>ValorAvisoRecebimento"`
	ValorValorDeclarado   CustomFloat `xml:"Servicos>cServico>ValorValorDeclarado"`
	EntregaDomiciliar     string      `xml:"Servicos>cServico>EntregaDomiciliar"`
	EntregaSabado         string      `xml:"Servicos>cServico>EntregaSabado"`
	Erro                  int         `xml:"Servicos>cServico>Erro"`
	MsgErro               string      `xml:"Servicos>cServico>MsgErro"`
	ValorSemAdicionais    CustomFloat `xml:"Servicos>cServico>ValorSemAdicionais"`
}

type (
	CustomFloat     float64
	validInt        int
	validPostalCode int
	validFloat      float64
	stringBool      string
	customBool      string
	service         string
)

const defaultService service = "04014"

var mapping map[string]string = map[string]string{
	"enterpriseContract":    "nCdEmpresa",
	"enterprisePassword":    "sDsSenha",
	"serviceType":           "nCdServico",
	"originPostalCode":      "sCepOrigem",
	"destinationPostalCode": "sCepDestino",
	"weight":                "nVlPeso",
	"shape":                 "nCdFormato",
	"length":                "nVlComprimento",
	"height":                "nVlAltura",
	"width":                 "nVlLargura",
	"diameter":              "nVlDiametro",
	"deliverInHands":        "sCdMaoPropria",
	"declaredValue":         "nVlValorDeclarado",
	"notifyReceiving":       "sCdAvisoRecebimento",
}

func New(contract, password, service string,
	notify, inHands bool,
	origin, destiny, weight, shape int,
	length, height, width, diameter, declaredValue float64) request {

	request := request{}

	request.enterpriseContract = contract
	request.enterprisePassword = password
	request.serviceType.Set(service)
	request.notifyReceiving.Set(notify)
	request.deliverInHands.Set(inHands)
	request.destinationPostalCode.Set(origin)
	request.originPostalCode.Set(destiny)
	request.weight.Set(weight)
	request.shape.Set(shape)
	request.length.Set(length)
	request.height.Set(height)
	request.width.Set(width)
	request.diameter.Set(diameter)
	request.declaredValue.Set(declaredValue)

	return request
}

func Calc(request request) (Response, error) {
	var parsedResponse Response
	parsedURL := parseQueryParams(request)
	resp, err := http.Get(parsedURL)

	if err != nil {
		log.Println("http.Get err: ", err)
		return parsedResponse, err
	}

	err = parseRequest(resp, &parsedResponse)

	if err != nil {
		return parsedResponse, err
	}

	return parsedResponse, nil
}

func parseQueryParams(request request) string {
	var fullUrl string

	fullUrl += baseUrl
	fullUrl += fmt.Sprintf("%v=%v&", mapping["enterpriseContract"], request.enterpriseContract)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["enterprisePassword"], request.enterprisePassword)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["serviceType"], request.serviceType)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["originPostalCode"], request.originPostalCode)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["destinationPostalCode"], request.destinationPostalCode)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["weight"], request.weight)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["shape"], request.shape)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["length"], request.length)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["height"], request.height)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["width"], request.width)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["diameter"], request.diameter)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["deliverInHands"], request.deliverInHands)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["declaredValue"], request.declaredValue)
	fullUrl += fmt.Sprintf("%v=%v&", mapping["notifyReceiving"], request.notifyReceiving)
	fullUrl += urlComplement

	log.Println(fullUrl)
	return fullUrl
}

func parseRequest(resp *http.Response, parsedResponse *Response) error {
	v, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("io.ReadAll err: ", err)
		return err
	}

	err = xml.Unmarshal(v, parsedResponse)

	if err != nil {
		log.Println("xml.Unmarshal err: ", err)
		return err
	}

	return nil
}

func (a *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.ReplaceAll(s, ",", ".")
	float, _ := strconv.ParseFloat(s, 64)
	*a = CustomFloat(float)

	return nil
}

func (a *request) EnterpriseContract() string {
	return (*a).enterpriseContract
}

func (a *request) EnterprisePassword() string {
	return (*a).enterprisePassword
}

func (a *request) ServiceType() string {
	return string((*a).serviceType)
}

func (a *request) OriginPostalCode() int {
	return int((*a).originPostalCode)
}

func (a *request) DestinationPostalCode() int {
	return int((*a).destinationPostalCode)
}

func (a *request) Weight() int {
	return int((*a).weight)
}

func (a *request) Shape() int {
	return int((*a).shape)
}

func (a *request) Length() float64 {
	return float64((*a).length)
}

func (a *request) Height() float64 {
	return float64((*a).height)
}

func (a *request) Width() float64 {
	return float64((*a).width)
}

func (a *request) Diameter() float64 {
	return float64((*a).diameter)
}

func (a *request) DeliverInHands() bool {
	v := (*a).deliverInHands
	if v == "S" {
		return true
	}
	return false
}

func (a *request) DeclaredValue() float64 {
	return float64((*a).declaredValue)
}

func (a *request) NotifyReceiving() bool {
	v := (*a).notifyReceiving
	if v == "S" {
		return true
	}
	return false
}

func (field *validInt) Set(v int) {
	if v <= 0 {
		*field = -1
	} else {
		*field = validInt(v)
	}
}

func (field *service) Set(v string) {
	if v == "" {
		*field = defaultService
	} else {
		*field = service(v)
	}
}

func (field *validFloat) Set(v float64) {
	if v < 0 {
		*field = -1
	} else {
		*field = validFloat(v)
	}
}

func (field *customBool) Set(v bool) {
	if v {
		*field = customBool("S")
	} else {
		*field = customBool("N")
	}
}

func (field *validPostalCode) Set(v int) {
	if v <= 0 || len(fmt.Sprint(v)) != 8 {
		*field = -1
	} else {
		*field = validPostalCode(v)
	}
}
