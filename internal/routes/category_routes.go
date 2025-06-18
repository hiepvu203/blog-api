package routes

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCategoryRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewCategoryRepository(db)
	service := services.NewCategoryService(repo)
	controller := controllers.NewCategoryController(service)

	adminGroup := r.Group("admin/categories").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminGroup.GET("", controller.AdminListCategories)
		adminGroup.POST("", controller.CreateCategory)
		adminGroup.PUT("/:id", controller.UpdateCategory)
		adminGroup.DELETE("/:id", controller.DeleteCategory)
	}

	publicGroup := r.Group("/categories")
    {
        publicGroup.GET("", controller.ListCategories)
    }
}