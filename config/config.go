package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	defaultAddress       = ":8000"
	defaultOperationMode = "debug"
	defaultPreferredAPI  = "viacep"
)

type Config struct {
	Address        string
	OperationMode  string
	PreferredAPI   string
	ViaCepUrl      string
	CepAbertoUrl   string
	CepAbertoToken string
}

// Loader defines the behaviour to load the values from
// multiple sources that will populate the config
type Loader interface {
	Load(*Config)
}

func NewConfig(loaders []Loader) (*Config, error) {
	if loaders == nil || len(loaders) <= 0 {
		return nil, fmt.Errorf("NewConfig %v", "no loader provided")
	}

	var cfg = new(Config)
	for _, l := range loaders {
		l.Load(cfg)
	}

	if len(cfg.Address) <= 0 {
		cfg.Address = defaultAddress
	}
	if len(cfg.OperationMode) <= 0 {
		cfg.OperationMode = defaultOperationMode
	}
	if len(cfg.PreferredAPI) <= 0 {
		cfg.PreferredAPI = defaultPreferredAPI
	}

	return cfg, nil
}

func (c *Config) GetGinOperationMode() string {
	switch c.OperationMode {
	case "test":
		return gin.TestMode
	case "debug":
		return gin.DebugMode
	default:
		return gin.ReleaseMode
	}
}
