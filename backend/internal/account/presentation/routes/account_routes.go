package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/account/presentation/handlers"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupAccountRoutes configures account routes.
func SetupAccountRoutes(router fiber.Router, accountHandler *handlers.AccountHandler, jwtService *services.JWTService) {
	accounts := router.Group("/accounts")

	// Apply authentication middleware to all account routes
	accounts.Use(middleware.AuthMiddleware(jwtService))

	{
		accounts.Post("/", accountHandler.Create)
		accounts.Get("/", accountHandler.List)
		accounts.Get("/:id", accountHandler.Get)
	}
}
