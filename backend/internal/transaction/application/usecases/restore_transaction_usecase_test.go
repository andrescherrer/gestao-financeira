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

// mockRestoreTransactionRepository is a mock implementation for restore tests
type mockRestoreTransactionRepository struct {
	*mockTransactionRepository
	restoreErr  error
	softDeleted map[string]bool // Track soft-deleted transactions
}

func newMockRestoreTransactionRepository() *mockRestoreTransactionRepository {
	return &mockRestoreTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
		softDeleted:               make(map[string]bool),
	}
}

func (m *mockRestoreTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	// If soft-deleted, return nil (simulating GORM behavior)
	if m.softDeleted[id.Value()] {
		return nil, nil
	}
	return m.mockTransactionRepository.FindByID(id)
}

func (m *mockRestoreTransactionRepository) Restore(id transactionvalueobjects.TransactionID) error {
	if m.restoreErr != nil {
		return m.restoreErr
	}
	if !m.softDeleted[id.Value()] {
		return errors.New("transaction is not deleted")
	}
	delete(m.softDeleted, id.Value())
	return nil
}

func (m *mockRestoreTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	return m.mockTransactionRepository.FindActiveRecurringTransactions()
}

func (m *mockRestoreTransactionRepository) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	return m.mockTransactionRepository.FindByParentIDAndDate(parentID, date)
}

func (m *mockRestoreTransactionRepository) FindByUserIDAndFiltersWithPagination(
	userID identityvalueobjects.UserID,
	accountID string,
	transactionType string,
	offset, limit int,
) ([]*entities.Transaction, int64, error) {
	return m.mockTransactionRepository.FindByUserIDAndFiltersWithPagination(userID, accountID, transactionType, offset, limit)
}

func TestRestoreTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	tests := []struct {
		name      string
		input     dtos.RestoreTransactionInput
		setupMock func(*mockRestoreTransactionRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.RestoreTransactionOutput)
	}{
		{
			name: "restore soft-deleted transaction",
			input: dtos.RestoreTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockRestoreTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				// Simulate soft delete
				m.softDeleted[transaction.ID().Value()] = true
				_ = m.Delete(transaction.ID())
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.RestoreTransactionOutput) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				if output.TransactionID == "" {
					t.Errorf("expected transaction ID to be set")
				}
				if output.Message != "Transaction restored successfully" {
					t.Errorf("expected message 'Transaction restored successfully', got %s", output.Message)
				}
			},
		},
		{
			name: "transaction not found",
			input: dtos.RestoreTransactionInput{
				TransactionID: transactionvalueobjects.GenerateTransactionID().Value(),
			},
			setupMock: func(m *mockRestoreTransactionRepository) {
				// No transactions in repository
			},
			wantError: true,
			errorMsg:  "failed to restore transaction",
		},
		{
			name: "transaction is not deleted",
			input: dtos.RestoreTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockRestoreTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				// Transaction is not soft-deleted
			},
			wantError: true,
			errorMsg:  "transaction is not deleted",
		},
		{
			name: "invalid transaction ID",
			input: dtos.RestoreTransactionInput{
				TransactionID: "invalid-uuid",
			},
			setupMock: func(m *mockRestoreTransactionRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid transaction ID",
		},
		{
			name: "repository restore error",
			input: dtos.RestoreTransactionInput{
				TransactionID: "",
			},
			setupMock: func(m *mockRestoreTransactionRepository) {
				transaction, _ := createTestTransaction(userID, accountID, "INCOME", 100.00, "BRL", "Salário", date)
				_ = m.Save(transaction)
				m.softDeleted[transaction.ID().Value()] = true
				m.restoreErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to restore transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockRestoreTransactionRepository()
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

			useCase := NewRestoreTransactionUseCase(mockRepo)
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
