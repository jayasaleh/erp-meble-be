package middleware

import (
	"strings"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Note: Import gin.H untuk rate limiter

// AuthMiddleware adalah middleware untuk autentikasi JWT
func AuthMiddleware() gin.HandlerFunc {
	logger := utils.GetLogger()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header",
				zap.String("path", c.Request.URL.Path),
			)
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format",
				zap.String("path", c.Request.URL.Path),
			)
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, utils.ErrTokenInvalid
			}
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		if err != nil {
			logger.Warn("Token parsing error",
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Warn("Invalid token",
				zap.String("path", c.Request.URL.Path),
			)
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Set user context
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(userID))
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}

		c.Next()
	}
}
