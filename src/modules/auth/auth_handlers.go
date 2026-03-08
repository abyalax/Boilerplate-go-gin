package auth

import (
	"net/http"

	api "github.com/abyalax/Boilerplate-go-gin/src/config/response"
	httpBind "github.com/abyalax/Boilerplate-go-gin/src/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *AuthService
	logger      *zap.Logger
}

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
	req := httpBind.MustGetBody[LoginRequest](c)

	signedIn, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Warn("login failed for email " + req.Email)
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

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	req := httpBind.MustGetBody[RegisterRequest](c)

	registered, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("registration failed", zap.Error(err))
		c.Error(err)
		return
	}

	resp := api.Response[RegisterResponse]{
		Message: "registration successfully",
		Data:    registered,
	}

	c.JSON(http.StatusCreated, resp)

}
