package users

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
	middlewares "github.com/abyalax/Boilerplate-go-gin/src/middleware"
)

type UserModule struct {
	handler *UserHandler
}

func NewUserModule(db DBTX, logger *zap.Logger) *UserModule {
	repo := New(db)
	service := NewUserService(repo)
	handler := NewUserHandler(service, logger)

	return &UserModule{
		handler: handler,
	}
}

func (m *UserModule) RegisterRoutes(r *gin.RouterGroup, logger *zap.Logger, cfg *env.Config) {
	users := r.Group("/users", middlewares.AuthMiddleware(logger, cfg))

	users.POST("", middlewares.BindJSON[CreateUserRequest](logger), m.handler.CreateUser)
	users.GET("", m.handler.ListUsers)
	users.GET("/:id", middlewares.BindURI[UserIDParams](logger), m.handler.GetUser)
	users.PUT("/:id", middlewares.BindURI[UserIDParams](logger), middlewares.BindJSON[UpdateUserRequest](logger), m.handler.UpdateUser)
	users.DELETE("/:id", middlewares.BindURI[UserIDParams](logger), m.handler.DeleteUser)
}
