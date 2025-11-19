package errors

import "fmt"

// AppError represents an application-specific error with a code
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error types
var (
	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}
	
	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized access",
	}
	
	ErrForbidden = &AppError{
		Code:    "FORBIDDEN",
		Message: "Access forbidden",
	}
	
	ErrValidation = &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input data",
	}
	
	ErrConflict = &AppError{
		Code:    "CONFLICT",
		Message: "Resource conflict",
	}
	
	ErrInternal = &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
	}
	
	ErrBadRequest = &AppError{
		Code:    "BAD_REQUEST",
		Message: "Bad request",
	}
	
	ErrDatabaseError = &AppError{
		Code:    "DATABASE_ERROR",
		Message: "Database operation failed",
	}
)

// New creates a new AppError
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NotFound creates a not found error with custom message
func NotFound(message string) *AppError {
	return &AppError{
		Code:    ErrNotFound.Code,
		Message: message,
	}
}

// Validation creates a validation error with custom message
func Validation(message string) *AppError {
	return &AppError{
		Code:    ErrValidation.Code,
		Message: message,
	}
}

// Unauthorized creates an unauthorized error with custom message
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    ErrUnauthorized.Code,
		Message: message,
	}
}

// Internal creates an internal error with wrapped error
func Internal(err error) *AppError {
	return &AppError{
		Code:    ErrInternal.Code,
		Message: ErrInternal.Message,
		Err:     err,
	}
}

// Database creates a database error with wrapped error
func Database(err error) *AppError {
	return &AppError{
		Code:    ErrDatabaseError.Code,
		Message: ErrDatabaseError.Message,
		Err:     err,
	}
}
