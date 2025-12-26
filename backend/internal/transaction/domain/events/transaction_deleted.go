package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// TransactionDeleted represents a domain event that occurs when a transaction is deleted.
// This event includes the transaction details to allow handlers to reverse balance updates.
type TransactionDeleted struct {
	events.BaseDomainEvent
	accountID       string
	transactionType string // INCOME or EXPENSE
	amount          int64  // Amount in cents
	currency        string
}

// NewTransactionDeleted creates a new TransactionDeleted event.
func NewTransactionDeleted(
	transactionID string,
	accountID string,
	transactionType string,
	amount sharedvalueobjects.Money,
) *TransactionDeleted {
	baseEvent := events.NewBaseDomainEvent(
		"TransactionDeleted",
		transactionID,
		"Transaction",
	)

	return &TransactionDeleted{
		BaseDomainEvent: baseEvent,
		accountID:       accountID,
		transactionType: transactionType,
		amount:          amount.Amount(),
		currency:        amount.Currency().Code(),
	}
}

// AccountID returns the account ID associated with this transaction.
func (e *TransactionDeleted) AccountID() string {
	return e.accountID
}

// TransactionType returns the type of transaction (INCOME or EXPENSE).
func (e *TransactionDeleted) TransactionType() string {
	return e.transactionType
}

// Amount returns the transaction amount in cents.
func (e *TransactionDeleted) Amount() int64 {
	return e.amount
}

// Currency returns the currency code.
func (e *TransactionDeleted) Currency() string {
	return e.currency
}
