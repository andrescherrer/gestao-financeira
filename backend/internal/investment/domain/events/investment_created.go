package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// InvestmentCreated represents a domain event that occurs when an investment is created.
type InvestmentCreated struct {
	events.BaseDomainEvent
	userID         string
	accountID      string
	investmentType string
	name           string
	ticker         string
	purchaseAmount string
	currency       string
	context        string
}

// NewInvestmentCreated creates a new InvestmentCreated event.
func NewInvestmentCreated(
	investmentID string,
	userID string,
	accountID string,
	investmentType string,
	name string,
	ticker string,
	purchaseAmount string,
	currency string,
	context string,
) *InvestmentCreated {
	baseEvent := events.NewBaseDomainEvent(
		"InvestmentCreated",
		investmentID,
		"Investment",
	)

	return &InvestmentCreated{
		BaseDomainEvent: baseEvent,
		userID:          userID,
		accountID:       accountID,
		investmentType:  investmentType,
		name:            name,
		ticker:          ticker,
		purchaseAmount:  purchaseAmount,
		currency:        currency,
		context:         context,
	}
}

// UserID returns the user ID who owns the investment.
func (e *InvestmentCreated) UserID() string {
	return e.userID
}

// AccountID returns the account ID associated with the investment.
func (e *InvestmentCreated) AccountID() string {
	return e.accountID
}

// InvestmentType returns the investment type.
func (e *InvestmentCreated) InvestmentType() string {
	return e.investmentType
}

// Name returns the investment name.
func (e *InvestmentCreated) Name() string {
	return e.name
}

// Ticker returns the investment ticker (if available).
func (e *InvestmentCreated) Ticker() string {
	return e.ticker
}

// PurchaseAmount returns the purchase amount as a string.
func (e *InvestmentCreated) PurchaseAmount() string {
	return e.purchaseAmount
}

// Currency returns the currency code.
func (e *InvestmentCreated) Currency() string {
	return e.currency
}

// Context returns the account context.
func (e *InvestmentCreated) Context() string {
	return e.context
}
