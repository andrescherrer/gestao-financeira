package usecases

import (
	"fmt"
	"math"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
)

// ListNotificationsUseCase handles listing notifications.
type ListNotificationsUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewListNotificationsUseCase creates a new ListNotificationsUseCase instance.
func NewListNotificationsUseCase(
	notificationRepository repositories.NotificationRepository,
) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs the listing of notifications.
func (uc *ListNotificationsUseCase) Execute(input dtos.ListNotificationsInput) (*dtos.ListNotificationsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Set default pagination values
	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 20 // Default page size
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Find notifications with filters and pagination
	notifications, total, err := uc.notificationRepository.FindByUserIDWithFiltersWithPagination(
		userID,
		input.Status,
		input.Type,
		offset,
		pageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	// Convert to output items
	items := make([]dtos.NotificationItem, 0, len(notifications))
	for _, notification := range notifications {
		var readAt *string
		if notification.ReadAt() != nil {
			readAtStr := notification.ReadAt().Format("2006-01-02T15:04:05Z07:00")
			readAt = &readAtStr
		}

		items = append(items, dtos.NotificationItem{
			NotificationID: notification.ID().Value(),
			Title:          notification.Title().Value(),
			Message:        notification.Message().Value(),
			Type:           notification.Type().Value(),
			Status:         notification.Status().Value(),
			ReadAt:         readAt,
			Metadata:       notification.Metadata(),
			CreatedAt:      notification.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Build output
	output := &dtos.ListNotificationsOutput{
		Notifications: items,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}

	return output, nil
}
