package routes

import (
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupReportRoutes mengatur routes untuk laporan penjualan
func SetupReportRoutes(api *gin.RouterGroup) {
	reportService := services.NewReportService(database.DB)
	reportHandler := handlers.NewReportHandler(reportService)

	reports := api.Group("/reports")
	reports.Use(middleware.AuthMiddleware())
	{
		// Sales Reports
		reports.GET("/sales", reportHandler.GetSalesReportByPeriod)               // ?tanggal_dari=&tanggal_sampai=&id_gudang=
		reports.GET("/sales/by-product", reportHandler.GetSalesReportByProduct)   // ?tanggal_dari=&tanggal_sampai=&id_gudang=
		reports.GET("/sales/by-customer", reportHandler.GetSalesReportByCustomer) // ?tanggal_dari=&tanggal_sampai=

		// Returns Report
		reports.GET("/returns", reportHandler.GetReturnReport) // ?tanggal_dari=&tanggal_sampai=

		// Stocks / Inventory Report
		reports.GET("/stocks", reportHandler.GetStockReport) // ?page=1&limit=10&search=&low_stock_only=true
	}
}
