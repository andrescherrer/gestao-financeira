package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// GoalProgressUpdated represents a domain event when a goal's progress is updated.
type GoalProgressUpdated struct {
	events.BaseDomainEvent
	oldAmount    string
	newAmount    string
	targetAmount string
	progress     float64
	currency     string
}

// NewGoalProgressUpdated creates a new GoalProgressUpdated event.
func NewGoalProgressUpdated(
	goalID string,
	oldAmount string,
	newAmount string,
	targetAmount string,
	progress float64,
	currency string,
) *GoalProgressUpdated {
	baseEvent := events.NewBaseDomainEvent(
		"GoalProgressUpdated",
		goalID,
		"Goal",
	)

	return &GoalProgressUpdated{
		BaseDomainEvent: baseEvent,
		oldAmount:       oldAmount,
		newAmount:       newAmount,
		targetAmount:    targetAmount,
		progress:        progress,
		currency:        currency,
	}
}

// OldAmount returns the old amount as a string.
func (e *GoalProgressUpdated) OldAmount() string {
	return e.oldAmount
}

// NewAmount returns the new amount as a string.
func (e *GoalProgressUpdated) NewAmount() string {
	return e.newAmount
}

// TargetAmount returns the target amount as a string.
func (e *GoalProgressUpdated) TargetAmount() string {
	return e.targetAmount
}

// Progress returns the progress percentage.
func (e *GoalProgressUpdated) Progress() float64 {
	return e.progress
}

// Currency returns the currency code.
func (e *GoalProgressUpdated) Currency() string {
	return e.currency
}
