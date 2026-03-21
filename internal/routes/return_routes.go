package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupReturnRoutes mengatur routes untuk retur penjualan dan retur pembelian
func SetupReturnRoutes(api *gin.RouterGroup, db *gorm.DB) {
	returnRepo := repositories.NewReturnRepository(db)
	stockRepo := repositories.NewStockRepository(db)
	batchRepo := repositories.NewStockBatchRepository(db)
	salesRepo := repositories.NewSalesRepository(db)

	returnService := services.NewReturnService(returnRepo, stockRepo, batchRepo, salesRepo)
	returnHandler := handlers.NewReturnHandler(returnService)

	// Retur Penjualan (Customer → Toko)
	salesReturns := api.Group("/sales-returns")
	salesReturns.Use(middleware.AuthMiddleware())
	{
		salesReturns.GET("", returnHandler.ListReturPenjualan)
		salesReturns.POST("", returnHandler.CreateReturPenjualan)
		salesReturns.GET("/:id", returnHandler.GetReturPenjualan)
		salesReturns.PATCH("/:id/approve", returnHandler.ApproveReturPenjualan) // Stok masuk kembali
	}

	// Retur Pembelian (Toko → Supplier/Vendor)
	purchaseReturns := api.Group("/purchase-returns")
	purchaseReturns.Use(middleware.AuthMiddleware())
	{
		purchaseReturns.GET("", returnHandler.ListReturPembelian)
		purchaseReturns.POST("", returnHandler.CreateReturPembelian)
		purchaseReturns.GET("/:id", returnHandler.GetReturPembelian)
		purchaseReturns.PATCH("/:id/approve", returnHandler.ApproveReturPembelian) // Stok keluar via FIFO
	}
}
