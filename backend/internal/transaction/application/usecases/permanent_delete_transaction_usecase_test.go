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

// mockPermanentDeleteTransactionRepository is a mock implementation for permanent delete tests
type mockPermanentDeleteTransactionRepository struct {
	*mockTransactionRepository
	permanentDeleteErr error
}

func newMockPermanentDeleteTransactionRepository() *mockPermanentDeleteTransactionRepository {
	return &mockPermanentDeleteTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
	}
}

func (m *mockPermanentDeleteTransactionRepository) PermanentDelete(id transactionvalueobjects.TransactionID) error {
	if m.permanentDeleteErr != nil {
		return m.permanentDeleteErr
	}
	// Remove from transactions map (hard delete)
	delete(m.transactions, id.Value())
	return nil
}

func (m *mockPermanentDeleteTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	return m.mockTransactionRepository.FindActiveRecurringTransactions()
}

func (m *mockPermanentDeleteTransactionRepository) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	return m.mockTransactionRepository.FindByParentIDAndDate(parentID, date)
}

func (m *mockPermanentDeleteTransactionRepository) FindByUserIDAndFiltersWithPagination(
	userID identityvalueobjects.UserID,
	accountID string,
	transactionType string,
	offset, limit int,
) ([]*entities.Transaction, int64, error) {
	return m.mockTransactionRepository.FindByUserIDAndFiltersWithPagination(userID, accountID, transactionType, offset, limit)
}

func TestPermanentDeleteTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	tests := []struct {
		name      string
		input     dtos.PermanentDeleteTransactionInput
		setupMock func(*mockPermanentDeleteTransactionRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.PermanentDeleteTransactionOutput)
	}{
		{
			name: "permanently delete existing transaction",
			input: dtos.PermanentDeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockPermanentDeleteTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.PermanentDeleteTransactionOutput) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				if output.TransactionID == "" {
					t.Errorf("expected transaction ID to be set")
				}
				if output.Message != "Transaction permanently deleted successfully" {
					t.Errorf("expected message 'Transaction permanently deleted successfully', got %s", output.Message)
				}
			},
		},
		{
			name: "transaction not found",
			input: dtos.PermanentDeleteTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(m *mockPermanentDeleteTransactionRepository) {
				// No transactions in repository
			},
			wantError: true,
			errorMsg:  "transaction not found",
		},
		{
			name: "invalid transaction ID",
			input: dtos.PermanentDeleteTransactionInput{
				TransactionID: "invalid-uuid",
			},
			setupMock: func(m *mockPermanentDeleteTransactionRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "repository permanent delete error",
			input: dtos.PermanentDeleteTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockPermanentDeleteTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				m.permanentDeleteErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to permanently delete transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockPermanentDeleteTransactionRepository()
			tt.setupMock(mockRepo)

			// For test cases that need TransactionID, set it after creating the transaction
			if tt.input.TransactionID == "" {
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

			useCase := NewPermanentDeleteTransactionUseCase(mockRepo)
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
