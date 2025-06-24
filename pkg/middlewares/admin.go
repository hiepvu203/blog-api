package middlewares

import (
	"blog-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")

		if !exists || role != "admin"{
			utils.SendFail(ctx, 403, "403", utils.ErrAdminAccessRequired, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}