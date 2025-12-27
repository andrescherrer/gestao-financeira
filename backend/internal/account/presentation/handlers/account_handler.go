package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/application/usecases"
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
// @Description Creates a new account for the authenticated user. Supports BANK, WALLET, INVESTMENT, and CREDIT_CARD account types.
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateAccountInput true "Account creation data (name, type, initial_balance, currency, context)"
// @Success 201 {object} map[string]interface{} "Account created successfully"
// @Success 201 {object} dtos.CreateAccountOutput "Account data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 409 {object} map[string]interface{} "Conflict - account already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Account created successfully",
		"data":    output,
	})
}

// List handles account listing requests.
// @Summary List accounts
// @Description Lists all accounts for the authenticated user. Optionally filter by context (PERSONAL or BUSINESS).
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param context query string false "Filter by context (PERSONAL or BUSINESS)"
// @Success 200 {object} map[string]interface{} "Accounts retrieved successfully"
// @Success 200 {object} dtos.ListAccountsOutput "List of accounts with count"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID or context"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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

	// Build input
	input := dtos.ListAccountsInput{
		UserID:  userID,
		Context: context,
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
// @Tags accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Account ID (UUID)"
// @Success 200 {object} map[string]interface{} "Account retrieved successfully"
// @Success 200 {object} dtos.GetAccountOutput "Account data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid account ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - account does not belong to user"
// @Failure 404 {object} map[string]interface{} "Not found - account does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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
func (h *AccountHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid user ID") {
		log.Warn().Err(err).Msg("Account operation failed: invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "invalid account") {
		log.Warn().Err(err).Msg("Account operation failed: invalid account data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account data",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "already exists") {
		log.Warn().Err(err).Msg("Account operation failed: account already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Account already exists",
			"code":  fiber.StatusConflict,
		})
	}

	// Generic error handling
	log.Error().Err(err).Msg("Account operation failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}

// handleGetAccountError handles errors from GetAccountUseCase and returns appropriate HTTP responses.
func (h *AccountHandler) handleGetAccountError(c *fiber.Ctx, err error, accountID string) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid account ID") {
		log.Warn().Err(err).Str("account_id", accountID).Msg("Get account failed: invalid account ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "account not found") {
		log.Warn().Err(err).Str("account_id", accountID).Msg("Get account failed: account not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Generic error handling
	log.Error().Err(err).Str("account_id", accountID).Msg("Get account failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}
