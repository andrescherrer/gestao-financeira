package entities

import (
	"errors"
	"time"

	budgetevents "gestao-financeira/backend/internal/budget/domain/events"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// Budget represents a budget aggregate root in the Budget context.
type Budget struct {
	id         valueobjects.BudgetID
	userID     identityvalueobjects.UserID
	categoryID categoryvalueobjects.CategoryID
	amount     sharedvalueobjects.Money
	period     valueobjects.BudgetPeriod
	context    sharedvalueobjects.AccountContext
	createdAt  time.Time
	updatedAt  time.Time
	isActive   bool

	// Domain events
	events []events.DomainEvent
}

// NewBudget creates a new Budget aggregate.
func NewBudget(
	userID identityvalueobjects.UserID,
	categoryID categoryvalueobjects.CategoryID,
	amount sharedvalueobjects.Money,
	period valueobjects.BudgetPeriod,
	context sharedvalueobjects.AccountContext,
) (*Budget, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if categoryID.IsEmpty() {
		return nil, errors.New("category ID cannot be empty")
	}

	if amount.IsNegative() {
		return nil, errors.New("budget amount cannot be negative")
	}

	if amount.IsZero() {
		return nil, errors.New("budget amount must be greater than zero")
	}

	now := time.Now()

	budget := &Budget{
		id:         valueobjects.GenerateBudgetID(),
		userID:     userID,
		categoryID: categoryID,
		amount:     amount,
		period:     period,
		context:    context,
		createdAt:  now,
		updatedAt:  now,
		isActive:   true,
		events:     []events.DomainEvent{},
	}

	// Add domain event
	budget.addEvent(budgetevents.NewBudgetCreated(
		budget.id.Value(),
		budget.userID.Value(),
		budget.categoryID.Value(),
		budget.amount.Float64(),
		budget.amount.Currency().Code(),
		string(budget.period.PeriodType()),
		budget.period.Year(),
		budget.period.Month(),
		budget.context.Value(),
	))

	return budget, nil
}

// BudgetFromPersistence reconstructs a Budget aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func BudgetFromPersistence(
	id valueobjects.BudgetID,
	userID identityvalueobjects.UserID,
	categoryID categoryvalueobjects.CategoryID,
	amount sharedvalueobjects.Money,
	period valueobjects.BudgetPeriod,
	context sharedvalueobjects.AccountContext,
	createdAt time.Time,
	updatedAt time.Time,
	isActive bool,
) (*Budget, error) {
	if id.IsEmpty() {
		return nil, errors.New("budget ID cannot be empty")
	}

	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if categoryID.IsEmpty() {
		return nil, errors.New("category ID cannot be empty")
	}

	return &Budget{
		id:         id,
		userID:     userID,
		categoryID: categoryID,
		amount:     amount,
		period:     period,
		context:    context,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		isActive:   isActive,
		events:     []events.DomainEvent{},
	}, nil
}

// ID returns the budget ID.
func (b *Budget) ID() valueobjects.BudgetID {
	return b.id
}

// UserID returns the user ID.
func (b *Budget) UserID() identityvalueobjects.UserID {
	return b.userID
}

// CategoryID returns the category ID.
func (b *Budget) CategoryID() categoryvalueobjects.CategoryID {
	return b.categoryID
}

// Amount returns the budget amount.
func (b *Budget) Amount() sharedvalueobjects.Money {
	return b.amount
}

// Period returns the budget period.
func (b *Budget) Period() valueobjects.BudgetPeriod {
	return b.period
}

// Context returns the account context.
func (b *Budget) Context() sharedvalueobjects.AccountContext {
	return b.context
}

// CreatedAt returns the creation timestamp.
func (b *Budget) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns the last update timestamp.
func (b *Budget) UpdatedAt() time.Time {
	return b.updatedAt
}

// IsActive returns whether the budget is active.
func (b *Budget) IsActive() bool {
	return b.isActive
}

// UpdateAmount updates the budget amount.
func (b *Budget) UpdateAmount(amount sharedvalueobjects.Money) error {
	if amount.IsNegative() {
		return errors.New("budget amount cannot be negative")
	}

	if amount.IsZero() {
		return errors.New("budget amount must be greater than zero")
	}

	// Currency must match
	if !b.amount.Currency().Equals(amount.Currency()) {
		return errors.New("cannot change currency of existing budget")
	}

	b.amount = amount
	b.updatedAt = time.Now()

	b.addEvent(events.NewBaseDomainEvent(
		"BudgetUpdated",
		b.id.Value(),
		"Budget",
	))

	return nil
}

// UpdatePeriod updates the budget period.
func (b *Budget) UpdatePeriod(period valueobjects.BudgetPeriod) error {
	b.period = period
	b.updatedAt = time.Now()

	b.addEvent(events.NewBaseDomainEvent(
		"BudgetUpdated",
		b.id.Value(),
		"Budget",
	))

	return nil
}

// Deactivate deactivates the budget.
func (b *Budget) Deactivate() error {
	if !b.isActive {
		return errors.New("budget is already inactive")
	}

	b.isActive = false
	b.updatedAt = time.Now()

	b.addEvent(events.NewBaseDomainEvent(
		"BudgetDeactivated",
		b.id.Value(),
		"Budget",
	))

	return nil
}

// Activate activates the budget.
func (b *Budget) Activate() error {
	if b.isActive {
		return errors.New("budget is already active")
	}

	b.isActive = true
	b.updatedAt = time.Now()

	b.addEvent(events.NewBaseDomainEvent(
		"BudgetActivated",
		b.id.Value(),
		"Budget",
	))

	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (b *Budget) GetEvents() []events.DomainEvent {
	return b.events
}

// ClearEvents clears all domain events from this aggregate.
func (b *Budget) ClearEvents() {
	b.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (b *Budget) addEvent(event events.DomainEvent) {
	b.events = append(b.events, event)
}
