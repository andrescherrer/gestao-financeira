package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

// AccountContext represents the context of an account (Personal or Business).
type AccountContext struct {
	value string
}

// Valid account context values
const (
	Personal = "PERSONAL"
	Business = "BUSINESS"
)

// ValidAccountContexts is a map of all supported account contexts.
var ValidAccountContexts = map[string]string{
	Personal: "Personal",
	Business: "Business",
}

// NewAccountContext creates a new AccountContext value object.
func NewAccountContext(value string) (AccountContext, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	if !IsValidAccountContext(value) {
		return AccountContext{}, fmt.Errorf("invalid account context: %s. Supported values: PERSONAL, BUSINESS", value)
	}

	return AccountContext{value: value}, nil
}

// MustAccountContext creates a new AccountContext value object and panics if the value is invalid.
// Use this only when you are certain the context value is valid.
func MustAccountContext(value string) AccountContext {
	context, err := NewAccountContext(value)
	if err != nil {
		panic(err)
	}
	return context
}

// IsValidAccountContext checks if an account context value is valid.
func IsValidAccountContext(value string) bool {
	value = strings.ToUpper(strings.TrimSpace(value))
	_, exists := ValidAccountContexts[value]
	return exists
}

// Value returns the account context value (PERSONAL or BUSINESS).
func (ac AccountContext) Value() string {
	return ac.value
}

// String returns the account context value as a string.
func (ac AccountContext) String() string {
	return ac.value
}

// DisplayName returns the human-readable name of the account context.
func (ac AccountContext) DisplayName() string {
	if name, exists := ValidAccountContexts[ac.value]; exists {
		return name
	}
	return ac.value
}

// IsPersonal checks if the account context is Personal.
func (ac AccountContext) IsPersonal() bool {
	return ac.value == Personal
}

// IsBusiness checks if the account context is Business.
func (ac AccountContext) IsBusiness() bool {
	return ac.value == Business
}

// Equals checks if two AccountContext values are equal.
func (ac AccountContext) Equals(other AccountContext) bool {
	return ac.value == other.value
}

// PersonalContext returns an AccountContext for Personal.
func PersonalContext() AccountContext {
	return AccountContext{value: Personal}
}

// BusinessContext returns an AccountContext for Business.
func BusinessContext() AccountContext {
	return AccountContext{value: Business}
}

// ParseAccountContext attempts to parse an account context from a string.
// It accepts both uppercase and lowercase values.
func ParseAccountContext(s string) (AccountContext, error) {
	if s == "" {
		return AccountContext{}, errors.New("account context cannot be empty")
	}

	return NewAccountContext(s)
}

// AllAccountContexts returns all valid account contexts.
func AllAccountContexts() []AccountContext {
	return []AccountContext{
		PersonalContext(),
		BusinessContext(),
	}
}

