package errors

import (
	"fmt"
	"strings"
)

// MapDomainError maps domain errors to AppError based on error message patterns.
// This is a transitional helper to migrate from string matching to typed errors.
// Eventually, all domain errors should return AppError directly.
func MapDomainError(err error) *AppError {
	if err == nil {
		return nil
	}

	// If already an AppError, return it
	if appErr, ok := AsAppError(err); ok {
		return appErr
	}

	errMsg := strings.ToLower(err.Error())

	// Validation errors
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "must be") ||
		strings.Contains(errMsg, "should be") {
		return NewValidationError(
			err.Error(),
			map[string]interface{}{
				"original_error": err.Error(),
			},
		)
	}

	// Not found errors
	if strings.Contains(errMsg, "not found") {
		// Try to extract resource and identifier from error message
		resource := "Resource"
		identifier := "unknown"

		// Common patterns: "{resource} not found: {identifier}"
		parts := strings.Split(errMsg, "not found")
		if len(parts) > 0 {
			resource = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				identifier = strings.TrimSpace(strings.TrimPrefix(parts[1], ":"))
			}
		}

		return NewNotFoundError(resource, identifier)
	}

	// Conflict errors
	if strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate") ||
		strings.Contains(errMsg, "conflict") {
		return NewConflictError(err.Error())
	}

	// Unauthorized errors
	if strings.Contains(errMsg, "unauthorized") ||
		strings.Contains(errMsg, "invalid credentials") ||
		strings.Contains(errMsg, "authentication failed") {
		return NewUnauthorizedError(err.Error())
	}

	// Forbidden errors
	if strings.Contains(errMsg, "forbidden") ||
		strings.Contains(errMsg, "permission denied") ||
		strings.Contains(errMsg, "access denied") {
		return NewForbiddenError(err.Error())
	}

	// Domain/business logic errors
	if strings.Contains(errMsg, "failed to") ||
		strings.Contains(errMsg, "unable to") ||
		strings.Contains(errMsg, "cannot") {
		return NewDomainError(err.Error(), err)
	}

	// Default to internal error for unknown errors
	return NewInternalError(
		"An unexpected error occurred",
		err,
	)
}

// MapToAppError is a convenience function that wraps MapDomainError
// and handles nil errors gracefully.
func MapToAppError(err error) error {
	if err == nil {
		return nil
	}
	return MapDomainError(err)
}

// WrapDomainError wraps a domain error with additional context.
func WrapDomainError(err error, errorType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	// If already an AppError, preserve it but update message if provided
	if appErr, ok := AsAppError(err); ok {
		if message != "" {
			appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		}
		return appErr
	}

	return WrapError(err, errorType, message)
}
