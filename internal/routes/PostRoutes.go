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
	userRepo := repositories.NewUserRepository(db)
	service := services.NewPostService(repo, categoryRepo, userRepo)
    controller := controllers.NewPostController(service)

    userGroup := r.Group("/posts").Use(middlewares.AuthMiddleware())
    {
        userGroup.POST("", controller.CreatePost)
        userGroup.PUT("/:id", middlewares.OwnerOrAdminMiddleware(db), controller.UpdatePost)
        userGroup.DELETE("/:id", middlewares.OwnerOrAdminMiddleware(db), controller.DeletePost) 
    }

    adminGroup := r.Group("/admin/posts").Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
    {
        adminGroup.DELETE("/:id", controller.DeletePost) 
    }

    publicGroup := r.Group("/posts")
    {
        publicGroup.GET("", controller.GetAllPosts)
        publicGroup.GET("/:post_id", controller.GetPostDetail)
    }
}