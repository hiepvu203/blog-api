package routes

import (
	"blog-api/internal/controllers"
	"blog-api/internal/repositories"
	"blog-api/internal/services"
	"blog-api/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCommentRoutes(r *gin.Engine, db *gorm.DB) {
	repo := repositories.NewCommentRepository(db)
	service := services.NewCommentService(repo)
	controller := controllers.NewCommentController(service)

	r.POST("/posts/:post_id/comments", middlewares.AuthMiddleware(), controller.CreateComment)
    r.PUT("/comments/:comment_id", middlewares.AuthMiddleware(), controller.UpdateComment)
    r.DELETE("/comments/:comment_id", middlewares.AuthMiddleware(), middlewares.CommentOwnerOrPostOwnerMiddleware(db), controller.DeleteComment)
    r.GET("/posts/:post_id/comments", controller.GetCommentsByPost)
}