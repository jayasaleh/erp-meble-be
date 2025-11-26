package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/utils"
	"real-erp-mebel/be/internal/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mengatur semua routes aplikasi
func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()

	// Health check
	r.GET("/health", healthCheck)

	// WebSocket endpoint
	r.GET("/ws", websocket.HandleWebSocket(hub))

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		setupPublicRoutes(api, authHandler)

		// Protected routes (require authentication)
		setupProtectedRoutes(api, userHandler)
	}
}

// setupPublicRoutes mengatur public routes
func setupPublicRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}
}

// setupProtectedRoutes mengatur protected routes
func setupProtectedRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetCurrentUser)
		}

		// Add more protected routes here
		// Example:
		// products := protected.Group("/products")
		// {
		// 	products.GET("", productHandler.List)
		// 	products.POST("", productHandler.Create)
		// 	products.GET("/:id", productHandler.GetByID)
		// 	products.PUT("/:id", productHandler.Update)
		// 	products.DELETE("/:id", productHandler.Delete)
		// }
	}
}

// healthCheck adalah handler untuk health check endpoint
func healthCheck(c *gin.Context) {
	utils.OK(c, "Server is running", gin.H{
		"status": "ok",
	})
}

