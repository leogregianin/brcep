package config

import (
	"fmt"
)

const (
	defaultAddress      = ":8000"
	defaultLogLevel     = "debug"
	defaultPreferredAPI = "viacep"
)

// Config hold configuration ..
type Config struct {
	Address        string
	LogLevel       string
	PreferredAPI   string
	ViaCepURL      string
	CepAbertoURL   string
	CepAbertoToken string
	CorreiosURL    string
}

// Loader defines the behaviour to load the values from
// multiple sources that will populate the config
type Loader interface {
	Load(*Config)
}

// NewConfig creates and returns new config structure ..
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
	if len(cfg.LogLevel) <= 0 {
		cfg.LogLevel = defaultLogLevel
	}
	if len(cfg.PreferredAPI) <= 0 {
		cfg.PreferredAPI = defaultPreferredAPI
	}

	return cfg, nil
}
