package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/metrics"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// AccountHandler handles account-related HTTP requests.
type AccountHandler struct {
	createAccountUseCase *usecases.CreateAccountUseCase
	listAccountsUseCase  *usecases.ListAccountsUseCase
	getAccountUseCase    *usecases.GetAccountUseCase
}

// NewAccountHandler creates a new AccountHandler instance.
func NewAccountHandler(
	createAccountUseCase *usecases.CreateAccountUseCase,
	listAccountsUseCase *usecases.ListAccountsUseCase,
	getAccountUseCase *usecases.GetAccountUseCase,
) *AccountHandler {
	return &AccountHandler{
		createAccountUseCase: createAccountUseCase,
		listAccountsUseCase:  listAccountsUseCase,
		getAccountUseCase:    getAccountUseCase,
	}
}

// Create handles account creation requests.
// @Summary Create a new account
// @Description Creates a new account for the authenticated user. Supports multiple account types and contexts.
//
// **Tipos de Conta Suportados**:
// - `BANK`: Conta bancária tradicional
// - `WALLET`: Carteira digital (ex: PayPal, Mercado Pago)
// - `INVESTMENT`: Conta de investimentos
// - `CREDIT_CARD`: Cartão de crédito
//
// **Contextos**:
// - `PERSONAL`: Conta pessoal
// - `BUSINESS`: Conta empresarial
//
// **Validações**:
// - Nome da conta é obrigatório e deve ser único para o usuário
// - Saldo inicial pode ser zero ou positivo
// - Moeda deve ser válida (ex: BRL, USD, EUR)
//
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateAccountInput true "Account creation data" example({"name":"Conta Corrente","type":"BANK","initial_balance":1000.00,"currency":"BRL","context":"PERSONAL"})
// @Success 201 {object} map[string]interface{} "Account created successfully" example({"message":"Account created successfully","data":{"account_id":"550e8400-e29b-41d4-a716-446655440000","user_id":"550e8400-e29b-41d4-a716-446655440000","name":"Conta Corrente","type":"BANK","balance":1000.00,"currency":"BRL","context":"PERSONAL","is_active":true,"created_at":"2025-12-29T10:00:00Z"}})
// @Success 201 {object} dtos.CreateAccountOutput "Account data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data or validation failed" example({"error":"Invalid account data","error_type":"VALIDATION_ERROR","code":400,"details":{"field":"name","message":"name cannot be empty"}})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 409 {object} map[string]interface{} "Conflict - account with this name already exists" example({"error":"Account already exists","error_type":"CONFLICT","code":409})
// @Failure 422 {object} map[string]interface{} "Unprocessable entity - domain validation failed" example({"error":"Invalid account type","error_type":"DOMAIN_ERROR","code":422})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
// @Router /accounts [post]
func (h *AccountHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse request body
	var input dtos.CreateAccountInput
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
	output, err := h.createAccountUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Record business metric
	metrics.BusinessMetrics.AccountsCreated.WithLabelValues(input.Type).Inc()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account created successfully",
		"data":    output,
	})
}

// List handles account listing requests.
// @Summary List accounts
// @Description Lists all accounts for the authenticated user. Supports filtering by context and pagination.
//
// **Filtros Disponíveis**:
// - `context`: Filtra por contexto (`PERSONAL` ou `BUSINESS`)
//
// **Paginação**:
// - `page`: Número da página (1-based, padrão: 1)
// - `limit`: Itens por página (padrão: 10, máximo: 100)
//
// **Ordenação**: Contas são ordenadas por data de criação (mais recentes primeiro).
//
// **Exemplo sem paginação**: Retorna todas as contas (compatibilidade retroativa)
// **Exemplo com paginação**: `GET /accounts?page=1&limit=20&context=PERSONAL`
//
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param context query string false "Filter by context (PERSONAL or BUSINESS)" Enums(PERSONAL, BUSINESS) example(PERSONAL)
// @Param page query string false "Page number (1-based, default: 1)" example(1)
// @Param limit query string false "Items per page (default: 10, max: 100)" example(20)
// @Success 200 {object} map[string]interface{} "Accounts retrieved successfully" example({"message":"Accounts retrieved successfully","data":{"accounts":[{"account_id":"550e8400-e29b-41d4-a716-446655440000","name":"Conta Corrente","type":"BANK","balance":1000.00,"currency":"BRL","context":"PERSONAL"}],"count":20,"pagination":{"page":1,"limit":20,"total":45,"total_pages":3,"has_next":true,"has_prev":false}}})
// @Success 200 {object} dtos.ListAccountsOutput "List of accounts with count and pagination metadata"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID, context, or pagination parameters" example({"error":"Invalid context value","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
// @Router /accounts [get]
func (h *AccountHandler) List(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get optional context filter from query parameter
	context := c.Query("context", "")
	// Get pagination parameters
	page := c.Query("page", "")
	limit := c.Query("limit", "")

	// Build input
	input := dtos.ListAccountsInput{
		UserID:  userID,
		Context: context,
		Page:    page,
		Limit:   limit,
	}

	// Execute use case
	output, err := h.listAccountsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Accounts retrieved successfully",
		"data":    output,
	})
}

// Get handles account retrieval requests.
// @Summary Get account by ID
// @Description Retrieves a specific account by its ID. Only returns accounts that belong to the authenticated user.
//
// **Segurança**: O endpoint valida que a conta pertence ao usuário autenticado. Tentativas de acessar contas de outros usuários retornam 404 (não 403) para evitar vazamento de informação.
//
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Account ID (UUID)" example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} map[string]interface{} "Account retrieved successfully" example({"message":"Account retrieved successfully","data":{"account_id":"550e8400-e29b-41d4-a716-446655440000","user_id":"550e8400-e29b-41d4-a716-446655440000","name":"Conta Corrente","type":"BANK","balance":1000.00,"currency":"BRL","context":"PERSONAL","is_active":true,"created_at":"2025-12-29T10:00:00Z","updated_at":"2025-12-29T10:00:00Z"}})
// @Success 200 {object} dtos.GetAccountOutput "Account data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid account ID format" example({"error":"Invalid account ID format","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - account does not exist or does not belong to user" example({"error":"Account not found","error_type":"NOT_FOUND","code":404,"details":{"resource":"Account","identifier":"550e8400-e29b-41d4-a716-446655440000"}})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
// @Router /accounts/{id} [get]
func (h *AccountHandler) Get(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get account ID from path parameter
	accountID := c.Params("id")
	if accountID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.GetAccountInput{
		AccountID: accountID,
	}

	// Execute use case
	output, err := h.getAccountUseCase.Execute(input)
	if err != nil {
		return h.handleGetAccountError(c, err, accountID)
	}

	// Verify that the account belongs to the authenticated user
	if output.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
			"code":  fiber.StatusForbidden,
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Account retrieved successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *AccountHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Account operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Account operation failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handleGetAccountError handles errors from GetAccountUseCase and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *AccountHandler) handleGetAccountError(c *fiber.Ctx, err error, accountID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add account ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["account_id"] = accountID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("account_id", accountID).Msg("Get account failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("account_id", accountID).Msg("Get account failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}
