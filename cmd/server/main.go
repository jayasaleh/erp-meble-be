package main

import (
	"strings"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/routes"
	"real-erp-mebel/be/internal/utils"
	"real-erp-mebel/be/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize logger
	if err := utils.InitLogger(config.AppConfig.GinMode); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	logger := utils.GetLogger()
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("mode", config.AppConfig.GinMode),
		zap.String("port", config.AppConfig.Port),
	)

	// Connect to database
	if err := database.Connect(); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	logger.Info("Database connected successfully")

	// Note: Database migration should be run separately using:
	// - go run cmd/migrate/main.go (normal migration)
	// - go run cmd/migrate-fresh/main.go (fresh migration)

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Initialize router
	r := gin.New()

	// Global middleware
	r.Use(middleware.ErrorRecovery()) // Recovery from panic
	r.Use(middleware.RequestLogger()) // Request logging

	// Initialize rate limiter
	if err := middleware.InitRateLimiter(); err != nil {
		logger.Warn("Failed to initialize rate limiter, continuing without it", zap.Error(err))
	} else {
		r.Use(middleware.RateLimitMiddleware()) // Rate limiting
	}

	// Setup CORS
	corsConfig := cors.DefaultConfig()
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
	routes.SetupRoutes(r, hub)

	// Start server
	port := ":" + config.AppConfig.Port
	logger.Info("Server starting", zap.String("port", port))
	if err := r.Run(port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
