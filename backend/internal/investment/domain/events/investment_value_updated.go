package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// InvestmentValueUpdated represents a domain event that occurs when an investment value is updated.
type InvestmentValueUpdated struct {
	events.BaseDomainEvent
	currentValue string
	currency     string
}

// NewInvestmentValueUpdated creates a new InvestmentValueUpdated event.
func NewInvestmentValueUpdated(
	investmentID string,
	currentValue string,
	currency string,
) *InvestmentValueUpdated {
	baseEvent := events.NewBaseDomainEvent(
		"InvestmentValueUpdated",
		investmentID,
		"Investment",
	)

	return &InvestmentValueUpdated{
		BaseDomainEvent: baseEvent,
		currentValue:    currentValue,
		currency:        currency,
	}
}

// CurrentValue returns the new current value as a string.
func (e *InvestmentValueUpdated) CurrentValue() string {
	return e.currentValue
}

// Currency returns the currency code.
func (e *InvestmentValueUpdated) Currency() string {
	return e.currency
}

