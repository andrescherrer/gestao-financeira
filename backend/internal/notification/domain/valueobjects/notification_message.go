package valueobjects

import (
	"errors"
	"strings"
)

const (
	// MinMessageLength is the minimum length for a notification message.
	MinMessageLength = 1
	// MaxMessageLength is the maximum length for a notification message.
	MaxMessageLength = 1000
)

// NotificationMessage represents a notification message value object.
type NotificationMessage struct {
	value string
}

// NewNotificationMessage creates a new NotificationMessage from a string.
func NewNotificationMessage(message string) (NotificationMessage, error) {
	trimmed := strings.TrimSpace(message)

	if trimmed == "" {
		return NotificationMessage{}, errors.New("notification message cannot be empty")
	}

	if len(trimmed) < MinMessageLength {
		return NotificationMessage{}, errors.New("notification message is too short")
	}

	if len(trimmed) > MaxMessageLength {
		return NotificationMessage{}, errors.New("notification message is too long")
	}

	return NotificationMessage{value: trimmed}, nil
}

// MustNotificationMessage creates a new NotificationMessage and panics if invalid.
// Use this only when you are certain the message is valid (e.g., in tests).
func MustNotificationMessage(message string) NotificationMessage {
	nm, err := NewNotificationMessage(message)
	if err != nil {
		panic(err)
	}
	return nm
}

// Value returns the notification message as a string.
func (nm NotificationMessage) Value() string {
	return nm.value
}

// String returns the notification message as a string (implements fmt.Stringer).
func (nm NotificationMessage) String() string {
	return nm.value
}

// Equals checks if two NotificationMessage values are equal.
func (nm NotificationMessage) Equals(other NotificationMessage) bool {
	return nm.value == other.value
}

// IsEmpty checks if the notification message is empty.
func (nm NotificationMessage) IsEmpty() bool {
	return nm.value == ""
}

