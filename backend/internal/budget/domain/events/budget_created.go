package events

import "gestao-financeira/backend/internal/shared/domain/events"

// BudgetCreated represents a domain event when a budget is created.
type BudgetCreated struct {
	events.BaseDomainEvent
	UserID     string
	CategoryID string
	Amount     float64
	Currency   string
	PeriodType string
	Year       int
	Month      *int
	Context    string
}

// NewBudgetCreated creates a new BudgetCreated event.
func NewBudgetCreated(
	budgetID string,
	userID string,
	categoryID string,
	amount float64,
	currency string,
	periodType string,
	year int,
	month *int,
	context string,
) *BudgetCreated {
	baseEvent := events.NewBaseDomainEvent(
		"BudgetCreated",
		budgetID,
		"Budget",
	)

	return &BudgetCreated{
		BaseDomainEvent: baseEvent,
		UserID:          userID,
		CategoryID:      categoryID,
		Amount:          amount,
		Currency:        currency,
		PeriodType:      periodType,
		Year:            year,
		Month:           month,
		Context:         context,
	}
}
