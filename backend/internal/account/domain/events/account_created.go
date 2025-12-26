package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
)

// AccountCreated represents a domain event that occurs when an account is created.
type AccountCreated struct {
	events.BaseDomainEvent
	userID      string
	name        string
	accountType string
	currency    string
	context     string
}

// NewAccountCreated creates a new AccountCreated event.
func NewAccountCreated(
	accountID string,
	userID string,
	name string,
	accountType string,
	currency string,
	context string,
) *AccountCreated {
	baseEvent := events.NewBaseDomainEvent(
		"AccountCreated",
		accountID,
		"Account",
	)

	return &AccountCreated{
		BaseDomainEvent: baseEvent,
		userID:          userID,
		name:            name,
		accountType:     accountType,
		currency:        currency,
		context:         context,
	}
}

// UserID returns the user ID who owns the account.
func (e *AccountCreated) UserID() string {
	return e.userID
}

// Name returns the account name.
func (e *AccountCreated) Name() string {
	return e.name
}

// AccountType returns the account type.
func (e *AccountCreated) AccountType() string {
	return e.accountType
}

// Currency returns the currency code.
func (e *AccountCreated) Currency() string {
	return e.currency
}

// Context returns the account context.
func (e *AccountCreated) Context() string {
	return e.context
}
