package entities

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

func TestNewTransaction(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL")) // 100.00 BRL
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	tests := []struct {
		name            string
		userID          identityvalueobjects.UserID
		accountID       accountvalueobjects.AccountID
		transactionType transactionvalueobjects.TransactionType
		amount          sharedvalueobjects.Money
		description     transactionvalueobjects.TransactionDescription
		date            time.Time
		wantError       bool
	}{
		{"valid transaction", userID, accountID, transactionType, amount, description, date, false},
		{"empty user ID", identityvalueobjects.UserID{}, accountID, transactionType, amount, description, date, true},
		{"empty account ID", userID, accountvalueobjects.AccountID{}, transactionType, amount, description, date, true},
		{"zero amount", userID, accountID, transactionType, sharedvalueobjects.Zero(sharedvalueobjects.MustCurrency("BRL")), description, date, true},
		{"negative amount", userID, accountID, transactionType, func() sharedvalueobjects.Money {
			m, _ := sharedvalueobjects.NewMoney(-1000, sharedvalueobjects.MustCurrency("BRL"))
			return m
		}(), description, date, true},
		{"empty description", userID, accountID, transactionType, amount, transactionvalueobjects.TransactionDescription{}, date, true},
		{"zero date", userID, accountID, transactionType, amount, description, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewTransaction(tt.userID, tt.accountID, tt.transactionType, tt.amount, tt.description, tt.date)
			if (err != nil) != tt.wantError {
				t.Errorf("NewTransaction() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if transaction.ID().IsEmpty() {
					t.Error("NewTransaction() returned transaction with empty ID")
				}
				// Check domain event
				events := transaction.GetEvents()
				if len(events) == 0 {
					t.Error("NewTransaction() should create TransactionCreated event")
				}
			}
		})
	}
}

func TestTransaction_UpdateAmount(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Update amount
	newAmount, _ := sharedvalueobjects.NewMoney(20000, sharedvalueobjects.MustCurrency("BRL"))
	err := transaction.UpdateAmount(newAmount)
	if err != nil {
		t.Errorf("Transaction.UpdateAmount() error = %v, want nil", err)
	}

	if !transaction.Amount().Equals(newAmount) {
		t.Errorf("Transaction.Amount() = %v, want %v", transaction.Amount(), newAmount)
	}

	// Try to update with zero amount
	zeroAmount, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	err = transaction.UpdateAmount(zeroAmount)
	if err == nil {
		t.Error("Transaction.UpdateAmount() should fail with zero amount")
	}

	// Try to update with different currency
	usdAmount, _ := sharedvalueobjects.NewMoney(1000, sharedvalueobjects.MustCurrency("USD"))
	err = transaction.UpdateAmount(usdAmount)
	if err == nil {
		t.Error("Transaction.UpdateAmount() should fail with different currency")
	}
}

func TestTransaction_UpdateDescription(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Update description
	newDescription, _ := transactionvalueobjects.NewTransactionDescription("Nova descrição")
	err := transaction.UpdateDescription(newDescription)
	if err != nil {
		t.Errorf("Transaction.UpdateDescription() error = %v, want nil", err)
	}

	if !transaction.Description().Equals(newDescription) {
		t.Errorf("Transaction.Description() = %v, want %v", transaction.Description(), newDescription)
	}

	// Try to update with empty description
	err = transaction.UpdateDescription(transactionvalueobjects.TransactionDescription{})
	if err == nil {
		t.Error("Transaction.UpdateDescription() should fail with empty description")
	}
}

func TestTransaction_UpdateDate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Update date
	newDate := time.Now().Add(24 * time.Hour)
	err := transaction.UpdateDate(newDate)
	if err != nil {
		t.Errorf("Transaction.UpdateDate() error = %v, want nil", err)
	}

	if !transaction.Date().Equal(newDate) {
		t.Errorf("Transaction.Date() = %v, want %v", transaction.Date(), newDate)
	}

	// Try to update with zero date
	err = transaction.UpdateDate(time.Time{})
	if err == nil {
		t.Error("Transaction.UpdateDate() should fail with zero date")
	}
}

