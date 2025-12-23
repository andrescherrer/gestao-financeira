package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// AccountID represents an account identifier value object.
type AccountID struct {
	value string
}

// NewAccountID creates a new AccountID from a string.
func NewAccountID(id string) (AccountID, error) {
	if id == "" {
		return AccountID{}, errors.New("account ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return AccountID{}, errors.New("invalid account ID format (must be UUID)")
	}

	return AccountID{value: id}, nil
}

// GenerateAccountID generates a new AccountID.
func GenerateAccountID() AccountID {
	return AccountID{value: uuid.New().String()}
}

// MustAccountID creates a new AccountID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustAccountID(id string) AccountID {
	aid, err := NewAccountID(id)
	if err != nil {
		panic(err)
	}
	return aid
}

// Value returns the account ID as a string.
func (aid AccountID) Value() string {
	return aid.value
}

// String returns the account ID as a string (implements fmt.Stringer).
func (aid AccountID) String() string {
	return aid.value
}

// Equals checks if two AccountID values are equal.
func (aid AccountID) Equals(other AccountID) bool {
	return aid.value == other.value
}

// IsEmpty checks if the account ID is empty.
func (aid AccountID) IsEmpty() bool {
	return aid.value == ""
}
