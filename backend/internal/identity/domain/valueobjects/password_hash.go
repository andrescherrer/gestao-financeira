package valueobjects

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHash represents a hashed password value object.
// It uses bcrypt for password hashing and verification.
type PasswordHash struct {
	hash string
}

const (
	// MinPasswordLength is the minimum length for a password.
	MinPasswordLength = 8
	// MaxPasswordLength is the maximum length for a password.
	MaxPasswordLength = 72 // bcrypt limit
	// DefaultCost is the default bcrypt cost factor.
	DefaultCost = 10
)

// NewPasswordHashFromPlain creates a new PasswordHash from a plain text password.
// It validates the password and hashes it using bcrypt.
func NewPasswordHashFromPlain(plainPassword string) (PasswordHash, error) {
	if err := validatePassword(plainPassword); err != nil {
		return PasswordHash{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), DefaultCost)
	if err != nil {
		return PasswordHash{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return PasswordHash{hash: string(hash)}, nil
}

// NewPasswordHashFromHash creates a new PasswordHash from an existing hash.
// Use this when loading a password hash from the database.
func NewPasswordHashFromHash(hash string) (PasswordHash, error) {
	if hash == "" {
		return PasswordHash{}, errors.New("password hash cannot be empty")
	}

	// Validate that it's a valid bcrypt hash format
	if len(hash) < 60 { // bcrypt hashes are always 60 characters
		return PasswordHash{}, errors.New("invalid bcrypt hash format")
	}

	return PasswordHash{hash: hash}, nil
}

// MustPasswordHashFromHash creates a new PasswordHash from an existing hash and panics if invalid.
// Use this only when you are certain the hash is valid (e.g., in tests).
func MustPasswordHashFromHash(hash string) PasswordHash {
	ph, err := NewPasswordHashFromHash(hash)
	if err != nil {
		panic(err)
	}
	return ph
}

// Verify checks if a plain text password matches the hash.
func (ph PasswordHash) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ph.hash), []byte(plainPassword))
	return err == nil
}

// Value returns the hash as a string.
func (ph PasswordHash) Value() string {
	return ph.hash
}

// String returns the hash as a string (implements fmt.Stringer).
// Note: In production, you might want to avoid logging the hash.
func (ph PasswordHash) String() string {
	return ph.hash
}

// Equals checks if two PasswordHash values are equal.
func (ph PasswordHash) Equals(other PasswordHash) bool {
	return ph.hash == other.hash
}

// IsEmpty checks if the hash is empty.
func (ph PasswordHash) IsEmpty() bool {
	return ph.hash == ""
}

// validatePassword validates a plain text password.
func validatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must be at most %d characters long", MaxPasswordLength)
	}

	return nil
}

// ValidatePasswordStrength performs additional password strength validation.
// This is a helper function that can be used before creating a PasswordHash.
// Currently validates only length; can be extended to check for uppercase, lowercase, digits, and special characters.
func ValidatePasswordStrength(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	// Optional: enforce stronger password requirements
	// For now, we'll just validate length and let the application layer decide on strength requirements
	// Future: can add checks for:
	// - hasUpper: at least one uppercase letter
	// - hasLower: at least one lowercase letter
	// - hasDigit: at least one digit
	// - hasSpecial: at least one special character

	return nil
}
