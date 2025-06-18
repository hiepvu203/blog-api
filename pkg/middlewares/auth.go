package middlewares

import (
	"blog-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		// get token from header
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse(utils.ErrNotToken))
			return
		}

		// Loại bỏ tiền tố "Bearer " 
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
		
		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse(utils.ErrInvalidToken))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("userID", claims["user_id"])
			ctx.Set("role", claims["role"])
		}else {
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse(utils.ErrInvalidTokenClaims))
			return
		}
		ctx.Next()
	}
}