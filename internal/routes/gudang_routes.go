package routes

import (
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"real-erp-mebel/be/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupGudangRoutes(r *gin.RouterGroup) {
	gudangRepo := repositories.NewGudangRepository(database.DB)
	gudangService := services.NewGudangService(gudangRepo)
	gudangHandler := handlers.NewGudangHandler(gudangService)

	gudang := r.Group("/warehouses")
	gudang.Use(middleware.AuthMiddleware())
	{
		gudang.POST("", gudangHandler.CreateGudang)
		gudang.GET("", gudangHandler.ListGudangs)
		gudang.GET("/:id", gudangHandler.GetGudang)
		gudang.PUT("/:id", gudangHandler.UpdateGudang)
		gudang.DELETE("/:id", gudangHandler.DeleteGudang)
	}
}
