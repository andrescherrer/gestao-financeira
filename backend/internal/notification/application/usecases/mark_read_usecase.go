package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/pkg/errors"
)

// MarkNotificationReadUseCase handles marking a notification as read.
type MarkNotificationReadUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewMarkNotificationReadUseCase creates a new MarkNotificationReadUseCase instance.
func NewMarkNotificationReadUseCase(
	notificationRepository repositories.NotificationRepository,
) *MarkNotificationReadUseCase {
	return &MarkNotificationReadUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs marking a notification as read.
func (uc *MarkNotificationReadUseCase) Execute(input dtos.MarkNotificationReadInput) (*dtos.MarkNotificationReadOutput, error) {
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

	// Mark as read
	if err := notification.MarkAsRead(); err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	// Save notification
	if err := uc.notificationRepository.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Build output
	readAtStr := notification.ReadAt().Format("2006-01-02T15:04:05Z07:00")
	output := &dtos.MarkNotificationReadOutput{
		NotificationID: notification.ID().Value(),
		Status:         notification.Status().Value(),
		ReadAt:         readAtStr,
	}

	return output, nil
}
