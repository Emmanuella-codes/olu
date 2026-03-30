package adminhandler

import (
	"errors"
	"net/http"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) DeactivateCandidate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid candidate ID"})
		return
	}

	if err := admin.AdminRepo.DeactivateCandidate(c.Request.Context(), id); err != nil {
		if errors.Is(err, admin.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
			return
		}
		log.Error().Err(err).Str("id", id.String()).Msg("failed to deactivate candidate")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate candidate"})
		return
	}

	log.Info().Str("id", id.String()).Msg("candidate deactivated")
	c.JSON(http.StatusOK, gin.H{"data": "candidate deactivated"})
}
