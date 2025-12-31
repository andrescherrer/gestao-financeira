package valueobjects

import (
	"errors"
	"strings"
)

const (
	// TypeInfo represents an informational notification.
	TypeInfo = "INFO"
	// TypeWarning represents a warning notification.
	TypeWarning = "WARNING"
	// TypeSuccess represents a success notification.
	TypeSuccess = "SUCCESS"
	// TypeError represents an error notification.
	TypeError = "ERROR"
)

// NotificationType represents a notification type value object.
type NotificationType struct {
	value string
}

// NewNotificationType creates a new NotificationType from a string.
func NewNotificationType(notificationType string) (NotificationType, error) {
	if notificationType == "" {
		return NotificationType{}, errors.New("notification type cannot be empty")
	}

	normalized := strings.ToUpper(strings.TrimSpace(notificationType))

	validTypes := map[string]bool{
		TypeInfo:    true,
		TypeWarning: true,
		TypeSuccess: true,
		TypeError:   true,
	}

	if !validTypes[normalized] {
		return NotificationType{}, errors.New("invalid notification type: must be INFO, WARNING, SUCCESS, or ERROR")
	}

	return NotificationType{value: normalized}, nil
}

// MustNotificationType creates a new NotificationType and panics if invalid.
// Use this only when you are certain the type is valid (e.g., in tests).
func MustNotificationType(notificationType string) NotificationType {
	nt, err := NewNotificationType(notificationType)
	if err != nil {
		panic(err)
	}
	return nt
}

// Value returns the notification type as a string.
func (nt NotificationType) Value() string {
	return nt.value
}

// String returns the notification type as a string (implements fmt.Stringer).
func (nt NotificationType) String() string {
	return nt.value
}

// Equals checks if two NotificationType values are equal.
func (nt NotificationType) Equals(other NotificationType) bool {
	return nt.value == other.value
}

// IsInfo checks if the notification type is INFO.
func (nt NotificationType) IsInfo() bool {
	return nt.value == TypeInfo
}

// IsWarning checks if the notification type is WARNING.
func (nt NotificationType) IsWarning() bool {
	return nt.value == TypeWarning
}

// IsSuccess checks if the notification type is SUCCESS.
func (nt NotificationType) IsSuccess() bool {
	return nt.value == TypeSuccess
}

// IsError checks if the notification type is ERROR.
func (nt NotificationType) IsError() bool {
	return nt.value == TypeError
}
