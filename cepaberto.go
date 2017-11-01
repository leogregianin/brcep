package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func getCepaberto(cep string) *CepAbertoResult {
	cepAberto := url.QueryEscape(cep)

	url := fmt.Sprintf("http://www.cepaberto.com/api/v2/ceps.json?cep=%s", cepAberto)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf(`Token token="%s"`, os.Getenv("cepabertoToken")))
	if err != nil {
		return nil
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var resultado CepAbertoResult
	err = json.Unmarshal(content, &resultado)
	if err != nil {
		return nil
	}

	return &resultado
}

func mapCepabertoJSON(resp *CepAbertoResult) string {
	//func (cepaberto CepAbertoResult) cepabertoJSON() string {

	var resultado brcepResult

	resultado.Cep = resp.Cep
	resultado.Endereco = resp.Logradouro
	resultado.Bairro = resp.Bairro
	resultado.Complemento = resp.Complemento
	resultado.Cidade = resp.Cidade
	resultado.Uf = resp.Estado
	resultado.Ibge = resp.Ibge
	resultado.Latitude = resp.Latitude
	resultado.Longitude = resp.Longitude

	/*
		resultado.Cep = cepaberto.Cep
		resultado.Endereco = cepaberto.Logradouro
		resultado.Bairro = cepaberto.Bairro
		resultado.Complemento = cepaberto.Complemento
		resultado.Cidade = cepaberto.Cidade
		resultado.Uf = cepaberto.Estado
		resultado.Ibge = cepaberto.Ibge
		resultado.Latitude = cepaberto.Latitude
		resultado.Longitude = cepaberto.Longitude
	*/

	return brcepAPI(&resultado)
}
