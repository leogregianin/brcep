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

	var resultado ViaCepResult
	err = json.Unmarshal(content, &resultado)
	if err != nil {
		return nil
	}

	return &resultado
}

func mapViacepJSON(resp *ViaCepResult) string {
	//func (viacep ViaCepResult) viacepJSON() string {

	var resultado brcepResult

	resultado.Cep = resp.Cep
	resultado.Endereco = resp.Logradouro
	resultado.Bairro = resp.Bairro
	resultado.Complemento = resp.Complemento
	resultado.Cidade = resp.Localidade
	resultado.Uf = resp.Uf
	resultado.Ibge = resp.Ibge
	resultado.Latitude = resp.Latitude
	resultado.Longitude = resp.Longitude

	/*
		resultado.Cep = viacep.Cep
		resultado.Endereco = viacep.Logradouro
		resultado.Bairro = viacep.Bairro
		resultado.Complemento = viacep.Complemento
		resultado.Cidade = viacep.Localidade
		resultado.Uf = viacep.Uf
		resultado.Ibge = viacep.Ibge
		resultado.Latitude = viacep.Latitude
		resultado.Longitude = viacep.Longitude
	*/

	return brcepAPI(&resultado)
}
