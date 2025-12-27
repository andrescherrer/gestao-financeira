package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/reporting/application/dtos"
	"gestao-financeira/backend/internal/reporting/application/usecases"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// ReportHandler handles report-related HTTP requests.
type ReportHandler struct {
	monthlyReportUseCase   *usecases.MonthlyReportUseCase
	annualReportUseCase    *usecases.AnnualReportUseCase
	categoryReportUseCase  *usecases.CategoryReportUseCase
	incomeVsExpenseUseCase *usecases.IncomeVsExpenseUseCase
}

// NewReportHandler creates a new ReportHandler instance.
func NewReportHandler(
	monthlyReportUseCase *usecases.MonthlyReportUseCase,
	annualReportUseCase *usecases.AnnualReportUseCase,
	categoryReportUseCase *usecases.CategoryReportUseCase,
	incomeVsExpenseUseCase *usecases.IncomeVsExpenseUseCase,
) *ReportHandler {
	return &ReportHandler{
		monthlyReportUseCase:   monthlyReportUseCase,
		annualReportUseCase:    annualReportUseCase,
		categoryReportUseCase:  categoryReportUseCase,
		incomeVsExpenseUseCase: incomeVsExpenseUseCase,
	}
}

// GetMonthlyReport handles monthly report requests.
// @Summary Get monthly report
// @Description Generates a monthly financial report for the authenticated user
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param year query int true "Year (e.g., 2025)"
// @Param month query int true "Month (1-12)"
// @Param currency query string false "Currency filter (BRL, USD, EUR)"
// @Success 200 {object} dtos.MonthlyReportOutput "Monthly report data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /reports/monthly [get]
func (h *ReportHandler) GetMonthlyReport(c *fiber.Ctx) error {
	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	currency := c.Query("currency")

	if yearStr == "" || monthStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "year and month are required",
			"code":  fiber.StatusBadRequest,
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid year format",
			"code":  fiber.StatusBadRequest,
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid month format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.MonthlyReportInput{
		UserID:   userID,
		Year:     year,
		Month:    month,
		Currency: currency,
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		return err
	}

	// Execute use case
	output, err := h.monthlyReportUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})
}

// GetAnnualReport handles annual report requests.
// @Summary Get annual report
// @Description Generates an annual financial report for the authenticated user
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param year query int true "Year (e.g., 2025)"
// @Param currency query string false "Currency filter (BRL, USD, EUR)"
// @Success 200 {object} dtos.AnnualReportOutput "Annual report data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /reports/annual [get]
func (h *ReportHandler) GetAnnualReport(c *fiber.Ctx) error {
	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	yearStr := c.Query("year")
	currency := c.Query("currency")

	if yearStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "year is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid year format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.AnnualReportInput{
		UserID:   userID,
		Year:     year,
		Currency: currency,
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		return err
	}

	// Execute use case
	output, err := h.annualReportUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})
}

// GetCategoryReport handles category report requests.
// @Summary Get category report
// @Description Generates a category-based financial report for the authenticated user
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_id query string false "Category ID filter (UUID)"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param currency query string false "Currency filter (BRL, USD, EUR)"
// @Success 200 {object} dtos.CategoryReportOutput "Category report data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /reports/category [get]
func (h *ReportHandler) GetCategoryReport(c *fiber.Ctx) error {
	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	categoryID := c.Query("category_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	currency := c.Query("currency")

	// Build input
	input := dtos.CategoryReportInput{
		UserID:     userID,
		CategoryID: categoryID,
		Currency:   currency,
	}

	// Parse dates if provided
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid start_date format (expected YYYY-MM-DD)",
				"code":  fiber.StatusBadRequest,
			})
		}
		input.StartDate = &startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid end_date format (expected YYYY-MM-DD)",
				"code":  fiber.StatusBadRequest,
			})
		}
		input.EndDate = &endDate
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		return err
	}

	// Execute use case
	output, err := h.categoryReportUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})
}

// GetIncomeVsExpense handles income vs expense report requests.
// @Summary Get income vs expense report
// @Description Generates an income vs expense comparison report for the authenticated user
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param start_date query string false "Start date filter (YYYY-MM-DD)"
// @Param end_date query string false "End date filter (YYYY-MM-DD)"
// @Param currency query string false "Currency filter (BRL, USD, EUR)"
// @Param group_by query string false "Group by period (day, week, month, year)"
// @Success 200 {object} dtos.IncomeVsExpenseOutput "Income vs expense report data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /reports/income-vs-expense [get]
func (h *ReportHandler) GetIncomeVsExpense(c *fiber.Ctx) error {
	// Get user ID from context
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	currency := c.Query("currency")
	groupBy := c.Query("group_by")

	// Build input
	input := dtos.IncomeVsExpenseInput{
		UserID:   userID,
		Currency: currency,
		GroupBy:  groupBy,
	}

	// Parse dates if provided
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid start_date format (expected YYYY-MM-DD)",
				"code":  fiber.StatusBadRequest,
			})
		}
		input.StartDate = &startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid end_date format (expected YYYY-MM-DD)",
				"code":  fiber.StatusBadRequest,
			})
		}
		input.EndDate = &endDate
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		return err
	}

	// Execute use case
	output, err := h.incomeVsExpenseUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})
}

// handleUseCaseError handles errors from use cases.
func (h *ReportHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	requestID := middleware.GetRequestID(c)

	// Log error
	log.Error().
		Err(err).
		Str("request_id", requestID).
		Str("path", c.Path()).
		Str("method", c.Method()).
		Msg("Use case error")

	// Return error response
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":      "Failed to generate report",
		"error_type": "INTERNAL_ERROR",
		"code":       fiber.StatusInternalServerError,
		"request_id": requestID,
	})
}
