package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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