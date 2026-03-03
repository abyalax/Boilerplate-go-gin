package http

import (
	"github.com/abyalax/Boilerplate-go-gin/src/conf/logger"
	"github.com/gin-gonic/gin"
)

func MustGetBody[T any](c *gin.Context) T {
	log := logger.GetLogger()
	val, exists := c.Get("body")
	if !exists {
		log.Warn("body not found in context")
	}
	return val.(T)
}

func MustGetURI[T any](c *gin.Context) T {
	log := logger.GetLogger()
	val, exists := c.Get("uri")
	if !exists {
		log.Warn("uri not found in context")
	}
	return val.(T)
}
