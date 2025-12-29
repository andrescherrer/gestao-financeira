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

// mockTransactionRepository is a mock implementation of TransactionRepository for testing.
type mockTransactionRepository struct {
	transactions map[string]*entities.Transaction
	saveErr      error
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{
		transactions: make(map[string]*entities.Transaction),
	}
}

func (m *mockTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	transaction, exists := m.transactions[id.Value()]
	if !exists {
		return nil, nil
	}
	return transaction, nil
}

func (m *mockTransactionRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.AccountID().Value() == accountID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() && transaction.AccountID().Value() == accountID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() && transaction.TransactionType().Value() == transactionType.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) Save(transaction *entities.Transaction) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.transactions[transaction.ID().Value()] = transaction
	return nil
}

func (m *mockTransactionRepository) Delete(id transactionvalueobjects.TransactionID) error {
	delete(m.transactions, id.Value())
	return nil
}

func (m *mockTransactionRepository) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	_, exists := m.transactions[id.Value()]
	return exists, nil
}

func (m *mockTransactionRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() {
			count++
		}
	}
	return count, nil
}

func (m *mockTransactionRepository) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	count := int64(0)
	for _, transaction := range m.transactions {
		if transaction.AccountID().Value() == accountID.Value() {
			count++
		}
	}
	return count, nil
}

func (m *mockTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		// Simple mock: return all transactions (real implementation would check isRecurring flag)
		result = append(result, transaction)
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	// Simple mock: return nil (not found)
	return nil, nil
}

func (m *mockTransactionRepository) FindByUserIDAndFiltersWithPagination(
	userID identityvalueobjects.UserID,
	accountID string,
	transactionType string,
	offset, limit int,
) ([]*entities.Transaction, int64, error) {
	var result []*entities.Transaction
	var total int64

	for _, transaction := range m.transactions {
		if transaction.UserID().Value() != userID.Value() {
			continue
		}
		if accountID != "" && transaction.AccountID().Value() != accountID {
			continue
		}
		if transactionType != "" && transaction.TransactionType().Value() != transactionType {
			continue
		}
		total++
		if len(result) < limit && len(result) >= offset {
			result = append(result, transaction)
		}
	}

	return result, total, nil
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	eventBus := eventbus.NewEventBus()
	date := time.Now().Format("2006-01-02")

	tests := []struct {
		name      string
		input     dtos.CreateTransactionInput
		setupMock func(*mockTransactionRepository)
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid transaction creation",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: false,
		},
		{
			name: "invalid user ID",
			input: dtos.CreateTransactionInput{
				UserID:      "invalid-uuid",
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid user ID",
		},
		{
			name: "invalid account ID",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   "invalid-uuid",
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid account ID",
		},
		{
			name: "invalid transaction type",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INVALID",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid transaction type",
		},
		{
			name: "invalid currency",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "INVALID",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid currency",
		},
		{
			name: "invalid description (too short)",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "AB",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid transaction description",
		},
		{
			name: "invalid date format",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        "invalid-date",
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid date format",
		},
		{
			name: "zero amount",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      0,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {},
			wantError: true,
			errorMsg:  "transaction amount cannot be zero",
		},
		{
			name: "repository save error",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Compra de supermercado",
				Date:        date,
			},
			setupMock: func(m *mockTransactionRepository) {
				m.saveErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to save transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockTransactionRepository()
			tt.setupMock(mockRepo)

			useCase := NewCreateTransactionUseCase(mockRepo, eventBus)
			output, err := useCase.Execute(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateTransactionUseCase.Execute() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				if tt.errorMsg != "" && err != nil && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("CreateTransactionUseCase.Execute() error = %v, want error containing %v", err, tt.errorMsg)
				}
				if output != nil {
					t.Errorf("CreateTransactionUseCase.Execute() output = %v, want nil", output)
				}
			} else {
				if err != nil {
					t.Errorf("CreateTransactionUseCase.Execute() unexpected error = %v", err)
					return
				}
				if output == nil {
					t.Errorf("CreateTransactionUseCase.Execute() output = nil, want CreateTransactionOutput")
					return
				}
				if output.TransactionID == "" {
					t.Error("CreateTransactionUseCase.Execute() output.TransactionID is empty")
				}
				if output.UserID != tt.input.UserID {
					t.Errorf("CreateTransactionUseCase.Execute() output.UserID = %v, want %v", output.UserID, tt.input.UserID)
				}
				if output.AccountID != tt.input.AccountID {
					t.Errorf("CreateTransactionUseCase.Execute() output.AccountID = %v, want %v", output.AccountID, tt.input.AccountID)
				}
				if output.Type != tt.input.Type {
					t.Errorf("CreateTransactionUseCase.Execute() output.Type = %v, want %v", output.Type, tt.input.Type)
				}
				if output.Description != tt.input.Description {
					t.Errorf("CreateTransactionUseCase.Execute() output.Description = %v, want %v", output.Description, tt.input.Description)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
