package handlers

import (
	"net/http"

	"github.com/abyalax/Boilerplate-go-gin/internal/users/services"
	"github.com/abyalax/Boilerplate-go-gin/internal/users/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserResponse is the response DTO for user queries
type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
	logger      *zap.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	userService *services.UserService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req validation.CreateUserRequest

	// Gin binding validation using struct tags
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.Error(err)
		return
	}

	// Delegate to service
	userID, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("user created", zap.Int64("id", userID))
	c.JSON(http.StatusCreated, UserResponse{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	})
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")

	// Delegate to service (handles path parameter validation)
	userDTO, err := h.userService.GetUser(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err), zap.String("id", idStr))
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    userDTO.ID,
		Name:  userDTO.Name,
		Email: userDTO.Email,
	})
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")

	var req validation.UpdateUserRequest

	// Gin binding validation using struct tags
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.Error(err)
		return
	}

	// Delegate to service
	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), idStr, &req)
	if err != nil {
		h.logger.Error("failed to update user", zap.Error(err), zap.String("id", idStr))
		c.Error(err)
		return
	}

	h.logger.Info("user updated", zap.String("id", idStr))
	c.JSON(http.StatusOK, UserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	// Delegate to service (handles path parameter validation)
	err := h.userService.DeleteUser(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error("failed to delete user", zap.Error(err), zap.String("id", idStr))
		c.Error(err)
		return
	}

	h.logger.Info("user deleted", zap.String("id", idStr))
	c.JSON(http.StatusNoContent, nil)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Delegate to service
	users, err := h.userService.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		c.Error(err)
		return
	}

	// Map to response
	response := make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
