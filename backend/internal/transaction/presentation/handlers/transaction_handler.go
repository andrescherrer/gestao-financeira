package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/application/usecases"
	"gestao-financeira/backend/pkg/middleware"
)

// TransactionHandler handles transaction-related HTTP requests.
type TransactionHandler struct {
	createTransactionUseCase *usecases.CreateTransactionUseCase
	listTransactionsUseCase  *usecases.ListTransactionsUseCase
	getTransactionUseCase    *usecases.GetTransactionUseCase
	updateTransactionUseCase *usecases.UpdateTransactionUseCase
	deleteTransactionUseCase *usecases.DeleteTransactionUseCase
}

// NewTransactionHandler creates a new TransactionHandler instance.
func NewTransactionHandler(
	createTransactionUseCase *usecases.CreateTransactionUseCase,
	listTransactionsUseCase *usecases.ListTransactionsUseCase,
	getTransactionUseCase *usecases.GetTransactionUseCase,
	updateTransactionUseCase *usecases.UpdateTransactionUseCase,
	deleteTransactionUseCase *usecases.DeleteTransactionUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUseCase: createTransactionUseCase,
		listTransactionsUseCase:  listTransactionsUseCase,
		getTransactionUseCase:    getTransactionUseCase,
		updateTransactionUseCase: updateTransactionUseCase,
		deleteTransactionUseCase: deleteTransactionUseCase,
	}
}

// Create handles transaction creation requests.
// @Summary Create a new transaction
// @Description Creates a new transaction (INCOME or EXPENSE) for the authenticated user. Links transaction to an account.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateTransactionInput true "Transaction creation data (account_id, type, amount, currency, description, date)"
// @Success 201 {object} map[string]interface{} "Transaction created successfully"
// @Success 201 {object} dtos.CreateTransactionOutput "Transaction data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data or account not found"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse request body
	var input dtos.CreateTransactionInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set user ID from context (override any user_id in request body for security)
	input.UserID = userID

	// Execute use case
	output, err := h.createTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaction created successfully",
		"data":    output,
	})
}

// List handles transaction listing requests.
// @Summary List transactions
// @Description Lists all transactions for the authenticated user. Optionally filter by account_id and/or type (INCOME or EXPENSE).
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param account_id query string false "Filter by account ID (UUID)"
// @Param type query string false "Filter by transaction type (INCOME or EXPENSE)"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Success 200 {object} dtos.ListTransactionsOutput "List of transactions with count"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID, account ID or type"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions [get]
func (h *TransactionHandler) List(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get optional filters from query parameters
	accountID := c.Query("account_id", "")
	transactionType := c.Query("type", "")

	// Build input
	input := dtos.ListTransactionsInput{
		UserID:    userID,
		AccountID: accountID,
		Type:      transactionType,
	}

	// Execute use case
	output, err := h.listTransactionsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transactions retrieved successfully",
		"data":    output,
	})
}

// Get handles transaction retrieval requests.
// @Summary Get transaction by ID
// @Description Retrieves a specific transaction by its ID. Only returns transactions that belong to the authenticated user.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction retrieved successfully"
// @Success 200 {object} dtos.GetTransactionOutput "Transaction data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - transaction does not belong to user"
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) Get(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get transaction ID from path parameter
	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.GetTransactionInput{
		TransactionID: transactionID,
	}

	// Execute use case
	output, err := h.getTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleGetTransactionError(c, err, transactionID)
	}

	// Verify that the transaction belongs to the authenticated user
	if output.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
			"code":  fiber.StatusForbidden,
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction retrieved successfully",
		"data":    output,
	})
}

// Update handles transaction update requests.
// @Summary Update a transaction
// @Description Updates an existing transaction. Supports partial updates - only provided fields will be updated. At least one field must be provided.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)"
// @Param request body dtos.UpdateTransactionInput true "Transaction update data (all fields optional: type, amount, currency, description, date)"
// @Success 200 {object} map[string]interface{} "Transaction updated successfully"
// @Success 200 {object} dtos.UpdateTransactionOutput "Updated transaction data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID, invalid data, or no fields provided"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - transaction does not belong to user"
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get transaction ID from path parameter
	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Parse request body
	var input dtos.UpdateTransactionInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set transaction ID from path parameter (override any transaction_id in request body)
	input.TransactionID = transactionID

	// Execute use case
	output, err := h.updateTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleUpdateTransactionError(c, err, transactionID)
	}

	// Verify that the transaction belongs to the authenticated user
	if output.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
			"code":  fiber.StatusForbidden,
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction updated successfully",
		"data":    output,
	})
}

// Delete handles transaction deletion requests.
// @Summary Delete a transaction
// @Description Deletes a transaction by its ID (soft delete). The transaction is marked as deleted but not permanently removed from the database.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Success 200 {object} dtos.DeleteTransactionOutput "Deletion confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Get transaction ID from path parameter
	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.DeleteTransactionInput{
		TransactionID: transactionID,
	}

	// Execute use case
	output, err := h.deleteTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleDeleteTransactionError(c, err, transactionID)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
func (h *TransactionHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid user ID") {
		log.Warn().Err(err).Msg("Transaction operation failed: invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "invalid transaction") {
		log.Warn().Err(err).Msg("Transaction operation failed: invalid transaction data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction data",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "account not found") {
		log.Warn().Err(err).Msg("Transaction operation failed: account not found")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account not found",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Generic error handling
	log.Error().Err(err).Msg("Transaction operation failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}

// handleGetTransactionError handles errors from GetTransactionUseCase and returns appropriate HTTP responses.
func (h *TransactionHandler) handleGetTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid transaction ID") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Get transaction failed: invalid transaction ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "transaction not found") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Get transaction failed: transaction not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Generic error handling
	log.Error().Err(err).Str("transaction_id", transactionID).Msg("Get transaction failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}

// handleUpdateTransactionError handles errors from UpdateTransactionUseCase and returns appropriate HTTP responses.
func (h *TransactionHandler) handleUpdateTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid transaction ID") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Update transaction failed: invalid transaction ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "transaction not found") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Update transaction failed: transaction not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
			"code":  fiber.StatusNotFound,
		})
	}

	if strings.Contains(errMsg, "invalid transaction") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Update transaction failed: invalid transaction data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction data",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Generic error handling
	log.Error().Err(err).Str("transaction_id", transactionID).Msg("Update transaction failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}

// handleDeleteTransactionError handles errors from DeleteTransactionUseCase and returns appropriate HTTP responses.
func (h *TransactionHandler) handleDeleteTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid transaction ID") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Delete transaction failed: invalid transaction ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "transaction not found") {
		log.Warn().Err(err).Str("transaction_id", transactionID).Msg("Delete transaction failed: transaction not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
			"code":  fiber.StatusNotFound,
		})
	}

	// Generic error handling
	log.Error().Err(err).Str("transaction_id", transactionID).Msg("Delete transaction failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}
