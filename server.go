package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

const helpMessage = "Welcome to brcep! Use https://brcep.herokuapp.com/cep/json"

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

func cepaberto(cep string) *CepAbertoResult {
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

func viacep(cep string) *ViaCepResult {

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

// return json brcep template
func apiWriteJSON(resp *brcepResult) string {

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		fmt.Printf("RegExp: %s", err)
	}
	cepClean := reg.ReplaceAllString(resp.Cep, "")

	jsonConvert := &brcepResult{
		Cep:         cepClean,
		Endereco:    resp.Endereco,
		Bairro:      resp.Bairro,
		Complemento: resp.Complemento,
		Cidade:      resp.Cidade,
		Uf:          resp.Uf,
		Ibge:        resp.Ibge,
		Latitude:    resp.Latitude,
		Longitude:   resp.Longitude,
	}

	conv, err := json.MarshalIndent(jsonConvert, "", "  ")
	if err != nil {
		fmt.Printf("apiWriteJSON: %s", err)
	}
	return string(conv)
}

// Mapping CepAbertoResult to brcepResult
func apiCepabertoJSON(resp *CepAbertoResult) string {

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

	return apiWriteJSON(&resultado)
}

// Mapping ViaCepResult to brcepResult
func apiViacepJSON(resp *ViaCepResult) string {

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

	return apiWriteJSON(&resultado)
}

// rewrite this please :D
func apiCep(c *gin.Context) {

	cep := c.Param("cep")
	c.Header("Content-Type", "application/json; charset=utf-8")

	resp := cepaberto(cep) // get CEPAberto
	if (resp != nil) && (resp.Cep != "") {
		c.String(200, apiCepabertoJSON(resp))
	} else {
		resp := viacep(cep) // get ViaCEP
		if (resp != nil) && (resp.Cep != "") {
			c.String(200, apiViacepJSON(resp))
		} else {
			c.JSON(500, gin.H{"status": "500"})
		}
	}
}

// 404 error showing start page
func startPage(c *gin.Context) {
	c.String(404, helpMessage)
}

func main() {

	fmt.Println(`   ___.                                  `)
	fmt.Println(`   \_ |_________   ____  ____ ______     `)
	fmt.Println(`    | __ \_  __ \_/ ___\/ __ \\____ \    `)
	fmt.Println(`    | \_\ \  | \/\  \__\  ___/|  |_> >   `)
	fmt.Println(`    |___  /__|    \___  >___  >   __/    `)
	fmt.Println(`        \/            \/    \/|__|       `)
	fmt.Printf("   %s\n\n", "http://github.com/leogregianin/brcep")

	gotenv.Load(".env")

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.Use(gin.ErrorLogger())

	router.NoRoute(startPage)
	router.GET("/:cep/json", apiCep)

	fmt.Println("starting server on", os.Getenv("PORT"))
	router.Run(":" + os.Getenv("PORT"))
}
