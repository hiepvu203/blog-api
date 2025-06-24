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
			utils.SendFail(ctx, 401, "401", utils.ErrNotToken, nil)
			ctx.Abort()
			return
		}

		// remove the 'Bearer' prefix
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
		
		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Println("Auth failed: token expired")
				utils.SendFail(ctx, 401, "401", "token expired", nil)
				ctx.Abort()
				return
			}
			log.Println("Auth failed: invalid token:", err)
			utils.SendFail(ctx, 401, "401", utils.ErrInvalidToken, nil)
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("userID", claims["user_id"])
			ctx.Set("role", claims["role"])
		} else {
			utils.SendFail(ctx, 401, "401", utils.ErrInvalidTokenClaims, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}