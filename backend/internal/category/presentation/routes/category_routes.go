package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/category/presentation/handlers"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupCategoryRoutes configures category routes.
func SetupCategoryRoutes(router fiber.Router, categoryHandler *handlers.CategoryHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	categories := router.Group("/categories")

	// Apply authentication middleware to all category routes
	categories.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		categories.Post("/", categoryHandler.Create)
		categories.Get("/", categoryHandler.List)
		categories.Get("/:id", categoryHandler.Get)
		categories.Put("/:id", categoryHandler.Update)
		categories.Delete("/:id", categoryHandler.Delete)
		categories.Post("/:id/restore", categoryHandler.Restore)
		categories.Delete("/:id/permanent", categoryHandler.PermanentDelete)
	}
}
