package events

import (
	"gestao-financeira/backend/internal/notification/application/dtos"
	notificationevents "gestao-financeira/backend/internal/notification/domain/events"
	"gestao-financeira/backend/internal/notification/infrastructure/websocket"
	"gestao-financeira/backend/internal/shared/domain/events"

	"github.com/rs/zerolog/log"
)

// NotificationEventHandler handles domain events related to notifications.
type NotificationEventHandler struct {
	notificationService *websocket.NotificationService
}

// NewNotificationEventHandler creates a new NotificationEventHandler instance.
func NewNotificationEventHandler(notificationService *websocket.NotificationService) *NotificationEventHandler {
	return &NotificationEventHandler{
		notificationService: notificationService,
	}
}

// HandleNotificationCreated handles NotificationCreated domain events.
func (h *NotificationEventHandler) HandleNotificationCreated(event events.DomainEvent) error {
	notifEvent, ok := event.(*notificationevents.NotificationCreated)
	if !ok {
		log.Warn().Str("event_type", event.EventType()).Msg("Invalid event type for NotificationCreated handler")
		return nil
	}

	// Create notification output DTO for WebSocket message
	notificationOutput := &dtos.CreateNotificationOutput{
		NotificationID: event.AggregateID(),
		UserID:         notifEvent.UserID(),
		Title:          notifEvent.Title(),
		Type:           notifEvent.Type(),
		Status:         "UNREAD",
		CreatedAt:      event.OccurredAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Send notification via WebSocket
	if err := h.notificationService.SendNotification(notifEvent.UserID(), notificationOutput); err != nil {
		log.Error().Err(err).
			Str("user_id", notifEvent.UserID()).
			Str("notification_id", event.AggregateID()).
			Msg("Failed to send notification via WebSocket")
		return err
	}

	log.Info().
		Str("user_id", notifEvent.UserID()).
		Str("notification_id", event.AggregateID()).
		Msg("Notification sent via WebSocket")

	return nil
}
