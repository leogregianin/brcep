package flag

import (
	"flag"

	"github.com/leogregianin/brcep/config"
)

// Loader ..
type Loader struct{}

// NewFlagLoader returns new loader ..
func NewFlagLoader() *Loader {
	return new(Loader)
}

// Load will load configuration from flags ..
func (f *Loader) Load(cfg *config.Config) {
	flag.StringVar(&cfg.Address, "address", "", "address to bind the server (default :8000)")
	flag.StringVar(&cfg.OperationMode, "mode", "", "mode which the server will operate (default debug)")
	flag.StringVar(&cfg.PreferredAPI, "preferred-api", "", "preferred API return (default viacep)")
	flag.StringVar(&cfg.ViaCepURL, "via-cep-url", "", "viacep url (default http://viacep.com.br/)")
	flag.StringVar(&cfg.CepAbertoURL, "cep-aberto-url", "", "cepaberto url (default http://www.cepaberto.com/)")
	flag.StringVar(&cfg.CepAbertoToken, "cep-aberto-token", "", "token to use cep aberto API (if not provided will not be enabled)")
	flag.StringVar(&cfg.CorreiosURL, "correios-url", "", "correios url (default https://apps.correios.com.br/)")

	flag.Parse()
}
