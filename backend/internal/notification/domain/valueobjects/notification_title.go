package valueobjects

import (
	"errors"
	"strings"
)

const (
	// MinTitleLength is the minimum length for a notification title.
	MinTitleLength = 1
	// MaxTitleLength is the maximum length for a notification title.
	MaxTitleLength = 200
)

// NotificationTitle represents a notification title value object.
type NotificationTitle struct {
	value string
}

// NewNotificationTitle creates a new NotificationTitle from a string.
func NewNotificationTitle(title string) (NotificationTitle, error) {
	trimmed := strings.TrimSpace(title)

	if trimmed == "" {
		return NotificationTitle{}, errors.New("notification title cannot be empty")
	}

	if len(trimmed) < MinTitleLength {
		return NotificationTitle{}, errors.New("notification title is too short")
	}

	if len(trimmed) > MaxTitleLength {
		return NotificationTitle{}, errors.New("notification title is too long")
	}

	return NotificationTitle{value: trimmed}, nil
}

// MustNotificationTitle creates a new NotificationTitle and panics if invalid.
// Use this only when you are certain the title is valid (e.g., in tests).
func MustNotificationTitle(title string) NotificationTitle {
	nt, err := NewNotificationTitle(title)
	if err != nil {
		panic(err)
	}
	return nt
}

// Value returns the notification title as a string.
func (nt NotificationTitle) Value() string {
	return nt.value
}

// String returns the notification title as a string (implements fmt.Stringer).
func (nt NotificationTitle) String() string {
	return nt.value
}

// Equals checks if two NotificationTitle values are equal.
func (nt NotificationTitle) Equals(other NotificationTitle) bool {
	return nt.value == other.value
}

// IsEmpty checks if the notification title is empty.
func (nt NotificationTitle) IsEmpty() bool {
	return nt.value == ""
}
