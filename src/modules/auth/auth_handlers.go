package auth

import (
	"net/http"

	api "github.com/abyalax/Boilerplate-go-gin/src/conf/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *AuthService
	logger      *zap.Logger
}

// NewUserHandler creates a new UserHandler
func NewAuthHandler(
	authService *AuthService,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.Error(err)
		return
	}

	signedIn, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("login failed", zap.Error(err))
		c.Error(err)
		return
	}

	resp := api.Response[LoginResponse]{
		Message: "login successfully",
		Data: &LoginResponse{
			User:  signedIn.User,
			Token: signedIn.Token,
		},
	}

	c.JSON(http.StatusAccepted, resp)

}
