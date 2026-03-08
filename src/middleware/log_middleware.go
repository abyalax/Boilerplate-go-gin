package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		sugar := logger.Sugar()
		sugar.Infof(
			"HTTP %s %s %d %s",
			method,
			path,
			statusCode,
			duration,
		)
	}
}
