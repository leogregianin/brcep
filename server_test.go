package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

var testCasesStatus = []struct {
	input  string
	status int
}{
	{
		"01311200",
		200,
	},
	{
		"08048000",
		500,
	},
	{
		"22070-011",
		200,
	},
}

var testCasesContent = []struct {
	input    string
	expected string
}{
	{
		"01311200",
		`{ 
  "cep": "01311200",
  "endereco": "Avenida Paulista, de 1047 a 1865 - lado ímpar",
  "bairro": "Bela Vista",
  "complemento": "",
  "cidade": "São Paulo",
  "uf": "SP",
  "latitude": "-23.5360299954",
  "longitude": "-46.622942654",
  "ddd": "",
  "unidade": "",
  "ibge": ""
}`,
	},
	{"22070011",
		`{ 
  "cep": "22070011",
  "endereco": "Avenida Nossa Senhora de Copacabana, de 1109 ao fim - lado ímpar",
  "bairro": "Copacabana",
  "complemento": "",
  "cidade": "Rio de Janeiro",
  "uf": "RJ",
  "latitude": "-22.9697777",
  "longitude": "-43.1868592",
  "ddd": "",
  "unidade": "",
  "ibge": ""
}`,
	},
}

func TestStatus(t *testing.T) {

	for _, tt := range testCasesStatus {

		gotenv.Load(".env")
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		t.Run(tt.input, func(t *testing.T) {

			router.GET("/:cep/json", apiCepJSON)

			req, err := http.NewRequest(http.MethodGet, "/"+tt.input+"/json", nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tt.status {
				t.Fatalf("Test status code - Expected to get status %d but instead got %d\n", tt.status, resp.Code)
			}
		})
	}
}

func TestContent(t *testing.T) {

	for _, tt := range testCasesContent {

		gotenv.Load(".env")
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		t.Run(tt.input, func(t *testing.T) {

			router.GET("/:cep/json", apiCepJSON)

			req, err := http.NewRequest(http.MethodGet, "/"+tt.input+"/json", nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			actual := resp.Body.String()
			if actual != tt.expected {
				fmt.Printf("---Request---\n")
				fmt.Printf(actual)
				fmt.Printf("---Expected---\n")
				fmt.Printf(tt.expected)
				fmt.Printf("\n\n")

				t.Fatalf("Test JSON Content - Expected to get %s but instead got %s\n", tt.expected, actual)
			}
		})
	}
}
