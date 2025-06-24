package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Status 	string 	    `json:"status"`
	Code 	string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Meta struct {
	Page   	 int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func SendSuccess(ctx *gin.Context, httpCode int, apiCode , message string, data interface{}) {
	ctx.JSON(httpCode, APIResponse{
		Status:  "success",
		Code:    apiCode,
		Message: message,
		Data:    data,
	})
}

func SendFail(ctx *gin.Context, httpCode int, apiCode, message string, data interface{}) {
	ctx.JSON(httpCode, APIResponse{
		Status:  "error",
		Code:    apiCode,
		Message: message,
		Data:    data,
	})
}

func BindAndValidate(ctx *gin.Context, obj interface{}) map[string]string {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return ParseValidationErrors(err)
	}
	return nil
}

func ParseValidationErrors(err error) map[string]string {
    var ve validator.ValidationErrors
    errs := make(map[string]string)

    if errors.As(err, &ve) {
        for _, fe := range ve {
            tag := fe.ActualTag()
            param := fe.Param()
            field := fe.Field()
            var msg string
            switch tag {
            case "required":
                msg = fmt.Sprintf("%s is required", field)
            case "min":
                msg = fmt.Sprintf("%s must have at least %s characters", field, param)
            case "max":
                msg = fmt.Sprintf("%s must not exceed %s characters", field, param)
            case "slug":
				msg = "slug can only contain lowercase letters, numbers, and hyphens. It cannot start or end with a hyphen, anh cannot have two consecutive hyphens."			
			case "username":
				msg = "username can only contain letters, numbers, underscores, and hyphens; no spaces or special characters."
			case "strongpwd":
				msg = "Password must have at least 8 characters, including uppercase letters, lowercase letters, numbers, and special characters."			
			default:
                msg = fmt.Sprintf("Invalid %s", field)
            }
            errs[field] = msg
        }
    }
    return errs
}

