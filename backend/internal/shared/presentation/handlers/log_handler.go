package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// LogHandler handles frontend log submissions
type LogHandler struct{}

// NewLogHandler creates a new log handler
func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// LogEntry represents a single log entry from the frontend
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
	User      struct {
		ID    string `json:"id,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"user,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	TraceID    string `json:"trace_id,omitempty"`
	SpanID     string `json:"span_id,omitempty"`
	URL        string `json:"url,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// LogsRequest represents a batch of logs from the frontend
type LogsRequest struct {
	Logs []LogEntry `json:"logs"`
}

// SubmitLogs handles batch log submissions from the frontend
// @Summary Submit frontend logs
// @Description Receives structured logs from the frontend for centralized logging
// @Tags logs
// @Accept json
// @Produce json
// @Param logs body LogsRequest true "Batch of log entries"
// @Success 200 {object} map[string]string "Logs received"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/logs [post]
func (h *LogHandler) SubmitLogs(c *fiber.Ctx) error {
	var req LogsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get request ID from context
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = c.Locals("request_id").(string)
	}

	// Get user context if available
	userID := c.Locals("user_id")
	userEmail := c.Locals("user_email")

	// Process each log entry
	for _, entry := range req.Logs {
		// Create logger with context
		logger := log.With().
			Str("source", "frontend").
			Str("frontend_request_id", entry.RequestID).
			Str("frontend_trace_id", entry.TraceID).
			Str("frontend_span_id", entry.SpanID).
			Str("frontend_url", entry.URL).
			Str("frontend_user_agent", entry.UserAgent).
			Str("frontend_environment", entry.Environment)

		// Add user context if available
		if entry.User.ID != "" {
			logger = logger.Str("frontend_user_id", entry.User.ID)
		}
		if entry.User.Email != "" {
			logger = logger.Str("frontend_user_email", entry.User.Email)
		}

		// Add backend request context
		if requestID != "" {
			logger = logger.Str("backend_request_id", requestID)
		}
		if userID != nil {
			logger = logger.Interface("backend_user_id", userID)
		}
		if userEmail != nil {
			logger = logger.Interface("backend_user_email", userEmail)
		}

		// Add custom context
		if entry.Context != nil {
			logger = logger.Interface("context", entry.Context)
		}

		// Log based on level
		switch entry.Level {
		case "debug":
			logger.Logger().Debug().Msg(entry.Message)
		case "info":
			logger.Logger().Info().Msg(entry.Message)
		case "warn":
			logger.Logger().Warn().Msg(entry.Message)
		case "error":
			logger.Logger().Error().Msg(entry.Message)
		default:
			logger.Logger().Info().Msg(entry.Message)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Logs received",
		"count":   len(req.Logs),
	})
}

