package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/application/usecases"
	"gestao-financeira/backend/pkg/metrics"
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
// Create handles budget creation requests.
// @Summary Create a new budget
// @Description Creates a new budget for the authenticated user. Supports MONTHLY and YEARLY period types.
//
// **Tipos de Período**:
// - `MONTHLY`: Orçamento mensal (requer `year` e `month`)
// - `YEARLY`: Orçamento anual (requer apenas `year`)
//
// **Validações**:
// - Category ID deve existir e pertencer ao usuário
// - Amount deve ser maior que zero
// - Currency deve ser válida (ex: BRL, USD, EUR)
// - Year deve estar entre 1900 e 3000
// - Month deve estar entre 1 e 12 (apenas para MONTHLY)
// - Não pode existir outro orçamento para a mesma categoria e período
//
// **Contextos**:
// - `PERSONAL`: Orçamento pessoal
// - `BUSINESS`: Orçamento empresarial
//
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateBudgetInput true "Budget creation data" example({"category_id":"550e8400-e29b-41d4-a716-446655440000","amount":1000.00,"currency":"BRL","period_type":"MONTHLY","year":2025,"month":12,"context":"PERSONAL"})
// @Success 201 {object} map[string]interface{} "Budget created successfully" example({"message":"Budget created successfully","data":{"budget_id":"550e8400-e29b-41d4-a716-446655440000","user_id":"550e8400-e29b-41d4-a716-446655440000","category_id":"550e8400-e29b-41d4-a716-446655440000","amount":1000.00,"currency":"BRL","period_type":"MONTHLY","year":2025,"month":12,"context":"PERSONAL","is_active":true,"created_at":"2025-12-29T10:00:00Z"}})
// @Success 201 {object} dtos.CreateBudgetOutput "Budget data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data or validation failed" example({"error":"Invalid budget data","error_type":"VALIDATION_ERROR","code":400,"details":{"field":"amount","message":"amount must be greater than 0"}})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist" example({"error":"Category not found","error_type":"NOT_FOUND","code":404})
// @Failure 409 {object} map[string]interface{} "Conflict - budget already exists for this category and period" example({"error":"Budget already exists for this category and period","error_type":"CONFLICT","code":409})
// @Failure 422 {object} map[string]interface{} "Unprocessable entity - domain validation failed" example({"error":"Invalid period type","error_type":"DOMAIN_ERROR","code":422})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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

	// Record business metric
	metrics.BusinessMetrics.BudgetsCreated.WithLabelValues(input.PeriodType).Inc()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Budget created successfully",
		"data":    output,
	})
}

// List handles budget listing requests.
// @Summary List budgets
// @Description Lists all budgets for the authenticated user. Supports multiple filters and pagination. Filters are applied in memory, then pagination is applied to the filtered results.
//
// **Filtros Disponíveis**:
// - `category_id`: Filtra por categoria específica (UUID)
// - `period_type`: Filtra por tipo de período (`MONTHLY` ou `YEARLY`)
// - `year`: Filtra por ano (ex: 2025)
// - `month`: Filtra por mês (1-12, apenas para orçamentos mensais)
// - `context`: Filtra por contexto (`PERSONAL` ou `BUSINESS`)
// - `is_active`: Filtra por status ativo/inativo (`true` ou `false`)
//
// **Paginação**:
// - `page`: Número da página (1-based, padrão: 1)
// - `limit`: Itens por página (padrão: 10, máximo: 100)
//
// **Nota**: Filtros são aplicados em memória após buscar todos os orçamentos do usuário. Para melhor performance com muitos orçamentos, considere usar filtros específicos.
//
// **Exemplo com filtros e paginação**: `GET /budgets?year=2025&period_type=MONTHLY&page=1&limit=20`
//
// @Tags budgets
// @Accept json
// @Produce json
// @Security Bearer
// @Param category_id query string false "Filter by category ID (UUID)" example(550e8400-e29b-41d4-a716-446655440000)
// @Param period_type query string false "Filter by period type (MONTHLY or YEARLY)" Enums(MONTHLY, YEARLY) example(MONTHLY)
// @Param year query int false "Filter by year" example(2025)
// @Param month query int false "Filter by month (1-12)" example(12)
// @Param context query string false "Filter by context (PERSONAL or BUSINESS)" Enums(PERSONAL, BUSINESS) example(PERSONAL)
// @Param is_active query bool false "Filter by active status (true or false)" example(true)
// @Param page query string false "Page number (1-based, default: 1)" example(1)
// @Param limit query string false "Items per page (default: 10, max: 100)" example(20)
// @Success 200 {object} map[string]interface{} "Budgets retrieved successfully" example({"message":"Budgets retrieved successfully","data":{"budgets":[{"budget_id":"550e8400-e29b-41d4-a716-446655440000","category_id":"550e8400-e29b-41d4-a716-446655440000","amount":1000.00,"currency":"BRL","period_type":"MONTHLY","year":2025,"month":12}],"total":12,"pagination":{"page":1,"limit":20,"total":12,"total_pages":1,"has_next":false,"has_prev":false}}})
// @Success 200 {object} dtos.ListBudgetsOutput "List of budgets with count and pagination metadata"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID, filters, or pagination parameters" example({"error":"Invalid year value","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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

	// Get pagination parameters
	input.Page = c.Query("page", "")
	input.Limit = c.Query("limit", "")

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
