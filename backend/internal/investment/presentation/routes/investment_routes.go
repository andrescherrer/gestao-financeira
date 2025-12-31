package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/investment/presentation/handlers"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupInvestmentRoutes configures investment routes.
func SetupInvestmentRoutes(router fiber.Router, investmentHandler *handlers.InvestmentHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	investments := router.Group("/investments")

	// Apply authentication middleware to all investment routes
	investments.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		investments.Post("/", investmentHandler.Create)
		investments.Get("/", investmentHandler.List)
		investments.Get("/:id", investmentHandler.Get)
		investments.Put("/:id", investmentHandler.Update)
		investments.Delete("/:id", investmentHandler.Delete)
	}
}
