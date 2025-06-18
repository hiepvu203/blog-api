package middlewares

import (
	"blog-api/internal/entities"
	"blog-api/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OwnerOrAdminMiddleware(db *gorm.DB) gin.HandlerFunc{
	return func(ctx *gin.Context) {
        userID, ok := ctx.Get("userID")
        if !ok {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
            return
        }
        role, _ := ctx.Get("role")
        if role == "admin" {
            ctx.Next()
            return
        }

        uidFloat, ok := userID.(float64)
        if !ok {
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
            return
        }
        uid := uint(uidFloat)

        postIDParam := ctx.Param("id")
        postID, err := strconv.ParseUint(postIDParam, 10, 64)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post id"))
            return
        }

        var post entities.Post
        if err := db.First(&post, postID).Error; err != nil {
            ctx.AbortWithStatusJSON(http.StatusNotFound, utils.ErrorResponse("Post not found"))
            return
        }

        if post.AuthorID != uid {
            ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse("You do not have permission to edit or delete this post."))
            return
        }
        ctx.Next()
    }
}