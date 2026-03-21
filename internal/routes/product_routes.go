package routes

import (
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes mengatur routes untuk produk
func SetupProductRoutes(api *gin.RouterGroup) {
	// Initialize dependencies
	db := database.DB
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Product routes (protected)
	products := api.Group("/products")
	products.Use(middleware.AuthMiddleware())
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
}
