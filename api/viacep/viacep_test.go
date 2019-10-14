package viacep

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&ViaCepSuite{})

type ViaCepSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *ViaCepSuite) TestNewViaCepApiSetDefaultUrl(c *gc.C) {
	var viaCepAPI = NewViaCepAPI("", nil)
	c.Check(viaCepAPI.url, gc.Equals, "http://viacep.com.br/")
	c.Check(viaCepAPI.client, gc.NotNil)
}

func (s *ViaCepSuite) TestFetchShouldFailWhenInvalidStatusCode(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var viaCepAPI = NewViaCepAPI("", httpClient)
	_, err := viaCepAPI.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *ViaCepSuite) TestFetchShouldFailWhenInvalidJSON(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var viaCepAPI = NewViaCepAPI("", httpClient)
	_, err := viaCepAPI.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *ViaCepSuite) TestFetchShouldSucceedWhenCorrectRemoteResponse(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "cep": "01001-000",
    "logradouro": "Praça da Sé",
    "complemento": "lado ímpar",
    "bairro": "Sé",
    "localidade": "São Paulo",
    "uf": "SP",
    "unidade": "",
    "ibge": "3550308",
    "gia": "1004"
}`))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var viaCepAPI = NewViaCepAPI("http://localhost:8080/", httpClient)
	result, err := viaCepAPI.Fetch("78048000")

	c.Check(err, gc.IsNil)
	c.Check(result, gc.NotNil)
	c.Check(result.Cep, gc.Equals, "01001-000")
	c.Check(result.Endereco, gc.Equals, "Praça da Sé")
	c.Check(result.Complemento, gc.Equals, "lado ímpar")
	c.Check(result.Bairro, gc.Equals, "Sé")
	c.Check(result.Cidade, gc.Equals, "São Paulo")
	c.Check(result.Uf, gc.Equals, "SP")
	c.Check(result.Ibge, gc.Equals, "3550308")
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}
