package handlers

import (
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/notification/infrastructure/websocket"
)

// WebSocketHandler handles WebSocket connections for notifications.
type WebSocketHandler struct {
	hub        *websocket.Hub
	jwtService *services.JWTService
}

// NewWebSocketHandler creates a new WebSocketHandler instance.
func NewWebSocketHandler(hub *websocket.Hub, jwtService *services.JWTService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtService: jwtService,
	}
}

// HandleWebSocket handles WebSocket upgrade requests.
// It authenticates the user using JWT token from query parameter or header.
func (h *WebSocketHandler) HandleWebSocket(c *fiber.Ctx) error {
	// Upgrade connection to WebSocket
	if !fiberws.IsWebSocketUpgrade(c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "WebSocket upgrade required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Get JWT token from query parameter or Authorization header
	token := c.Query("token")
	if token == "" {
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication token required",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Validate token and extract user ID
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid WebSocket authentication token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authentication token",
			"code":  fiber.StatusUnauthorized,
		})
	}

	userID := claims.UserID
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in token",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Upgrade to WebSocket and serve
	return fiberws.New(func(conn *fiberws.Conn) {
		websocket.ServeWebSocket(conn, h.hub, userID)
	})(c)
}
