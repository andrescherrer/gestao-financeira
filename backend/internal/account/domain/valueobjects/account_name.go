package valueobjects

import (
	"errors"
	"strings"
)

// AccountName represents an account name value object.
type AccountName struct {
	value string
}

// NewAccountName creates a new AccountName value object.
func NewAccountName(name string) (AccountName, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return AccountName{}, errors.New("account name cannot be empty")
	}

	if len(trimmed) < 3 {
		return AccountName{}, errors.New("account name must have at least 3 characters")
	}

	if len(trimmed) > 100 {
		return AccountName{}, errors.New("account name must have at most 100 characters")
	}

	return AccountName{value: trimmed}, nil
}

// MustAccountName creates a new AccountName and panics if invalid.
// Use this only when you are certain the name is valid (e.g., in tests).
func MustAccountName(name string) AccountName {
	accountName, err := NewAccountName(name)
	if err != nil {
		panic(err)
	}
	return accountName
}

// Value returns the account name as a string.
func (an AccountName) Value() string {
	return an.value
}

// String returns the account name as a string (implements fmt.Stringer).
func (an AccountName) String() string {
	return an.value
}

// Equals checks if two AccountName values are equal.
func (an AccountName) Equals(other AccountName) bool {
	return an.value == other.value
}

// IsEmpty checks if the account name is empty.
func (an AccountName) IsEmpty() bool {
	return an.value == ""
}
