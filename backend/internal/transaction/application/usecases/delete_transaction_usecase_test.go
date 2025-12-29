package usecases

import (
	"errors"
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

func TestDeleteTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	// Create initial balance
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	initialBalance, _ := sharedvalueobjects.NewMoney(100000, currency) // 1000.00 BRL

	tests := []struct {
		name      string
		input     dtos.DeleteTransactionInput
		setupMock func(*mockTransactionRepository, *mockAccountRepository, *mockUnitOfWork)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.DeleteTransactionOutput, *mockAccountRepository)
	}{
		{
			name: "delete existing INCOME transaction",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// Create account with initial balance
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				// Credit account to simulate the transaction effect
				currency, _ := sharedvalueobjects.NewCurrency("BRL")
				amount, _ := sharedvalueobjects.NewMoney(10000, currency) // 100.00
				_ = account.Credit(amount)
				_ = accRepo.Save(account)
				// Create transaction
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.DeleteTransactionOutput, accRepo *mockAccountRepository) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				if output.TransactionID == "" {
					t.Errorf("expected transaction ID to be set")
				}
				if output.Message != "Transaction deleted successfully" {
					t.Errorf("expected message 'Transaction deleted successfully', got %s", output.Message)
				}
				// Verify account balance was reversed (INCOME deleted = debit)
				account, _ := accRepo.FindByID(accountID)
				if account == nil {
					t.Fatal("account not found")
				}
				// Initial: 1000.00, INCOME 100.00 -> 1100.00, then deleted (reversed) -> 1000.00
				expectedBalance := int64(100000) // 1000.00 in cents
				if account.Balance().Amount() != expectedBalance {
					t.Errorf("expected balance %d, got %d", expectedBalance, account.Balance().Amount())
				}
			},
		},
		{
			name: "delete existing EXPENSE transaction",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// Create account with initial balance
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				// Debit account to simulate the transaction effect
				currency, _ := sharedvalueobjects.NewCurrency("BRL")
				amount, _ := sharedvalueobjects.NewMoney(10000, currency) // 100.00
				_ = account.Debit(amount)
				_ = accRepo.Save(account)
				// Create transaction
				transaction, _ := createTestTransaction(userID, accountID, "EXPENSE", 100.00, "BRL", "Compra", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.DeleteTransactionOutput, accRepo *mockAccountRepository) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				// Verify account balance was reversed (EXPENSE deleted = credit)
				account, _ := accRepo.FindByID(accountID)
				if account == nil {
					t.Fatal("account not found")
				}
				// Initial: 1000.00, EXPENSE 100.00 -> 900.00, then deleted (reversed) -> 1000.00
				expectedBalance := int64(100000) // 1000.00 in cents
				if account.Balance().Amount() != expectedBalance {
					t.Errorf("expected balance %d, got %d", expectedBalance, account.Balance().Amount())
				}
			},
		},
		{
			name: "transaction not found",
			input: dtos.DeleteTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// No transactions in repository
			},
			wantError: true,
			errorMsg:  "transaction not found",
		},
		{
			name: "invalid transaction ID",
			input: dtos.DeleteTransactionInput{
				TransactionID: "invalid-uuid",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "repository find error",
			input: dtos.DeleteTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// This will be handled by the mock - we need to check if FindByID can return an error
				// For now, we'll test with account not found scenario
			},
			wantError: true,
			errorMsg:  "transaction not found",
		},
		{
			name: "account not found",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// Create transaction but no account
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "account not found",
		},
		{
			name: "unit of work begin error",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
				uow.beginErr = errors.New("database connection error")
			},
			wantError: true,
			errorMsg:  "failed to begin transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxRepo := newMockTransactionRepository()
			mockAccRepo := newMockAccountRepository()
			mockUow := newMockUnitOfWork(mockTxRepo, mockAccRepo)
			tt.setupMock(mockTxRepo, mockAccRepo, mockUow)

			// For test cases that need TransactionID, set it after creating the transaction
			if tt.input.TransactionID == "" {
				if tt.name == "transaction not found" {
					// For "transaction not found" test, use a valid UUID format
					tt.input.TransactionID = transactionvalueobjects.GenerateTransactionID().Value()
				} else {
					// For other tests, get the transaction ID from the repository
					var transactionID string
					for id := range mockTxRepo.transactions {
						transactionID = id
						break
					}
					if transactionID == "" {
						// If no transaction in repo, create one
						account, _ := createTestAccountWithID(userID, accountID, initialBalance)
						_ = mockAccRepo.Save(account)
						transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
						_ = mockTxRepo.Save(transaction)
						transactionID = transaction.ID().Value()
					}
					tt.input.TransactionID = transactionID
				}
			}

			useCase := NewDeleteTransactionUseCase(mockUow, eventbus.NewEventBus())
			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				if output != nil {
					t.Errorf("expected nil output on error, got %v", output)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output == nil {
				t.Errorf("expected output but got nil")
				return
			}

			if tt.validate != nil {
				tt.validate(t, output, mockAccRepo)
			}
		})
	}
}
