package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// GoalHandler handles goal-related HTTP requests.
type GoalHandler struct {
	createGoalUseCase      *usecases.CreateGoalUseCase
	listGoalsUseCase       *usecases.ListGoalsUseCase
	getGoalUseCase         *usecases.GetGoalUseCase
	addContributionUseCase *usecases.AddContributionUseCase
	updateProgressUseCase  *usecases.UpdateProgressUseCase
	cancelGoalUseCase      *usecases.CancelGoalUseCase
	deleteGoalUseCase      *usecases.DeleteGoalUseCase
}

// NewGoalHandler creates a new GoalHandler instance.
func NewGoalHandler(
	createGoalUseCase *usecases.CreateGoalUseCase,
	listGoalsUseCase *usecases.ListGoalsUseCase,
	getGoalUseCase *usecases.GetGoalUseCase,
	addContributionUseCase *usecases.AddContributionUseCase,
	updateProgressUseCase *usecases.UpdateProgressUseCase,
	cancelGoalUseCase *usecases.CancelGoalUseCase,
	deleteGoalUseCase *usecases.DeleteGoalUseCase,
) *GoalHandler {
	return &GoalHandler{
		createGoalUseCase:      createGoalUseCase,
		listGoalsUseCase:       listGoalsUseCase,
		getGoalUseCase:         getGoalUseCase,
		addContributionUseCase: addContributionUseCase,
		updateProgressUseCase:  updateProgressUseCase,
		cancelGoalUseCase:      cancelGoalUseCase,
		deleteGoalUseCase:      deleteGoalUseCase,
	}
}

// Create handles goal creation requests.
// @Summary Create a new goal
// @Description Creates a new goal for the authenticated user.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateGoalInput true "Goal creation data"
// @Success 201 {object} map[string]interface{} "Goal created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals [post]
func (h *GoalHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	var input dtos.CreateGoalInput
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

	output, err := h.createGoalUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Goal created successfully",
		"data":    output,
	})
}

// List handles goal listing requests.
// @Summary List goals
// @Description Lists all goals for the authenticated user.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param context query string false "Filter by context"
// @Param status query string false "Filter by status"
// @Param page query string false "Page number"
// @Param limit query string false "Items per page"
// @Success 200 {object} map[string]interface{} "Goals retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals [get]
func (h *GoalHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	context := c.Query("context", "")
	status := c.Query("status", "")
	page := c.Query("page", "")
	limit := c.Query("limit", "")

	input := dtos.ListGoalsInput{
		UserID:  userID,
		Context: context,
		Status:  status,
		Page:    page,
		Limit:   limit,
	}

	output, err := h.listGoalsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Goals retrieved successfully",
		"data":    output,
	})
}

// Get handles goal retrieval requests.
// @Summary Get goal by ID
// @Description Retrieves a specific goal by its ID.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Goal ID"
// @Success 200 {object} map[string]interface{} "Goal retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals/{id} [get]
func (h *GoalHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	goalID := c.Params("id")
	if goalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.GetGoalInput{
		GoalID: goalID,
		UserID: userID,
	}

	output, err := h.getGoalUseCase.Execute(input)
	if err != nil {
		return h.handleGetGoalError(c, err, goalID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Goal retrieved successfully",
		"data":    output,
	})
}

// AddContribution handles adding a contribution to a goal.
// @Summary Add contribution to goal
// @Description Adds a contribution to a goal.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Goal ID"
// @Param request body dtos.AddContributionInput true "Contribution data"
// @Success 200 {object} map[string]interface{} "Contribution added successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals/{id}/contribute [post]
func (h *GoalHandler) AddContribution(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	goalID := c.Params("id")
	if goalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	var input dtos.AddContributionInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	input.GoalID = goalID
	input.UserID = userID

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.addContributionUseCase.Execute(input)
	if err != nil {
		return h.handleGetGoalError(c, err, goalID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Contribution added successfully",
		"data":    output,
	})
}

// UpdateProgress handles updating goal progress.
// @Summary Update goal progress
// @Description Updates the progress of a goal.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Goal ID"
// @Param request body dtos.UpdateProgressInput true "Progress update data"
// @Success 200 {object} map[string]interface{} "Progress updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals/{id}/progress [put]
func (h *GoalHandler) UpdateProgress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	goalID := c.Params("id")
	if goalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	var input dtos.UpdateProgressInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	input.GoalID = goalID
	input.UserID = userID

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.updateProgressUseCase.Execute(input)
	if err != nil {
		return h.handleGetGoalError(c, err, goalID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Progress updated successfully",
		"data":    output,
	})
}

// Cancel handles goal cancellation requests.
// @Summary Cancel goal
// @Description Cancels a goal.
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Goal ID"
// @Success 200 {object} map[string]interface{} "Goal cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals/{id}/cancel [post]
func (h *GoalHandler) Cancel(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	goalID := c.Params("id")
	if goalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.CancelGoalInput{
		GoalID: goalID,
		UserID: userID,
	}

	output, err := h.cancelGoalUseCase.Execute(input)
	if err != nil {
		return h.handleGetGoalError(c, err, goalID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Goal cancelled successfully",
		"data":    output,
	})
}

// Delete handles goal deletion requests.
// @Summary Delete goal
// @Description Deletes a goal (soft delete).
// @Tags goals
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Goal ID"
// @Success 200 {object} map[string]interface{} "Goal deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /goals/{id} [delete]
func (h *GoalHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	goalID := c.Params("id")
	if goalID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Goal ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.DeleteGoalInput{
		GoalID: goalID,
		UserID: userID,
	}

	output, err := h.deleteGoalUseCase.Execute(input)
	if err != nil {
		return h.handleGetGoalError(c, err, goalID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Goal deleted successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases.
func (h *GoalHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	appErr := apperrors.MapDomainError(err)

	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Goal operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Goal operation failed")
	}

	return appErr
}

// handleGetGoalError handles errors when getting a goal.
func (h *GoalHandler) handleGetGoalError(c *fiber.Ctx, err error, goalID string) error {
	if err == nil {
		return nil
	}

	appErr := apperrors.MapDomainError(err)

	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["goal_id"] = goalID

	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("goal_id", goalID).Msg("Get goal failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("goal_id", goalID).Msg("Get goal failed")
	}

	return appErr
}
