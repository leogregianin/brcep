package main

import (
	"net/http"
	"os"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/leogregianin/brcep/api"
	"github.com/leogregianin/brcep/api/cepaberto"
	"github.com/leogregianin/brcep/api/correios"
	"github.com/leogregianin/brcep/api/viacep"
	"github.com/leogregianin/brcep/config"
	"github.com/leogregianin/brcep/config/env"
	"github.com/leogregianin/brcep/config/flag"
	"github.com/leogregianin/brcep/handler"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	cfg, err := config.NewConfig([]config.Loader{
		flag.NewFlagLoader(),
		env.NewEnvLoader(),
	})

	if err != nil {
		log.Fatal(err)
	}

	logLevel, err := log.ParseLevel(cfg.CorreiosURL)
	if err == nil {
		log.SetLevel(logLevel)
	}

	var (
		cepApis = map[string]api.API{
			viacep.ID:   viacep.NewViaCepAPI(cfg.ViaCepURL, http.DefaultClient),
			correios.ID: correios.NewCorreiosAPI(cfg.CorreiosURL, http.DefaultClient),
		}
	)

	if len(cfg.CepAbertoToken) > 0 {
		cepApis[cepaberto.ID] = cepaberto.NewCepAbertoAPI(
			cfg.CepAbertoURL,
			cfg.CepAbertoToken,
			http.DefaultClient)
	}

	var cepHandler = &handler.CepHandler{
		PreferredAPI: cfg.PreferredAPI,
		CepApis:      cepApis,
		Cache:        cache.New(5*time.Minute, 10*time.Minute),
	}

	router := http.NewServeMux()
	router.HandleFunc("/", cepHandler.Handle)

	server := &http.Server{
		Addr:           cfg.Address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Infof("starting server at %s", cfg.Address)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
