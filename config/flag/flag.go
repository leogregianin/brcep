package flag

import (
	"flag"

	"github.com/leogregianin/brcep/config"
)

type Loader struct{}

func NewFlagLoader() *Loader {
	return new(Loader)
}

func (f *Loader) Load(cfg *config.Config) {
	flag.StringVar(&cfg.Address, "address", "", "address to bind the server (default :8000)")
	flag.StringVar(&cfg.OperationMode, "mode", "", "mode which the server will operate (default debug)")
	flag.StringVar(&cfg.PreferredAPI, "preferred-api", "", "preferred API return (default viacep)")
	flag.StringVar(&cfg.CepAbertoToken, "cep-aberto-token", "", "token to use cep aberto API (if not provided will not be enabled)")

	flag.Parse()
}
