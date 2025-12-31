package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// NotificationHandler handles notification-related HTTP requests.
type NotificationHandler struct {
	createNotificationUseCase  *usecases.CreateNotificationUseCase
	listNotificationsUseCase   *usecases.ListNotificationsUseCase
	getNotificationUseCase     *usecases.GetNotificationUseCase
	markReadUseCase            *usecases.MarkNotificationReadUseCase
	markUnreadUseCase          *usecases.MarkNotificationUnreadUseCase
	archiveNotificationUseCase *usecases.ArchiveNotificationUseCase
	deleteNotificationUseCase  *usecases.DeleteNotificationUseCase
}

// NewNotificationHandler creates a new NotificationHandler instance.
func NewNotificationHandler(
	createNotificationUseCase *usecases.CreateNotificationUseCase,
	listNotificationsUseCase *usecases.ListNotificationsUseCase,
	getNotificationUseCase *usecases.GetNotificationUseCase,
	markReadUseCase *usecases.MarkNotificationReadUseCase,
	markUnreadUseCase *usecases.MarkNotificationUnreadUseCase,
	archiveNotificationUseCase *usecases.ArchiveNotificationUseCase,
	deleteNotificationUseCase *usecases.DeleteNotificationUseCase,
) *NotificationHandler {
	return &NotificationHandler{
		createNotificationUseCase:  createNotificationUseCase,
		listNotificationsUseCase:   listNotificationsUseCase,
		getNotificationUseCase:     getNotificationUseCase,
		markReadUseCase:            markReadUseCase,
		markUnreadUseCase:          markUnreadUseCase,
		archiveNotificationUseCase: archiveNotificationUseCase,
		deleteNotificationUseCase:  deleteNotificationUseCase,
	}
}

// Create handles notification creation requests.
// @Summary Create a new notification
// @Description Creates a new notification for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dtos.CreateNotificationInput true "Notification creation data"
// @Success 201 {object} map[string]interface{} "Notification created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications [post]
func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	var input dtos.CreateNotificationInput
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

	output, err := h.createNotificationUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Notification created successfully",
		"data":    output,
	})
}

// List handles notification listing requests.
// @Summary List notifications
// @Description Lists all notifications for the authenticated user with optional filters.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by status (UNREAD, READ, ARCHIVED)"
// @Param type query string false "Filter by type (INFO, WARNING, SUCCESS, ERROR)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} map[string]interface{} "Notifications retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications [get]
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	var input dtos.ListNotificationsInput
	input.UserID = userID
	input.Status = c.Query("status")
	input.Type = c.Query("type")

	// Parse pagination
	if page := c.QueryInt("page"); page > 0 {
		input.Page = page
	}
	if pageSize := c.QueryInt("page_size"); pageSize > 0 {
		input.PageSize = pageSize
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.listNotificationsUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notifications retrieved successfully",
		"data":    output,
	})
}

// Get handles getting a single notification.
// @Summary Get notification
// @Description Gets a single notification by ID for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications/{id} [get]
func (h *NotificationHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.GetNotificationInput{
		NotificationID: notificationID,
		UserID:         userID,
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.getNotificationUseCase.Execute(input)
	if err != nil {
		return h.handleGetNotificationError(c, err, notificationID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification retrieved successfully",
		"data":    output,
	})
}

// MarkRead handles marking a notification as read.
// @Summary Mark notification as read
// @Description Marks a notification as read for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification marked as read"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.MarkNotificationReadInput{
		NotificationID: notificationID,
		UserID:         userID,
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.markReadUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification marked as read",
		"data":    output,
	})
}

// MarkUnread handles marking a notification as unread.
// @Summary Mark notification as unread
// @Description Marks a notification as unread for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification marked as unread"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications/{id}/unread [post]
func (h *NotificationHandler) MarkUnread(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.MarkNotificationUnreadInput{
		NotificationID: notificationID,
		UserID:         userID,
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.markUnreadUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification marked as unread",
		"data":    output,
	})
}

// Archive handles archiving a notification.
// @Summary Archive notification
// @Description Archives a notification for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification archived"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications/{id}/archive [post]
func (h *NotificationHandler) Archive(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.ArchiveNotificationInput{
		NotificationID: notificationID,
		UserID:         userID,
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.archiveNotificationUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification archived",
		"data":    output,
	})
}

// Delete handles notification deletion requests.
// @Summary Delete notification
// @Description Deletes a notification for the authenticated user.
// @Tags notifications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			"code":  fiber.StatusUnauthorized,
		})
	}

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	input := dtos.DeleteNotificationInput{
		NotificationID: notificationID,
		UserID:         userID,
	}

	if err := validator.Validate(&input); err != nil {
		return err
	}

	output, err := h.deleteNotificationUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification deleted successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
func (h *NotificationHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Notification operation failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Notification operation failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// handleGetNotificationError handles errors from GetNotificationUseCase and returns appropriate HTTP responses.
func (h *NotificationHandler) handleGetNotificationError(c *fiber.Ctx, err error, notificationID string) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Add notification ID to error details if available
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["notification_id"] = notificationID

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeNotFound {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Str("notification_id", notificationID).Msg("Get notification failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Str("notification_id", notificationID).Msg("Get notification failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}
