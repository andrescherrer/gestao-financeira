package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/presentation/handlers"
)

// SetupAuthRoutes configures authentication routes.
func SetupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	auth := app.Group("/api/v1/auth")
	{
		auth.Post("/register", authHandler.Register)
		// Login route will be added in ID-010
	}
}
