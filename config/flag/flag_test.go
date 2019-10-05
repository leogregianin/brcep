package flag

import (
	"os"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/leogregianin/brcep/config"
)

var _ = gc.Suite(&FlagLoaderSuite{})

type FlagLoaderSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *FlagLoaderSuite) TestNewFlagLoaderShouldLoadValuesIntoConfig(c *gc.C) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"brcep",
		"-address=:8080",
		"-mode=test",
		"-preferred-api=cep-aberto",
		"-cep-aberto-token=token-sample",
	}

	var (
		cfg    = &config.Config{}
		loader = NewFlagLoader()
	)

	loader.Load(cfg)
	c.Check(cfg.Address, gc.Equals, ":8080")
	c.Check(cfg.OperationMode, gc.Equals, "test")
	c.Check(cfg.PreferredAPI, gc.Equals, "cep-aberto")
	c.Check(cfg.CepAbertoToken, gc.Equals, "token-sample")
}
