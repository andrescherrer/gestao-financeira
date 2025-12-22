package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// UserID represents a user identifier value object.
type UserID struct {
	value string
}

// NewUserID creates a new UserID from a string.
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, errors.New("invalid user ID format (must be UUID)")
	}

	return UserID{value: id}, nil
}

// GenerateUserID generates a new UserID.
func GenerateUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// MustUserID creates a new UserID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustUserID(id string) UserID {
	uid, err := NewUserID(id)
	if err != nil {
		panic(err)
	}
	return uid
}

// Value returns the user ID as a string.
func (uid UserID) Value() string {
	return uid.value
}

// String returns the user ID as a string (implements fmt.Stringer).
func (uid UserID) String() string {
	return uid.value
}

// Equals checks if two UserID values are equal.
func (uid UserID) Equals(other UserID) bool {
	return uid.value == other.value
}

// IsEmpty checks if the user ID is empty.
func (uid UserID) IsEmpty() bool {
	return uid.value == ""
}
