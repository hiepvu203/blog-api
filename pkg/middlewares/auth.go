package middlewares

import (
	"blog-api/pkg/utils"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		// get token from header
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("Auth failed: missing token")
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse("token",utils.ErrNotToken))
			return
		}

		// Loại bỏ tiền tố "Bearer " 
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
		
		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Println("Auth failed: token expired")
				ctx.AbortWithStatusJSON(401, utils.ErrorResponse("token","Token expired"))
				return
			}
			log.Println("Auth failed: invalid token:", err)
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse("token",utils.ErrInvalidToken))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("userID", claims["user_id"])
			ctx.Set("role", claims["role"])
		}else {
			ctx.AbortWithStatusJSON(401, utils.ErrorResponse("token",utils.ErrInvalidTokenClaims))
			return
		}
		ctx.Next()
	}
}