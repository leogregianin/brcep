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

	url := fmt.Sprintf("http://www.cepaberto.com/api/v3/cep?cep=%s", cepAberto)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf(`Token token=%s`, os.Getenv("cepabertoToken")))
	if err != nil {
		fmt.Println("Get error")
		return nil
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Do request error")
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("200 error")
		return nil
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Real error")
		return nil
	}

	var resultado CepAbertoResult
	err = json.Unmarshal(content, &resultado)
	if err != nil {
		fmt.Println("json error")
		return nil
	}

	return &resultado
}

func mapCepabertoJSON(resp *CepAbertoResult) string {

	var resultado brcepResult

	resultado.Cep = resp.Cep
	resultado.Endereco = resp.Logradouro
	resultado.Bairro = resp.Bairro
	resultado.Complemento = resp.Complemento
	resultado.Cidade = resp.Cidade.Nome
	resultado.Uf = resp.Estado.Sigla
	resultado.Latitude = resp.Latitude
	resultado.Longitude = resp.Longitude
	resultado.DDD = resp.UfDdd.DDD
	resultado.Unidade = resp.Unidade
	resultado.Ibge = resp.CodigoIbge.Ibge

	return brcepAPI(&resultado)
}
