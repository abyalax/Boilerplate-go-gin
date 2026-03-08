package users

import (
	"net/http"

	api "github.com/abyalax/Boilerplate-go-gin/src/config/response"
	httpBind "github.com/abyalax/Boilerplate-go-gin/src/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *UserService
	logger      *zap.Logger
}

func NewUserHandler(
	userService *UserService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	req := httpBind.MustGetBody[CreateUserRequest](c)

	userID, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		c.Error(err)
		return
	}

	resp := api.Response[UserDTO]{
		Message: "user created successfully",
		Data: &UserDTO{
			ID:    userID,
			Name:  req.Name,
			Email: req.Email,
		},
	}

	h.logger.Info("user created", zap.Int32("id", userID))
	c.JSON(http.StatusCreated, resp)
}

// GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	params := httpBind.MustGetURI[UserIDParams](c)

	userDTO, err := h.userService.GetUser(c.Request.Context(), params.ID)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err), zap.Int32("id", params.ID))
		c.Error(err)
		return
	}

	resp := api.Response[UserDTO]{
		Message: "get user data successfully",
		Data: &UserDTO{
			ID:    userDTO.ID,
			Name:  userDTO.Name,
			Email: userDTO.Email,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	params := httpBind.MustGetURI[UserIDParams](c)
	req := httpBind.MustGetBody[UpdateUserRequest](c)

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), params.ID, &req)
	if err != nil {
		h.logger.Error("failed to update user", zap.Error(err), zap.Int32("id", params.ID))
		c.Error(err)
		return
	}

	resp := api.Response[UserDTO]{
		Message: "user updated successfully",
		Data: &UserDTO{
			ID:    updatedUser.ID,
			Name:  updatedUser.Name,
			Email: updatedUser.Email,
		},
	}

	h.logger.Info("user updated", zap.Int32("id", params.ID))
	c.JSON(http.StatusOK, resp)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	params := httpBind.MustGetURI[UserIDParams](c)

	err := h.userService.DeleteUser(c.Request.Context(), params.ID)
	if err != nil {
		h.logger.Error("failed to delete user", zap.Error(err), zap.Int32("id", params.ID))
		c.Error(err)
		return
	}

	h.logger.Info("user deleted", zap.Int32("id", params.ID))
	c.JSON(http.StatusNoContent, nil)
}

// GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		c.Error(err)
		return
	}

	listUsers := make([]UserDTO, len(users))
	for i, user := range users {
		listUsers[i] = UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	resp := api.Response[[]UserDTO]{
		Message: "get data users succcessfully",
		Data:    &listUsers,
	}

	c.JSON(http.StatusOK, resp)
}
