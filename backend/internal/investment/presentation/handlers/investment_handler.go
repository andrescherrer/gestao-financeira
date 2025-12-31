package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/metrics"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// InvestmentHandler handles investment-related HTTP requests.
type InvestmentHandler struct {
	createInvestmentUseCase *usecases.CreateInvestmentUseCase
	listInvestmentsUseCase  *usecases.ListInvestmentsUseCase
	getInvestmentUseCase    *usecases.GetInvestmentUseCase
	updateInvestmentUseCase *usecases.UpdateInvestmentUseCase
	deleteInvestmentUseCase *usecases.DeleteInvestmentUseCase
}

// NewInvestmentHandler creates a new InvestmentHandler instance.
func NewInvestmentHandler(
	createInvestmentUseCase *usecases.CreateInvestmentUseCase,
	listInvestmentsUseCase *usecases.ListInvestmentsUseCase,
	getInvestmentUseCase *usecases.GetInvestmentUseCase,
	updateInvestmentUseCase *usecases.UpdateInvestmentUseCase,
	deleteInvestmentUseCase *usecases.DeleteInvestmentUseCase,
) *InvestmentHandler {
	return &InvestmentHandler{
		createInvestmentUseCase: createInvestmentUseCase,
		listInvestmentsUseCase:  listInvestmentsUseCase,
		getInvestmentUseCase:    getInvestmentUseCase,
		updateInvestmentUseCase: updateInvestmentUseCase,
		deleteInvestmentUseCase: deleteInvestmentUseCase,
	}
}

// Create handles investment creation requests.
// @Summary Create a new investment
// @Description Creates a new investment for the authenticated user.
// @Tags investments
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateInvestmentInput true "Investment creation data"
// @Success 201 {object} map[string]interface{} "Investment created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /investments [post]
func (h *InvestmentHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	var input dtos.CreateInvestmentInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	input.UserID = userID

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.createInvestmentUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	metrics.BusinessMetrics.InvestmentsCreated.WithLabelValues(input.Type).Inc()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Investment created successfully",
		"data":    output,
	})
}

// List handles investment listing requests.
// @Summary List investments
// @Description Lists all investments for the authenticated user.
// @Tags investments
// @Accept json
// @Produce json
// @Security Bearer
// @Param context query string false "Filter by context"
// @Param type query string false "Filter by type"
// @Param page query string false "Page number"
// @Param limit query string false "Items per page"
// @Success 200 {object} map[string]interface{} "Investments retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /investments [get]
func (h *InvestmentHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	context := c.Query("context", "")
	investmentType := c.Query("type", "")
	page := c.Query("page", "")
	limit := c.Query("limit", "")

	input := dtos.ListInvestmentsInput{
		UserID:  userID,
		Context: context,
		Type:    investmentType,
		Page:    page,
		Limit:   limit,
	}

	output, err := h.listInvestmentsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Investments retrieved successfully",
		"data":    output,
	})
}

// Get handles investment retrieval requests.
// @Summary Get investment by ID
// @Description Retrieves a specific investment by its ID.
// @Tags investments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Investment ID"
// @Success 200 {object} map[string]interface{} "Investment retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /investments/{id} [get]
func (h *InvestmentHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	investmentID := c.Params("id")
	if investmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.GetInvestmentInput{
		InvestmentID: investmentID,
	}

	output, err := h.getInvestmentUseCase.Execute(input)
	if err != nil {
		return h.handleGetInvestmentError(c, err, investmentID)
	}

	if output.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
			"code":  fiber.StatusForbidden,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Investment retrieved successfully",
		"data":    output,
	})
}

// Update handles investment update requests.
// @Summary Update investment
// @Description Updates an investment (current value or quantity).
// @Tags investments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Investment ID"
// @Param request body dtos.UpdateInvestmentInput true "Investment update data"
// @Success 200 {object} map[string]interface{} "Investment updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /investments/{id} [put]
func (h *InvestmentHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	investmentID := c.Params("id")
	if investmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	var input dtos.UpdateInvestmentInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	input.InvestmentID = investmentID

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.updateInvestmentUseCase.Execute(input)
	if err != nil {
		return h.handleGetInvestmentError(c, err, investmentID)
	}

	if output.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
			"code":  fiber.StatusForbidden,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Investment updated successfully",
		"data":    output,
	})
}

// Delete handles investment deletion requests.
// @Summary Delete investment
// @Description Deletes an investment (soft delete).
// @Tags investments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Investment ID"
// @Success 200 {object} map[string]interface{} "Investment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /investments/{id} [delete]
func (h *InvestmentHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	investmentID := c.Params("id")
	if investmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Investment ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.DeleteInvestmentInput{
		InvestmentID: investmentID,
	}

	output, err := h.deleteInvestmentUseCase.Execute(input)
	if err != nil {
		return h.handleGetInvestmentError(c, err, investmentID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Investment deleted successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases.
func (h *InvestmentHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	appErr := apperrors.MapDomainError(err)

	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Investment operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Investment operation failed")
	}

	return appErr
}

// handleGetInvestmentError handles errors from GetInvestmentUseCase.
func (h *InvestmentHandler) handleGetInvestmentError(c *fiber.Ctx, err error, investmentID string) error {
	if err == nil {
		return nil
	}

	appErr := apperrors.MapDomainError(err)

	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["investment_id"] = investmentID

	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("investment_id", investmentID).Msg("Get investment failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("investment_id", investmentID).Msg("Get investment failed")
	}

	return appErr
}
