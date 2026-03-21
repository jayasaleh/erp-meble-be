package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupSalesRoutes mengatur semua routes untuk modul penjualan (Mode 1: POS)
func SetupSalesRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Initialize dependencies
	salesRepo := repositories.NewSalesRepository(db)
	stockRepo := repositories.NewStockRepository(db)
	batchRepo := repositories.NewStockBatchRepository(db)

	salesService := services.NewSalesService(salesRepo, stockRepo, batchRepo)
	salesHandler := handlers.NewSalesHandler(salesService)

	sales := api.Group("/sales")
	sales.Use(middleware.AuthMiddleware())
	{
		// Daftar & buat transaksi penjualan
		sales.GET("", salesHandler.ListSales)
		sales.POST("", salesHandler.CreateSale) // multipart/form-data: field "data" (JSON) + "bukti_bayar" (file, opsional)

		// Detail & invoice
		sales.GET("/:id", salesHandler.GetSale)
		sales.GET("/:id/invoice", salesHandler.GetInvoice)

		// Upload bukti bayar (untuk transaksi transfer yang belum upload saat transaksi)
		sales.POST("/:id/bukti-bayar", salesHandler.UploadBuktiBayar)
	}
}
