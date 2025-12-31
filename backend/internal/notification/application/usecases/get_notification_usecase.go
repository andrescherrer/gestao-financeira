package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/pkg/errors"
)

// GetNotificationUseCase handles getting a single notification.
type GetNotificationUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewGetNotificationUseCase creates a new GetNotificationUseCase instance.
func NewGetNotificationUseCase(
	notificationRepository repositories.NotificationRepository,
) *GetNotificationUseCase {
	return &GetNotificationUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs getting a notification.
func (uc *GetNotificationUseCase) Execute(input dtos.GetNotificationInput) (*dtos.GetNotificationOutput, error) {
	// Create notification ID value object
	notificationID, err := notificationvalueobjects.NewNotificationID(input.NotificationID)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find notification
	notification, err := uc.notificationRepository.FindByID(notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}

	if notification == nil {
		return nil, errors.NewNotFoundError("notification", input.NotificationID)
	}

	// Verify ownership
	if !notification.UserID().Equals(userID) {
		return nil, errors.NewForbiddenError("notification does not belong to user")
	}

	// Build output
	var readAt *string
	if notification.ReadAt() != nil {
		readAtStr := notification.ReadAt().Format("2006-01-02T15:04:05Z07:00")
		readAt = &readAtStr
	}

	output := &dtos.GetNotificationOutput{
		NotificationID: notification.ID().Value(),
		UserID:         notification.UserID().Value(),
		Title:          notification.Title().Value(),
		Message:        notification.Message().Value(),
		Type:           notification.Type().Value(),
		Status:         notification.Status().Value(),
		ReadAt:         readAt,
		Metadata:       notification.Metadata(),
		CreatedAt:      notification.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      notification.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
