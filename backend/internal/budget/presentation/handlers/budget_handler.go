package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/application/usecases"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// BudgetHandler handles budget-related HTTP requests.
type BudgetHandler struct {
	createBudgetUseCase      *usecases.CreateBudgetUseCase
	listBudgetsUseCase       *usecases.ListBudgetsUseCase
	getBudgetUseCase         *usecases.GetBudgetUseCase
	updateBudgetUseCase      *usecases.UpdateBudgetUseCase
	deleteBudgetUseCase      *usecases.DeleteBudgetUseCase
	getBudgetProgressUseCase *usecases.GetBudgetProgressUseCase
}

// NewBudgetHandler creates a new BudgetHandler instance.
func NewBudgetHandler(
	createBudgetUseCase *usecases.CreateBudgetUseCase,
	listBudgetsUseCase *usecases.ListBudgetsUseCase,
	getBudgetUseCase *usecases.GetBudgetUseCase,
	updateBudgetUseCase *usecases.UpdateBudgetUseCase,
	deleteBudgetUseCase *usecases.DeleteBudgetUseCase,
	getBudgetProgressUseCase *usecases.GetBudgetProgressUseCase,
) *BudgetHandler {
	return &BudgetHandler{
		createBudgetUseCase:      createBudgetUseCase,
		listBudgetsUseCase:       listBudgetsUseCase,
		getBudgetUseCase:         getBudgetUseCase,
		updateBudgetUseCase:      updateBudgetUseCase,
		deleteBudgetUseCase:      deleteBudgetUseCase,
		getBudgetProgressUseCase: getBudgetProgressUseCase,
	}
}

// Create handles budget creation requests.
// @Summary Create a new budget
// @Description Creates a new budget for the authenticated user. Supports MONTHLY and YEARLY period types.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateBudgetInput true "Budget creation data (category_id, amount, currency, period_type, year, month, context)"
// @Success 201 {object} map[string]interface{} "Budget created successfully"
// @Success 201 {object} dtos.CreateBudgetOutput "Budget data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 409 {object} map[string]interface{} "Conflict - budget already exists for this category and period"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets [post]
func (h *BudgetHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse request body
	var input dtos.CreateBudgetInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set user ID from context (override any user_id in request body for security)
	input.UserID = userID

	// Validate input
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.createBudgetUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Budget created successfully",
		"data":    output,
	})
}

// List handles budget listing requests.
// @Summary List budgets
// @Description Lists all budgets for the authenticated user. Optionally filter by category_id, period_type, year, month, context, or is_active.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_id query string false "Filter by category ID (UUID)"
// @Param period_type query string false "Filter by period type (MONTHLY or YEARLY)"
// @Param year query int false "Filter by year"
// @Param month query int false "Filter by month (1-12)"
// @Param context query string false "Filter by context (PERSONAL or BUSINESS)"
// @Param is_active query bool false "Filter by active status (true or false)"
// @Success 200 {object} map[string]interface{} "Budgets retrieved successfully"
// @Success 200 {object} dtos.ListBudgetsOutput "List of budgets with count"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID or filters"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets [get]
func (h *BudgetHandler) List(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	input := dtos.ListBudgetsInput{
		UserID: userID,
	}

	// Get optional filters from query parameters
	if categoryID := c.Query("category_id"); categoryID != "" {
		input.CategoryID = categoryID
	}
	if periodType := c.Query("period_type"); periodType != "" {
		input.PeriodType = periodType
	}
	if yearStr := c.Query("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil && year > 0 {
			input.Year = &year
		}
	}
	if monthStr := c.Query("month"); monthStr != "" {
		if month, err := strconv.Atoi(monthStr); err == nil && month > 0 {
			input.Month = &month
		}
	}
	if context := c.Query("context"); context != "" {
		input.Context = context
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if strings.ToLower(isActiveStr) == "true" {
			active := true
			input.IsActive = &active
		} else if strings.ToLower(isActiveStr) == "false" {
			active := false
			input.IsActive = &active
		}
	}

	// Execute use case
	output, err := h.listBudgetsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Budgets retrieved successfully",
		"data":    output,
	})
}

// Get handles budget retrieval requests.
// @Summary Get budget by ID
// @Description Retrieves a specific budget by its ID. Only returns budgets that belong to the authenticated user.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Budget ID (UUID)"
// @Success 200 {object} map[string]interface{} "Budget retrieved successfully"
// @Success 200 {object} dtos.GetBudgetOutput "Budget data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid budget ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - budget does not belong to user"
// @Failure 404 {object} map[string]interface{} "Not found - budget does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id} [get]
func (h *BudgetHandler) Get(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get budget ID from path parameter
	budgetID := c.Params("id")
	if budgetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Budget ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.GetBudgetInput{
		BudgetID: budgetID,
		UserID:   userID,
	}

	// Execute use case
	output, err := h.getBudgetUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Budget retrieved successfully",
		"data":    output,
	})
}

// Update handles budget update requests.
// @Summary Update budget
// @Description Updates a budget. Only provided fields will be updated.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Budget ID (UUID)"
// @Param request body dtos.UpdateBudgetInput true "Budget update data (amount, period_type, year, month, is_active - at least one required)"
// @Success 200 {object} map[string]interface{} "Budget updated successfully"
// @Success 200 {object} dtos.UpdateBudgetOutput "Updated budget data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - budget does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id} [put]
func (h *BudgetHandler) Update(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get budget ID from path parameter
	budgetID := c.Params("id")
	if budgetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Budget ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Parse request body
	var input dtos.UpdateBudgetInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set budget ID and user ID from path parameter and context
	input.BudgetID = budgetID
	input.UserID = userID

	// Validate input
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.updateBudgetUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Budget updated successfully",
		"data":    output,
	})
}

// Delete handles budget deletion requests.
// @Summary Delete budget
// @Description Deletes a budget.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Budget ID (UUID)"
// @Success 200 {object} map[string]interface{} "Budget deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid budget ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - budget does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id} [delete]
func (h *BudgetHandler) Delete(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get budget ID from path parameter
	budgetID := c.Params("id")
	if budgetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Budget ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.DeleteBudgetInput{
		BudgetID: budgetID,
		UserID:   userID,
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	err := h.deleteBudgetUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Budget deleted successfully",
	})
}

// GetProgress handles budget progress calculation requests.
// @Summary Get budget progress
// @Description Calculates the progress of a budget, including spent amount, remaining amount, and percentage used.
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Budget ID (UUID)"
// @Success 200 {object} map[string]interface{} "Budget progress retrieved successfully"
// @Success 200 {object} dtos.GetBudgetProgressOutput "Budget progress data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid budget ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - budget does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /budgets/{id}/progress [get]
func (h *BudgetHandler) GetProgress(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get budget ID from path parameter
	budgetID := c.Params("id")
	if budgetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Budget ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.GetBudgetProgressInput{
		BudgetID: budgetID,
		UserID:   userID,
	}

	// Validate input
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.getBudgetProgressUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Budget progress retrieved successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
func (h *BudgetHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "not found") {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errMsg,
			"code":  fiber.StatusNotFound,
		})
	}

	if strings.Contains(errMsg, "already exists") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Já existe um orçamento para esta categoria e período",
			"code":  fiber.StatusConflict,
		})
	}

	if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "cannot be empty") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg,
			"code":  fiber.StatusBadRequest,
		})
	}

	// Default to 500 Internal Server Error
	log.Error().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Budget use case error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
		"code":  fiber.StatusInternalServerError,
	})
}
