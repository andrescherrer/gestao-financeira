package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/notification/presentation/handlers"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupNotificationRoutes configures notification routes.
func SetupNotificationRoutes(router fiber.Router, notificationHandler *handlers.NotificationHandler, websocketHandler *handlers.WebSocketHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	notifications := router.Group("/notifications")

	// Apply authentication middleware to all notification routes
	notifications.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		notifications.Post("/", notificationHandler.Create)
		notifications.Get("/", notificationHandler.List)
		notifications.Get("/:id", notificationHandler.Get)
		notifications.Post("/:id/read", notificationHandler.MarkRead)
		notifications.Post("/:id/unread", notificationHandler.MarkUnread)
		notifications.Post("/:id/archive", notificationHandler.Archive)
		notifications.Delete("/:id", notificationHandler.Delete)
	}

	// WebSocket route (no auth middleware needed, authentication is handled in the handler)
	router.Get("/ws/notifications", websocketHandler.HandleWebSocket)
}
