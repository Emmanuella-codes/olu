package adminhandler

import (
	"errors"
	"net/http"

	"github.com/emmanuella-codes/olu/dtos"
	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/emmanuella-codes/olu/validator"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) CreateCandidate(c *gin.Context) {
	var req dtos.CreateCandidateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Code = validator.NormalizeCode(req.Code)
	if !validator.IsValidCandidateCode(req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid candidate code",
		})
		return
	}

	req.Party = validator.NormalizeParty(req.Party)
	if !validator.IsValidParty(req.Party) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid party code",
			"code":  "INVALID_PARTY",
		})
		return
	}

	candidate, err := admin.AdminRepo.CreateCandidate(c.Request.Context(), req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "candidate code already exists",
				"code":  "DUPLICATE_CODE",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create candidate",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	log.Info().Str("code", candidate.Code).Str("name", candidate.Name).Msg("candidate created")
	c.JSON(http.StatusCreated, gin.H{"data": candidate})
}
