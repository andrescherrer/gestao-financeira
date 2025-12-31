package entities

import (
	"errors"
	"fmt"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	investmentevents "gestao-financeira/backend/internal/investment/domain/events"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// Investment represents an investment aggregate root in the Investment context.
type Investment struct {
	id             investmentvalueobjects.InvestmentID
	userID         identityvalueobjects.UserID
	accountID      accountvalueobjects.AccountID
	investmentType investmentvalueobjects.InvestmentType
	name           investmentvalueobjects.InvestmentName
	purchaseDate   time.Time
	purchaseAmount sharedvalueobjects.Money
	currentValue   sharedvalueobjects.Money
	quantity       *float64 // Optional, depends on investment type
	context        sharedvalueobjects.AccountContext
	createdAt      time.Time
	updatedAt      time.Time

	// Domain events
	events []events.DomainEvent
}

// NewInvestment creates a new Investment aggregate.
func NewInvestment(
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	investmentType investmentvalueobjects.InvestmentType,
	name investmentvalueobjects.InvestmentName,
	purchaseDate time.Time,
	purchaseAmount sharedvalueobjects.Money,
	quantity *float64,
	context sharedvalueobjects.AccountContext,
) (*Investment, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if accountID.IsEmpty() {
		return nil, errors.New("account ID cannot be empty")
	}

	if name.IsEmpty() {
		return nil, errors.New("investment name cannot be empty")
	}

	if purchaseDate.IsZero() {
		return nil, errors.New("purchase date cannot be zero")
	}

	if purchaseDate.After(time.Now()) {
		return nil, errors.New("purchase date cannot be in the future")
	}

	if purchaseAmount.IsZero() {
		return nil, errors.New("purchase amount cannot be zero")
	}

	if purchaseAmount.IsNegative() {
		return nil, errors.New("purchase amount cannot be negative")
	}

	// Validate quantity based on investment type
	if investmentType.RequiresQuantity() {
		if quantity == nil || *quantity <= 0 {
			return nil, errors.New("quantity is required and must be positive for this investment type")
		}
	}

	// Currency validation should be done at the use case level
	// to ensure it matches the account's currency

	now := time.Now()

	investment := &Investment{
		id:             investmentvalueobjects.GenerateInvestmentID(),
		userID:         userID,
		accountID:      accountID,
		investmentType: investmentType,
		name:           name,
		purchaseDate:   purchaseDate,
		purchaseAmount: purchaseAmount,
		currentValue:   purchaseAmount, // Initially, current value equals purchase amount
		quantity:       quantity,
		context:        context,
		createdAt:      now,
		updatedAt:      now,
		events:         []events.DomainEvent{},
	}

	// Add domain event
	ticker := ""
	if name.HasTicker() {
		ticker = *name.Ticker()
	}
	investment.addEvent(investmentevents.NewInvestmentCreated(
		investment.id.Value(),
		investment.userID.Value(),
		investment.accountID.Value(),
		investment.investmentType.Value(),
		investment.name.Name(),
		ticker,
		fmt.Sprintf("%.2f", purchaseAmount.Float64()),
		purchaseAmount.Currency().Code(),
		investment.context.Value(),
	))

	return investment, nil
}

// InvestmentFromPersistence reconstructs an Investment aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func InvestmentFromPersistence(
	id investmentvalueobjects.InvestmentID,
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	investmentType investmentvalueobjects.InvestmentType,
	name investmentvalueobjects.InvestmentName,
	purchaseDate time.Time,
	purchaseAmount sharedvalueobjects.Money,
	currentValue sharedvalueobjects.Money,
	quantity *float64,
	context sharedvalueobjects.AccountContext,
	createdAt time.Time,
	updatedAt time.Time,
) (*Investment, error) {
	if id.IsEmpty() {
		return nil, errors.New("investment ID cannot be empty")
	}
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}
	if accountID.IsEmpty() {
		return nil, errors.New("account ID cannot be empty")
	}
	if name.IsEmpty() {
		return nil, errors.New("investment name cannot be empty")
	}

	return &Investment{
		id:             id,
		userID:         userID,
		accountID:      accountID,
		investmentType: investmentType,
		name:           name,
		purchaseDate:   purchaseDate,
		purchaseAmount: purchaseAmount,
		currentValue:   currentValue,
		quantity:       quantity,
		context:        context,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		events:         []events.DomainEvent{},
	}, nil
}

// ID returns the investment ID.
func (i *Investment) ID() investmentvalueobjects.InvestmentID {
	return i.id
}

// UserID returns the user ID.
func (i *Investment) UserID() identityvalueobjects.UserID {
	return i.userID
}

