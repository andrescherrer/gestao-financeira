package valueobjects

import (
	"errors"
	"strings"
)

const (
	// StatusUnread represents an unread notification.
	StatusUnread = "UNREAD"
	// StatusRead represents a read notification.
	StatusRead = "READ"
	// StatusArchived represents an archived notification.
	StatusArchived = "ARCHIVED"
)

// NotificationStatus represents a notification status value object.
type NotificationStatus struct {
	value string
}

// NewNotificationStatus creates a new NotificationStatus from a string.
func NewNotificationStatus(status string) (NotificationStatus, error) {
	if status == "" {
		return NotificationStatus{}, errors.New("notification status cannot be empty")
	}

	normalized := strings.ToUpper(strings.TrimSpace(status))

	validStatuses := map[string]bool{
		StatusUnread:   true,
		StatusRead:     true,
		StatusArchived: true,
	}

	if !validStatuses[normalized] {
		return NotificationStatus{}, errors.New("invalid notification status: must be UNREAD, READ, or ARCHIVED")
	}

	return NotificationStatus{value: normalized}, nil
}

// MustNotificationStatus creates a new NotificationStatus and panics if invalid.
// Use this only when you are certain the status is valid (e.g., in tests).
func MustNotificationStatus(status string) NotificationStatus {
	ns, err := NewNotificationStatus(status)
	if err != nil {
		panic(err)
	}
	return ns
}

// Value returns the notification status as a string.
func (ns NotificationStatus) Value() string {
	return ns.value
}

// String returns the notification status as a string (implements fmt.Stringer).
func (ns NotificationStatus) String() string {
	return ns.value
}

// Equals checks if two NotificationStatus values are equal.
func (ns NotificationStatus) Equals(other NotificationStatus) bool {
	return ns.value == other.value
}

// IsUnread checks if the notification status is UNREAD.
func (ns NotificationStatus) IsUnread() bool {
	return ns.value == StatusUnread
}

// IsRead checks if the notification status is READ.
func (ns NotificationStatus) IsRead() bool {
	return ns.value == StatusRead
}

// IsArchived checks if the notification status is ARCHIVED.
func (ns NotificationStatus) IsArchived() bool {
	return ns.value == StatusArchived
}

