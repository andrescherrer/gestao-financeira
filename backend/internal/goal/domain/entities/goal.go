package entities

import (
	"errors"
	"fmt"
	"time"

	goalevents "gestao-financeira/backend/internal/goal/domain/events"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// Goal represents a goal aggregate root in the Goal context.
type Goal struct {
	id            goalvalueobjects.GoalID
	userID        identityvalueobjects.UserID
	name          goalvalueobjects.GoalName
	targetAmount  sharedvalueobjects.Money
	currentAmount sharedvalueobjects.Money
	deadline      time.Time
	context       sharedvalueobjects.AccountContext
	status        goalvalueobjects.GoalStatus
	createdAt     time.Time
	updatedAt     time.Time

	// Domain events
	events []events.DomainEvent
}

// NewGoal creates a new Goal aggregate.
func NewGoal(
	userID identityvalueobjects.UserID,
	name goalvalueobjects.GoalName,
	targetAmount sharedvalueobjects.Money,
	deadline time.Time,
	context sharedvalueobjects.AccountContext,
) (*Goal, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if name.Name() == "" {
		return nil, errors.New("goal name cannot be empty")
	}

	if targetAmount.IsZero() {
		return nil, errors.New("target amount cannot be zero")
	}

	if targetAmount.IsNegative() {
		return nil, errors.New("target amount cannot be negative")
	}

	if deadline.IsZero() {
		return nil, errors.New("deadline cannot be zero")
	}

	if deadline.Before(time.Now()) {
		return nil, errors.New("deadline cannot be in the past")
	}

	now := time.Now()

	goal := &Goal{
		id:            goalvalueobjects.GenerateGoalID(),
		userID:        userID,
		name:          name,
		targetAmount:  targetAmount,
		currentAmount: sharedvalueobjects.Zero(targetAmount.Currency()), // Start with zero
		deadline:      deadline,
		context:       context,
		status:        goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusInProgress),
		createdAt:     now,
		updatedAt:     now,
		events:        []events.DomainEvent{},
	}

	// Add domain event
	goal.addEvent(goalevents.NewGoalCreated(
		goal.id.Value(),
		goal.userID.Value(),
		goal.name.Name(),
		fmt.Sprintf("%.2f", targetAmount.Float64()),
		targetAmount.Currency().Code(),
		deadline.Format(time.RFC3339),
		goal.context.Value(),
	))

	return goal, nil
}

// GoalFromPersistence reconstructs a Goal aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func GoalFromPersistence(
	id goalvalueobjects.GoalID,
	userID identityvalueobjects.UserID,
	name goalvalueobjects.GoalName,
	targetAmount sharedvalueobjects.Money,
	currentAmount sharedvalueobjects.Money,
	deadline time.Time,
	context sharedvalueobjects.AccountContext,
	status goalvalueobjects.GoalStatus,
	createdAt time.Time,
	updatedAt time.Time,
) (*Goal, error) {
	if id.IsEmpty() {
		return nil, errors.New("goal ID cannot be empty")
	}
	if id.IsEmpty() {
		return nil, errors.New("goal ID cannot be empty")
	}
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}
	if name.Name() == "" {
		return nil, errors.New("goal name cannot be empty")
	}

	return &Goal{
		id:            id,
		userID:        userID,
		name:          name,
		targetAmount:  targetAmount,
		currentAmount: currentAmount,
		deadline:      deadline,
		context:       context,
		status:        status,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		events:        []events.DomainEvent{},
	}, nil
}

// ID returns the goal ID.
func (g *Goal) ID() goalvalueobjects.GoalID {
	return g.id
}

// UserID returns the user ID.
func (g *Goal) UserID() identityvalueobjects.UserID {
	return g.userID
}

// Name returns the goal name.
func (g *Goal) Name() goalvalueobjects.GoalName {
	return g.name
}

// TargetAmount returns the target amount.
func (g *Goal) TargetAmount() sharedvalueobjects.Money {
	return g.targetAmount
}

// CurrentAmount returns the current amount.
func (g *Goal) CurrentAmount() sharedvalueobjects.Money {
	return g.currentAmount
}

// Deadline returns the deadline.
func (g *Goal) Deadline() time.Time {
	return g.deadline
}

// Context returns the account context.
func (g *Goal) Context() sharedvalueobjects.AccountContext {
	return g.context
}

// Status returns the goal status.
func (g *Goal) Status() goalvalueobjects.GoalStatus {
	return g.status
}

// CreatedAt returns the creation timestamp.
func (g *Goal) CreatedAt() time.Time {
	return g.createdAt
}

// UpdatedAt returns the last update timestamp.
func (g *Goal) UpdatedAt() time.Time {
	return g.updatedAt
}

// AddContribution adds a contribution to the goal.
func (g *Goal) AddContribution(amount sharedvalueobjects.Money) error {
	if amount.IsZero() {
		return errors.New("contribution amount cannot be zero")
	}

	if amount.IsNegative() {
		return errors.New("contribution amount cannot be negative")
	}

	// Ensure currencies match
	if !amount.Currency().Equals(g.targetAmount.Currency()) {
		return errors.New("contribution currency must match target amount currency")
	}

	oldAmount := g.currentAmount
	newAmount, err := g.currentAmount.Add(amount)
	if err != nil {
		return err
	}

	// Ensure current amount doesn't exceed target
	greater, err := newAmount.GreaterThan(g.targetAmount)
	if err != nil {
		return err
	}
	if greater {
		return errors.New("contribution would exceed target amount")
	}

	g.currentAmount = newAmount
	g.updatedAt = time.Now()

	// Add domain event
	g.addEvent(goalevents.NewGoalProgressUpdated(
		g.id.Value(),
		fmt.Sprintf("%.2f", oldAmount.Float64()),
		fmt.Sprintf("%.2f", newAmount.Float64()),
		fmt.Sprintf("%.2f", g.targetAmount.Float64()),
		g.CalculateProgress(),
		g.targetAmount.Currency().Code(),
	))

	// Check and update status
	g.checkAndUpdateStatus()

	return nil
}

