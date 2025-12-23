package usecases

import (
	"errors"
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// Extended mock repository with error support for GetTransactionUseCase
type mockGetTransactionRepository struct {
	*mockTransactionRepository
	findByIDErr error
}

func newMockGetTransactionRepository() *mockGetTransactionRepository {
	return &mockGetTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
	}
}

func (m *mockGetTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.mockTransactionRepository.FindByID(id)
}

func TestGetTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	tests := []struct {
		name      string
		input     dtos.GetTransactionInput
		setupMock func(*mockGetTransactionRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.GetTransactionOutput)
	}{
		{
			name: "get existing transaction",
			input: dtos.GetTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockGetTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.50, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.GetTransactionOutput) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				if output.TransactionID == "" {
					t.Errorf("expected transaction ID to be set")
				}
				if output.UserID != userID.Value() {
					t.Errorf("expected user ID %s, got %s", userID.Value(), output.UserID)
				}
				if output.AccountID != accountID.Value() {
					t.Errorf("expected account ID %s, got %s", accountID.Value(), output.AccountID)
				}
				if output.Type != "INCOME" {
					t.Errorf("expected type 'INCOME', got %s", output.Type)
				}
				if output.Amount != 100.50 {
					t.Errorf("expected amount 100.50, got %f", output.Amount)
				}
				if output.Currency != "BRL" {
					t.Errorf("expected currency 'BRL', got %s", output.Currency)
				}
				if output.Description != "Salário" {
					t.Errorf("expected description 'Salário', got %s", output.Description)
				}
			},
		},
		{
			name: "transaction not found",
			input: dtos.GetTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(m *mockGetTransactionRepository) {
				// No transactions in repository
			},
			wantError: true,
			errorMsg:  "transaction not found",
		},
		{
			name: "invalid transaction ID",
			input: dtos.GetTransactionInput{
				TransactionID: "invalid-uuid",
			},
			setupMock: func(m *mockGetTransactionRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "repository error",
			input: dtos.GetTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(m *mockGetTransactionRepository) {
				m.findByIDErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to find transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockGetTransactionRepository()
			tt.setupMock(mockRepo)

			// For the first test case, we need to set the TransactionID after creating the transaction
			if tt.name == "get existing transaction" {
				// Create transaction first
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.50, "BRL", "Salário", date)
				_ = mockRepo.Save(transaction)
				tt.input.TransactionID = transaction.ID().Value()
			}

			useCase := NewGetTransactionUseCase(mockRepo)
			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
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

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return contains(s, substr)
}
