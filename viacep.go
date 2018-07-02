package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getViacep(cep string) *ViaCepResult {

	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get error")
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

	var resultado ViaCepResult
	err = json.Unmarshal(content, &resultado)
	if err != nil {
		fmt.Println("json error")
		return nil
	}

	return &resultado
}

func mapViacepJSON(resp *ViaCepResult) string {

	var resultado brcepResult

	resultado.Cep = resp.Cep
	resultado.Endereco = resp.Logradouro
	resultado.Bairro = resp.Bairro
	resultado.Complemento = resp.Complemento
	resultado.Cidade = resp.Cidade
	resultado.Uf = resp.Estado
	resultado.Latitude = resp.Latitude
	resultado.Longitude = resp.Longitude
	resultado.DDD = resp.Ibge
	resultado.Unidade = resp.Unidade
	resultado.Ibge = resp.Ibge

	return brcepAPI(&resultado)
}
