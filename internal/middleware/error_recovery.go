package middleware

import (
	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorRecovery adalah middleware untuk recovery dari panic
func ErrorRecovery() gin.HandlerFunc {
	logger := utils.GetLogger()

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		utils.InternalServerError(c, "Internal server error", nil)
		c.Abort()
	})
}
