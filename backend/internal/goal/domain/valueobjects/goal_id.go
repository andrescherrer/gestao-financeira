package valueobjects

import (
	"errors"

	"github.com/google/uuid"
)

// GoalID represents a goal identifier value object.
type GoalID struct {
	value string
}

// NewGoalID creates a new GoalID from a string.
func NewGoalID(id string) (GoalID, error) {
	if id == "" {
		return GoalID{}, errors.New("goal ID cannot be empty")
	}

	// Validate UUID format
	_, err := uuid.Parse(id)
	if err != nil {
		return GoalID{}, errors.New("invalid goal ID format (must be UUID)")
	}

	return GoalID{value: id}, nil
}

// GenerateGoalID generates a new GoalID.
func GenerateGoalID() GoalID {
	return GoalID{value: uuid.New().String()}
}

// MustGoalID creates a new GoalID and panics if invalid.
// Use this only when you are certain the ID is valid (e.g., in tests).
func MustGoalID(id string) GoalID {
	gid, err := NewGoalID(id)
	if err != nil {
		panic(err)
	}
	return gid
}

// Value returns the goal ID as a string.
func (gid GoalID) Value() string {
	return gid.value
}

// String returns the goal ID as a string (implements fmt.Stringer).
func (gid GoalID) String() string {
	return gid.value
}

// Equals checks if two GoalID values are equal.
func (gid GoalID) Equals(other GoalID) bool {
	return gid.value == other.value
}

// IsEmpty checks if the goal ID is empty.
func (gid GoalID) IsEmpty() bool {
	return gid.value == ""
}
