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
		logger.Debug("Auth middleware",
			zap.String("path", c.Request.URL.Path),
			zap.String("auth_header", authHeader),
		)

		if authHeader == "" {
			logger.Warn("Missing authorization header",
				zap.String("path", c.Request.URL.Path),
			)
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" or just token
		var tokenString string
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			// Format: "Bearer <token>"
			tokenString = parts[1]
		} else if len(parts) == 1 {
			// Format: just token (for backward compatibility)
			tokenString = parts[0]
		} else {
			logger.Warn("Invalid authorization header format",
				zap.String("path", c.Request.URL.Path),
				zap.String("header", authHeader),
			)
			utils.Unauthorized(c, "Invalid authorization header format. Use 'Bearer <token>' or just '<token>'")
			c.Abort()
			return
		}

		if tokenString == "" {
			logger.Warn("Empty token",
				zap.String("path", c.Request.URL.Path),
			)
			utils.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

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
				zap.String("error_type", err.Error()),
			)
			// Provide more specific error message
			errMsg := "Invalid or expired token"
			if strings.Contains(err.Error(), "expired") {
				errMsg = "Token has expired. Please login again."
			} else if strings.Contains(err.Error(), "signature") {
				errMsg = "Invalid token signature"
			}
			utils.Unauthorized(c, errMsg)
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
