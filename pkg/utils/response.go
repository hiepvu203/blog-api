package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"` // omitempty: bỏ qua nếu rỗng
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

func ParseValidationError(err error) string {
	var sb strings.Builder
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				sb.WriteString(field + " is required. ")
			case "email":
				sb.WriteString(field + " must be a valid email. ")
			case "min":
				sb.WriteString(field + " must be at least " + e.Param() + " characters. ")
			case "max":
				sb.WriteString(field + " must be at most " + e.Param() + " characters. ")
			case "oneof":
				sb.WriteString(field + " must be one of: " + e.Param() + ". ")
			default:
				sb.WriteString(field + " is invalid. ")
			}
		}
		return strings.TrimSpace(sb.String())
	}
	return err.Error()
}