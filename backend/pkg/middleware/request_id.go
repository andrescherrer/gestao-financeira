package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDContextKey is the context key for request ID
	RequestIDContextKey = "request_id"
)

// RequestIDMiddleware creates a middleware that generates and adds a request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get(RequestIDHeader)

		// If not present, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context
		c.Locals(RequestIDContextKey, requestID)

		// Add to response header
		c.Set(RequestIDHeader, requestID)

		// Continue to next handler
		return c.Next()
	}
}

// GetRequestID extracts the request ID from the request context
func GetRequestID(c *fiber.Ctx) string {
	if requestID, ok := c.Locals(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
