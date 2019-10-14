package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/leogregianin/brcep/api"
	"github.com/leogregianin/brcep/cache"
)

type (
	// CepHandler provides a handler for a http.Server
	// used to interface the API implementations and
	// the server request
	CepHandler struct {
		PreferredAPI string
		CepApis      map[string]api.API
		Cache        cache.Cache
	}
	channelResponse struct {
		ApiID  string
		Result *api.BrCepResult
		Error  error
	}
	responseError struct {
		Error string `json:"error"`
	}
)

func (r *responseError) String() string {
	return r.Error
}

// Handle handles a request to /:cep/ which will extract the CEP
// from the URL.Path and fetch on multiple API implementations
// and return a common JSON result
func (h *CepHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	log.WithFields(log.Fields{
		"remote_address": r.RemoteAddr,
		"method":         r.Method,
		"url":            r.URL.String(),
	}).Info("request received")

	cep, respErr := h.parseCepFromPath(r.URL.Path)
	if respErr != nil {
		h.renderJSON(w, http.StatusBadRequest, nil, respErr)
		return
	}

	if h.Cache != nil {
		if cached, found := h.Cache.Get(cep); found {
			h.renderJSON(w, http.StatusNotModified, cached.(*api.BrCepResult), nil)
			return
		}
	}

	var (
		ch         = make(chan *channelResponse)
		apiResults = make(map[string]*api.BrCepResult)
	)

	for id, currentAPI := range h.CepApis {
		go h.fetch(cep, id, currentAPI, ch)
	}

	for range h.CepApis {
		var current = <-ch
		if current.Error != nil {
			log.WithFields(log.Fields{
				"api_id": current.ApiID,
				"error":  current.Error.Error(),
			}).Error("API.Fetch error")
		} else {
			apiResults[current.ApiID] = current.Result
		}
	}

	if len(apiResults) <= 0 {
		h.renderJSON(w, http.StatusInternalServerError, nil, &responseError{"no API responded successfully"})
		return
	}

	result, ok := apiResults[h.PreferredAPI]
	if !ok {
		for _, currentResult := range apiResults {
			result = currentResult
			break
		}
	}

	result.Sanitize()

	if h.Cache != nil {
		h.Cache.Set(cep, result, 1*time.Hour)
	}

	h.renderJSON(w, http.StatusOK, result, nil)
}

func (h *CepHandler) fetch(cep string, id string, api api.API, ch chan<- *channelResponse) {
	log.WithFields(log.Fields{
		"api_id": id,
		"cep":    cep,
	}).Debug("dispatching API request")

	result, err := api.Fetch(cep)
	ch <- &channelResponse{
		ApiID:  id,
		Result: result,
		Error:  err,
	}
}

func (h *CepHandler) renderJSON(w http.ResponseWriter, code int, data interface{}, respErr *responseError) {
	w.WriteHeader(code)

	logFields := log.Fields{
		"response_code": code,
		"preferred_api": h.PreferredAPI,
	}

	if respErr != nil {
		logFields["error"] = respErr.String()

		if code >= http.StatusInternalServerError {
			log.WithFields(logFields).Error("response emit")
		} else if code >= http.StatusBadRequest {
			log.WithFields(logFields).Warn("response emit")
		}

		data = respErr
	} else {
		log.WithFields(logFields).Debug("response emit")
	}

	content, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	_, err = w.Write(content)
	if err != nil {
		log.Error(err)
	}

}

func (h *CepHandler) parseCepFromPath(path string) (string, *responseError) {
	var pathParts = strings.Split(path, "/")
	if len(pathParts) < 3 {
		return "", &responseError{"invalid CEP provided"}
	}

	return pathParts[1], nil
}
