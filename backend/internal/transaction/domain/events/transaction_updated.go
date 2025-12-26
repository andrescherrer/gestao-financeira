package events

import (
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// TransactionUpdated represents a domain event that occurs when a transaction is updated.
// This event includes the old values to allow handlers to reverse previous balance updates.
type TransactionUpdated struct {
	events.BaseDomainEvent
	accountID       string
	transactionType string // INCOME or EXPENSE
	oldAmount       int64  // Old amount in cents
	oldType         string // Old transaction type
	newAmount       int64  // New amount in cents
	newType         string // New transaction type
	currency        string
}

// NewTransactionUpdated creates a new TransactionUpdated event.
func NewTransactionUpdated(
	transactionID string,
	accountID string,
	oldType string,
	oldAmount sharedvalueobjects.Money,
	newType string,
	newAmount sharedvalueobjects.Money,
) *TransactionUpdated {
	baseEvent := events.NewBaseDomainEvent(
		"TransactionUpdated",
		transactionID,
		"Transaction",
	)

	return &TransactionUpdated{
		BaseDomainEvent: baseEvent,
		accountID:       accountID,
		oldType:         oldType,
		oldAmount:       oldAmount.Amount(),
		newType:         newType,
		newAmount:       newAmount.Amount(),
		currency:        newAmount.Currency().Code(),
		transactionType: newType, // For convenience, use new type
	}
}

// AccountID returns the account ID associated with this transaction.
func (e *TransactionUpdated) AccountID() string {
	return e.accountID
}

// TransactionType returns the new type of transaction (INCOME or EXPENSE).
func (e *TransactionUpdated) TransactionType() string {
	return e.transactionType
}

// Amount returns the new transaction amount in cents.
func (e *TransactionUpdated) Amount() int64 {
	return e.newAmount
}

// Currency returns the currency code.
func (e *TransactionUpdated) Currency() string {
	return e.currency
}

// OldAmount returns the old transaction amount in cents.
func (e *TransactionUpdated) OldAmount() int64 {
	return e.oldAmount
}

// OldType returns the old transaction type.
func (e *TransactionUpdated) OldType() string {
	return e.oldType
}

// NewAmount returns the new transaction amount in cents.
func (e *TransactionUpdated) NewAmount() int64 {
	return e.newAmount
}

// NewType returns the new transaction type.
func (e *TransactionUpdated) NewType() string {
	return e.newType
}
