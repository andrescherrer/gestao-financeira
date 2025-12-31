package valueobjects

import (
	"errors"
)

// GoalStatus represents the status of a goal.
type GoalStatus struct {
	value string
}

const (
	// StatusInProgress represents a goal that is in progress.
	StatusInProgress = "IN_PROGRESS"
	// StatusCompleted represents a goal that has been completed.
	StatusCompleted = "COMPLETED"
	// StatusOverdue represents a goal that is overdue.
	StatusOverdue = "OVERDUE"
	// StatusCancelled represents a goal that has been cancelled.
	StatusCancelled = "CANCELLED"
)

// Valid statuses
var validStatuses = []string{
	StatusInProgress,
	StatusCompleted,
	StatusOverdue,
	StatusCancelled,
}

// NewGoalStatus creates a new GoalStatus.
func NewGoalStatus(status string) (GoalStatus, error) {
	if status == "" {
		return GoalStatus{}, errors.New("goal status cannot be empty")
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return GoalStatus{value: status}, nil
		}
	}

	return GoalStatus{}, errors.New("invalid goal status")
}

// MustGoalStatus creates a new GoalStatus and panics if invalid.
// Use this only when you are certain the status is valid (e.g., in tests).
func MustGoalStatus(status string) GoalStatus {
	gs, err := NewGoalStatus(status)
	if err != nil {
		panic(err)
	}
	return gs
}

// Value returns the goal status as a string.
func (gs GoalStatus) Value() string {
	return gs.value
}

// String returns the goal status as a string (implements fmt.Stringer).
func (gs GoalStatus) String() string {
	return gs.value
}

// Equals checks if two GoalStatus values are equal.
func (gs GoalStatus) Equals(other GoalStatus) bool {
	return gs.value == other.value
}

// IsInProgress checks if the goal status is IN_PROGRESS.
func (gs GoalStatus) IsInProgress() bool {
	return gs.value == StatusInProgress
}

// IsCompleted checks if the goal status is COMPLETED.
func (gs GoalStatus) IsCompleted() bool {
	return gs.value == StatusCompleted
}

// IsOverdue checks if the goal status is OVERDUE.
func (gs GoalStatus) IsOverdue() bool {
	return gs.value == StatusOverdue
}

// IsCancelled checks if the goal status is CANCELLED.
func (gs GoalStatus) IsCancelled() bool {
	return gs.value == StatusCancelled
}

// CanBeCancelled checks if the goal can be cancelled.
func (gs GoalStatus) CanBeCancelled() bool {
	return gs.value == StatusInProgress || gs.value == StatusOverdue
}

