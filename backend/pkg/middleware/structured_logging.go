package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// StructuredLoggingMiddleware creates a middleware that enhances logging with request context
// This middleware should be used after RequestIDMiddleware and AuthMiddleware
func StructuredLoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request context
		requestID := GetRequestID(c)
		userID := GetUserID(c)
		userEmail := GetUserEmail(c)

		// Create logger with request context
		requestLogger := log.With().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent"))

		// Add user context if available
		if userID != "" {
			requestLogger = requestLogger.Str("user_id", userID)
		}
		if userEmail != "" {
			requestLogger = requestLogger.Str("user_email", userEmail)
		}

		// Store enhanced logger in context
		c.Locals("logger", requestLogger.Logger())

		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request completion
		logEvent := requestLogger.Logger().Info().
			Int("status", c.Response().StatusCode()).
			Dur("duration_ms", duration).
			Int("bytes_sent", len(c.Response().Body()))

		// Add error information if present
		if err != nil {
			logEvent = logEvent.Err(err)
		}

		logEvent.Msg("Request completed")

		return err
	}
}

// GetRequestLogger extracts the enhanced logger from request context
// Returns the global logger if not found
func GetRequestLogger(c *fiber.Ctx) zerolog.Logger {
	if logger, ok := c.Locals("logger").(zerolog.Logger); ok {
		return logger
	}
	// Fallback to global logger with request ID
	requestID := GetRequestID(c)
	if requestID != "" {
		return log.With().Str("request_id", requestID).Logger()
	}
	return log.Logger
}

