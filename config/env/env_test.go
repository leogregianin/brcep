package env

import (
	"os"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/leogregianin/brcep/config"
)

var _ = gc.Suite(&EnvLoaderSuite{})

type EnvLoaderSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *EnvLoaderSuite) TestNewEnvLoaderShouldLoadValuesIntoConfig(c *gc.C) {
	os.Setenv("BRCEP_ADDRESS", ":8080")
	os.Setenv("BRCEP_MODE", "test")
	os.Setenv("BRCEP_PREFERRED_API", "cep-aberto")
	os.Setenv("BRCEP_CEP_ABERTO_TOKEN", "token-sample")

	var (
		cfg    = &config.Config{}
		loader = NewEnvLoader()
	)

	loader.Load(cfg)
	c.Check(cfg.Address, gc.Equals, ":8080")
	c.Check(cfg.OperationMode, gc.Equals, "test")
	c.Check(cfg.PreferredAPI, gc.Equals, "cep-aberto")
	c.Check(cfg.CepAbertoToken, gc.Equals, "token-sample")
}
