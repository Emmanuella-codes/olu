package handlers

import (
	"errors"
	"net/http"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
	"github.com/emmanuella-codes/olu/validator"
	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	voteSvc *services.VoteService
}

func NewVoteHandler(voteSvc *services.VoteService) *VoteHandler {
	return &VoteHandler{voteSvc: voteSvc}
}

type castVoteRequest struct {
	CandidateCode string `json:"candidate_code" binding:"required"`
}

func (h *VoteHandler) Cast(c *gin.Context) {
	var req castVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "candidate_code is required",
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	req.CandidateCode = validator.NormalizeCode(req.CandidateCode)
	if !validator.IsValidCandidateCode(req.CandidateCode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid candidate code",
			"code":  "INVALID_CANDIDATE_CODE",
		})
		return
	}

	phoneValue, ok := c.Get("phone")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing verified phone in token",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	phone, ok := phoneValue.(string)
	if !ok || phone == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid verified phone in token",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	result, err := h.voteSvc.CastVote(c.Request.Context(), services.CastVoteInput{
		Phone:         phone,
		CandidateCode: req.CandidateCode,
		Channel:       models.WebVoteChannel,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCandidate):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "candidate not found",
				"code":  "INVALID_CANDIDATE",
			})
		case errors.Is(err, services.ErrAlreadyVoted):
			c.JSON(http.StatusConflict, gin.H{
				"error": "voter has already voted",
				"code":  "ALREADY_VOTED",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to cast vote",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
