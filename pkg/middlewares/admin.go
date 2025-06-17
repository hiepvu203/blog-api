package middlewares

import "github.com/gin-gonic/gin"

func AdminMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")

		if !exists || role != "admin"{
			ctx.AbortWithStatusJSON(403, gin.H{"error" : "Admin access required"})
			return
		}
		ctx.Next()
	}
}