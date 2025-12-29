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

func TestUpdateTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	eventBus := eventbus.NewEventBus()
	date := time.Now()
	newDateStr := date.Add(24 * time.Hour).Format("2006-01-02")

	// Create initial balance
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	initialBalance, _ := sharedvalueobjects.NewMoney(100000, currency) // 1000.00 BRL

	tests := []struct {
		name      string
		input     dtos.UpdateTransactionInput
		setupMock func(*mockTransactionRepository, *mockAccountRepository, *mockUnitOfWork)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.UpdateTransactionOutput, *mockAccountRepository)
	}{
		{
			name: "update transaction type",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput, accRepo *mockAccountRepository) {
				if output.Type != "EXPENSE" {
					t.Errorf("expected type 'EXPENSE', got %s", output.Type)
				}
				// Verify account balance was updated (reversed INCOME, applied EXPENSE)
				account, _ := accRepo.FindByID(accountID)
				if account == nil {
					t.Fatal("account not found")
				}
				// Initial: 1000.00, INCOME 100.00 -> 1100.00, then reversed -> 1000.00, then EXPENSE 100.00 -> 900.00
				expectedBalance := int64(90000) // 900.00 in cents
				if account.Balance().Amount() != expectedBalance {
					t.Errorf("expected balance %d, got %d", expectedBalance, account.Balance().Amount())
				}
			},
		},
		{
			name: "update transaction amount",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Amount:        floatPtr(200.50),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput, accRepo *mockAccountRepository) {
				if output.Amount != 200.50 {
					t.Errorf("expected amount 200.50, got %f", output.Amount)
				}
				// Verify account balance was updated
				account, _ := accRepo.FindByID(accountID)
				if account == nil {
					t.Fatal("account not found")
				}
				// Initial: 1000.00, INCOME 100.00 -> 1100.00, then reversed -> 1000.00, then INCOME 200.50 -> 1200.50
				expectedBalance := int64(120050) // 1200.50 in cents
				if account.Balance().Amount() != expectedBalance {
					t.Errorf("expected balance %d, got %d", expectedBalance, account.Balance().Amount())
				}
			},
		},
		{
			name: "update transaction description",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Description:   stringPtr("Nova descrição"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput, accRepo *mockAccountRepository) {
				if output.Description != "Nova descrição" {
					t.Errorf("expected description 'Nova descrição', got %s", output.Description)
				}
				// Description change doesn't affect balance
			},
		},
		{
			name: "update transaction date",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Date:          stringPtr(newDateStr),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput, accRepo *mockAccountRepository) {
				if output.Date != newDateStr {
					t.Errorf("expected date %s, got %s", newDateStr, output.Date)
				}
				// Date change doesn't affect balance
			},
		},
		{
			name: "update multiple fields",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
				Amount:        floatPtr(150.75),
				Description:   stringPtr("Atualizado"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput, accRepo *mockAccountRepository) {
				if output.Type != "EXPENSE" {
					t.Errorf("expected type 'EXPENSE', got %s", output.Type)
				}
				if output.Amount != 150.75 {
					t.Errorf("expected amount 150.75, got %f", output.Amount)
				}
				if output.Description != "Atualizado" {
					t.Errorf("expected description 'Atualizado', got %s", output.Description)
				}
			},
		},
		{
			name: "transaction not found",
			input: dtos.UpdateTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
				Type:          stringPtr("EXPENSE"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				// No transactions in repository
			},
			wantError: true,
			errorMsg:  "transaction not found",
		},
		{
			name: "invalid transaction ID",
			input: dtos.UpdateTransactionInput{
				TransactionID: "invalid-uuid",
				Type:          stringPtr("EXPENSE"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "invalid transaction type",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("INVALID"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "invalid transaction type",
		},
		{
			name: "invalid amount (zero)",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Amount:        floatPtr(0.0),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "transaction amount cannot be zero",
		},
		{
			name: "invalid description (too short)",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Description:   stringPtr("AB"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "invalid transaction description",
		},
		{
			name: "invalid date format",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Date:          stringPtr("invalid-date"),
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "invalid date format",
		},
		{
			name: "no fields provided for update",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWork) {
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = txRepo.Save(transaction)
			},
			wantError: true,
			errorMsg:  "at least one field must be provided",
		},
		{
			name: "account not found",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
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
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
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

			useCase := NewUpdateTransactionUseCase(mockUow, eventBus)
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

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
