package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// NotificationCreated represents a domain event when a notification is created.
type NotificationCreated struct {
	events.BaseDomainEvent
	userID    string
	title     string
	notifType string
}

// NewNotificationCreated creates a new NotificationCreated event.
func NewNotificationCreated(
	notificationID string,
	userID string,
	title string,
	notifType string,
) *NotificationCreated {
	baseEvent := events.NewBaseDomainEvent(
		"NotificationCreated",
		notificationID,
		"Notification",
	)

	return &NotificationCreated{
		BaseDomainEvent: baseEvent,
		userID:          userID,
		title:           title,
		notifType:       notifType,
	}
}

// UserID returns the user ID associated with this notification.
func (e *NotificationCreated) UserID() string {
	return e.userID
}

// Title returns the notification title.
func (e *NotificationCreated) Title() string {
	return e.title
}

// Type returns the notification type.
func (e *NotificationCreated) Type() string {
	return e.notifType
}
