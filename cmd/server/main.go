package main

import (
	"log"
	"strings"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Auto-migrate models
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Initialize router
	r := gin.Default()

	// Setup CORS
	corsConfig := cors.DefaultConfig()
	// Split origins by comma
	origins := strings.Split(config.AppConfig.CORS.AllowOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	corsConfig.AllowOrigins = origins
	corsConfig.AllowMethods = strings.Split(config.AppConfig.CORS.AllowMethods, ",")
	corsConfig.AllowHeaders = strings.Split(config.AppConfig.CORS.AllowHeaders, ",")
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Setup routes
	setupRoutes(r, hub)

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// WebSocket endpoint
	r.GET("/ws", websocket.HandleWebSocket(hub))

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/register", handlers.Register)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Add your protected routes here
			protected.GET("/users/me", handlers.GetCurrentUser)
		}
	}
}
