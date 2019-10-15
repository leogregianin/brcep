package correios

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&CorreiosSuite{})

type CorreiosSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *CorreiosSuite) TestNewCorreiosApiSetDefaultUrl(c *gc.C) {
	var correiosAPI = NewCorreiosAPI("", nil)
	c.Check(correiosAPI.url, gc.Equals, "https://apps.correios.com.br/")
	c.Check(correiosAPI.client, gc.NotNil)
}

func (s *CorreiosSuite) TestFetchShouldFailWhenInvalidStatusCode(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("", httpClient)
	_, err := correiosAPI.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *CorreiosSuite) TestFetchShouldFailWhenInvalidJSON(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("", httpClient)
	_, err := correiosAPI.Fetch("78048000")

	c.Check(err, gc.NotNil)
}

func (s *CorreiosSuite) TestFetchShouldSucceedWhenCorrectRemoteResponse(c *gc.C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
	<ns2:consultaCEPResponse xmlns:ns2="http://cliente.bean.master.sigep.bsb.correios.com.br/">
		<return>
			<bairro>Alvorada</bairro>
			<cep>78048000</cep>
			<cidade>Cuiabá</cidade>
			<complemento2>- de 5686 a 6588 - lado par</complemento2>
			<end>Avenida Miguel Sutil</end>
			<uf>MT</uf>
		</return>
	</ns2:consultaCEPResponse>
</soap:Body></soap:Envelope>`))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var correiosAPI = NewCorreiosAPI("http://localhost:8080/", httpClient)
	result, err := correiosAPI.Fetch("78048000")

	c.Check(err, gc.IsNil)
	c.Check(result, gc.NotNil)
	c.Check(result.Cep, gc.Equals, "78048000")
	c.Check(result.Endereco, gc.Equals, "Avenida Miguel Sutil")
	c.Check(result.Complemento, gc.Equals, "- de 5686 a 6588 - lado par")
	c.Check(result.Bairro, gc.Equals, "Alvorada")
	c.Check(result.Cidade, gc.Equals, "CuiabÃ¡")
	c.Check(result.Uf, gc.Equals, "MT")
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
