package cepaberto

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&CepAbertoSuite{})

type CepAbertoSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *CepAbertoSuite) TestNewCepAbertoApiSetDefaultUrl(c *gc.C) {
	var cepAbertoApi = NewCepAbertoApi("", "token-example", nil)

	c.Check(cepAbertoApi.url, gc.Equals, "http://www.cepaberto.com/")
	c.Check(cepAbertoApi.token, gc.Equals, "token-example")
	c.Check(cepAbertoApi.client, gc.NotNil)
}

func (s *CepAbertoSuite) TestFetchShouldFailWhenInvalidStatusCode(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var cepAbertoApi = NewCepAbertoApi("", "token-example", httpClient)
	_, err := cepAbertoApi.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *CepAbertoSuite) TestFetchShouldFailWhenInvalidJSON(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var cepAbertoApi = NewCepAbertoApi("", "token-example", httpClient)
	_, err := cepAbertoApi.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *CepAbertoSuite) TestFetchShouldSucceedWhenCorrectRemoteResponse(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "altitude": 7.0,
    "cep": "40010000",
    "latitude": "-12.967192",
    "longitude": "-38.5101976",
    "logradouro": "Avenida da França",
    "bairro": "Comércio",
    "cidade": {
        "ddd": 71,
        "ibge": "2927408",
        "nome": "Salvador"
    },
    "estado": {
        "sigla": "BA"
    }
}`))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var cepAbertoApi = NewCepAbertoApi("", "token-example", httpClient)
	result, err := cepAbertoApi.Fetch("78048000")

	c.Check(err, gc.IsNil)
	c.Check(result, gc.NotNil)
	c.Check(result.Cep, gc.Equals, "40010000")
	c.Check(result.Endereco, gc.Equals, "Avenida da França")
	c.Check(result.Bairro, gc.Equals, "Comércio")
	c.Check(result.Cidade, gc.Equals, "Salvador")
	c.Check(result.Uf, gc.Equals, "BA")
	c.Check(result.Latitude, gc.Equals, "-12.967192")
	c.Check(result.Longitude, gc.Equals, "-38.5101976")
	c.Check(result.DDD, gc.Equals, "71")
	c.Check(result.Ibge, gc.Equals, "2927408")
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
