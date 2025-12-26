package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// TransactionCreated represents a domain event that occurs when a transaction is created.
type TransactionCreated struct {
	events.BaseDomainEvent
	accountID       string
	transactionType string // INCOME or EXPENSE
	amount          int64  // Amount in cents
	currency        string
}

// NewTransactionCreated creates a new TransactionCreated event.
func NewTransactionCreated(
	transactionID string,
	accountID string,
	transactionType string,
	amount sharedvalueobjects.Money,
) *TransactionCreated {
	baseEvent := events.NewBaseDomainEvent(
		"TransactionCreated",
		transactionID,
		"Transaction",
	)

	return &TransactionCreated{
		BaseDomainEvent: baseEvent,
		accountID:       accountID,
		transactionType: transactionType,
		amount:          amount.Amount(),
		currency:        amount.Currency().Code(),
	}
}

// AccountID returns the account ID associated with this transaction.
func (e *TransactionCreated) AccountID() string {
	return e.accountID
}

// TransactionType returns the type of transaction (INCOME or EXPENSE).
func (e *TransactionCreated) TransactionType() string {
	return e.transactionType
}

// Amount returns the transaction amount in cents.
func (e *TransactionCreated) Amount() int64 {
	return e.amount
}

// Currency returns the currency code.
func (e *TransactionCreated) Currency() string {
	return e.currency
}
