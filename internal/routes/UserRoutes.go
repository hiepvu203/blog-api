package routes

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	UserController := controllers.NewUserController(authService, userService)

	public := r.Group("/users")
	{
		public.POST("/register", UserController.Register)
		public.POST("/login", UserController.Login) 
		// public.POST("/forgot-password", UserController.ForgotPassword)
        // public.POST("/reset-password", UserController.ResetPassword)
	}

	authGroup := r.Group("/users").Use(middlewares.AuthMiddleware())
	{
		authGroup.GET("/me", UserController.GetMe)
		authGroup.PUT("/change-password", UserController.ChangePassword)
		authGroup.DELETE("/me", UserController.DeleteMe)
	}

	adminGroup := r.Group("/admin/users").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminGroup.GET("", UserController.ListUsers)
		adminGroup.GET("/:id", UserController.GetUserDetail)
		adminGroup.PUT("/:id/role", UserController.ChangeUserRole)
		adminGroup.DELETE("/:id", UserController.DeleteUser)
		adminGroup.PUT("/:id/ban-post", UserController.UpdateCanPost)
	}
}