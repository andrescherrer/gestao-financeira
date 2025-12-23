package entities

import (
	"errors"
	"time"

	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// Account represents an account aggregate root in the Account Management context.
type Account struct {
	id          valueobjects.AccountID
	userID      identityvalueobjects.UserID
	name        valueobjects.AccountName
	accountType valueobjects.AccountType
	balance     sharedvalueobjects.Money
	context     sharedvalueobjects.AccountContext
	createdAt   time.Time
	updatedAt   time.Time
	isActive    bool

	// Domain events
	events []events.DomainEvent
}

// NewAccount creates a new Account aggregate.
func NewAccount(
	userID identityvalueobjects.UserID,
	name valueobjects.AccountName,
	accountType valueobjects.AccountType,
	initialBalance sharedvalueobjects.Money,
	context sharedvalueobjects.AccountContext,
) (*Account, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if name.IsEmpty() {
		return nil, errors.New("account name cannot be empty")
	}

	if accountType.Value() == "" {
		return nil, errors.New("account type cannot be empty")
	}

	// Validate that initial balance currency matches context (if needed)
	// For now, we'll just ensure balance is not negative
	if initialBalance.IsNegative() {
		return nil, errors.New("initial balance cannot be negative")
	}

	now := time.Now()

	account := &Account{
		id:          valueobjects.GenerateAccountID(),
		userID:      userID,
		name:        name,
		accountType: accountType,
		balance:     initialBalance,
		context:     context,
		createdAt:   now,
		updatedAt:   now,
		isActive:    true,
		events:      []events.DomainEvent{},
	}

	// Add domain event
	account.addEvent(events.NewBaseDomainEvent(
		"AccountCreated",
		account.id.Value(),
		"Account",
	))

	return account, nil
}

// FromPersistence reconstructs an Account aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func AccountFromPersistence(
	id valueobjects.AccountID,
	userID identityvalueobjects.UserID,
	name valueobjects.AccountName,
	accountType valueobjects.AccountType,
	balance sharedvalueobjects.Money,
	context sharedvalueobjects.AccountContext,
	createdAt time.Time,
	updatedAt time.Time,
	isActive bool,
) (*Account, error) {
	if id.IsEmpty() {
		return nil, errors.New("account ID cannot be empty")
	}
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}
	if name.IsEmpty() {
		return nil, errors.New("account name cannot be empty")
	}

	return &Account{
		id:          id,
		userID:      userID,
		name:        name,
		accountType: accountType,
		balance:     balance,
		context:     context,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		isActive:    isActive,
		events:      []events.DomainEvent{},
	}, nil
}

// ID returns the account ID.
func (a *Account) ID() valueobjects.AccountID {
	return a.id
}

// UserID returns the user ID.
func (a *Account) UserID() identityvalueobjects.UserID {
	return a.userID
}

// Name returns the account name.
func (a *Account) Name() valueobjects.AccountName {
	return a.name
}

// AccountType returns the account type.
func (a *Account) AccountType() valueobjects.AccountType {
	return a.accountType
}

// Balance returns the account balance.
func (a *Account) Balance() sharedvalueobjects.Money {
	return a.balance
}

// Context returns the account context.
func (a *Account) Context() sharedvalueobjects.AccountContext {
	return a.context
}

// IsActive returns whether the account is active.
func (a *Account) IsActive() bool {
	return a.isActive
}

// CreatedAt returns the creation timestamp.
func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

// UpdatedAt returns the last update timestamp.
func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

// Credit adds money to the account balance.
func (a *Account) Credit(amount sharedvalueobjects.Money) error {
	if !a.isActive {
		return errors.New("cannot credit inactive account")
	}

	// Ensure currencies match
	if !a.balance.Currency().Equals(amount.Currency()) {
		return errors.New("cannot credit with different currency")
	}

	newBalance, err := a.balance.Add(amount)
	if err != nil {
		return err
	}

	a.balance = newBalance
	a.updatedAt = time.Now()

	a.addEvent(events.NewBaseDomainEvent(
		"AccountBalanceUpdated",
		a.id.Value(),
		"Account",
	))

	return nil
}

// Debit subtracts money from the account balance.
func (a *Account) Debit(amount sharedvalueobjects.Money) error {
	if !a.isActive {
		return errors.New("cannot debit inactive account")
	}

	// Ensure currencies match
	if !a.balance.Currency().Equals(amount.Currency()) {
		return errors.New("cannot debit with different currency")
	}

	// Check if balance would become negative
	newBalance, err := a.balance.Subtract(amount)
	if err != nil {
		return err
	}

	if newBalance.IsNegative() {
		return errors.New("insufficient balance")
	}

	a.balance = newBalance
	a.updatedAt = time.Now()

	a.addEvent(events.NewBaseDomainEvent(
		"AccountBalanceUpdated",
		a.id.Value(),
		"Account",
	))

	return nil
}

// UpdateName updates the account name.
func (a *Account) UpdateName(name valueobjects.AccountName) error {
	if !a.isActive {
		return errors.New("cannot update name for inactive account")
	}

	if name.IsEmpty() {
		return errors.New("account name cannot be empty")
	}

	a.name = name
	a.updatedAt = time.Now()

	a.addEvent(events.NewBaseDomainEvent(
		"AccountNameChanged",
		a.id.Value(),
		"Account",
	))

	return nil
}

// Deactivate deactivates the account.
func (a *Account) Deactivate() error {
	if !a.isActive {
		return errors.New("account is already inactive")
	}

	a.isActive = false
	a.updatedAt = time.Now()

	a.addEvent(events.NewBaseDomainEvent(
		"AccountDeactivated",
		a.id.Value(),
		"Account",
	))

	return nil
}

// Activate activates the account.
func (a *Account) Activate() error {
	if a.isActive {
		return errors.New("account is already active")
	}

	a.isActive = true
	a.updatedAt = time.Now()

	a.addEvent(events.NewBaseDomainEvent(
		"AccountActivated",
		a.id.Value(),
		"Account",
	))

	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (a *Account) GetEvents() []events.DomainEvent {
	return a.events
}

// ClearEvents clears all domain events from this aggregate.
func (a *Account) ClearEvents() {
	a.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (a *Account) addEvent(event events.DomainEvent) {
	a.events = append(a.events, event)
}
