package middleware

import (
	"time"

	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger adalah middleware untuk logging request
func RequestLogger() gin.HandlerFunc {
	logger := utils.GetLogger()

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ip),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add user_id if available
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		if status >= 500 {
			logger.Error("HTTP Request", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP Request", fields...)
		} else {
			logger.Info("HTTP Request", fields...)
		}
	}
}

