package viacep

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
)

const (
	viaCepApi = "viacep"
	viaCepApiUrl = "http://viacep.com.br/ws/%s/json/"
)

type ViaCepApi struct {}

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

func (api *ViaCepApi) Fetch(_ *gin.Context, cep string) (*api.BrCepResult, error) {
	resp, err := http.Get(fmt.Sprintf(viaCepApiUrl, cep))
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

func (api *ViaCepApi) Name() string {
	return viaCepApi
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
