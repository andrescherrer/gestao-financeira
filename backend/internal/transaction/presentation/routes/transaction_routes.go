package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/transaction/presentation/handlers"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupTransactionRoutes configures transaction routes.
func SetupTransactionRoutes(app *fiber.App, transactionHandler *handlers.TransactionHandler, jwtService *services.JWTService) {
	transactions := app.Group("/api/v1/transactions")

	// Apply authentication middleware to all transaction routes
	transactions.Use(middleware.AuthMiddleware(jwtService))

	{
		transactions.Post("/", transactionHandler.Create)
		transactions.Get("/", transactionHandler.List)
		transactions.Get("/:id", transactionHandler.Get)
		transactions.Put("/:id", transactionHandler.Update)
		transactions.Delete("/:id", transactionHandler.Delete)
	}
}
