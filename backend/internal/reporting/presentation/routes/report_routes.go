package routes

import (
	"github.com/gofiber/fiber/v2"

	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/reporting/presentation/handlers"
	"gestao-financeira/backend/pkg/middleware"
)

// SetupReportRoutes configures report routes.
func SetupReportRoutes(router fiber.Router, reportHandler *handlers.ReportHandler, jwtService *services.JWTService) {
	reports := router.Group("/reports")

	// Apply authentication middleware to all report routes
	reports.Use(middleware.AuthMiddleware(jwtService))

	{
		reports.Get("/monthly", reportHandler.GetMonthlyReport)
		reports.Get("/annual", reportHandler.GetAnnualReport)
		reports.Get("/category", reportHandler.GetCategoryReport)
		reports.Get("/income-vs-expense", reportHandler.GetIncomeVsExpense)
	}
}
