package routes

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	UserController := controllers.NewUserController(authService)

	r.POST("/register", UserController.Register)
}