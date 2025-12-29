package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// CategoryHandler handles category-related HTTP requests.
type CategoryHandler struct {
	createCategoryUseCase          *usecases.CreateCategoryUseCase
	listCategoriesUseCase          *usecases.ListCategoriesUseCase
	getCategoryUseCase             *usecases.GetCategoryUseCase
	updateCategoryUseCase          *usecases.UpdateCategoryUseCase
	deleteCategoryUseCase          *usecases.DeleteCategoryUseCase
	restoreCategoryUseCase         *usecases.RestoreCategoryUseCase
	permanentDeleteCategoryUseCase *usecases.PermanentDeleteCategoryUseCase
}

// NewCategoryHandler creates a new CategoryHandler instance.
func NewCategoryHandler(
	createCategoryUseCase *usecases.CreateCategoryUseCase,
	listCategoriesUseCase *usecases.ListCategoriesUseCase,
	getCategoryUseCase *usecases.GetCategoryUseCase,
	updateCategoryUseCase *usecases.UpdateCategoryUseCase,
	deleteCategoryUseCase *usecases.DeleteCategoryUseCase,
	restoreCategoryUseCase *usecases.RestoreCategoryUseCase,
	permanentDeleteCategoryUseCase *usecases.PermanentDeleteCategoryUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		createCategoryUseCase:          createCategoryUseCase,
		listCategoriesUseCase:          listCategoriesUseCase,
		getCategoryUseCase:             getCategoryUseCase,
		updateCategoryUseCase:          updateCategoryUseCase,
		deleteCategoryUseCase:          deleteCategoryUseCase,
		restoreCategoryUseCase:         restoreCategoryUseCase,
		permanentDeleteCategoryUseCase: permanentDeleteCategoryUseCase,
	}
}

// Create handles category creation requests.
// @Summary Create a new category
// @Description Creates a new category for the authenticated user.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateCategoryInput true "Category creation data (name, description)"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Success 201 {object} dtos.CreateCategoryOutput "Category data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse request body
	var input dtos.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set user ID from context
	input.UserID = userID

	// Validate input
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.createCategoryUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
		"data":    output,
	})
}

// List handles category listing requests.
// @Summary List categories
// @Description Lists all categories for the authenticated user. Optionally filter by active status.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param is_active query bool false "Filter by active status (true or false)"
// @Success 200 {object} map[string]interface{} "Categories retrieved successfully"
// @Success 200 {object} dtos.ListCategoriesOutput "List of categories with count"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	// Parse query parameters
	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if strings.ToLower(isActiveStr) == "true" {
			active := true
			isActive = &active
		} else if strings.ToLower(isActiveStr) == "false" {
			active := false
			isActive = &active
		}
	}

	input := dtos.ListCategoriesInput{
		UserID:   userID,
		IsActive: isActive,
	}

	// Execute use case
	output, err := h.listCategoriesUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Categories retrieved successfully",
		"data":    output,
	})
}

// Get handles category retrieval requests.
// @Summary Get category by ID
// @Description Retrieves a specific category by its ID for the authenticated user.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} map[string]interface{} "Category retrieved successfully"
// @Success 200 {object} dtos.GetCategoryOutput "Category data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Execute use case
	output, err := h.getCategoryUseCase.Execute(categoryID)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category retrieved successfully",
		"data":    output,
	})
}

// Update handles category update requests.
// @Summary Update category
// @Description Updates a category. Only provided fields will be updated.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID (UUID)"
// @Param request body dtos.UpdateCategoryInput true "Category update data (name, description - at least one required)"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Success 200 {object} dtos.UpdateCategoryOutput "Updated category data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Parse request body
	var input dtos.UpdateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Set category ID from URL parameter
	input.CategoryID = categoryID

	// Execute use case
	output, err := h.updateCategoryUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
		"data":    output,
	})
}

// Delete handles category deletion requests.
// @Summary Delete category
// @Description Deletes a category (soft delete).
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Success 200 {object} dtos.DeleteCategoryOutput "Deletion confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.DeleteCategoryInput{
		CategoryID: categoryID,
	}

	// Execute use case
	output, err := h.deleteCategoryUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *CategoryHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Category operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Category operation failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// Restore handles category restoration requests.
// @Summary Restore a soft-deleted category
// @Description Restores a soft-deleted category by setting deleted_at to NULL.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} map[string]interface{} "Category restored successfully"
// @Success 200 {object} dtos.RestoreCategoryOutput "Restoration confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid category ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist or is not deleted"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id}/restore [post]
func (h *CategoryHandler) Restore(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.RestoreCategoryInput{
		CategoryID: categoryID,
	}

	// Execute use case
	output, err := h.restoreCategoryUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}

// PermanentDelete handles permanent category deletion requests (admin only).
// @Summary Permanently delete a category
// @Description Permanently deletes a category from the database (hard delete). This action cannot be undone.
// @Tags categories
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} map[string]interface{} "Category permanently deleted successfully"
// @Success 200 {object} dtos.PermanentDeleteCategoryOutput "Deletion confirmation"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid category ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized - missing or invalid JWT token"
// @Failure 403 {object} map[string]interface{} "Forbidden - admin access required"
// @Failure 404 {object} map[string]interface{} "Not found - category does not exist"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id}/permanent [delete]
func (h *CategoryHandler) PermanentDelete(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.PermanentDeleteCategoryInput{
		CategoryID: categoryID,
	}

	// Execute use case
	output, err := h.permanentDeleteCategoryUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": output.Message,
		"data":    output,
	})
}
