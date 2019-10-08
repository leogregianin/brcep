package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
)

var _ = gc.Suite(&HandlerSuite{})

type HandlerSuite struct{}

type MockAPI struct {
	shouldErr    error
	shouldReturn *api.BrCepResult
}

func (a *MockAPI) Fetch(cep string) (*api.BrCepResult, error) {
	if a.shouldErr != nil {
		return nil, a.shouldErr
	}
	return a.shouldReturn, nil
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func setupRouter(h *CepHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/:cep/json", h.Handle)
	return r
}

func (s *HandlerSuite) TestHandleShouldReturnErrorIfNoPreferredAPIFound(c *gc.C) {
	var cepHandler = &CepHandler{
		PreferredAPI: "non-existent",
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/789000/json", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 500)
}

func (s *HandlerSuite) TestHandleShouldReturnErrorIfFetchReturnsError(c *gc.C) {
	var cepHandler = &CepHandler{
		PreferredAPI: "mock",
		CepApis: map[string]api.API{
			"mock": &MockAPI{
				shouldErr: errors.New("unknown error"),
			},
		},
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/789000/json", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 500)
}

func (s *HandlerSuite) TestHandleShouldSucceed(c *gc.C) {
	var cepHandler = &CepHandler{
		PreferredAPI: "mock",
		CepApis: map[string]api.API{
			"mock": &MockAPI{
				shouldErr: nil,
				shouldReturn: &api.BrCepResult{
					Cep:         "01001-000",
					Endereco:    "Praça da Sé",
					Complemento: "lado ímpar",
					Cidade:      "São Paulo",
					Uf:          "SP",
					Bairro:      "Sé",
					Ibge:        "3550308",
				},
			},
		},
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/789000/json", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 200)
	c.Check(w.Body.String(), gc.Equals, "{\"cep\":\"01001000\",\"endereco\":\"Praça da Sé\",\"bairro\":\"Sé\",\"complemento\":\"lado ímpar\",\"cidade\":\"São Paulo\",\"uf\":\"SP\",\"latitude\":\"\",\"longitude\":\"\",\"ddd\":\"\",\"unidade\":\"\",\"ibge\":\"3550308\"}\n")
}
