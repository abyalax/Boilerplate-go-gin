package middleware

import (
	"errors"
	"net/http"

	user_error "github.com/abyalax/Boilerplate-go-gin/internal/users/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse is the standard error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler middleware untuk handle error dengan proper
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during request processing
		lastErr := c.Errors.Last()
		if lastErr != nil {
			handleError(c, logger, lastErr.Err)
		}
	}
}

// handleError maps domain/application errors ke HTTP response
func handleError(c *gin.Context, logger *zap.Logger, err error) {
	if err == nil {
		return
	}

	// Domain errors
	switch {
	case errors.Is(err, user_error.ErrUserNotFound):
		logger.Warn("user not found", zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
			Status:  http.StatusNotFound,
		})

	case errors.Is(err, user_error.ErrUserAlreadyExists):
		logger.Warn("user already exists", zap.Error(err))
		c.JSON(http.StatusConflict, ErrorResponse{
			Code:    "USER_ALREADY_EXISTS",
			Message: "User with this email already exists",
			Status:  http.StatusConflict,
		})

	case errors.Is(err, user_error.ErrInvalidEmail):
		logger.Warn("invalid email", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_EMAIL",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})

	case errors.Is(err, user_error.ErrInvalidName):
		logger.Warn("invalid name", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_NAME",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})

	case errors.Is(err, user_error.ErrInvalidPassword):
		logger.Warn("invalid password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_PASSWORD",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})

	default:
		// Validation or parameter errors
		logger.Error("request error", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
	}
}
