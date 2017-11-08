package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

const helpMessage = `Bem-vindo ao brcep!

Utilize desta forma: https://brcep.herokuapp.com/cep/json
                 ou: https://brcep.herokuapp.com/cep/xml

Por exemplo: https://brcep.herokuapp.com/78048000/json

Resultado: 
					
{
	"cep": "78048000",
	"endereco": "Avenida Miguel Sutil, de 5799/5800 a 7887/7888",
	"bairro": "Consil",
	"complemento": "",
	"cidade": "Cuiabá",
	"uf": "MT",
	"ibge": "5103403",
	"latitude": "-15.5786867",
	"longitude": "-56.0952081"
  }
`

// return json brcep template
func brcepAPI(resp *brcepResult) string {

	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		fmt.Printf("RegExp: %s", err)
	}

	jsonConvert := &brcepResult{
		Cep:         reg.ReplaceAllString(resp.Cep, ""),
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

func apiCepJSON(c *gin.Context) {

	cep := c.Param("cep")
	c.Header("Content-Type", "application/json; charset=utf-8")

	resp := getCepaberto(cep) // get CEPAberto
	if (resp != nil) && (resp.Cep != "") {
		c.String(200, mapCepabertoJSON(resp))
	} else {
		resp := getViacep(cep) // get ViaCEP
		if (resp != nil) && (resp.Cep != "") {
			c.String(200, mapViacepJSON(resp))
		} else {
			c.JSON(500, gin.H{"status": "500"})
		}
	}
}

func apiCepXML(c *gin.Context) {

	cep := c.Param("cep")
	c.Header("Content-Type", "application/xml; charset=utf-8")

	resp := getCepaberto(cep)
	if (resp != nil) && (resp.Cep != "") {

		data := mapCepabertoJSON(resp)
		var cepParser = brcep{}
		err := json.Unmarshal([]byte(data), &cepParser)
		if err != nil {
			fmt.Println("error:", err)
		}

		c.XML(200, gin.H{
			"Cep":         cepParser.Cep,
			"Endereco":    cepParser.Endereco,
			"Bairro":      cepParser.Bairro,
			"Complemento": cepParser.Complemento,
			"Cidade":      cepParser.Cidade,
			"Uf":          cepParser.Uf,
			"Ibge":        cepParser.Ibge,
			"Latitude":    cepParser.Latitude,
			"Longitude":   cepParser.Longitude,
		})

	} else {
		resp := getViacep(cep) // get ViaCEP

		if (resp != nil) && (resp.Cep != "") {
			data := mapViacepJSON(resp)
			var cepParser = brcep{}
			err := json.Unmarshal([]byte(data), &cepParser)
			if err != nil {
				fmt.Println("error:", err)
			}

			c.XML(200, gin.H{
				"Cep":         cepParser.Cep,
				"Endereco":    cepParser.Endereco,
				"Bairro":      cepParser.Bairro,
				"Complemento": cepParser.Complemento,
				"Cidade":      cepParser.Cidade,
				"Uf":          cepParser.Uf,
				"Ibge":        cepParser.Ibge,
				"Latitude":    cepParser.Latitude,
				"Longitude":   cepParser.Longitude,
			})

		} else {
			c.XML(500, gin.H{"status": "500"})
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
	router.GET("/:cep/json", apiCepJSON)
	router.GET("/:cep/xml", apiCepXML)
	//router.GET("/:cep/graphql", apiCepGraphQL)

	port := os.Getenv("PORT")
	fmt.Println("starting server on", port)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
