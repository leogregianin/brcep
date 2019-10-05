package cepaberto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
)

const (
	cepAbertoApi = "cep-aberto"
	cepAbertoApiUrl = "http://www.cepaberto.com/api/v3/cep?cep=%s"
)

type CepAbertoApi struct {
	token string
}

// CepAbertoResult holds the result from cepaberto.com API
type CepAbertoResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Complemento string ``
	Cidade      struct {
		Nome string `json:"nome"`
	}
	Estado struct {
		Sigla string `json:"sigla"`
	}
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	UfDdd     struct {
		DDD string `json:"ddd"`
	}
	Unidade    string ``
	CodigoIbge struct {
		Ibge string `json:"ibge"`
	}
}

func NewCepAbertoApi(token string) *CepAbertoApi {
	return &CepAbertoApi{token}
}

func (api *CepAbertoApi) Fetch(_ *gin.Context, cep string) (*api.BrCepResult, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(cepAbertoApiUrl, url.QueryEscape(cep)), nil)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf(`Token token=%s`, api.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v %d", "invalid status code", resp.StatusCode)
	}

	var result CepAbertoResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("CepAbertoApi.Fetch %v", err)
	}

	return result.toBrCepResult(), nil
}

func (api *CepAbertoApi) Name() string {
	return cepAbertoApi
}

func (r CepAbertoResult) toBrCepResult() *api.BrCepResult {
	var result = new(api.BrCepResult)

	result.Cep = r.Cep
	result.Endereco = r.Logradouro
	result.Bairro = r.Bairro
	result.Complemento = r.Complemento
	result.Cidade = r.Cidade.Nome
	result.Uf = r.Estado.Sigla
	result.Latitude = r.Latitude
	result.Longitude = r.Longitude
	result.DDD = r.UfDdd.DDD
	result.Unidade = r.Unidade
	result.Ibge = r.CodigoIbge.Ibge

	return result
}
