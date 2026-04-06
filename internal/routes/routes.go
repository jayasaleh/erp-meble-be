package routes

import (
	"real-erp-mebel/be/internal/database"
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

	// Serve Static Files for Uploads
	r.Static("/uploads", "./uploads")

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		SetupAuthRoutes(api)

		// Protected routes (require authentication)
		// Routes dipisahkan per modul untuk kemudahan maintenance
		SetupUserRoutes(api)

		// Add more module routes here:
		SetupProductRoutes(api)
		SetupStockRoutes(api, database.DB)  // Registered Stock Routes
		SetupPemasokRoutes(api)             // Registered Supplier Routes
		SetupGudangRoutes(api)              // Registered Warehouse Routes
		SetupSalesRoutes(api, database.DB)  // Registered Sales Routes (Mode 1: POS)
		SetupReturnRoutes(api, database.DB) // Registered Return Routes (Sales Return + Purchase Return)
		SetupReportRoutes(api)              // Registered Report Routes (Sales by Period/Product/Customer)
		// SetupPurchaseOrderRoutes(api)
		// SetupFinanceRoutes(api)
	}
}

// healthCheck adalah handler untuk health check endpoint
// @Summary      Health check
// @Description  Check status server
// @Tags         health
// @Produce      json
// @Success      200  {object}  utils.Response
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	utils.OK(c, "Server is running", gin.H{
		"status": "ok",
	})
}
