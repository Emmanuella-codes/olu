package adminhandler

import (
	"net/http"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := admin.AdminRepo.GetAllStats(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
