package middleware

import (
	api "github.com/abyalax/Boilerplate-go-gin/src/config/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
