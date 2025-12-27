package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeDomain represents domain/business logic errors
	ErrorTypeDomain ErrorType = "DOMAIN_ERROR"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeConflict represents conflict errors (e.g., duplicate)
	ErrorTypeConflict ErrorType = "CONFLICT"
	// ErrorTypeUnauthorized represents unauthorized errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden represents forbidden errors
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
)

// AppError represents an application error with type and HTTP status code
type AppError struct {
	Type    ErrorType
	Message string
	Code    int
	Details map[string]interface{}
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (original: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusBadRequest,
		Details: details,
	}
}

// NewDomainError creates a new domain error
func NewDomainError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDomain,
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string, identifier string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Code:    http.StatusNotFound,
		Details: map[string]interface{}{
			"resource":   resource,
			"identifier": identifier,
		},
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    http.StatusConflict,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Code:    http.StatusForbidden,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

// WrapError wraps an existing error with an AppError
func WrapError(err error, errorType ErrorType, message string) *AppError {
	code := http.StatusInternalServerError
	switch errorType {
	case ErrorTypeValidation:
		code = http.StatusBadRequest
	case ErrorTypeDomain:
		code = http.StatusBadRequest
	case ErrorTypeNotFound:
		code = http.StatusNotFound
	case ErrorTypeConflict:
		code = http.StatusConflict
	case ErrorTypeUnauthorized:
		code = http.StatusUnauthorized
	case ErrorTypeForbidden:
		code = http.StatusForbidden
	}

	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError if possible
func AsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
