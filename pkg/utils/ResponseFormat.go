package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct  {
	Field 	string `json:"Field"`
	Tag 	string `json:"Tag"`
	Param 	string `json:"Param"`
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
                msg = fmt.Sprintf("%s là bắt buộc", field)
            case "min":
                msg = fmt.Sprintf("%s phải có ít nhất %s ký tự", field, param)
            case "max":
                msg = fmt.Sprintf("%s không được vượt quá %s ký tự", field, param)
            case "slug":
				msg = "Slug chỉ được chứa chữ thường, số, dấu gạch ngang. Không bắt đầu/kết thúc bằng dấu gạch ngang. Không có hai dấu gạch ngang liên tiếp"			
			case "username":
				msg = "Username chỉ được chứa chữ cái, số, dấu gạch dưới hoặc gạch ngang, không khoảng trắng, không ký tự đặc biệt"
			case "strongpwd":
				msg = "Password phải có ít nhất 8 ký tự, gồm chữ hoa, chữ thường, số và ký tự đặc biệt"			
			default:
                msg = fmt.Sprintf("%s không hợp lệ", field)
            }
            errs = append(errs, FieldError{
                Field:   field,
                Tag:     tag,
                Param:   param,
                Message: msg,
            })
        }
    }
    return errs
}