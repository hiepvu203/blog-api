package middlewares

import (
	"blog-api/internal/entities"
	"blog-api/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func OwnerOrAdminMiddleware(db *gorm.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
        userID, ok := ctx.Get("userID")
        if !ok {
            utils.SendFail(ctx, http.StatusUnauthorized, "401", utils.ErrUnauthorized, nil)
            ctx.Abort()
            return
        }
        role, _ := ctx.Get("role")
        if role == "admin" {
            ctx.Next()
            return
        }

        uidFloat, ok := userID.(float64)
        if !ok {
            utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrInvalidUserIDType, nil)
            ctx.Abort()
            return
        }
        uid := uint(uidFloat)

        postIDParam := ctx.Param("id")
        postID, err := strconv.ParseUint(postIDParam, 10, 64)
        if err != nil {
            utils.SendFail(ctx, http.StatusBadRequest, "400", utils.ErrInvalidPostID, nil)
            ctx.Abort()
            return
        }

        var post entities.Post
        if err := db.First(&post, postID).Error; err != nil {
            utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrPostNotFound, nil)
            ctx.Abort()
            return
        }

        if post.AuthorID != uid {
            utils.SendFail(ctx, http.StatusForbidden, "403", utils.ErrNoPermissionPost, nil)
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}

func CommentOwnerOrPostOwnerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("userID")
		if !ok {
			utils.SendFail(ctx, http.StatusUnauthorized, "401", utils.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		uid, ok := userID.(float64)
		if !ok {
			utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrInvalidUserIDType, nil)
			ctx.Abort()
			return
		}
		commentIDParam := ctx.Param("comment_id")
		var commentID uint
		_, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
		if err != nil {
			utils.SendFail(ctx, http.StatusBadRequest, "400", utils.ErrInvalidCommentID, nil)
			ctx.Abort()
			return
		}

		var comment entities.Comment
		if err := db.Preload("Post").First(&comment, commentID).Error; err != nil {
			utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrCommentNotFound, nil)
			ctx.Abort()
			return
		}

		if comment.UserID != uint(uid) && comment.Post.AuthorID != uint(uid) {
			utils.SendFail(ctx, http.StatusForbidden, "403", utils.ErrNoPermissionComment, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}