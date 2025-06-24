package middlewares

import (
	"blog-api/pkg/utils"
	"github.com/gin-gonic/gin"
    "errors"
	"log"
)

func OwnerOrAdminMiddleware(db *gorm.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
        userID, ok := ctx.Get("userID")
        if !ok {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("authorization",utils.ErrUnauthorized))
            return
        }
        role, _ := ctx.Get("role")
        if role == "admin" {
            ctx.Next()
            return
        }

        uidFloat, ok := userID.(float64)
        if !ok {
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse("userID",utils.ErrInvalidUserIDType))
            return
        }
        uid := uint(uidFloat)

        postIDParam := ctx.Param("id")
        postID, err := strconv.ParseUint(postIDParam, 10, 64)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse("postID",utils.ErrInvalidPostID))
            return
        }

        var post entities.Post
        if err := db.First(&post, postID).Error; err != nil {
            ctx.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse("",utils.ErrPostNotFound))
            return
        }

        if post.AuthorID != uid {
            ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("",utils.ErrNoPermissionPost))
            return
        }
        ctx.Next()
    }
}

func CommentOwnerOrPostOwnerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := ctx.Get("userID")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("authorization",utils.ErrUnauthorized))
			return
		}
		uid, ok := userID.(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse("userID",utils.ErrInvalidUserIDType))
			return
		}
		commentIDParam := ctx.Param("comment_id")
		var commentID uint
		_, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse("commentID",utils.ErrInvalidCommentID))
			return
		}

		var comment entities.Comment
		if err := db.Preload("Post").First(&comment, commentID).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse("",utils.ErrCommentNotFound))
			return
		}

		if comment.UserID != uint(uid) && comment.Post.AuthorID != uint(uid) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("",utils.ErrNoPermissionComment))
			return
		}
		ctx.Next()
	}
}