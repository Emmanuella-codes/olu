package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// rejects requests without X-Webhook-Secret header
func WebhookSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		incoming := c.GetHeader("X-Webhook-Secret")
		if incoming == "" || subtle.ConstantTimeCompare([]byte(incoming), []byte(secret)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid webhook secret",
			})
			return
		}
		c.Next()

	}
}
