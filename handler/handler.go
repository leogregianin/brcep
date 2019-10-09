package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/leogregianin/brcep/api"
)

// CepHandler ..
type CepHandler struct {
	PreferredAPI string
	CepApis      map[string]api.API
}

type responseError struct {
	Error string `json:"error"`
}

func renderJSON(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(j)
}

// Handle handles the request ..
func (h *CepHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	cep, valid := h.parseCepFromPath(r.URL.Path)

	if !valid {
		renderJSON(w, http.StatusBadRequest, invalidCep)
		return
	}

	preferredApi, ok := h.CepApis[h.PreferredApi]
	if !ok {
		renderJSON(w, http.StatusInternalServerError, &responseError{"preferred api not available"})
		return
	}

	result, err := preferredAPI.Fetch(cep)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, &responseError{Error: err.Error()})
		return
	}

	result.Sanitize()

	renderJSON(w, http.StatusOK, result)
}

var invalidCep = responseError{"Invalid CEP provided"}

func (h *CepHandler) parseCepFromPath(path string) (string, bool) {
	if len(path) <= 1 {
		return "", false
	}

	var pathParts = strings.Split(path, "/")

	if len(pathParts) < 3 {
		return "", false
	}

	return pathParts[1], true

}
