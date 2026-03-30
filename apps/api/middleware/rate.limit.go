package middleware

import (
	"net/http"
	"time"

	"github.com/emmanuella-codes/olu/cache"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, prefix string, max int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := cache.RateLimitKey(prefix, c.ClientIP())
		count, err := cache.IncreaseRateLimit(c.Request.Context(), rdb, key, window)
		if err != nil {
			// fail open — don't block requests if Redis is unavailable
			c.Next()
			return
		}
		if count > max {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
				"code":  "RATE_LIMITED",
			})
			return
		}
		c.Next()
	}
}
