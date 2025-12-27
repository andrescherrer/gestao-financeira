package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/pkg/errors"
)

// ErrorHandlerMiddleware creates a middleware that handles errors globally
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request ID from context (set by request ID middleware)
		requestID := GetRequestID(c)

		// Call next handler
		err := c.Next()

		// If no error, return
		if err == nil {
			return nil
		}

		// Handle AppError
		if appErr, ok := errors.AsAppError(err); ok {
			// Log error with request ID
			logFields := log.Error().
				Str("request_id", requestID).
				Str("error_type", string(appErr.Type)).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Int("status", appErr.Code)

			if appErr.Err != nil {
				logFields = logFields.Err(appErr.Err)
			}

			if appErr.Details != nil {
				logFields = logFields.Interface("details", appErr.Details)
			}

			logFields.Msg(appErr.Message)

			// Build response
			response := fiber.Map{
				"error":      appErr.Message,
				"error_type": string(appErr.Type),
				"code":       appErr.Code,
			}

			// Add request ID if available
			if requestID != "" {
				response["request_id"] = requestID
			}

			// Add details if available (only in development)
			if appErr.Details != nil && c.Locals("env") == "development" {
				response["details"] = appErr.Details
			}

			return c.Status(appErr.Code).JSON(response)
		}

		// Handle Fiber errors
		if fiberErr, ok := err.(*fiber.Error); ok {
			log.Warn().
				Str("request_id", requestID).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Int("status", fiberErr.Code).
				Msg(fiberErr.Message)

			response := fiber.Map{
				"error": fiberErr.Message,
				"code":  fiberErr.Code,
			}

			if requestID != "" {
				response["request_id"] = requestID
			}

			return c.Status(fiberErr.Code).JSON(response)
		}

		// Handle domain errors (check error message patterns)
		errMsg := err.Error()
		code := fiber.StatusInternalServerError
		errorType := "INTERNAL_ERROR"

		// Check for common error patterns
		if strings.Contains(errMsg, "not found") {
			code = fiber.StatusNotFound
			errorType = "NOT_FOUND"
		} else if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "duplicate") {
			code = fiber.StatusConflict
			errorType = "CONFLICT"
		} else if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "cannot be empty") {
			code = fiber.StatusBadRequest
			errorType = "VALIDATION_ERROR"
		} else if strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "forbidden") {
			code = fiber.StatusUnauthorized
			errorType = "UNAUTHORIZED"
		}

		// Log error
		log.Error().
			Str("request_id", requestID).
			Err(err).
			Str("error_type", errorType).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Unhandled error")

		// Build response
		response := fiber.Map{
			"error":      "An error occurred processing your request",
			"error_type": errorType,
			"code":       code,
		}

		if requestID != "" {
			response["request_id"] = requestID
		}

		// In development, include error message
		if c.Locals("env") == "development" {
			response["message"] = errMsg
		}

		return c.Status(code).JSON(response)
	}
}
