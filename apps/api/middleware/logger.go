package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Str("status", strconv.Itoa(c.Writer.Status())).
			Int("bytes", c.Writer.Size()).
			Dur("duration", time.Since(start)).
			Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
			Msg("HTTP request")
	}
}
