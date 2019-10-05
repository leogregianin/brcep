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
	cfg.CepAbertoToken = os.Getenv("BRCEP_CEP_ABERTO_TOKEN")
}
