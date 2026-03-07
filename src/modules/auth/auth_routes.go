package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	middlewares "github.com/abyalax/Boilerplate-go-gin/src/middleware"
)

type AuthModule struct {
	handler *AuthHandler
}

func NewAuthModule(db DBTX, logger *zap.Logger) *AuthModule {
	repo := New(db)
	service := NewAuthService(repo)
	handler := NewAuthHandler(service, logger)

	return &AuthModule{
		handler: handler,
	}
}

func (m *AuthModule) RegisterRoutes(r *gin.RouterGroup, logger *zap.Logger) {
	users := r.Group("/auth")

	users.POST("/login", middlewares.BindJSON[LoginRequest](logger), m.handler.Login)
	users.POST("/register", middlewares.BindJSON[RegisterRequest](logger), m.handler.Register)
}
