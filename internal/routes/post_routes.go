package routes

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPostRoutes(r *gin.Engine, db *gorm.DB) {
    repo := repositories.NewPostRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	service := services.NewPostService(repo, categoryRepo)
    controller := controllers.NewPostController(service)

    // Route cho user (cần đăng nhập)
    userGroup := r.Group("/posts").Use(middlewares.AuthMiddleware())
    {
        userGroup.POST("", controller.CreatePost)
        userGroup.PUT("/:id", middlewares.OwnerOrAdminMiddleware(db), controller.UpdatePost)
        userGroup.DELETE("/:id", middlewares.OwnerOrAdminMiddleware(db), controller.DeletePost) 
    }

    // Route cho admin (quản lý mọi post)
    adminGroup := r.Group("/admin/posts").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
    {
        adminGroup.DELETE("/:id", controller.DeletePost) 
    }

	// Public routes (ai cũng xem được)
    publicGroup := r.Group("/posts")
    {
        publicGroup.GET("", controller.GetAllPosts)
        publicGroup.GET("/:post_id", controller.GetPostDetail)
    }
}