package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/goal/presentation/handlers"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupGoalRoutes configures goal routes.
func SetupGoalRoutes(router fiber.Router, goalHandler *handlers.GoalHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	goals := router.Group("/goals")

	// Apply authentication middleware to all goal routes
	goals.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		goals.Post("/", goalHandler.Create)
		goals.Get("/", goalHandler.List)
		goals.Get("/:id", goalHandler.Get)
		goals.Post("/:id/contribute", goalHandler.AddContribution)
		goals.Put("/:id/progress", goalHandler.UpdateProgress)
		goals.Post("/:id/cancel", goalHandler.Cancel)
		goals.Delete("/:id", goalHandler.Delete)
	}
}
