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

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", zap.Any("error", err))
				c.JSON(500, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
