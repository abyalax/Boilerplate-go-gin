package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
