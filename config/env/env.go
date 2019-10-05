package env

import (
	"os"

	"github.com/leogregianin/brcep/config"
)

type Loader struct{}

func NewEnvLoader() *Loader {
	return new(Loader)
}

func (l *Loader) Load(cfg *config.Config) {
	cfg.Address = os.Getenv("BRCEP_ADDRESS")
	cfg.OperationMode = os.Getenv("BRCEP_MODE")
	cfg.PreferredAPI = os.Getenv("BRCEP_PREFERRED_API")
	cfg.ViaCepUrl = os.Getenv("BRCEP_VIACEP_URL")
	cfg.CepAbertoUrl = os.Getenv("BRCEP_CEPABERTO_URL")
	cfg.CepAbertoToken = os.Getenv("BRCEP_CEPABERTO_TOKEN")
}
