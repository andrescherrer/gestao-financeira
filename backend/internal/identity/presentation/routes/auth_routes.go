package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/presentation/handlers"
)

// SetupAuthRoutes configures authentication routes.
func SetupAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
	}
}
