package valueobjects

import (
	"errors"
	"strings"
)

// GoalName represents a goal name value object.
type GoalName struct {
	name string
}

// NewGoalName creates a new GoalName.
func NewGoalName(name string) (GoalName, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return GoalName{}, errors.New("goal name cannot be empty")
	}

	if len(trimmed) < 3 {
		return GoalName{}, errors.New("goal name must have at least 3 characters")
	}

	if len(trimmed) > 200 {
		return GoalName{}, errors.New("goal name must have at most 200 characters")
	}

	return GoalName{name: trimmed}, nil
}

// MustGoalName creates a new GoalName and panics if invalid.
// Use this only when you are certain the name is valid (e.g., in tests).
func MustGoalName(name string) GoalName {
	gn, err := NewGoalName(name)
	if err != nil {
		panic(err)
	}
	return gn
}

// Name returns the goal name as a string.
func (gn GoalName) Name() string {
	return gn.name
}

// String returns the goal name as a string (implements fmt.Stringer).
func (gn GoalName) String() string {
	return gn.name
}

// Equals checks if two GoalName values are equal.
func (gn GoalName) Equals(other GoalName) bool {
	return gn.name == other.name
}

