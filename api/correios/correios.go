package correios

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/leogregianin/brcep/api"
)

const (
	// ID holds the identifier of this implementation
	ID                    = "correios"
	defaultCorreiosAPIURL = "https://apps.correios.com.br/"
)

// API holds the API implementation for correios.com.br
type API struct {
	url    string
	client *http.Client
}

type (
	responseEnvelope struct {
		Body responseBody `xml:"Body"`
	}
	responseBody struct {
		CepResponse responseReturn `xml:"consultaCEPResponse"`
	}
	responseReturn struct {
		Return responseResult `xml:"return"`
	}
	responseResult struct {
		Bairro       string `xml:"bairro"`
		Cep          string `xml:"cep"`
		Cidade       string `xml:"cidade"`
		Complemento2 string `xml:"complemento2"`
		End          string `xml:"end"`
		Id           int    `xml:"id"`
		Uf           string `xml:"uf"`
	}
)

// NewCorreiosAPI creates and return a new API struct
func NewCorreiosAPI(url string, client *http.Client) *API {
	if len(url) <= 0 {
		url = defaultCorreiosAPIURL
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &API{url, client}
}

// Fetch implements API.Fetch by fetching the correios.com.br and
// returning a BrCepResult
func (api *API) Fetch(cep string) (*api.BrCepResult, error) {
	envelope := `<soapenv:Envelope
xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
xmlns:cli="http://cliente.bean.master.sigep.bsb.correios.com.br/">
	<soapenv:Header/>
	<soapenv:Body>
		<cli:consultaCEP>
			<cep>` + cep + `s</cep>
		</cli:consultaCEP>
	</soapenv:Body>
</soapenv:Envelope>`

	req, err := http.NewRequest("POST", api.url+"SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl", bytes.NewBufferString(envelope))
	if err != nil {
		return nil, fmt.Errorf("CorreiosApi.Fetch %v", err)
	}

	req.Header.Add("Content-Type", "application/soap+xml; charset=iso-8859-1")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CorreiosApi.Fetch %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("CorreiosApi.Fetch %v %d", "invalid status code", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CorreiosApi.Fetch %v", err)
	}

	var result responseEnvelope
	err = xml.Unmarshal(api.toUtf8(body), &result)
	if err != nil {
		return nil, fmt.Errorf("CorreiosApi.Fetch %v", err)
	}

	return result.toBrCepResult(), nil
}

func (api *API) toUtf8(iso8859buff []byte) []byte {
	buf := make([]rune, len(iso8859buff))
	for i, b := range iso8859buff {
		buf[i] = rune(b)
	}
	return []byte(string(buf))
}

func (r responseEnvelope) toBrCepResult() *api.BrCepResult {
	var result = new(api.BrCepResult)

	result.Cep = r.Body.CepResponse.Return.Cep
	result.Endereco = r.Body.CepResponse.Return.End
	result.Bairro = r.Body.CepResponse.Return.Bairro
	result.Complemento = r.Body.CepResponse.Return.Complemento2
	result.Cidade = r.Body.CepResponse.Return.Cidade
	result.Uf = r.Body.CepResponse.Return.Uf

	return result
}
