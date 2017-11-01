package main

// CepAbertoResult : Retorno da API CepAberto.com
type CepAbertoResult struct {
	Cep         string  `json:"cep"`
	Logradouro  string  `json:"logradouro"`
	Bairro      string  `json:"bairro"`
	Complemento string  `json:"complemento"`
	Cidade      string  `json:"cidade"`
	Estado      string  `json:"estado"`
	Latitude    string  `json:"latitude"`
	Longitude   string  `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	DDD         int     `json:"ddd"`
	Unidade     string  `json:"unidade"`
	Ibge        string  `json:"ibge"`
}

type cepaberto interface {
	cepabertoJSON() CepAbertoResult
}

// ViaCepResult : Retorno da API ViaCep.com.br
type ViaCepResult struct {
	Cep         string  `json:"cep"`
	Logradouro  string  `json:"logradouro"`
	Bairro      string  `json:"bairro"`
	Complemento string  `json:"complemento"`
	Localidade  string  `json:"localidade"`
	Uf          string  `json:"uf"`
	Latitude    string  `json:"latitude"`
	Longitude   string  `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	DDD         int     `json:"ddd"`
	Unidade     string  `json:"unidade"`
	Ibge        string  `json:"ibge"`
}

type viacep interface {
	viacepJSON() ViaCepResult
}

// brcepResult : Padronização do JSON do brcep
type brcepResult struct {
	Cep         string `json:"cep"`
	Endereco    string `json:"endereco"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"cidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

type brcep interface {
	brcepJSON() brcepResult
}
