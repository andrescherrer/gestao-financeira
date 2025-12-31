package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// GoalCompleted represents a domain event when a goal is completed.
type GoalCompleted struct {
	events.BaseDomainEvent
	currentAmount string
	targetAmount  string
	currency      string
}

// NewGoalCompleted creates a new GoalCompleted event.
func NewGoalCompleted(
	goalID string,
	currentAmount string,
	targetAmount string,
	currency string,
) *GoalCompleted {
	baseEvent := events.NewBaseDomainEvent(
		"GoalCompleted",
		goalID,
		"Goal",
	)

	return &GoalCompleted{
		BaseDomainEvent: baseEvent,
		currentAmount:   currentAmount,
		targetAmount:    targetAmount,
		currency:        currency,
	}
}

// CurrentAmount returns the current amount as a string.
func (e *GoalCompleted) CurrentAmount() string {
	return e.currentAmount
}

// TargetAmount returns the target amount as a string.
func (e *GoalCompleted) TargetAmount() string {
	return e.targetAmount
}

// Currency returns the currency code.
func (e *GoalCompleted) Currency() string {
	return e.currency
}
