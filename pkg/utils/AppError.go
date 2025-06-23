package utils

// AppError represents a custom error with a specific field and message.
type AppError struct {
    Field   string
    Message string
}

// Error returns the error message to satisfy the error interface.
func (e *AppError) Error() string {
    return e.Message
}

// NewAppError creates a new AppError.
func NewAppError(field, message string) *AppError {
    return &AppError{Field: field, Message: message}
} 