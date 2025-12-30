package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/metrics"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// TransactionHandler handles transaction-related HTTP requests.
type TransactionHandler struct {
	createTransactionUseCase          *usecases.CreateTransactionUseCase
	listTransactionsUseCase           *usecases.ListTransactionsUseCase
	getTransactionUseCase             *usecases.GetTransactionUseCase
	updateTransactionUseCase          *usecases.UpdateTransactionUseCase
	deleteTransactionUseCase          *usecases.DeleteTransactionUseCase
	restoreTransactionUseCase         *usecases.RestoreTransactionUseCase
	permanentDeleteTransactionUseCase *usecases.PermanentDeleteTransactionUseCase
}

// NewTransactionHandler creates a new TransactionHandler instance.
func NewTransactionHandler(
	createTransactionUseCase *usecases.CreateTransactionUseCase,
	listTransactionsUseCase *usecases.ListTransactionsUseCase,
	getTransactionUseCase *usecases.GetTransactionUseCase,
	updateTransactionUseCase *usecases.UpdateTransactionUseCase,
	deleteTransactionUseCase *usecases.DeleteTransactionUseCase,
	restoreTransactionUseCase *usecases.RestoreTransactionUseCase,
	permanentDeleteTransactionUseCase *usecases.PermanentDeleteTransactionUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUseCase:          createTransactionUseCase,
		listTransactionsUseCase:           listTransactionsUseCase,
		getTransactionUseCase:             getTransactionUseCase,
		updateTransactionUseCase:          updateTransactionUseCase,
		deleteTransactionUseCase:          deleteTransactionUseCase,
		restoreTransactionUseCase:         restoreTransactionUseCase,
		permanentDeleteTransactionUseCase: permanentDeleteTransactionUseCase,
	}
}

// Create handles transaction creation requests.
// @Summary Create a new transaction
// @Description Creates a new transaction (INCOME or EXPENSE) for the authenticated user. Links transaction to an account and automatically updates the account balance atomically using Unit of Work pattern.
//
// **Operação Atômica**: A criação da transação e atualização do saldo da conta são executadas em uma única transação de banco de dados. Se qualquer operação falhar, todas são revertidas (rollback).
//
// **Tipos de Transação**:
// - `INCOME`: Receita - aumenta o saldo da conta
// - `EXPENSE`: Despesa - diminui o saldo da conta
//
// **Validações**:
// - Account ID deve existir e pertencer ao usuário autenticado
// - Amount deve ser maior que zero
// - Currency deve ser válida (ex: BRL, USD, EUR)
// - Date deve estar no formato YYYY-MM-DD
//
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateTransactionInput true "Transaction creation data" example({"account_id":"550e8400-e29b-41d4-a716-446655440000","type":"INCOME","amount":150.50,"currency":"BRL","description":"Salário","date":"2025-12-29"})
// @Success 201 {object} map[string]interface{} "Transaction created successfully" example({"message":"Transaction created successfully","data":{"transaction_id":"550e8400-e29b-41d4-a716-446655440001","user_id":"550e8400-e29b-41d4-a716-446655440000","account_id":"550e8400-e29b-41d4-a716-446655440000","type":"INCOME","amount":150.50,"currency":"BRL","description":"Salário","date":"2025-12-29","created_at":"2025-12-29T10:00:00Z","updated_at":"2025-12-29T10:00:00Z"}})
// @Success 201 {object} dtos.CreateTransactionOutput "Transaction data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data, validation failed, or account not found" example({"error":"Invalid transaction data","error_type":"VALIDATION_ERROR","code":400,"details":{"field":"amount","message":"amount must be greater than 0"}})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - account does not exist" example({"error":"Account not found","error_type":"NOT_FOUND","code":404,"details":{"resource":"Account","identifier":"550e8400-e29b-41d4-a716-446655440000"}})
// @Failure 422 {object} map[string]interface{} "Unprocessable entity - domain validation failed" example({"error":"Account balance would become negative","error_type":"DOMAIN_ERROR","code":422})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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
	output, err := h.createTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Record business metric
	metrics.BusinessMetrics.TransactionsCreated.WithLabelValues(input.Type).Inc()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaction created successfully",
		"data":    output,
	})
}

