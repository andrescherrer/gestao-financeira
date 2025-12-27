package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// BudgetID represents a budget identifier value object.
type BudgetID struct {
	value string
}

// NewBudgetID creates a new BudgetID from a string.
func NewBudgetID(id string) (BudgetID, error) {
	if id == "" {
		return BudgetID{}, errors.New("budget ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return BudgetID{}, errors.New("invalid budget ID format (must be UUID)")
	}

	return BudgetID{value: id}, nil
}

// GenerateBudgetID generates a new BudgetID.
func GenerateBudgetID() BudgetID {
	return BudgetID{value: uuid.New().String()}
}

// MustBudgetID creates a new BudgetID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustBudgetID(id string) BudgetID {
	bid, err := NewBudgetID(id)
	if err != nil {
		panic(err)
	}
	return bid
}

// Value returns the budget ID as a string.
func (bid BudgetID) Value() string {
	return bid.value
}

// String returns the budget ID as a string (implements fmt.Stringer).
func (bid BudgetID) String() string {
	return bid.value
}

// Equals checks if two BudgetID values are equal.
func (bid BudgetID) Equals(other BudgetID) bool {
	return bid.value == other.value
}

// IsEmpty checks if the budget ID is empty.
func (bid BudgetID) IsEmpty() bool {
	return bid.value == ""
}
