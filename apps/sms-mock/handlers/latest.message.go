package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LatestMessage(c *gin.Context) {
	phone := c.Param("phone")
	msg := h.store.Latest(phone)
	if msg == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no messages found for this number",
			"phone": phone,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"phone":   phone,
		"message": msg,
	})
}
