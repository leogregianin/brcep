package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/leogregianin/brcep/api"
)

type (
	// CepHandler provides a handler for a http.Server
	// used to interface the API implementations and
	// the server request
	CepHandler struct {
		PreferredAPI string
		CepApis      map[string]api.API
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

	preferredAPI, ok := h.CepApis[h.PreferredAPI]
	if !ok {
		h.renderJSON(w, http.StatusInternalServerError, nil, &responseError{"preferred api not available"})
		return
	}

	result, err := preferredAPI.Fetch(cep)
	if err != nil {
		h.renderJSON(w, http.StatusInternalServerError, nil, &responseError{err.Error()})
		return
	}

	result.Sanitize()

	h.renderJSON(w, http.StatusOK, result, nil)
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
