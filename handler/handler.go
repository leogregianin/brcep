package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
)

type CepHandler struct {
	PreferredApi string
	CepApis      map[string]api.Api
}

func (h *CepHandler) Handle(ctx *gin.Context) {
	cep := ctx.Param("cep")
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	preferredApi, ok := h.CepApis[h.PreferredApi]
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "preferred api not available"})
		return
	}

	result, err := preferredApi.Fetch(cep)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result.Sanitize()

	ctx.JSON(http.StatusOK, result)
}