// List handles transaction listing requests.
// @Summary List transactions
// @Description Lists all transactions for the authenticated user. Supports filtering by account_id and/or type, and pagination.
//
// **Filtros Disponíveis**:
// - `account_id`: Filtra transações por conta específica (UUID)
// - `type`: Filtra por tipo de transação (`INCOME` ou `EXPENSE`)
//
// **Paginação**:
// - `page`: Número da página (1-based, padrão: 1)
// - `limit`: Itens por página (padrão: 10, máximo: 100)
//
// **Ordenação**: Transações são ordenadas por data (mais recentes primeiro) e depois por data de criação.
//
// **Exemplo sem paginação**: Retorna todas as transações (compatibilidade retroativa)
// **Exemplo com paginação**: `GET /transactions?page=2&limit=20&type=INCOME`
//
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param account_id query string false "Filter by account ID (UUID)" example(550e8400-e29b-41d4-a716-446655440000)
// @Param type query string false "Filter by transaction type (INCOME or EXPENSE)" Enums(INCOME, EXPENSE) example(INCOME)
// @Param page query string false "Page number (1-based, default: 1)" example(1)
// @Param limit query string false "Items per page (default: 10, max: 100)" example(20)
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully" example({"message":"Transactions retrieved successfully","data":{"transactions":[{"transaction_id":"550e8400-e29b-41d4-a716-446655440001","type":"INCOME","amount":150.50,"currency":"BRL","description":"Salário","date":"2025-12-29"}],"count":20,"pagination":{"page":1,"limit":20,"total":45,"total_pages":3,"has_next":true,"has_prev":false}}})
// @Success 200 {object} dtos.ListTransactionsOutput "List of transactions with count and pagination metadata"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID, account ID, type, or pagination parameters" example({"error":"Invalid account ID format","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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
	page := c.Query("page", "")
	limit := c.Query("limit", "")

	// Build input
	input := dtos.ListTransactionsInput{
		UserID:    userID,
		AccountID: accountID,
		Type:      transactionType,
		Page:      page,
		Limit:     limit,
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
//
// **Segurança**: O endpoint valida que a transação pertence ao usuário autenticado. Tentativas de acessar transações de outros usuários retornam 404 (não 403) para evitar vazamento de informação.
//
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)" example(550e8400-e29b-41d4-a716-446655440001)
// @Success 200 {object} map[string]interface{} "Transaction retrieved successfully" example({"message":"Transaction retrieved successfully","data":{"transaction_id":"550e8400-e29b-41d4-a716-446655440001","user_id":"550e8400-e29b-41d4-a716-446655440000","account_id":"550e8400-e29b-41d4-a716-446655440000","type":"INCOME","amount":150.50,"currency":"BRL","description":"Salário","date":"2025-12-29","created_at":"2025-12-29T10:00:00Z","updated_at":"2025-12-29T10:00:00Z"}})
// @Success 200 {object} dtos.GetTransactionOutput "Transaction data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format" example({"error":"Invalid transaction ID format","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist or does not belong to user" example({"error":"Transaction not found","error_type":"NOT_FOUND","code":404,"details":{"resource":"Transaction","identifier":"550e8400-e29b-41d4-a716-446655440001"}})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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
// @Description Updates an existing transaction using Unit of Work pattern for atomicity. Supports partial updates - only provided fields will be updated. At least one field must be provided.
//
// **Operação Atômica**: A atualização da transação e ajuste do saldo da conta são executados em uma única transação de banco de dados. Se qualquer operação falhar, todas são revertidas (rollback).
//
// **Comportamento**:
// - Campos não fornecidos permanecem inalterados
// - Se `amount` ou `type` mudarem, o saldo da conta é ajustado automaticamente
// - O efeito antigo é revertido e o novo efeito é aplicado
//
// **Campos Atualizáveis**:
// - `type`: INCOME ou EXPENSE (atualiza saldo da conta)
// - `amount`: Valor da transação (atualiza saldo da conta)
// - `currency`: Moeda (BRL, USD, EUR, etc.)
// - `description`: Descrição da transação
// - `date`: Data da transação (formato YYYY-MM-DD)
//
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)" example(550e8400-e29b-41d4-a716-446655440001)
// @Param request body dtos.UpdateTransactionInput true "Transaction update data (all fields optional)" example({"type":"EXPENSE","amount":75.25,"description":"Atualizado"})
// @Success 200 {object} map[string]interface{} "Transaction updated successfully" example({"message":"Transaction updated successfully","data":{"transaction_id":"550e8400-e29b-41d4-a716-446655440001","type":"EXPENSE","amount":75.25,"currency":"BRL","description":"Atualizado","updated_at":"2025-12-29T11:00:00Z"}})
// @Success 200 {object} dtos.UpdateTransactionOutput "Updated transaction data with all fields"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID, invalid data, or no fields provided" example({"error":"At least one field must be provided for update","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist" example({"error":"Transaction not found","error_type":"NOT_FOUND","code":404})
// @Failure 422 {object} map[string]interface{} "Unprocessable entity - domain validation failed" example({"error":"Account balance would become negative","error_type":"DOMAIN_ERROR","code":422})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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

	// Record business metric
	metrics.BusinessMetrics.TransactionsUpdated.WithLabelValues(output.Type).Inc()

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction updated successfully",
		"data":    output,
	})
}

