package handlers

import (
	"net/http"

	"github.com/emmanuella-codes/sms-mock/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store         *store.Store
	webhookSecret string
	apiBaseURL    string
}

func New(s *store.Store, webhookSecret, apiBaseURL string) *Handler {
	return &Handler{store: s, webhookSecret: webhookSecret, apiBaseURL: apiBaseURL}
}

type mockSMSRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
	SMS  string `json:"sms"`
}

func (h *Handler) Send(c *gin.Context) {
	var req mockSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.To == "" || req.SMS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to and sms are required"})
		return
	}

	from := req.From
	if from == "" {
		from = "Olu"
	}

	msg := h.store.Add(req.To, from, req.SMS)

	c.JSON(http.StatusOK, gin.H{
		"message_id": msg.ID,
		"message":    "Successfully Sent",
		"balance":    9999,
		"user":       "mock",
	})
}
