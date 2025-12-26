package valueobjects

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// CategoryName represents a category name value object.
type CategoryName struct {
	value string
}

const (
	// MinCategoryNameLength is the minimum length for a category name.
	MinCategoryNameLength = 2
	// MaxCategoryNameLength is the maximum length for a category name.
	MaxCategoryNameLength = 100
)

// NewCategoryName creates a new CategoryName value object.
func NewCategoryName(name string) (CategoryName, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return CategoryName{}, errors.New("category name cannot be empty")
	}

	if len(name) < MinCategoryNameLength {
		return CategoryName{}, fmt.Errorf("category name must be at least %d characters long", MinCategoryNameLength)
	}

	if len(name) > MaxCategoryNameLength {
		return CategoryName{}, fmt.Errorf("category name is too long (max %d characters)", MaxCategoryNameLength)
	}

	// Check for valid characters (letters, numbers, spaces, hyphens, apostrophes)
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != ' ' && char != '-' && char != '\'' {
			return CategoryName{}, errors.New("category name contains invalid characters")
		}
	}

	return CategoryName{value: name}, nil
}

// MustCategoryName creates a new CategoryName and panics if invalid.
// Use this only when you are certain the name is valid (e.g., in tests).
func MustCategoryName(name string) CategoryName {
	cn, err := NewCategoryName(name)
	if err != nil {
		panic(err)
	}
	return cn
}

// Value returns the category name as a string.
func (cn CategoryName) Value() string {
	return cn.value
}

// String returns the category name as a string (implements fmt.Stringer).
func (cn CategoryName) String() string {
	return cn.value
}

// Equals checks if two CategoryName values are equal.
func (cn CategoryName) Equals(other CategoryName) bool {
	return strings.EqualFold(cn.value, other.value)
}

// IsEmpty checks if the category name is empty.
func (cn CategoryName) IsEmpty() bool {
	return cn.value == ""
}
