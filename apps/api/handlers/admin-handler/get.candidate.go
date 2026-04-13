package adminhandler

import (
	"net/http"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) GetCandidate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid candidate id"})
		return
	}

	candidate, err := admin.AdminRepo.GetCandidateByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("id", id.String()).Msg("failed to retrieve candidate")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve candidate"})
		return
	}
	if candidate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": candidate})
}
