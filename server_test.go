package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

type testbrcepResult struct {
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
		"São Paulo",
	},
	{
		"22070011",
		"Rio de Janeiro",
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
			messageJSON := resp.Body.String()
			fmt.Printf(messageJSON)

			data := &testbrcepResult{}
			json.Unmarshal([]byte(messageJSON), data)
			actual := data.Cidade

			if actual != tt.expected {
				fmt.Printf("\n\n")
				fmt.Printf("Request: " + data.Cidade + "\n")
				fmt.Printf("Expected: " + tt.expected + "\n")
				fmt.Printf("\n\n")

				t.Fatalf("Test JSON Content - Expected to get %s but instead got %s\n", tt.expected, actual)
			}
		})
	}
}
