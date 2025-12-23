package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// TransactionID represents a transaction identifier value object.
type TransactionID struct {
	value string
}

// NewTransactionID creates a new TransactionID from a string.
func NewTransactionID(id string) (TransactionID, error) {
	if id == "" {
		return TransactionID{}, errors.New("transaction ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return TransactionID{}, errors.New("invalid transaction ID format (must be UUID)")
	}

	return TransactionID{value: id}, nil
}

// GenerateTransactionID generates a new TransactionID.
func GenerateTransactionID() TransactionID {
	return TransactionID{value: uuid.New().String()}
}

// MustTransactionID creates a new TransactionID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustTransactionID(id string) TransactionID {
	tid, err := NewTransactionID(id)
	if err != nil {
		panic(err)
	}
	return tid
}

// Value returns the transaction ID as a string.
func (tid TransactionID) Value() string {
	return tid.value
}

// String returns the transaction ID as a string (implements fmt.Stringer).
func (tid TransactionID) String() string {
	return tid.value
}

// Equals checks if two TransactionID values are equal.
func (tid TransactionID) Equals(other TransactionID) bool {
	return tid.value == other.value
}

// IsEmpty checks if the transaction ID is empty.
func (tid TransactionID) IsEmpty() bool {
	return tid.value == ""
}
