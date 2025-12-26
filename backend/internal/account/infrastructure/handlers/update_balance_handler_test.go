package handlers

import (
	"testing"

	"gestao-financeira/backend/internal/account/domain/entities"
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	transactionevents "gestao-financeira/backend/internal/transaction/domain/events"
)

// Mock account repository for testing
type mockAccountRepositoryForBalanceHandler struct {
	accounts    map[string]*entities.Account
	findByIDErr error
	saveErr     error
}

func newMockAccountRepositoryForBalanceHandler() *mockAccountRepositoryForBalanceHandler {
	return &mockAccountRepositoryForBalanceHandler{
		accounts: make(map[string]*entities.Account),
	}
}

func (m *mockAccountRepositoryForBalanceHandler) FindByID(id accountvalueobjects.AccountID) (*entities.Account, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.accounts[id.Value()], nil
}

func (m *mockAccountRepositoryForBalanceHandler) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	return nil, nil
}

func (m *mockAccountRepositoryForBalanceHandler) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	return nil, nil
}

func (m *mockAccountRepositoryForBalanceHandler) Save(account *entities.Account) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.accounts[account.ID().Value()] = account
	return nil
}

func (m *mockAccountRepositoryForBalanceHandler) Delete(id accountvalueobjects.AccountID) error {
	return nil
}

func (m *mockAccountRepositoryForBalanceHandler) Exists(id accountvalueobjects.AccountID) (bool, error) {
	return false, nil
}

func (m *mockAccountRepositoryForBalanceHandler) Count(userID identityvalueobjects.UserID) (int64, error) {
	return 0, nil
}

func TestUpdateBalanceHandler_HandleTransactionCreated(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  float64
		transactionType string
		amount          float64
		expectedBalance float64
		wantError       bool
	}{
		{
			name:            "INCOME transaction increases balance",
			initialBalance:  100.00,
			transactionType: "INCOME",
			amount:          50.00,
			expectedBalance: 150.00,
			wantError:       false,
		},
		{
			name:            "EXPENSE transaction decreases balance",
			initialBalance:  100.00,
			transactionType: "EXPENSE",
			amount:          30.00,
			expectedBalance: 70.00,
			wantError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := newMockAccountRepositoryForBalanceHandler()

			// Create account with initial balance
			userID := identityvalueobjects.GenerateUserID()
			accountName, _ := accountvalueobjects.NewAccountName("Test Account")
			accountType, _ := accountvalueobjects.NewAccountType("BANK")
			currency, _ := sharedvalueobjects.NewCurrency("BRL")
			initialBalance, _ := sharedvalueobjects.NewMoneyFromFloat(tt.initialBalance, currency)
			context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

			account, err := entities.NewAccount(userID, accountName, accountType, initialBalance, context)
			if err != nil {
				t.Fatalf("Failed to create account: %v", err)
			}

			// Save account to mock repository
			if err := mockRepo.Save(account); err != nil {
				t.Fatalf("Failed to save account: %v", err)
			}

			// Create handler
			handler := NewUpdateBalanceHandler(mockRepo)

			// Create transaction event
			amount, _ := sharedvalueobjects.NewMoneyFromFloat(tt.amount, currency)
			event := transactionevents.NewTransactionCreated(
				"transaction-123",
				account.ID().Value(),
				tt.transactionType,
				amount,
			)

			// Handle event
			err = handler.HandleTransactionCreated(event)
			if (err != nil) != tt.wantError {
				t.Errorf("HandleTransactionCreated() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify account balance was updated
				updatedAccount, err := mockRepo.FindByID(account.ID())
				if err != nil {
					t.Fatalf("Failed to find account: %v", err)
				}

				actualBalance := updatedAccount.Balance().Float64()
				if actualBalance != tt.expectedBalance {
					t.Errorf("Expected balance %f, got %f", tt.expectedBalance, actualBalance)
				}
			}
		})
	}
}