func TestTransaction_UpdateType(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Update type
	newType := transactionvalueobjects.ExpenseType()
	err := transaction.UpdateType(newType)
	if err != nil {
		t.Errorf("Transaction.UpdateType() error = %v, want nil", err)
	}

	if !transaction.TransactionType().Equals(newType) {
		t.Errorf("Transaction.TransactionType() = %v, want %v", transaction.TransactionType(), newType)
	}
}

func TestTransactionFromPersistence(t *testing.T) {
	transactionID := transactionvalueobjects.GenerateTransactionID()
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()
	createdAt := time.Now()
	updatedAt := time.Now()

	transaction, err := TransactionFromPersistence(transactionID, userID, accountID, transactionType, amount, description, date, createdAt, updatedAt)
	if err != nil {
		t.Errorf("TransactionFromPersistence() error = %v, want nil", err)
	}

	if !transaction.ID().Equals(transactionID) {
		t.Errorf("Transaction.ID() = %v, want %v", transaction.ID(), transactionID)
	}

	// Should not have events
	events := transaction.GetEvents()
	if len(events) != 0 {
		t.Error("TransactionFromPersistence() should not create domain events")
	}
}

func TestTransaction_GetEvents(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Should have TransactionCreated event
	events := transaction.GetEvents()
	if len(events) == 0 {
		t.Error("Transaction.GetEvents() should return events")
	}

	// Clear events
	transaction.ClearEvents()
	events = transaction.GetEvents()
	if len(events) != 0 {
		t.Error("Transaction.ClearEvents() should clear all events")
	}
}

func TestTransaction_Getters(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	date := time.Now()

	transaction, _ := NewTransaction(userID, accountID, transactionType, amount, description, date)

	// Test all getters
	if transaction.ID().IsEmpty() {
		t.Error("Transaction.ID() should not be empty")
	}

	if transaction.UserID().IsEmpty() {
		t.Error("Transaction.UserID() should not be empty")
	}

	if transaction.AccountID().IsEmpty() {
		t.Error("Transaction.AccountID() should not be empty")
	}

	if !transaction.TransactionType().Equals(transactionType) {
		t.Errorf("Transaction.TransactionType() = %v, want %v", transaction.TransactionType(), transactionType)
	}

	if !transaction.Amount().Equals(amount) {
		t.Errorf("Transaction.Amount() = %v, want %v", transaction.Amount(), amount)
	}

	if !transaction.Description().Equals(description) {
		t.Errorf("Transaction.Description() = %v, want %v", transaction.Description(), description)
	}

	if transaction.CreatedAt().IsZero() {
		t.Error("Transaction.CreatedAt() should not be zero")
	}

	if transaction.UpdatedAt().IsZero() {
		t.Error("Transaction.UpdatedAt() should not be zero")
	}
}

func TestNewTransactionWithRecurrence(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	date := time.Now()
	recurrenceFrequency := transactionvalueobjects.MonthlyFrequency()
	endDate := date.AddDate(1, 0, 0) // 1 year from now

	// Test valid recurring transaction
	transaction, err := NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionType,
		amount,
		description,
		date,
		true,
		&recurrenceFrequency,
		&endDate,
		nil,
	)
	if err != nil {
		t.Fatalf("NewTransactionWithRecurrence() error = %v, want nil", err)
	}

	if !transaction.IsRecurring() {
		t.Error("Transaction.IsRecurring() should be true")
	}

	if transaction.RecurrenceFrequency() == nil {
		t.Error("Transaction.RecurrenceFrequency() should not be nil")
	} else if !transaction.RecurrenceFrequency().Equals(recurrenceFrequency) {
		t.Errorf("Transaction.RecurrenceFrequency() = %v, want %v", transaction.RecurrenceFrequency(), recurrenceFrequency)
	}

	if transaction.RecurrenceEndDate() == nil {
		t.Error("Transaction.RecurrenceEndDate() should not be nil")
	} else if !transaction.RecurrenceEndDate().Equal(endDate) {
		t.Errorf("Transaction.RecurrenceEndDate() = %v, want %v", transaction.RecurrenceEndDate(), endDate)
	}

	// Test non-recurring transaction (should not have recurrence fields)
	nonRecurring, err := NewTransaction(userID, accountID, transactionType, amount, description, date)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v, want nil", err)
	}

	if nonRecurring.IsRecurring() {
		t.Error("Non-recurring transaction should have IsRecurring() = false")
	}

	if nonRecurring.RecurrenceFrequency() != nil {
		t.Error("Non-recurring transaction should have nil RecurrenceFrequency()")
	}

	if nonRecurring.RecurrenceEndDate() != nil {
		t.Error("Non-recurring transaction should have nil RecurrenceEndDate()")
	}
}

