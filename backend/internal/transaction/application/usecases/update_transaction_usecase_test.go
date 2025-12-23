package usecases

import (
	"errors"
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// Extended mock repository with error support for UpdateTransactionUseCase
type mockUpdateTransactionRepository struct {
	*mockTransactionRepository
	findByIDErr error
	saveErr     error
}

func newMockUpdateTransactionRepository() *mockUpdateTransactionRepository {
	return &mockUpdateTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
	}
}

func (m *mockUpdateTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.mockTransactionRepository.FindByID(id)
}

func (m *mockUpdateTransactionRepository) Save(transaction *entities.Transaction) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	return m.mockTransactionRepository.Save(transaction)
}

func TestUpdateTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	eventBus := eventbus.NewEventBus()
	date := time.Now()
	newDateStr := date.Add(24 * time.Hour).Format("2006-01-02")

	tests := []struct {
		name      string
		input     dtos.UpdateTransactionInput
		setupMock func(*mockUpdateTransactionRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.UpdateTransactionOutput)
	}{
		{
			name: "update transaction type",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput) {
				if output.Type != "EXPENSE" {
					t.Errorf("expected type 'EXPENSE', got %s", output.Type)
				}
			},
		},
		{
			name: "update transaction amount",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Amount:        floatPtr(200.50),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput) {
				if output.Amount != 200.50 {
					t.Errorf("expected amount 200.50, got %f", output.Amount)
				}
			},
		},
		{
			name: "update transaction description",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Description:   stringPtr("Nova descrição"),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput) {
				if output.Description != "Nova descrição" {
					t.Errorf("expected description 'Nova descrição', got %s", output.Description)
				}
			},
		},
		{
			name: "update transaction date",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Date:          stringPtr(newDateStr),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput) {
				if output.Date != newDateStr {
					t.Errorf("expected date %s, got %s", newDateStr, output.Date)
				}
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
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.UpdateTransactionOutput) {
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
			setupMock: func(m *mockUpdateTransactionRepository) {
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
			setupMock: func(m *mockUpdateTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "invalid transaction type",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("INVALID"),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
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
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
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
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
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
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: true,
			errorMsg:  "invalid date format",
		},
		{
			name: "no fields provided for update",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: true,
			errorMsg:  "at least one field must be provided",
		},
		{
			name: "repository save error",
			input: dtos.UpdateTransactionInput{
				TransactionID: "",
				Type:          stringPtr("EXPENSE"),
			},
			setupMock: func(m *mockUpdateTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				// Set saveErr after transaction is saved, so FindByID works but Save fails
				m.saveErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to save transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUpdateTransactionRepository()
			tt.setupMock(mockRepo)

			// For test cases that need TransactionID, set it after creating the transaction
			if tt.input.TransactionID == "" {
				if tt.name == "transaction not found" {
					// For "transaction not found" test, use a valid UUID format
					tt.input.TransactionID = transactionvalueobjects.GenerateTransactionID().Value()
				} else {
					// For other tests, get the transaction ID from the repository
					// (transaction was already created in setupMock)
					var transactionID string
					for id := range mockRepo.transactions {
						transactionID = id
						break
					}
					if transactionID == "" {
						// If no transaction in repo, create one
						transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
						_ = mockRepo.Save(transaction)
						transactionID = transaction.ID().Value()
					}
					tt.input.TransactionID = transactionID
				}
			}

			useCase := NewUpdateTransactionUseCase(mockRepo, eventBus)
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
				tt.validate(t, output)
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
