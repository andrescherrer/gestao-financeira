package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/pkg/errors"
)

// ArchiveNotificationUseCase handles archiving a notification.
type ArchiveNotificationUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewArchiveNotificationUseCase creates a new ArchiveNotificationUseCase instance.
func NewArchiveNotificationUseCase(
	notificationRepository repositories.NotificationRepository,
) *ArchiveNotificationUseCase {
	return &ArchiveNotificationUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs archiving a notification.
func (uc *ArchiveNotificationUseCase) Execute(input dtos.ArchiveNotificationInput) (*dtos.ArchiveNotificationOutput, error) {
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

	// Archive notification
	if err := notification.Archive(); err != nil {
		return nil, fmt.Errorf("failed to archive notification: %w", err)
	}

	// Save notification
	if err := uc.notificationRepository.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Build output
	output := &dtos.ArchiveNotificationOutput{
		NotificationID: notification.ID().Value(),
		Status:         notification.Status().Value(),
	}

	return output, nil
}
