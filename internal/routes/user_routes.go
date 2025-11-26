package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes mengatur routes untuk user management (protected routes)
func SetupUserRoutes(api *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler()

	// User routes (require authentication)
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		// Current user routes
		users.GET("/me", userHandler.GetCurrentUser)           // Get current user
		users.PUT("/me/password", userHandler.ChangePassword)  // Change password

		// User management routes
		users.GET("", userHandler.ListUsers)                    // List all users (owner/admin only)
		users.POST("", userHandler.CreateUser)                 // Create user (owner/admin only)
		users.GET("/:id", userHandler.GetUserByID)             // Get user by ID
		users.PUT("/:id", userHandler.UpdateUser)              // Update user
		users.DELETE("/:id", userHandler.DeleteUser)           // Delete user (owner/admin only)
	}
}

