package main

// CepAbertoResult : Retorno da API CepAberto.com
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

// Icepaberto : Interface da struct CepAbertoResult
type Icepaberto interface {
	cepabertoJSON() CepAbertoResult
}

// ViaCepResult : Retorno da API ViaCep.com.br
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

// Iviacep : Interface da struct ViaCepResult
type Iviacep interface {
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
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	DDD         string `json:"ddd"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
}

// Ibrcep : Interface da struct brcepResult
type Ibrcep interface {
	brcepJSON() brcepResult
}
