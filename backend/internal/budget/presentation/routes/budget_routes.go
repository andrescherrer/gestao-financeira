package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/budget/presentation/handlers"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupBudgetRoutes configures budget routes.
func SetupBudgetRoutes(router fiber.Router, budgetHandler *handlers.BudgetHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	budgets := router.Group("/budgets")

	// Apply authentication middleware to all budget routes
	budgets.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		budgets.Post("/", budgetHandler.Create)
		budgets.Get("/", budgetHandler.List)
		budgets.Get("/:id", budgetHandler.Get)
		budgets.Get("/:id/progress", budgetHandler.GetProgress)
		budgets.Put("/:id", budgetHandler.Update)
		budgets.Delete("/:id", budgetHandler.Delete)
	}
}
