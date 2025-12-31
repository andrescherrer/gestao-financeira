package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// InvestmentReturnCalculated represents a domain event that occurs when an investment return is calculated.
type InvestmentReturnCalculated struct {
	events.BaseDomainEvent
	absoluteReturn string
	percentage     float64
	currency       string
}

// NewInvestmentReturnCalculated creates a new InvestmentReturnCalculated event.
func NewInvestmentReturnCalculated(
	investmentID string,
	absoluteReturn string,
	percentage float64,
	currency string,
) *InvestmentReturnCalculated {
	baseEvent := events.NewBaseDomainEvent(
		"InvestmentReturnCalculated",
		investmentID,
		"Investment",
	)

	return &InvestmentReturnCalculated{
		BaseDomainEvent: baseEvent,
		absoluteReturn:  absoluteReturn,
		percentage:      percentage,
		currency:        currency,
	}
}

// AbsoluteReturn returns the absolute return as a string.
func (e *InvestmentReturnCalculated) AbsoluteReturn() string {
	return e.absoluteReturn
}

// Percentage returns the percentage return.
func (e *InvestmentReturnCalculated) Percentage() float64 {
	return e.percentage
}

// Currency returns the currency code.
func (e *InvestmentReturnCalculated) Currency() string {
	return e.currency
}