func TestUpdateBalanceHandler_HandleTransactionDeleted(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  float64
		transactionType string
		amount          float64
		expectedBalance float64
		wantError       bool
	}{
		{
			name:            "Deleting INCOME transaction decreases balance",
			initialBalance:  150.00,
			transactionType: "INCOME",
			amount:          50.00,
			expectedBalance: 100.00,
			wantError:       false,
		},
		{
			name:            "Deleting EXPENSE transaction increases balance",
			initialBalance:  70.00,
			transactionType: "EXPENSE",
			amount:          30.00,
			expectedBalance: 100.00,
			wantError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := newMockAccountRepositoryForBalanceHandler()

			// Create account with initial balance
			userID := identityvalueobjects.GenerateUserID()
			accountName, _ := accountvalueobjects.NewAccountName("Test Account")
			accountType, _ := accountvalueobjects.NewAccountType("BANK")
			currency, _ := sharedvalueobjects.NewCurrency("BRL")
			initialBalance, _ := sharedvalueobjects.NewMoneyFromFloat(tt.initialBalance, currency)
			context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

			account, err := entities.NewAccount(userID, accountName, accountType, initialBalance, context)
			if err != nil {
				t.Fatalf("Failed to create account: %v", err)
			}

			// Save account to mock repository
			if err := mockRepo.Save(account); err != nil {
				t.Fatalf("Failed to save account: %v", err)
			}

			// Create handler
			handler := NewUpdateBalanceHandler(mockRepo)

			// Create transaction deleted event
			amount, _ := sharedvalueobjects.NewMoneyFromFloat(tt.amount, currency)
			event := transactionevents.NewTransactionDeleted(
				"transaction-123",
				account.ID().Value(),
				tt.transactionType,
				amount,
			)

			// Handle event
			err = handler.HandleTransactionDeleted(event)
			if (err != nil) != tt.wantError {
				t.Errorf("HandleTransactionDeleted() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify account balance was updated
				updatedAccount, err := mockRepo.FindByID(account.ID())
				if err != nil {
					t.Fatalf("Failed to find account: %v", err)
				}

				actualBalance := updatedAccount.Balance().Float64()
				if actualBalance != tt.expectedBalance {
					t.Errorf("Expected balance %f, got %f", tt.expectedBalance, actualBalance)
				}
			}
		})
	}
}

func TestUpdateBalanceHandler_HandleTransactionUpdated(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  float64
		oldType         string
		oldAmount       float64
		newType         string
		newAmount       float64
		expectedBalance float64
		wantError       bool
	}{
		{
			name:            "Update INCOME amount",
			initialBalance:  100.00,
			oldType:         "INCOME",
			oldAmount:       50.00,
			newType:         "INCOME",
			newAmount:       75.00,
			expectedBalance: 175.00, // 100 (initial) + 50 (old INCOME) = 150, then reverse 50 = 100, then add 75 = 175
			wantError:       false,
		},
		{
			name:            "Change from INCOME to EXPENSE",
			initialBalance:  100.00,
			oldType:         "INCOME",
			oldAmount:       50.00,
			newType:         "EXPENSE",
			newAmount:       30.00,
			expectedBalance: 70.00, // 100 (initial) + 50 (old INCOME) = 150, then reverse 50 = 100, then subtract 30 = 70
			wantError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := newMockAccountRepositoryForBalanceHandler()

			// Create account with initial balance
			// Note: initial balance should reflect the old transaction already applied
			userID := identityvalueobjects.GenerateUserID()
			accountName, _ := accountvalueobjects.NewAccountName("Test Account")
			accountType, _ := accountvalueobjects.NewAccountType("BANK")
			currency, _ := sharedvalueobjects.NewCurrency("BRL")

			// Calculate initial balance: if old type was INCOME, add it; if EXPENSE, subtract it
			var balanceAfterOldTransaction float64
			if tt.oldType == "INCOME" {
				balanceAfterOldTransaction = tt.initialBalance + tt.oldAmount
			} else {
				balanceAfterOldTransaction = tt.initialBalance - tt.oldAmount
			}

			initialBalance, _ := sharedvalueobjects.NewMoneyFromFloat(balanceAfterOldTransaction, currency)
			context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

			account, err := entities.NewAccount(userID, accountName, accountType, initialBalance, context)
			if err != nil {
				t.Fatalf("Failed to create account: %v", err)
			}

			// Save account to mock repository
			if err := mockRepo.Save(account); err != nil {
				t.Fatalf("Failed to save account: %v", err)
			}

			// Create handler
			handler := NewUpdateBalanceHandler(mockRepo)

			// Create transaction updated event
			oldAmount, _ := sharedvalueobjects.NewMoneyFromFloat(tt.oldAmount, currency)
			newAmount, _ := sharedvalueobjects.NewMoneyFromFloat(tt.newAmount, currency)
			event := transactionevents.NewTransactionUpdated(
				"transaction-123",
				account.ID().Value(),
				tt.oldType,
				oldAmount,
				tt.newType,
				newAmount,
			)

			// Handle event
			err = handler.HandleTransactionUpdated(event)
			if (err != nil) != tt.wantError {
				t.Errorf("HandleTransactionUpdated() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify account balance was updated
				updatedAccount, err := mockRepo.FindByID(account.ID())
				if err != nil {
					t.Fatalf("Failed to find account: %v", err)
				}

				actualBalance := updatedAccount.Balance().Float64()
				// Allow small floating point differences
				diff := actualBalance - tt.expectedBalance
				if diff < 0 {
					diff = -diff
				}
				if diff > 0.01 {
					t.Errorf("Expected balance %f, got %f (diff: %f)", tt.expectedBalance, actualBalance, diff)
				}
			}
		})
	}
}
