package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// InvestmentID represents an investment identifier value object.
type InvestmentID struct {
	value string
}

// NewInvestmentID creates a new InvestmentID from a string.
func NewInvestmentID(id string) (InvestmentID, error) {
	if id == "" {
		return InvestmentID{}, errors.New("investment ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return InvestmentID{}, errors.New("invalid investment ID format (must be UUID)")
	}

	return InvestmentID{value: id}, nil
}

// GenerateInvestmentID generates a new InvestmentID.
func GenerateInvestmentID() InvestmentID {
	return InvestmentID{value: uuid.New().String()}
}

// MustInvestmentID creates a new InvestmentID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustInvestmentID(id string) InvestmentID {
	iid, err := NewInvestmentID(id)
	if err != nil {
		panic(err)
	}
	return iid
}

// Value returns the investment ID as a string.
func (iid InvestmentID) Value() string {
	return iid.value
}

// String returns the investment ID as a string (implements fmt.Stringer).
func (iid InvestmentID) String() string {
	return iid.value
}

// Equals checks if two InvestmentID values are equal.
func (iid InvestmentID) Equals(other InvestmentID) bool {
	return iid.value == other.value
}

// IsEmpty checks if the investment ID is empty.
func (iid InvestmentID) IsEmpty() bool {
	return iid.value == ""
}
