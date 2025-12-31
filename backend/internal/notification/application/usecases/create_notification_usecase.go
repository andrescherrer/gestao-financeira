package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/application/dtos"
	"gestao-financeira/backend/internal/notification/domain/entities"
	"gestao-financeira/backend/internal/notification/domain/repositories"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateNotificationUseCase handles notification creation.
type CreateNotificationUseCase struct {
	notificationRepository repositories.NotificationRepository
	eventBus               *eventbus.EventBus
}

// NewCreateNotificationUseCase creates a new CreateNotificationUseCase instance.
func NewCreateNotificationUseCase(
	notificationRepository repositories.NotificationRepository,
	eventBus *eventbus.EventBus,
) *CreateNotificationUseCase {
	return &CreateNotificationUseCase{
		notificationRepository: notificationRepository,
		eventBus:               eventBus,
	}
}

// Execute performs the notification creation.
// It validates the input, creates value objects, creates a new notification entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateNotificationUseCase) Execute(input dtos.CreateNotificationInput) (*dtos.CreateNotificationOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create notification title value object
	title, err := notificationvalueobjects.NewNotificationTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("invalid notification title: %w", err)
	}

	// Create notification message value object
	message, err := notificationvalueobjects.NewNotificationMessage(input.Message)
	if err != nil {
		return nil, fmt.Errorf("invalid notification message: %w", err)
	}

	// Create notification type value object
	notifType, err := notificationvalueobjects.NewNotificationType(input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid notification type: %w", err)
	}

	// Create notification entity
	notification, err := entities.NewNotification(
		userID,
		title,
		message,
		notifType,
		input.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Save notification
	if err := uc.notificationRepository.Save(notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Publish domain events
	domainEvents := notification.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the notification creation
			_ = err
		}
	}
	notification.ClearEvents()

	// Build output
	output := &dtos.CreateNotificationOutput{
		NotificationID: notification.ID().Value(),
		UserID:         notification.UserID().Value(),
		Title:          notification.Title().Value(),
		Message:        notification.Message().Value(),
		Type:           notification.Type().Value(),
		Status:         notification.Status().Value(),
		Metadata:       notification.Metadata(),
		CreatedAt:      notification.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
