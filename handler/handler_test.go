package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	gc "gopkg.in/check.v1"

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

func setupRouter(h *CepHandler) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/", h.Handle)
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

func (s *HandlerSuite) TestHandleShouldReturnErrorIfURLIsInvalid(c *gc.C) {
	var cepHandler = &CepHandler{
		PreferredAPI: "mock",
		CepApis: map[string]api.API{
			"mock": &MockAPI{
				shouldErr: nil,
			},
		},
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 400)
}

func (s *HandlerSuite) TestHandleShouldReturnErrorIfWithoutCepAndFormat(c *gc.C) {
	var cepHandler = &CepHandler{
		PreferredAPI: "mock",
		CepApis: map[string]api.API{
			"mock": &MockAPI{
				shouldErr: nil,
			},
		},
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/234423", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 400)
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
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/789000/json", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 200)
	c.Check(w.Body.String(), gc.Equals, "{\"cep\":\"01001000\",\"endereco\":\"Praça da Sé\",\"bairro\":\"Sé\",\"complemento\":\"lado ímpar\",\"cidade\":\"São Paulo\",\"uf\":\"SP\",\"latitude\":\"\",\"longitude\":\"\",\"ddd\":\"\",\"unidade\":\"\",\"ibge\":\"3550308\"}")

	cached, found := cepHandler.Cache.Get("789000")
	c.Check(found, gc.Equals, true)
	c.Check(cached, gc.NotNil)
}

func (s *HandlerSuite) TestHandleCachedShouldHitCacheAndSucceed(c *gc.C) {
	var cached = &api.BrCepResult{
		Cep:         "01001000",
		Endereco:    "Praça da Sé",
		Complemento: "lado ímpar",
		Cidade:      "São Paulo",
		Uf:          "SP",
		Bairro:      "Sé",
		Ibge:        "3550308",
	}

	var cepHandler = &CepHandler{
		PreferredAPI: "mock",
		CepApis: map[string]api.API{
			"mock": &MockAPI{
				shouldErr:    nil,
				shouldReturn: cached,
			},
		},
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}

	cepHandler.Cache.Set("789000", cached, 1*time.Hour)

	router := setupRouter(cepHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/789000/json", nil)
	router.ServeHTTP(w, req)

	c.Check(w.Code, gc.Equals, 304)
	c.Check(w.Body.String(), gc.Equals, "{\"cep\":\"01001000\",\"endereco\":\"Praça da Sé\",\"bairro\":\"Sé\",\"complemento\":\"lado ímpar\",\"cidade\":\"São Paulo\",\"uf\":\"SP\",\"latitude\":\"\",\"longitude\":\"\",\"ddd\":\"\",\"unidade\":\"\",\"ibge\":\"3550308\"}")
}
