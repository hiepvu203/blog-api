package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct  {
	Field 	string `json:"Field"`
	Message string `json:"Message"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(field , message string) Response {
	return Response{
		Success: false,
		Error:   FieldError{
			Field: field,
			Message: message,
		},
	}
}

func BindAndValidate(ctx *gin.Context, obj interface{}) []FieldError {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return ParseValidationErrors(err)
	}
	return nil
}

func ParseValidationErrors(err error) []FieldError {
    var ve validator.ValidationErrors
    var errs []FieldError

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
            errs = append(errs, FieldError{
                Field:   field,
                Message: msg,
            })
        }
    }
    return errs
}