// AccountID returns the account ID.
func (i *Investment) AccountID() accountvalueobjects.AccountID {
	return i.accountID
}

// InvestmentType returns the investment type.
func (i *Investment) InvestmentType() investmentvalueobjects.InvestmentType {
	return i.investmentType
}

// Name returns the investment name.
func (i *Investment) Name() investmentvalueobjects.InvestmentName {
	return i.name
}

// PurchaseDate returns the purchase date.
func (i *Investment) PurchaseDate() time.Time {
	return i.purchaseDate
}

// PurchaseAmount returns the purchase amount.
func (i *Investment) PurchaseAmount() sharedvalueobjects.Money {
	return i.purchaseAmount
}

// CurrentValue returns the current value.
func (i *Investment) CurrentValue() sharedvalueobjects.Money {
	return i.currentValue
}

// Quantity returns the quantity (if available).
func (i *Investment) Quantity() *float64 {
	return i.quantity
}

// Context returns the account context.
func (i *Investment) Context() sharedvalueobjects.AccountContext {
	return i.context
}

// CreatedAt returns the creation timestamp.
func (i *Investment) CreatedAt() time.Time {
	return i.createdAt
}

// UpdatedAt returns the last update timestamp.
func (i *Investment) UpdatedAt() time.Time {
	return i.updatedAt
}

// UpdateCurrentValue updates the current value of the investment.
func (i *Investment) UpdateCurrentValue(value sharedvalueobjects.Money) error {
	if value.IsNegative() {
		return errors.New("current value cannot be negative")
	}

	// Ensure currencies match
	if !value.Currency().Equals(i.purchaseAmount.Currency()) {
		return errors.New("current value currency must match purchase amount currency")
	}

	oldValue := i.currentValue
	i.currentValue = value
	i.updatedAt = time.Now()

	// Add domain event
	i.addEvent(investmentevents.NewInvestmentValueUpdated(
		i.id.Value(),
		fmt.Sprintf("%.2f", value.Float64()),
		value.Currency().Code(),
	))

	// Calculate and emit return if value changed
	if !oldValue.Equals(value) {
		returnObj := i.calculateReturn()
		i.addEvent(investmentevents.NewInvestmentReturnCalculated(
			i.id.Value(),
			fmt.Sprintf("%.2f", returnObj.Absolute().Float64()),
			returnObj.Percentage(),
			returnObj.Absolute().Currency().Code(),
		))
	}

	return nil
}

// CalculateReturn calculates the return of the investment.
func (i *Investment) CalculateReturn() investmentvalueobjects.InvestmentReturn {
	return i.calculateReturn()
}

// calculateReturn is the internal method that calculates the return.
func (i *Investment) calculateReturn() investmentvalueobjects.InvestmentReturn {
	// Calculate absolute return
	absoluteReturn, err := i.currentValue.Subtract(i.purchaseAmount)
	if err != nil {
		// This should not happen as currencies should match, but handle it gracefully
		absoluteReturn = sharedvalueobjects.Zero(i.purchaseAmount.Currency())
	}

	// Calculate percentage return
	var percentage float64
	if !i.purchaseAmount.IsZero() {
		percentage = (float64(absoluteReturn.Amount()) / float64(i.purchaseAmount.Amount())) * 100.0
	}

	return investmentvalueobjects.NewInvestmentReturn(absoluteReturn, percentage)
}

// CalculateReturnPercentage calculates the return percentage.
func (i *Investment) CalculateReturnPercentage() float64 {
	return i.calculateReturn().Percentage()
}

// AddQuantity adds quantity to the investment.
func (i *Investment) AddQuantity(quantity float64) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if !i.investmentType.RequiresQuantity() {
		return errors.New("this investment type does not require quantity")
	}

	if i.quantity == nil {
		i.quantity = &quantity
	} else {
		newQuantity := *i.quantity + quantity
		i.quantity = &newQuantity
	}

	i.updatedAt = time.Now()
	return nil
}

// RemoveQuantity removes quantity from the investment.
func (i *Investment) RemoveQuantity(quantity float64) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if !i.investmentType.RequiresQuantity() {
		return errors.New("this investment type does not require quantity")
	}

	if i.quantity == nil || *i.quantity < quantity {
		return errors.New("insufficient quantity")
	}

	newQuantity := *i.quantity - quantity
	if newQuantity == 0 {
		i.quantity = nil
	} else {
		i.quantity = &newQuantity
	}

	i.updatedAt = time.Now()
	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (i *Investment) GetEvents() []events.DomainEvent {
	return i.events
}

// ClearEvents clears all domain events from this aggregate.
func (i *Investment) ClearEvents() {
	i.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (i *Investment) addEvent(event events.DomainEvent) {
	i.events = append(i.events, event)
}
