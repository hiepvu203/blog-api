package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(ctx *gin.Context) (uint, bool) {
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse("Authorization",ErrUnauthorized))
		return 0, false
	}
	uid, ok := userID.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse("userID",ErrInvalidUserIDType))
		return 0, false
	}
	return uint(uid), true
}