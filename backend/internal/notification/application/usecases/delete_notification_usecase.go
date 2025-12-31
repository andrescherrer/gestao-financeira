package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/pkg/errors"
)

// DeleteNotificationUseCase handles deleting a notification.
type DeleteNotificationUseCase struct {
	notificationRepository repositories.NotificationRepository
}

// NewDeleteNotificationUseCase creates a new DeleteNotificationUseCase instance.
func NewDeleteNotificationUseCase(
	notificationRepository repositories.NotificationRepository,
) *DeleteNotificationUseCase {
	return &DeleteNotificationUseCase{
		notificationRepository: notificationRepository,
	}
}

// Execute performs deleting a notification.
func (uc *DeleteNotificationUseCase) Execute(input dtos.DeleteNotificationInput) (*dtos.DeleteNotificationOutput, error) {
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

	// Find notification to verify ownership
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

	// Delete notification (soft delete)
	if err := uc.notificationRepository.Delete(notificationID); err != nil {
		return nil, fmt.Errorf("failed to delete notification: %w", err)
	}

	// Build output
	output := &dtos.DeleteNotificationOutput{
		NotificationID: notificationID.Value(),
		Deleted:        true,
	}

	return output, nil
}
