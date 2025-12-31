package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/pkg/errors"
)

// MarkNotificationUnreadUseCase handles marking a notification as unread.
type MarkNotificationUnreadUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewMarkNotificationUnreadUseCase creates a new MarkNotificationUnreadUseCase instance.
func NewMarkNotificationUnreadUseCase(
	notificationRepository repositories.NotificationRepository,
) *MarkNotificationUnreadUseCase {
	return &MarkNotificationUnreadUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs marking a notification as unread.
func (uc *MarkNotificationUnreadUseCase) Execute(input dtos.MarkNotificationUnreadInput) (*dtos.MarkNotificationUnreadOutput, error) {
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

	// Mark as unread
	if err := notification.MarkAsUnread(); err != nil {
		return nil, fmt.Errorf("failed to mark notification as unread: %w", err)
	}

	// Save notification
	if err := uc.notificationRepository.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Build output
	output := &dtos.MarkNotificationUnreadOutput{
		NotificationID: notification.ID().Value(),
		Status:         notification.Status().Value(),
	}

	return output, nil
}
