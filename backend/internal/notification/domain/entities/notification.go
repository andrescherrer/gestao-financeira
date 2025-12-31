package entities

import (
	"errors"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
)

// Notification represents a notification aggregate root in the Notification context.
type Notification struct {
	id        notificationvalueobjects.NotificationID
	userID    identityvalueobjects.UserID
	title     notificationvalueobjects.NotificationTitle
	message   notificationvalueobjects.NotificationMessage
	notifType notificationvalueobjects.NotificationType
	status    notificationvalueobjects.NotificationStatus
	readAt    *time.Time
	metadata  map[string]interface{} // Additional data (e.g., related entity ID, action URL)
	createdAt time.Time
	updatedAt time.Time

	// Domain events
	events []events.DomainEvent
}

// NewNotification creates a new Notification aggregate.
func NewNotification(
	userID identityvalueobjects.UserID,
	title notificationvalueobjects.NotificationTitle,
	message notificationvalueobjects.NotificationMessage,
	notifType notificationvalueobjects.NotificationType,
	metadata map[string]interface{},
) (*Notification, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if title.IsEmpty() {
		return nil, errors.New("notification title cannot be empty")
	}

	if message.IsEmpty() {
		return nil, errors.New("notification message cannot be empty")
	}

	now := time.Now()

	notification := &Notification{
		id:        notificationvalueobjects.GenerateNotificationID(),
		userID:    userID,
		title:     title,
		message:   message,
		notifType: notifType,
		status:    notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusUnread),
		readAt:    nil,
		metadata:  metadata,
		createdAt: now,
		updatedAt: now,
		events:    []events.DomainEvent{},
	}

	return notification, nil
}

// NotificationFromPersistence reconstructs a Notification aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func NotificationFromPersistence(
	id notificationvalueobjects.NotificationID,
	userID identityvalueobjects.UserID,
	title notificationvalueobjects.NotificationTitle,
	message notificationvalueobjects.NotificationMessage,
	notifType notificationvalueobjects.NotificationType,
	status notificationvalueobjects.NotificationStatus,
	readAt *time.Time,
	metadata map[string]interface{},
	createdAt time.Time,
	updatedAt time.Time,
) *Notification {
	return &Notification{
		id:        id,
		userID:    userID,
		title:     title,
		message:   message,
		notifType: notifType,
		status:    status,
		readAt:    readAt,
		metadata:  metadata,
		createdAt: createdAt,
		updatedAt: updatedAt,
		events:    []events.DomainEvent{},
	}
}

// ID returns the notification ID.
func (n *Notification) ID() notificationvalueobjects.NotificationID {
	return n.id
}

// UserID returns the user ID.
func (n *Notification) UserID() identityvalueobjects.UserID {
	return n.userID
}

// Title returns the notification title.
func (n *Notification) Title() notificationvalueobjects.NotificationTitle {
	return n.title
}

// Message returns the notification message.
func (n *Notification) Message() notificationvalueobjects.NotificationMessage {
	return n.message
}

// Type returns the notification type.
func (n *Notification) Type() notificationvalueobjects.NotificationType {
	return n.notifType
}

// Status returns the notification status.
func (n *Notification) Status() notificationvalueobjects.NotificationStatus {
	return n.status
}

// ReadAt returns the read timestamp (nil if not read).
func (n *Notification) ReadAt() *time.Time {
	return n.readAt
}

// Metadata returns the notification metadata.
func (n *Notification) Metadata() map[string]interface{} {
	return n.metadata
}

// CreatedAt returns the creation timestamp.
func (n *Notification) CreatedAt() time.Time {
	return n.createdAt
}

// UpdatedAt returns the last update timestamp.
func (n *Notification) UpdatedAt() time.Time {
	return n.updatedAt
}

// MarkAsRead marks the notification as read.
func (n *Notification) MarkAsRead() error {
	if n.status.IsRead() {
		return errors.New("notification is already read")
	}

	if n.status.IsArchived() {
		return errors.New("cannot mark archived notification as read")
	}

	now := time.Now()
	n.status = notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusRead)
	n.readAt = &now
	n.updatedAt = now

	return nil
}

// MarkAsUnread marks the notification as unread.
func (n *Notification) MarkAsUnread() error {
	if n.status.IsUnread() {
		return errors.New("notification is already unread")
	}

	if n.status.IsArchived() {
		return errors.New("cannot mark archived notification as unread")
	}

	n.status = notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusUnread)
	n.readAt = nil
	n.updatedAt = time.Now()

	return nil
}

// Archive archives the notification.
func (n *Notification) Archive() error {
	if n.status.IsArchived() {
		return errors.New("notification is already archived")
	}

	n.status = notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusArchived)
	n.updatedAt = time.Now()

	return nil
}

// Unarchive unarchives the notification.
func (n *Notification) Unarchive() error {
	if !n.status.IsArchived() {
		return errors.New("notification is not archived")
	}

	// If it was read before, restore to read status, otherwise unread
	if n.readAt != nil {
		n.status = notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusRead)
	} else {
		n.status = notificationvalueobjects.MustNotificationStatus(notificationvalueobjects.StatusUnread)
	}
	n.updatedAt = time.Now()

	return nil
}

// UpdateMetadata updates the notification metadata.
func (n *Notification) UpdateMetadata(metadata map[string]interface{}) {
	n.metadata = metadata
	n.updatedAt = time.Now()
}

// GetEvents returns all domain events that occurred in this aggregate.
func (n *Notification) GetEvents() []events.DomainEvent {
	return n.events
}

// ClearEvents clears all domain events from this aggregate.
func (n *Notification) ClearEvents() {
	n.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (n *Notification) addEvent(event events.DomainEvent) {
	n.events = append(n.events, event)
}
