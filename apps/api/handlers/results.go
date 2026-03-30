package handlers

import (
	"net/http"

	"github.com/emmanuella-codes/olu/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ResultsHandler struct {
	svc *services.ResultsService
}

func NewResultsHandler(svc *services.ResultsService) *ResultsHandler {
	return &ResultsHandler{svc: svc}
}

func (h *ResultsHandler) GetResults(c *gin.Context) {
	results, err := h.svc.GetResults(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve results")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve results",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
