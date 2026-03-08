package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abyalax/Boilerplate-go-gin/src/config/app"
	reject "github.com/abyalax/Boilerplate-go-gin/src/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse is the standard error response structure
type ErrorResponse struct {
	Message string `json:"message"`
}

var errorStatusMap = map[error]int{
	reject.UserNotFound:      http.StatusNotFound,
	reject.UserAlreadyExists: http.StatusConflict,

	reject.AuthEmailAlreadyExists: http.StatusConflict,
	reject.AuthEmailNotFound:      http.StatusNotFound,
	reject.AuthInvalidPassword:    http.StatusUnauthorized,

	reject.InvalidEmail:    http.StatusBadRequest,
	reject.InvalidName:     http.StatusBadRequest,
	reject.InvalidPassword: http.StatusBadRequest,

	reject.JWTFailedGenerateToken: http.StatusInternalServerError,
}

func ErrorMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()
		if lastErr != nil {
			handleError(c, logger, lastErr.Err)
		}
	}
}

func handleError(c *gin.Context, logger *zap.Logger, err error) {
	if err == nil {
		return
	}

	var appErr *app.AppError

	if errors.As(err, &appErr) {
		for targetErr, statusCode := range errorStatusMap {
			if errors.Is(appErr.Base, targetErr) {

				logFunc := logger.Warn
				if statusCode >= 500 {
					logFunc = logger.Error
				}

				logFunc(fmt.Sprintf("%d ", statusCode)+appErr.Base.Error(),
					zap.Error(appErr.Cause),
				)

				c.JSON(statusCode, ErrorResponse{
					Message: appErr.Base.Error(),
				})
				return
			}
		}
	}
	// Default Fallback (Internal Server Error)
	logger.Error("Unexpected system failure", zap.Error(err))
	c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal server error"})
}
