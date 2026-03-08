package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func validateBearerScheme(c *gin.Context, authHeader string) []string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		c.Abort()
	}
	return parts
}

func validateExpiration(c *gin.Context, claims jwt.MapClaims) {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}
	}
}

func validateToken(c *gin.Context, cfg *env.Config, logger *zap.Logger, tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		logger.Error(err.Error())
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		}
		c.Abort()
	}

	return token
}

func AuthMiddleware(logger *zap.Logger, cfg *env.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := validateBearerScheme(c, authHeader)

		tokenString := parts[1]

		token := validateToken(c, cfg, logger, tokenString)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		validateExpiration(c, claims)

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])

		c.Next()
	}
}
