package middleware

import (
	"time"

	api "github.com/abyalax/Boilerplate-go-gin/src/conf/response"
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

func BindJSON[T any](logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn(api.InvalidRequestBody, zap.Error(err))
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("body", req)
		c.Next()
	}
}

func BindURI[T any](logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params T

		if err := c.ShouldBindUri(&params); err != nil {
			logger.Warn(api.InvalidRequestParams, zap.Error(err))
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("uri", params)
		c.Next()
	}
}
