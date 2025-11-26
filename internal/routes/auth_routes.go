package routes

import (
	"real-erp-mebel/be/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes mengatur routes untuk autentikasi (public routes)
func SetupAuthRoutes(api *gin.RouterGroup) {
	authHandler := handlers.NewAuthHandler()

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}
}

