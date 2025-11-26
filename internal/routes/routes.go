package routes

import (
	"real-erp-mebel/be/internal/utils"
	"real-erp-mebel/be/internal/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mengatur semua routes aplikasi
// Best Practice: Routes dipisahkan per modul untuk maintainability
func SetupRoutes(r *gin.Engine, hub *websocket.Hub) {
	// Health check (public)
	r.GET("/health", healthCheck)

	// WebSocket endpoint (public)
	r.GET("/ws", websocket.HandleWebSocket(hub))

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		SetupAuthRoutes(api)

		// Protected routes (require authentication)
		// Routes dipisahkan per modul untuk kemudahan maintenance
		SetupUserRoutes(api)

		// Add more module routes here:
		// SetupProductRoutes(api)
		// SetupStockRoutes(api)
		// SetupSalesRoutes(api)
		// SetupPurchaseOrderRoutes(api)
		// SetupFinanceRoutes(api)
		// SetupReportRoutes(api)
	}
}

// healthCheck adalah handler untuk health check endpoint
func healthCheck(c *gin.Context) {
	utils.OK(c, "Server is running", gin.H{
		"status": "ok",
	})
}
