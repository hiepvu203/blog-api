package middlewares

import (
	"blog-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")

		if !exists || role != "admin"{
			ctx.AbortWithStatusJSON(403, utils.ErrorResponse("Admin access required"))
			return
		}
		ctx.Next()
	}
}