func TestNewTransactionWithRecurrence_InvalidInput(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	date := time.Now()
	recurrenceFrequency := transactionvalueobjects.MonthlyFrequency()
	endDate := date.AddDate(1, 0, 0)

	tests := []struct {
		name                string
		isRecurring         bool
		recurrenceFrequency *transactionvalueobjects.RecurrenceFrequency
		recurrenceEndDate   *time.Time
		wantError           bool
	}{
		{
			name:                "recurring without frequency",
			isRecurring:         true,
			recurrenceFrequency: nil,
			recurrenceEndDate:   nil,
			wantError:           true,
		},
		{
			name:                "recurring with end date before transaction date",
			isRecurring:         true,
			recurrenceFrequency: &recurrenceFrequency,
			recurrenceEndDate:   func() *time.Time { d := date.AddDate(-1, 0, 0); return &d }(),
			wantError:           true,
		},
		{
			name:                "non-recurring with frequency",
			isRecurring:         false,
			recurrenceFrequency: &recurrenceFrequency,
			recurrenceEndDate:   nil,
			wantError:           true,
		},
		{
			name:                "non-recurring with end date",
			isRecurring:         false,
			recurrenceFrequency: nil,
			recurrenceEndDate:   &endDate,
			wantError:           true,
		},
		{
			name:                "valid recurring transaction",
			isRecurring:         true,
			recurrenceFrequency: &recurrenceFrequency,
			recurrenceEndDate:   &endDate,
			wantError:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransactionWithRecurrence(
				userID,
				accountID,
				transactionType,
				amount,
				description,
				date,
				tt.isRecurring,
				tt.recurrenceFrequency,
				tt.recurrenceEndDate,
				nil,
			)
			if (err != nil) != tt.wantError {
				t.Errorf("NewTransactionWithRecurrence() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestTransactionFromPersistenceWithRecurrence(t *testing.T) {
	transactionID := transactionvalueobjects.GenerateTransactionID()
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.IncomeType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	date := time.Now()
	createdAt := time.Now()
	updatedAt := time.Now()
	recurrenceFrequency := transactionvalueobjects.MonthlyFrequency()
	endDate := date.AddDate(1, 0, 0)

	transaction, err := TransactionFromPersistenceWithRecurrence(
		transactionID,
		userID,
		accountID,
		transactionType,
		amount,
		description,
		date,
		createdAt,
		updatedAt,
		true,
		&recurrenceFrequency,
		&endDate,
		nil,
	)
	if err != nil {
		t.Fatalf("TransactionFromPersistenceWithRecurrence() error = %v, want nil", err)
	}

	if !transaction.IsRecurring() {
		t.Error("Transaction.IsRecurring() should be true")
	}

	if transaction.RecurrenceFrequency() == nil {
		t.Error("Transaction.RecurrenceFrequency() should not be nil")
	}

	if transaction.RecurrenceEndDate() == nil {
		t.Error("Transaction.RecurrenceEndDate() should not be nil")
	}

	// Should not have events
	events := transaction.GetEvents()
	if len(events) != 0 {
		t.Error("TransactionFromPersistenceWithRecurrence() should not create domain events")
	}
}
