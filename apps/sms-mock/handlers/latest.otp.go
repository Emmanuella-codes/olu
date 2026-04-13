package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LatestOTP(c *gin.Context) {
	phone := c.Param("phone")
	code := h.store.LatestOTP(phone)
	if code == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no OTP found for this number",
			"phone": phone,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"phone": phone,
		"otp":   code,
	})
}
