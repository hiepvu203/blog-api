package utils

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetUintIDParam(ctx *gin.Context, paramName string, invalidMsg string) (uint, bool) {
	idParam := ctx.Param(paramName)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse("", invalidMsg))
		return 0, false
	}
	return uint(id), true
}

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

func GetPaginationParams(ctx *gin.Context) (page int, pageSize int, ok bool) {
    page = 1
    pageSize = 10
    ok = true

    if p := ctx.Query("page"); p != "" {
        if v, err := strconv.Atoi(p); err == nil && v > 0 {
            page = v
        } else {
            ctx.JSON(http.StatusBadRequest, ErrorResponse("Page param",ErrInvalidPageParam))
            return 0, 0, false
        }
    }
    if ps := ctx.Query("page_size"); ps != "" {
        if v, err := strconv.Atoi(ps); err == nil && v > 0 {
            pageSize = v
        } else {
            ctx.JSON(http.StatusBadRequest, ErrorResponse("Page size",ErrInvalidPageSizeParam))
            return 0, 0, false
        }
    }
    return page, pageSize, true
}