package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/budget/presentation/handlers"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupBudgetRoutes configures budget routes.
func SetupBudgetRoutes(router fiber.Router, budgetHandler *handlers.BudgetHandler, jwtService *services.JWTService) {
	budgets := router.Group("/budgets")

	// Apply authentication middleware to all budget routes
	budgets.Use(middleware.AuthMiddleware(jwtService))

	{
		budgets.Post("/", budgetHandler.Create)
		budgets.Get("/", budgetHandler.List)
		budgets.Get("/:id", budgetHandler.Get)
		budgets.Get("/:id/progress", budgetHandler.GetProgress)
		budgets.Put("/:id", budgetHandler.Update)
		budgets.Delete("/:id", budgetHandler.Delete)
	}
}
