package utils

import (
	"github.com/gin-gonic/gin"
)

func BindAndValidate(ctx *gin.Context, obj interface{}) []FieldError {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return ParseValidationErrors(err)
	}
	return nil
}