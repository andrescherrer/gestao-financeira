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

// Extended mock repository with error support for DeleteTransactionUseCase
type mockDeleteTransactionRepository struct {
	*mockTransactionRepository
	findByIDErr error
	deleteErr   error
}

func newMockDeleteTransactionRepository() *mockDeleteTransactionRepository {
	return &mockDeleteTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
	}
}

func (m *mockDeleteTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.mockTransactionRepository.FindByID(id)
}

func (m *mockDeleteTransactionRepository) Delete(id transactionvalueobjects.TransactionID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return m.mockTransactionRepository.Delete(id)
}

func TestDeleteTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	tests := []struct {
		name      string
		input     dtos.DeleteTransactionInput
		setupMock func(*mockDeleteTransactionRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.DeleteTransactionOutput)
	}{
		{
			name: "delete existing transaction",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockDeleteTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.DeleteTransactionOutput) {
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
			},
		},
		{
			name: "transaction not found",
			input: dtos.DeleteTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(m *mockDeleteTransactionRepository) {
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
			setupMock: func(m *mockDeleteTransactionRepository) {
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
			setupMock: func(m *mockDeleteTransactionRepository) {
				m.findByIDErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to find transaction",
		},
		{
			name: "repository delete error",
			input: dtos.DeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockDeleteTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				m.deleteErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to delete transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockDeleteTransactionRepository()
			tt.setupMock(mockRepo)

			// For test cases that need TransactionID, set it after creating the transaction
			if tt.input.TransactionID == "" {
				if tt.name == "transaction not found" {
					// For "transaction not found" test, use a valid UUID format
					tt.input.TransactionID = transactionvalueobjects.GenerateTransactionID().Value()
				} else {
					// For other tests, get the transaction ID from the repository
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

			useCase := NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
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
