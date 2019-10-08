package viacep

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leogregianin/brcep/api"
)

const (
	ID = "viacep"

	defaultViaCepApiUrl = "http://viacep.com.br/"
)

type ViaCepApi struct {
	url    string
	client *http.Client
}

// ViaCepResult holds the result from viacep.com.br API
type ViaCepResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"localidade"`
	Estado      string `json:"uf"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	DDD         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

func NewViaCepApi(url string, client *http.Client) *ViaCepApi {
	if len(url) <= 0 {
		url = defaultViaCepApiUrl
	}
	if client == nil {
		client = http.DefaultClient
	}

	return &ViaCepApi{url, client}
}

func (api *ViaCepApi) Fetch(cep string) (*api.BrCepResult, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(api.url+"ws/%s/json/", cep), nil)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v %d", "invalid status code", resp.StatusCode)
	}

	var result ViaCepResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("ViaCepApi.Fetch %v", err)
	}

	return result.toBrCepResult(), nil
}

func (r ViaCepResult) toBrCepResult() *api.BrCepResult {
	var result = new(api.BrCepResult)

	result.Cep = r.Cep
	result.Endereco = r.Logradouro
	result.Bairro = r.Bairro
	result.Complemento = r.Complemento
	result.Cidade = r.Cidade
	result.Uf = r.Estado
	result.Latitude = r.Latitude
	result.Longitude = r.Longitude
	result.DDD = r.Ibge
	result.Unidade = r.Unidade
	result.Ibge = r.Ibge

	return result
}
