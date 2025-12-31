package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// GoalCreated represents a domain event when a goal is created.
type GoalCreated struct {
	events.BaseDomainEvent
	userID       string
	name         string
	targetAmount string
	currency     string
	deadline     string
	context      string
}

// NewGoalCreated creates a new GoalCreated event.
func NewGoalCreated(
	goalID string,
	userID string,
	name string,
	targetAmount string,
	currency string,
	deadline string,
	context string,
) *GoalCreated {
	baseEvent := events.NewBaseDomainEvent(
		"GoalCreated",
		goalID,
		"Goal",
	)

	return &GoalCreated{
		BaseDomainEvent: baseEvent,
		userID:          userID,
		name:            name,
		targetAmount:    targetAmount,
		currency:        currency,
		deadline:        deadline,
		context:         context,
	}
}

// UserID returns the user ID who owns the goal.
func (e *GoalCreated) UserID() string {
	return e.userID
}

// Name returns the goal name.
func (e *GoalCreated) Name() string {
	return e.name
}

// TargetAmount returns the target amount as a string.
func (e *GoalCreated) TargetAmount() string {
	return e.targetAmount
}

// Currency returns the currency code.
func (e *GoalCreated) Currency() string {
	return e.currency
}

// Deadline returns the deadline as a string.
func (e *GoalCreated) Deadline() string {
	return e.deadline
}

// Context returns the account context.
func (e *GoalCreated) Context() string {
	return e.context
}
