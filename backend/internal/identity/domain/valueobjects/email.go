package valueobjects

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Email represents an email address value object.
// It validates the email format and normalizes it (lowercase).
type Email struct {
	value string
}

// emailRegex is a regular expression for validating email addresses.
// This is a simplified regex; for production, consider using a more robust validation library.
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// NewEmail creates a new Email value object with validation.
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	// Normalize: trim whitespace and convert to lowercase
	email = strings.ToLower(strings.TrimSpace(email))

	// Validate length
	if len(email) > 254 {
		return Email{}, errors.New("email is too long (max 254 characters)")
	}

	// Validate format
	if !emailRegex.MatchString(email) {
		return Email{}, fmt.Errorf("invalid email format: %s", email)
	}

	// Additional validation: check for consecutive dots
	if strings.Contains(email, "..") {
		return Email{}, errors.New("email cannot contain consecutive dots")
	}

	// Additional validation: check that @ appears only once
	atCount := strings.Count(email, "@")
	if atCount != 1 {
		return Email{}, errors.New("email must contain exactly one @ symbol")
	}

	// Split to validate local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return Email{}, errors.New("invalid email format")
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Validate local part
	if len(localPart) == 0 {
		return Email{}, errors.New("email local part cannot be empty")
	}
	if len(localPart) > 64 {
		return Email{}, errors.New("email local part is too long (max 64 characters)")
	}

	// Validate domain part
	if len(domainPart) == 0 {
		return Email{}, errors.New("email domain part cannot be empty")
	}
	if len(domainPart) > 253 {
		return Email{}, errors.New("email domain part is too long (max 253 characters)")
	}

	// Check that domain has at least one dot
	if !strings.Contains(domainPart, ".") {
		return Email{}, errors.New("email domain must contain at least one dot")
	}

	return Email{value: email}, nil
}

// MustEmail creates a new Email value object and panics if the email is invalid.
// Use this only when you are certain the email is valid (e.g., in tests).
func MustEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

// Value returns the email address as a string.
func (e Email) Value() string {
	return e.value
}

// String returns the email address as a string (implements fmt.Stringer).
func (e Email) String() string {
	return e.value
}

// Equals checks if two Email values are equal.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty checks if the email is empty (should not happen for valid emails).
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Domain returns the domain part of the email address.
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// LocalPart returns the local part (before @) of the email address.
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}