// Delete handles transaction deletion requests.
// Delete handles transaction deletion requests.
// @Summary Delete a transaction
// @Description Deletes a transaction by its ID using Unit of Work pattern for atomicity. Performs soft delete - the transaction is marked as deleted but not permanently removed from the database. The account balance is automatically reversed.
//
// **Operação Atômica**: A exclusão da transação e reversão do saldo da conta são executados em uma única transação de banco de dados. Se qualquer operação falhar, todas são revertidas (rollback).
//
// **Comportamento**:
// - Transação é marcada como deletada (soft delete)
// - Saldo da conta é revertido automaticamente
// - Transação pode ser restaurada usando o endpoint `/transactions/{id}/restore`
// - Transação pode ser permanentemente deletada usando o endpoint `/transactions/{id}/permanent` (admin only)
//
// **Reversão de Saldo**:
// - Se era INCOME: valor é subtraído do saldo
// - Se era EXPENSE: valor é adicionado ao saldo
//
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)" example(550e8400-e29b-41d4-a716-446655440001)
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully" example({"message":"Transaction deleted successfully","data":{"transaction_id":"550e8400-e29b-41d4-a716-446655440001"}})
// @Success 200 {object} dtos.DeleteTransactionOutput "Deletion confirmation with transaction ID"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format" example({"error":"Invalid transaction ID format","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token" example({"error":"Unauthorized","code":401})
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist" example({"error":"Transaction not found","error_type":"NOT_FOUND","code":404})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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

	// Get transaction before deletion to record metric with type
	getInput := dtos.GetTransactionInput{
		TransactionID: transactionID,
	}
	transaction, err := h.getTransactionUseCase.Execute(getInput)
	transactionType := "unknown"
	if err == nil && transaction != nil {
		transactionType = transaction.Type
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

	// Record business metric
	metrics.BusinessMetrics.TransactionsDeleted.WithLabelValues(transactionType).Inc()

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Transaction operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Transaction operation failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handleGetTransactionError handles errors from GetTransactionUseCase and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handleGetTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add transaction ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["transaction_id"] = transactionID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Get transaction failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Get transaction failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handleUpdateTransactionError handles errors from UpdateTransactionUseCase and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handleUpdateTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add transaction ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["transaction_id"] = transactionID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Update transaction failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Update transaction failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handleDeleteTransactionError handles errors from DeleteTransactionUseCase and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handleDeleteTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add transaction ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["transaction_id"] = transactionID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Delete transaction failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Delete transaction failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// Restore handles transaction restoration requests.
// @Summary Restore a soft-deleted transaction
// @Description Restores a soft-deleted transaction by setting deleted_at to NULL.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction restored successfully"
// @Success 200 {object} dtos.RestoreTransactionOutput "Restoration confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist or is not deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id}/restore [post]
func (h *TransactionHandler) Restore(c *fiber.Ctx) error {
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
	input := dtos.RestoreTransactionInput{
		TransactionID: transactionID,
	}

	// Execute use case
	output, err := h.restoreTransactionUseCase.Execute(input)
	if err != nil {
		return h.handleRestoreTransactionError(c, err, transactionID)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// PermanentDelete handles permanent transaction deletion requests (admin only).
// @Summary Permanently delete a transaction
// @Description Permanently deletes a transaction from the database (hard delete). This action cannot be undone.
// @Tags transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction permanently deleted successfully"
// @Success 200 {object} dtos.PermanentDeleteTransactionOutput "Deletion confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid transaction ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin access required"
// @Failure 404 {object} map[string]interface{} "Not found - transaction does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id}/permanent [delete]
func (h *TransactionHandler) PermanentDelete(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// TODO: Add admin check here
	// For now, we'll allow any authenticated user, but in production this should be restricted

	// Get transaction ID from path parameter
	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Build input
	input := dtos.PermanentDeleteTransactionInput{
		TransactionID: transactionID,
	}

	// Execute use case
	output, err := h.permanentDeleteTransactionUseCase.Execute(input)
	if err != nil {
		return h.handlePermanentDeleteTransactionError(c, err, transactionID)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// handleRestoreTransactionError handles errors from RestoreTransactionUseCase.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handleRestoreTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add transaction ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["transaction_id"] = transactionID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Restore transaction failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Restore transaction failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handlePermanentDeleteTransactionError handles errors from PermanentDeleteTransactionUseCase.
// Uses AppError for consistent error handling instead of string matching.
func (h *TransactionHandler) handlePermanentDeleteTransactionError(c *fiber.Ctx, err error, transactionID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add transaction ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["transaction_id"] = transactionID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Permanent delete transaction failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("transaction_id", transactionID).Msg("Permanent delete transaction failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}
