package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserIDValidity extracts userID from context safely
func GetUserIDValidity(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	// Depending on how it's stored (float64 or uint)
	if id, ok := userID.(uint); ok {
		return id
	}
	if id, ok := userID.(float64); ok {
		return uint(id)
	}
	return 0
}
