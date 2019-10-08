package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
)

// CepHandler ..
type CepHandler struct {
	PreferredAPI string
	CepApis      map[string]api.API
}

// Handle handles the request ..
func (h *CepHandler) Handle(ctx *gin.Context) {
	cep := ctx.Param("cep")
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	preferredAPI, ok := h.CepApis[h.PreferredAPI]
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "preferred api not available"})
		return
	}

	result, err := preferredAPI.Fetch(cep)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result.Sanitize()

	ctx.JSON(http.StatusOK, result)
}
