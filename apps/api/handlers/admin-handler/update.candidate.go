package adminhandler

import (
	"errors"
	"net/http"

	"github.com/emmanuella-codes/olu/dtos"
	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/emmanuella-codes/olu/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

func (h *AdminHandler) UpdateCandidate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid candidate ID"})
		return
	}

	var req dtos.UpdateCandidateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Code != nil {
		normalized := validator.NormalizeCode(*req.Code)
		if !validator.IsValidCandidateCode(normalized) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid candidate code",
				"code":  "INVALID_CANDIDATE_CODE",
			})
			return
		}
		req.Code = &normalized
	}

	if req.Party != nil {
		normalized := validator.NormalizeParty(*req.Party)
		if !validator.IsValidParty(normalized) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid party code",
				"code":  "INVALID_PARTY",
			})
			return
		}
		req.Party = &normalized
	}

	candidate, err := admin.AdminRepo.UpdateCandidate(c.Request.Context(), id, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "candidate code already exists",
				"code":  "DUPLICATE_CODE",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update candidate"})
		return
	}
	if candidate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	log.Info().Str("code", candidate.Code).Str("name", candidate.Name).Msg("candidate updated")
	c.JSON(http.StatusOK, gin.H{"data": candidate})
}
