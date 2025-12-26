package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/category/presentation/handlers"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupCategoryRoutes configures category routes.
func SetupCategoryRoutes(router fiber.Router, categoryHandler *handlers.CategoryHandler, jwtService *services.JWTService) {
	categories := router.Group("/categories")

	// Apply authentication middleware to all category routes
	categories.Use(middleware.AuthMiddleware(jwtService))

	{
		categories.Post("/", categoryHandler.Create)
		categories.Get("/", categoryHandler.List)
		categories.Get("/:id", categoryHandler.Get)
		categories.Put("/:id", categoryHandler.Update)
		categories.Delete("/:id", categoryHandler.Delete)
	}
}
