package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupStockRoutes(r *gin.RouterGroup, db *gorm.DB) {
	stockRepo := repositories.NewStockRepository(db)
	batchRepo := repositories.NewStockBatchRepository(db)
	stockService := services.NewStockService(stockRepo, batchRepo)
	stockHandler := handlers.NewStockHandler(stockService)

	stocks := r.Group("/stocks")
	{
		stocks.GET("", stockHandler.GetStocks)
		stocks.GET("/history", stockHandler.GetStockHistory)
		stocks.POST("/in", stockHandler.CreateStockIn)
		stocks.POST("/out", stockHandler.CreateStockOut)
		stocks.POST("/adjustment", stockHandler.CreateStockOpname)
		stocks.POST("/transfer", stockHandler.CreateStockTransfer)
	}
}
