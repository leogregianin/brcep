package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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

// Handle handles a request to /:cep/ which will extract the CEP
// from the URL.Path and fetch on multiple API implementations
// and return a common JSON result
func (h *CepHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	cep, respErr := h.parseCepFromPath(r.URL.Path)
	if respErr != nil {
		h.renderJSON(w, http.StatusBadRequest, respErr)
		return
	}

	preferredAPI, ok := h.CepApis[h.PreferredAPI]
	if !ok {
		h.renderJSON(w, http.StatusInternalServerError, &responseError{"preferred api not available"})
		return
	}

	result, err := preferredAPI.Fetch(cep)
	if err != nil {
		h.renderJSON(w, http.StatusInternalServerError, &responseError{Error: err.Error()})
		return
	}

	result.Sanitize()

	h.renderJSON(w, http.StatusOK, result)
}

func (h *CepHandler) renderJSON(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(j)
}

func (h *CepHandler) parseCepFromPath(path string) (string, *responseError) {
	var pathParts = strings.Split(path, "/")
	if len(pathParts) < 3 {
		return "", &responseError{"Invalid CEP provided"}
	}

	return pathParts[1], nil
}
