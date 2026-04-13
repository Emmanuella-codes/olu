package handlers

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
	"github.com/emmanuella-codes/olu/validator"
	"github.com/gin-gonic/gin"
)

var inboundSMSVotePattern = regexp.MustCompile(`^(?:VOTE\s+)?([A-Z]\d{1,2})$`)

type SMSWebhookHandler struct {
	voteSvc *services.VoteService
	smsSvc  *services.SMSService
}

type inboundSMSRequest struct {
	From string `json:"from" binding:"required"`
	Text string `json:"text" binding:"required"`
}

func NewSMSWebhookHandler(voteSvc *services.VoteService, smsSvc *services.SMSService) *SMSWebhookHandler {
	return &SMSWebhookHandler{voteSvc: voteSvc, smsSvc: smsSvc}
}

func (h *SMSWebhookHandler) InboundVote(c *gin.Context) {
	var req inboundSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and text are required",
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	if !validator.IsValidPhone(req.From) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid phone number",
			"code":  "INVALID_PHONE",
		})
		return
	}

	phone := validator.ToE164(req.From)
	candidateCode, err := parseInboundVoteCode(req.Text)
	if err != nil {
		h.queueSMSFormatRejection(c.Request.Context(), phone)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid sms vote format; use 'VOTE A1'",
			"code":  "INVALID_SMS_FORMAT",
		})
		return
	}

	result, err := h.voteSvc.CastVote(c.Request.Context(), services.CastVoteInput{
		Phone:         phone,
		CandidateCode: candidateCode,
		Channel:       models.SMSVoteChannel,
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
				"error": "failed to cast sms vote",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"candidate_code":  candidateCode,
			"candidate_name":  result.CandidateName,
			"confirmation_id": result.ConfirmationID,
			"phone":           maskPhone(phone),
			"channel":         models.SMSVoteChannel,
		},
	})
}

func parseInboundVoteCode(text string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(text))
	matches := inboundSMSVotePattern.FindStringSubmatch(normalized)
	if len(matches) != 2 {
		return "", errors.New("invalid sms vote format")
	}

	code := validator.NormalizeCode(matches[1])
	if !validator.IsValidCandidateCode(code) {
		return "", errors.New("invalid candidate code")
	}

	return code, nil
}

func (h *SMSWebhookHandler) queueSMSFormatRejection(ctx context.Context, phone string) {
	if h.smsSvc == nil {
		return
	}

	_ = h.smsSvc.QueueVoteRejection(ctx, phone, "Invalid vote format. Send 'VOTE A1' with your candidate code.")
}
