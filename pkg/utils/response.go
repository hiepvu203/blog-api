package utils

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"` // omitempty: bỏ qua nếu rỗng
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// successResponse
func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(message string) Response{
	return Response{
		Success: false,
		Error: message,
	}
}