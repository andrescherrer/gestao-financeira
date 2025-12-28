package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/transaction/presentation/handlers"
	"gestao-financeira/backend/pkg/cache"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupTransactionRoutes configures transaction routes.
func SetupTransactionRoutes(router fiber.Router, transactionHandler *handlers.TransactionHandler, jwtService *services.JWTService, userRepository repositories.UserRepository, cacheService *cache.CacheService) {
	transactions := router.Group("/transactions")

	// Apply authentication middleware to all transaction routes
	transactions.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   cacheService,
	}))

	{
		transactions.Post("/", transactionHandler.Create)
		transactions.Get("/", transactionHandler.List)
		transactions.Get("/:id", transactionHandler.Get)
		transactions.Put("/:id", transactionHandler.Update)
		transactions.Delete("/:id", transactionHandler.Delete)
	}
}
