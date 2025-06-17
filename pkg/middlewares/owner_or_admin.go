package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func OwnerOrAdminMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(uint)
		requestID := ctx.Param("id")

		if fmt.Sprintf("%d", userID) != requestID && ctx.MustGet("role") != "admin" {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "You can only edit your own data"})
			return
		}
		ctx.Next()
	}
}