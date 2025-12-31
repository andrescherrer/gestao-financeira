package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// NotificationID represents a notification identifier value object.
type NotificationID struct {
	value string
}

// NewNotificationID creates a new NotificationID from a string.
func NewNotificationID(id string) (NotificationID, error) {
	if id == "" {
		return NotificationID{}, errors.New("notification ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return NotificationID{}, errors.New("invalid notification ID format (must be UUID)")
	}

	return NotificationID{value: id}, nil
}

// GenerateNotificationID generates a new NotificationID.
func GenerateNotificationID() NotificationID {
	return NotificationID{value: uuid.New().String()}
}

// MustNotificationID creates a new NotificationID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustNotificationID(id string) NotificationID {
	nid, err := NewNotificationID(id)
	if err != nil {
		panic(err)
	}
	return nid
}

// Value returns the notification ID as a string.
func (nid NotificationID) Value() string {
	return nid.value
}

// String returns the notification ID as a string (implements fmt.Stringer).
func (nid NotificationID) String() string {
	return nid.value
}

// Equals checks if two NotificationID values are equal.
func (nid NotificationID) Equals(other NotificationID) bool {
	return nid.value == other.value
}

// IsEmpty checks if the notification ID is empty.
func (nid NotificationID) IsEmpty() bool {
	return nid.value == ""
}

