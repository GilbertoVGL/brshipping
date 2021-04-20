package rate

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseUrl     = "http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx/CalcPrecoPrazo"
	returnParam = "StrRetorno"
	returnType  = "xml"
)

type request struct {
	enterpriseContract    string
	enterprisePassword    string
	serviceType           string
	originPostalCode      uint64
	destinationPostalCode uint64
	weight                uint64
	shape                 uint64
	length                float64
	height                float64
	width                 float64
	diameter              float64
	deliverInHands        string
	declaredValue         float64
	notifyReceiving       string
}

type Response struct {
	Code                 int         `xml:"Servicos>cServico>Codigo"`
	Value                StringFloat `xml:"Servicos>cServico>Valor"`
	Deadline             int         `xml:"Servicos>cServico>PrazoEntrega"`
	DeliverInHandsValue  StringFloat `xml:"Servicos>cServico>ValorMaoPropria"`
	NotifyReceivingValue StringFloat `xml:"Servicos>cServico>ValorAvisoRecebimento"`
	DeclaredValue        StringFloat `xml:"Servicos>cServico>ValorValorDeclarado"`
	HomeDeliver          string      `xml:"Servicos>cServico>EntregaDomiciliar"`
	SaturdayDeliver      string      `xml:"Servicos>cServico>EntregaSabado"`
	Error                int         `xml:"Servicos>cServico>Erro"`
	ErrorMessage         string      `xml:"Servicos>cServico>MsgErro"`
	PlainValue           StringFloat `xml:"Servicos>cServico>ValorSemAdicionais"`
}

type StringFloat float64

const defaultService string = "04014"

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

func new(contract, password, service string,
	notify, inHands bool,
	origin, destiny, weight, shape uint64,
	length, height, width, diameter, declaredValue float64) request {

	return request{
		enterpriseContract:    contract,
		enterprisePassword:    password,
		serviceType:           validService(service),
		notifyReceiving:       boolStringField(notify),
		deliverInHands:        boolStringField(inHands),
		destinationPostalCode: origin,
		originPostalCode:      destiny,
		weight:                weight,
		shape:                 shape,
		length:                length,
		height:                height,
		width:                 width,
		diameter:              diameter,
		declaredValue:         declaredValue,
	}
}

func Calc(contract, password, service string,
	notify, inHands bool,
	origin, destiny, weight, shape uint64,
	length, height, width, diameter, declaredValue float64) (Response, error) {
	var parsedResponse Response

	request := new(contract, password, service, notify, inHands, origin, destiny, weight, shape, length, height, width, diameter, declaredValue)

	url, err := buildURL(request)
	if err != nil {
		return parsedResponse, err
	}

	res, err := fetch(url.String())
	if err != nil {
		return parsedResponse, err
	}

	err = parseRequest(res, &parsedResponse)

	return parsedResponse, err
}

func buildURL(request request) (*url.URL, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	buildQuery(u, request)
	log.Println(u.String())

	return u, nil
}

func buildQuery(u *url.URL, request request) {
	q := u.Query()
	q.Set(mapping["enterpriseContract"], request.enterpriseContract)
	q.Set(mapping["enterprisePassword"], request.enterprisePassword)
	q.Set(mapping["serviceType"], request.serviceType)
	q.Set(mapping["originPostalCode"], strconv.FormatUint(request.originPostalCode, 10))
	q.Set(mapping["destinationPostalCode"], strconv.FormatUint(request.destinationPostalCode, 10))
	q.Set(mapping["weight"], strconv.FormatUint(request.weight, 10))
	q.Set(mapping["shape"], strconv.FormatUint(request.shape, 10))
	q.Set(mapping["length"], strconv.FormatFloat(request.length, 'f', -1, 64))
	q.Set(mapping["height"], strconv.FormatFloat(request.height, 'f', -1, 64))
	q.Set(mapping["width"], strconv.FormatFloat(request.width, 'f', -1, 64))
	q.Set(mapping["diameter"], strconv.FormatFloat(request.diameter, 'f', -1, 64))
	q.Set(mapping["deliverInHands"], request.deliverInHands)
	q.Set(mapping["declaredValue"], strconv.FormatFloat(request.declaredValue, 'f', -1, 64))
	q.Set(mapping["notifyReceiving"], request.notifyReceiving)
	q.Set(returnParam, returnType)
	u.RawQuery = q.Encode()
}

func fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

func parseRequest(resp *http.Response, parsedResponse *Response) error {
	v, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = xml.Unmarshal(v, parsedResponse)

	return err
}

func (a *StringFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.ReplaceAll(s, ",", ".")
	float, _ := strconv.ParseFloat(s, 64)
	*a = StringFloat(float)

	return nil
}

func validService(v string) string {
	if v == "" {
		return defaultService
	} else {
		return v
	}
}

func boolStringField(v bool) string {
	if v {
		return "S"
	} else {
		return "N"
	}
}
