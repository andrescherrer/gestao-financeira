package repositories

import (
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/notification/domain/entities"
	notificationvalueobjects "gestao-financeira/backend/internal/notification/domain/valueobjects"
)

// NotificationRepository defines the interface for notification persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type NotificationRepository interface {
	// FindByID finds a notification by its ID.
	// Returns nil if the notification is not found.
	FindByID(id notificationvalueobjects.NotificationID) (*entities.Notification, error)

	// FindByUserID finds all notifications for a given user.
	// Returns an empty slice if no notifications are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Notification, error)

	// FindByUserIDWithFiltersWithPagination finds notifications for a given user with filters and pagination.
	FindByUserIDWithFiltersWithPagination(
		userID identityvalueobjects.UserID,
		status string,
		notifType string,
		offset, limit int,
	) ([]*entities.Notification, int64, error)

	// FindUnreadByUserID finds all unread notifications for a given user.
	FindUnreadByUserID(userID identityvalueobjects.UserID) ([]*entities.Notification, error)

	// Save saves or updates a notification.
	// If the notification already exists (by ID), it updates it.
	// If the notification doesn't exist, it creates a new one.
	Save(notification *entities.Notification) error

	// Delete deletes a notification by its ID (soft delete).
	Delete(id notificationvalueobjects.NotificationID) error

	// Exists checks if a notification with the given ID already exists.
	// Returns true if the notification exists, false otherwise.
	Exists(id notificationvalueobjects.NotificationID) (bool, error)

	// Count returns the total number of notifications for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)

	// CountUnread returns the number of unread notifications for a given user.
	CountUnread(userID identityvalueobjects.UserID) (int64, error)
}