// UpdateProgress updates the progress of the goal.
func (g *Goal) UpdateProgress(amount sharedvalueobjects.Money) error {
	if amount.IsNegative() {
		return errors.New("progress amount cannot be negative")
	}

	// Ensure currencies match
	if !amount.Currency().Equals(g.targetAmount.Currency()) {
		return errors.New("progress currency must match target amount currency")
	}

	// Ensure amount doesn't exceed target
	greater, err := amount.GreaterThan(g.targetAmount)
	if err != nil {
		return err
	}
	if greater {
		return errors.New("progress amount cannot exceed target amount")
	}

	oldAmount := g.currentAmount
	g.currentAmount = amount
	g.updatedAt = time.Now()

	// Add domain event
	g.addEvent(goalevents.NewGoalProgressUpdated(
		g.id.Value(),
		fmt.Sprintf("%.2f", oldAmount.Float64()),
		fmt.Sprintf("%.2f", amount.Float64()),
		fmt.Sprintf("%.2f", g.targetAmount.Float64()),
		g.CalculateProgress(),
		g.targetAmount.Currency().Code(),
	))

	// Check and update status
	g.checkAndUpdateStatus()

	return nil
}

// CheckStatus checks and returns the current status of the goal.
func (g *Goal) CheckStatus() goalvalueobjects.GoalStatus {
	g.checkAndUpdateStatus()
	return g.status
}

// checkAndUpdateStatus checks the goal status and updates it if necessary.
func (g *Goal) checkAndUpdateStatus() {
	now := time.Now()

	// If already completed or cancelled, don't change
	if g.status.IsCompleted() || g.status.IsCancelled() {
		return
	}

	// Check if goal is completed
	if g.isCompleted() {
		if !g.status.IsCompleted() {
			g.status = goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusCompleted)
			g.updatedAt = now

			// Add domain event
			g.addEvent(goalevents.NewGoalCompleted(
				g.id.Value(),
				fmt.Sprintf("%.2f", g.currentAmount.Float64()),
				fmt.Sprintf("%.2f", g.targetAmount.Float64()),
				g.targetAmount.Currency().Code(),
			))
		}
		return
	}

	// Check if goal is overdue
	if g.isOverdue() {
		if !g.status.IsOverdue() {
			g.status = goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusOverdue)
			g.updatedAt = now

			// Add domain event
			g.addEvent(goalevents.NewGoalOverdue(
				g.id.Value(),
				g.deadline.Format(time.RFC3339),
				fmt.Sprintf("%.2f", g.currentAmount.Float64()),
				fmt.Sprintf("%.2f", g.targetAmount.Float64()),
				g.targetAmount.Currency().Code(),
			))
		}
		return
	}

	// Otherwise, ensure it's IN_PROGRESS
	if !g.status.IsInProgress() {
		g.status = goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusInProgress)
		g.updatedAt = now
	}
}

// CalculateProgress calculates the progress percentage of the goal.
func (g *Goal) CalculateProgress() float64 {
	if g.targetAmount.IsZero() {
		return 0.0
	}

	progress := (float64(g.currentAmount.Amount()) / float64(g.targetAmount.Amount())) * 100.0
	if progress > 100.0 {
		return 100.0
	}
	return progress
}

// CalculateRemainingDays calculates the number of days remaining until the deadline.
func (g *Goal) CalculateRemainingDays() int {
	now := time.Now()
	if g.deadline.Before(now) {
		return 0
	}

	diff := g.deadline.Sub(now)
	return int(diff.Hours() / 24)
}

// IsCompleted checks if the goal is completed.
func (g *Goal) IsCompleted() bool {
	return g.isCompleted()
}

// isCompleted is the internal method that checks if the goal is completed.
func (g *Goal) isCompleted() bool {
	greaterOrEqual, err := g.currentAmount.GreaterThanOrEqual(g.targetAmount)
	if err != nil {
		return false
	}
	return greaterOrEqual
}

// IsOverdue checks if the goal is overdue.
func (g *Goal) IsOverdue() bool {
	return g.isOverdue()
}

// isOverdue is the internal method that checks if the goal is overdue.
func (g *Goal) isOverdue() bool {
	now := time.Now()
	return g.deadline.Before(now) && !g.isCompleted()
}

// Cancel cancels the goal.
func (g *Goal) Cancel() error {
	if !g.status.CanBeCancelled() {
		return errors.New("goal cannot be cancelled in its current status")
	}

	g.status = goalvalueobjects.MustGoalStatus(goalvalueobjects.StatusCancelled)
	g.updatedAt = time.Now()

	// Note: We could add a GoalCancelled event here if needed

	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (g *Goal) GetEvents() []events.DomainEvent {
	return g.events
}

// ClearEvents clears all domain events from this aggregate.
func (g *Goal) ClearEvents() {
	g.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (g *Goal) addEvent(event events.DomainEvent) {
	g.events = append(g.events, event)
}
