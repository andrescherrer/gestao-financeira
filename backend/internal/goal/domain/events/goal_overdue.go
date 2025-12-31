package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// GoalOverdue represents a domain event when a goal becomes overdue.
type GoalOverdue struct {
	events.BaseDomainEvent
	deadline      string
	currentAmount string
	targetAmount  string
	currency      string
}

// NewGoalOverdue creates a new GoalOverdue event.
func NewGoalOverdue(
	goalID string,
	deadline string,
	currentAmount string,
	targetAmount string,
	currency string,
) *GoalOverdue {
	baseEvent := events.NewBaseDomainEvent(
		"GoalOverdue",
		goalID,
		"Goal",
	)

	return &GoalOverdue{
		BaseDomainEvent: baseEvent,
		deadline:        deadline,
		currentAmount:   currentAmount,
		targetAmount:    targetAmount,
		currency:        currency,
	}
}

// Deadline returns the deadline as a string.
func (e *GoalOverdue) Deadline() string {
	return e.deadline
}

// CurrentAmount returns the current amount as a string.
func (e *GoalOverdue) CurrentAmount() string {
	return e.currentAmount
}

// TargetAmount returns the target amount as a string.
func (e *GoalOverdue) TargetAmount() string {
	return e.targetAmount
}

// Currency returns the currency code.
func (e *GoalOverdue) Currency() string {
	return e.currency
}

