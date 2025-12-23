package entities

import (
	"errors"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// Transaction represents a transaction aggregate root in the Transaction context.
type Transaction struct {
	id              transactionvalueobjects.TransactionID
	userID          identityvalueobjects.UserID
	accountID       accountvalueobjects.AccountID
	transactionType transactionvalueobjects.TransactionType
	amount          sharedvalueobjects.Money
	description     transactionvalueobjects.TransactionDescription
	date            time.Time
	createdAt       time.Time
	updatedAt       time.Time

	// Domain events
	events []events.DomainEvent
}

// NewTransaction creates a new Transaction aggregate.
func NewTransaction(
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	transactionType transactionvalueobjects.TransactionType,
	amount sharedvalueobjects.Money,
	description transactionvalueobjects.TransactionDescription,
	date time.Time,
) (*Transaction, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if accountID.IsEmpty() {
		return nil, errors.New("account ID cannot be empty")
	}

	if transactionType.Value() == "" {
		return nil, errors.New("transaction type cannot be empty")
	}

	if amount.IsZero() {
		return nil, errors.New("transaction amount cannot be zero")
	}

	if amount.IsNegative() {
		return nil, errors.New("transaction amount cannot be negative")
	}

	if description.IsEmpty() {
		return nil, errors.New("transaction description cannot be empty")
	}

	if date.IsZero() {
		return nil, errors.New("transaction date cannot be zero")
	}

	now := time.Now()

	transaction := &Transaction{
		id:              transactionvalueobjects.GenerateTransactionID(),
		userID:          userID,
		accountID:       accountID,
		transactionType: transactionType,
		amount:          amount,
		description:     description,
		date:            date,
		createdAt:       now,
		updatedAt:       now,
		events:          []events.DomainEvent{},
	}

	// Add domain event
	transaction.addEvent(events.NewBaseDomainEvent(
		"TransactionCreated",
		transaction.id.Value(),
		"Transaction",
	))

	return transaction, nil
}

// TransactionFromPersistence reconstructs a Transaction aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func TransactionFromPersistence(
	id transactionvalueobjects.TransactionID,
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	transactionType transactionvalueobjects.TransactionType,
	amount sharedvalueobjects.Money,
	description transactionvalueobjects.TransactionDescription,
	date time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*Transaction, error) {
	if id.IsEmpty() {
		return nil, errors.New("transaction ID cannot be empty")
	}
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}
	if accountID.IsEmpty() {
		return nil, errors.New("account ID cannot be empty")
	}
	if transactionType.Value() == "" {
		return nil, errors.New("transaction type cannot be empty")
	}
	if description.IsEmpty() {
		return nil, errors.New("transaction description cannot be empty")
	}

	return &Transaction{
		id:              id,
		userID:          userID,
		accountID:       accountID,
		transactionType: transactionType,
		amount:          amount,
		description:     description,
		date:            date,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
		events:          []events.DomainEvent{},
	}, nil
}

// ID returns the transaction ID.
func (t *Transaction) ID() transactionvalueobjects.TransactionID {
	return t.id
}

// UserID returns the user ID.
func (t *Transaction) UserID() identityvalueobjects.UserID {
	return t.userID
}

// AccountID returns the account ID.
func (t *Transaction) AccountID() accountvalueobjects.AccountID {
	return t.accountID
}

// TransactionType returns the transaction type.
func (t *Transaction) TransactionType() transactionvalueobjects.TransactionType {
	return t.transactionType
}

// Amount returns the transaction amount.
func (t *Transaction) Amount() sharedvalueobjects.Money {
	return t.amount
}

// Description returns the transaction description.
func (t *Transaction) Description() transactionvalueobjects.TransactionDescription {
	return t.description
}

// Date returns the transaction date.
func (t *Transaction) Date() time.Time {
	return t.date
}

// CreatedAt returns the creation timestamp.
func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update timestamp.
func (t *Transaction) UpdatedAt() time.Time {
	return t.updatedAt
}

// UpdateAmount updates the transaction amount.
func (t *Transaction) UpdateAmount(amount sharedvalueobjects.Money) error {
	if amount.IsZero() {
		return errors.New("transaction amount cannot be zero")
	}

	if amount.IsNegative() {
		return errors.New("transaction amount cannot be negative")
	}

	// Ensure currencies match
	if !t.amount.Currency().Equals(amount.Currency()) {
		return errors.New("cannot update amount with different currency")
	}

	t.amount = amount
	t.updatedAt = time.Now()

	t.addEvent(events.NewBaseDomainEvent(
		"TransactionAmountUpdated",
		t.id.Value(),
		"Transaction",
	))

	return nil
}

// UpdateDescription updates the transaction description.
func (t *Transaction) UpdateDescription(description transactionvalueobjects.TransactionDescription) error {
	if description.IsEmpty() {
		return errors.New("transaction description cannot be empty")
	}

	t.description = description
	t.updatedAt = time.Now()

	t.addEvent(events.NewBaseDomainEvent(
		"TransactionDescriptionUpdated",
		t.id.Value(),
		"Transaction",
	))

	return nil
}

// UpdateDate updates the transaction date.
func (t *Transaction) UpdateDate(date time.Time) error {
	if date.IsZero() {
		return errors.New("transaction date cannot be zero")
	}

	t.date = date
	t.updatedAt = time.Now()

	t.addEvent(events.NewBaseDomainEvent(
		"TransactionDateUpdated",
		t.id.Value(),
		"Transaction",
	))

	return nil
}

// UpdateType updates the transaction type.
func (t *Transaction) UpdateType(transactionType transactionvalueobjects.TransactionType) error {
	if transactionType.Value() == "" {
		return errors.New("transaction type cannot be empty")
	}

	t.transactionType = transactionType
	t.updatedAt = time.Now()

	t.addEvent(events.NewBaseDomainEvent(
		"TransactionTypeUpdated",
		t.id.Value(),
		"Transaction",
	))

	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (t *Transaction) GetEvents() []events.DomainEvent {
	return t.events
}

// ClearEvents clears all domain events from this aggregate.
func (t *Transaction) ClearEvents() {
	t.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (t *Transaction) addEvent(event events.DomainEvent) {
	t.events = append(t.events, event)
}
