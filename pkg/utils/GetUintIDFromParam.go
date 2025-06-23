package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUintIDParam(ctx *gin.Context, paramName string, invalidMsg string) (uint, bool){
	idParam := ctx.Param(paramName)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse("",invalidMsg))
		return 0, false
	}
	return uint(id), true
}