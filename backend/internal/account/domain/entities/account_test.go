package entities

import (
	"testing"
	"time"

	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

func TestNewAccount(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()

	tests := []struct {
		name           string
		userID         identityvalueobjects.UserID
		accountName    valueobjects.AccountName
		accountType    valueobjects.AccountType
		initialBalance sharedvalueobjects.Money
		context        sharedvalueobjects.AccountContext
		wantError      bool
	}{
		{"valid account", userID, accountName, accountType, initialBalance, context, false},
		{"empty user ID", identityvalueobjects.UserID{}, accountName, accountType, initialBalance, context, true},
		{"empty account name", userID, valueobjects.AccountName{}, accountType, initialBalance, context, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount(tt.userID, tt.accountName, tt.accountType, tt.initialBalance, tt.context)
			if (err != nil) != tt.wantError {
				t.Errorf("NewAccount() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if account.ID().IsEmpty() {
					t.Error("NewAccount() returned account with empty ID")
				}
				if !account.IsActive() {
					t.Error("NewAccount() returned inactive account")
				}
				// Check domain event
				events := account.GetEvents()
				if len(events) == 0 {
					t.Error("NewAccount() should create AccountCreated event")
				}
			}
		})
	}
}

func TestAccount_Credit(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL")) // 100.00 BRL
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)

	// Credit valid amount
	creditAmount, _ := sharedvalueobjects.NewMoney(5000, sharedvalueobjects.MustCurrency("BRL")) // 50.00 BRL
	err := account.Credit(creditAmount)
	if err != nil {
		t.Errorf("Account.Credit() error = %v, want nil", err)
	}

	// Check balance
	expectedBalance, _ := sharedvalueobjects.NewMoney(15000, sharedvalueobjects.MustCurrency("BRL")) // 150.00 BRL
	if !account.Balance().Equals(expectedBalance) {
		t.Errorf("Account.Balance() = %v, want %v", account.Balance(), expectedBalance)
	}

	// Try to credit inactive account
	account.Deactivate()
	err = account.Credit(creditAmount)
	if err == nil {
		t.Error("Account.Credit() should fail for inactive account")
	}
}

func TestAccount_Debit(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL")) // 100.00 BRL
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)

	// Debit valid amount
	debitAmount, _ := sharedvalueobjects.NewMoney(5000, sharedvalueobjects.MustCurrency("BRL")) // 50.00 BRL
	err := account.Debit(debitAmount)
	if err != nil {
		t.Errorf("Account.Debit() error = %v, want nil", err)
	}

	// Check balance
	expectedBalance, _ := sharedvalueobjects.NewMoney(5000, sharedvalueobjects.MustCurrency("BRL")) // 50.00 BRL
	if !account.Balance().Equals(expectedBalance) {
		t.Errorf("Account.Balance() = %v, want %v", account.Balance(), expectedBalance)
	}

	// Try to debit more than balance (try to debit 60.00, but only 50.00 left)
	largeDebitAmount, _ := sharedvalueobjects.NewMoney(6000, sharedvalueobjects.MustCurrency("BRL")) // 60.00 BRL
	err = account.Debit(largeDebitAmount)
	if err == nil {
		t.Error("Account.Debit() should fail when insufficient balance")
	}

	// Try to debit inactive account
	account2, _ := NewAccount(userID, accountName, accountType, initialBalance, context)
	account2.Deactivate()
	err = account2.Debit(debitAmount)
	if err == nil {
		t.Error("Account.Debit() should fail for inactive account")
	}
}

func TestAccount_UpdateName(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)

	// Update name
	newName, _ := valueobjects.NewAccountName("Conta Poupança")
	err := account.UpdateName(newName)
	if err != nil {
		t.Errorf("Account.UpdateName() error = %v, want nil", err)
	}

	if !account.Name().Equals(newName) {
		t.Errorf("Account.Name() = %v, want %v", account.Name(), newName)
	}

	// Try to update name for inactive account
	account.Deactivate()
	err = account.UpdateName(newName)
	if err == nil {
		t.Error("Account.UpdateName() should fail for inactive account")
	}
}

func TestAccount_Deactivate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)

	// Deactivate
	err := account.Deactivate()
	if err != nil {
		t.Errorf("Account.Deactivate() error = %v, want nil", err)
	}

	if account.IsActive() {
		t.Error("Account.IsActive() = true, want false")
	}

	// Try to deactivate again
	err = account.Deactivate()
	if err == nil {
		t.Error("Account.Deactivate() should fail if already inactive")
	}
}

func TestAccount_Activate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)
	account.Deactivate()

	// Activate
	err := account.Activate()
	if err != nil {
		t.Errorf("Account.Activate() error = %v, want nil", err)
	}

	if !account.IsActive() {
		t.Error("Account.IsActive() = false, want true")
	}

	// Try to activate again
	err = account.Activate()
	if err == nil {
		t.Error("Account.Activate() should fail if already active")
	}
}

func TestAccountFromPersistence(t *testing.T) {
	accountID := valueobjects.GenerateAccountID()
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	balance, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()
	createdAt := time.Now()
	updatedAt := time.Now()

	account, err := AccountFromPersistence(accountID, userID, accountName, accountType, balance, context, createdAt, updatedAt, true)
	if err != nil {
		t.Errorf("AccountFromPersistence() error = %v, want nil", err)
	}

	if account.ID() != accountID {
		t.Errorf("Account.ID() = %v, want %v", account.ID(), accountID)
	}

	// Should not have events
	events := account.GetEvents()
	if len(events) != 0 {
		t.Error("AccountFromPersistence() should not create domain events")
	}
}

func TestAccount_GetEvents(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType := valueobjects.BankType()
	initialBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
	context := sharedvalueobjects.PersonalContext()

	account, _ := NewAccount(userID, accountName, accountType, initialBalance, context)

	// Should have AccountCreated event
	events := account.GetEvents()
	if len(events) == 0 {
		t.Error("Account.GetEvents() should return events")
	}

	// Clear events
	account.ClearEvents()
	events = account.GetEvents()
	if len(events) != 0 {
		t.Error("Account.ClearEvents() should clear all events")
	}
}
