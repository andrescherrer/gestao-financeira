package valueobjects

import (
	"errors"
	"strings"
)

const (
	// MinDescriptionLength is the minimum length for a transaction description.
	MinDescriptionLength = 3
	// MaxDescriptionLength is the maximum length for a transaction description.
	MaxDescriptionLength = 500
)

// TransactionDescription represents a transaction description value object.
type TransactionDescription struct {
	value string
}

// NewTransactionDescription creates a new TransactionDescription value object.
func NewTransactionDescription(description string) (TransactionDescription, error) {
	trimmed := strings.TrimSpace(description)

	if trimmed == "" {
		return TransactionDescription{}, errors.New("transaction description cannot be empty")
	}

	if len(trimmed) < MinDescriptionLength {
		return TransactionDescription{}, errors.New("transaction description must have at least 3 characters")
	}

	if len(trimmed) > MaxDescriptionLength {
		return TransactionDescription{}, errors.New("transaction description must have at most 500 characters")
	}

	return TransactionDescription{value: trimmed}, nil
}

// MustTransactionDescription creates a new TransactionDescription and panics if invalid.
// Use this only when you are certain the description is valid (e.g., in tests).
func MustTransactionDescription(description string) TransactionDescription {
	td, err := NewTransactionDescription(description)
	if err != nil {
		panic(err)
	}
	return td
}

// Value returns the transaction description as a string.
func (td TransactionDescription) Value() string {
	return td.value
}

// String returns the transaction description as a string (implements fmt.Stringer).
func (td TransactionDescription) String() string {
	return td.value
}

// Equals checks if two TransactionDescription values are equal.
func (td TransactionDescription) Equals(other TransactionDescription) bool {
	return td.value == other.value
}

// IsEmpty checks if the transaction description is empty.
func (td TransactionDescription) IsEmpty() bool {
	return td.value == ""
}
