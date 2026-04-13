package adminhandler

import (
	"net/http"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) AllCandidates(c *gin.Context) {
	candidates, err := admin.AdminRepo.GetAllCandidates(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve candidates")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve candidates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  candidates,
		"count": len(candidates),
	})
}
