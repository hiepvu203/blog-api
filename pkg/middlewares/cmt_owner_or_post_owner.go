package middlewares

import (
	"blog-api/internal/entities"
	"blog-api/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommentOwnerOrPostOwnerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("userID")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
			return
		}
		uid, ok := userID.(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
			return
		}
		commentIDParam := ctx.Param("comment_id")
		var commentID uint
		_, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse("Invalid comment_id"))
			return
		}

		var comment entities.Comment
		if err := db.Preload("Post").First(&comment, commentID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse("Comment not found"))
			return
		}

		if comment.UserID != uint(uid) && comment.Post.AuthorID != uint(uid) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("You do not have permission to delete this comment"))
			return
		}
		ctx.Next()
	}
}