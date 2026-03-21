package routes

import (
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/handlers"
	"real-erp-mebel/be/internal/middleware"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupPemasokRoutes(api *gin.RouterGroup) {
	repo := repositories.NewPemasokRepository(database.DB)
	service := services.NewPemasokService(repo)
	handler := handlers.NewPemasokHandler(service)

	suppliers := api.Group("/suppliers")
	suppliers.Use(middleware.AuthMiddleware())
	{
		suppliers.POST("", handler.CreatePemasok)
		suppliers.GET("", handler.ListPemasok)
		suppliers.GET("/:id", handler.GetPemasok)
		suppliers.PUT("/:id", handler.UpdatePemasok)
		suppliers.DELETE("/:id", handler.DeletePemasok)
	}
}
