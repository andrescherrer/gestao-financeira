package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// CategoryID represents a category identifier value object.
type CategoryID struct {
	value string
}

// NewCategoryID creates a new CategoryID from a string.
func NewCategoryID(id string) (CategoryID, error) {
	if id == "" {
		return CategoryID{}, errors.New("category ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return CategoryID{}, errors.New("invalid category ID format (must be UUID)")
	}

	return CategoryID{value: id}, nil
}

// GenerateCategoryID generates a new CategoryID.
func GenerateCategoryID() CategoryID {
	return CategoryID{value: uuid.New().String()}
}

// MustCategoryID creates a new CategoryID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustCategoryID(id string) CategoryID {
	cid, err := NewCategoryID(id)
	if err != nil {
		panic(err)
	}
	return cid
}

// Value returns the category ID as a string.
func (cid CategoryID) Value() string {
	return cid.value
}

// String returns the category ID as a string (implements fmt.Stringer).
func (cid CategoryID) String() string {
	return cid.value
}

// Equals checks if two CategoryID values are equal.
func (cid CategoryID) Equals(other CategoryID) bool {
	return cid.value == other.value
}

// IsEmpty checks if the category ID is empty.
func (cid CategoryID) IsEmpty() bool {
	return cid.value == ""
}
