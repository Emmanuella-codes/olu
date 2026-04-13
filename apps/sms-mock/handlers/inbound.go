package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type inboundSMSRequest struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func (h *Handler) Inbound(c *gin.Context) {
	var req inboundSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.From) == "" || strings.TrimSpace(req.Text) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and text are required"})
		return
	}

	h.store.AddInbound(req.From, req.Text)

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode request"})
		return
	}

	apiURL := strings.TrimRight(h.apiBaseURL, "/") + "/api/v1/webhooks/sms/inbound"
	httpReq, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		apiURL,
		bytes.NewReader(body),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upstream request"})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Webhook-Secret", h.webhookSecret)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to reach api webhook: %v", err)})
		return
	}
	defer resp.Body.Close()

	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "api webhook returned invalid response"})
		return
	}

	c.JSON(resp.StatusCode, gin.H{
		"forwarded": true,
		"payload":   payload,
	})
}
