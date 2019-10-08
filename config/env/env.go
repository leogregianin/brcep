package env

import (
	"os"

	"github.com/leogregianin/brcep/config"
)

// Loader ..
type Loader struct{}

// NewEnvLoader returns new loader ..
func NewEnvLoader() *Loader {
	return new(Loader)
}

// Load will load configuration from environment variables ..
func (l *Loader) Load(cfg *config.Config) {
	cfg.Address = os.Getenv("BRCEP_ADDRESS")
	cfg.OperationMode = os.Getenv("BRCEP_MODE")
	cfg.PreferredAPI = os.Getenv("BRCEP_PREFERRED_API")
	cfg.ViaCepURL = os.Getenv("BRCEP_VIACEP_URL")
	cfg.CepAbertoURL = os.Getenv("BRCEP_CEPABERTO_URL")
	cfg.CepAbertoToken = os.Getenv("BRCEP_CEPABERTO_TOKEN")
}
