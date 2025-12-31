package usecases

import (
	"errors"
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// Helper function to create test transaction
func createTestTransaction(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID, transactionType string, amount float64, currency string, description string, date time.Time) (*entities.Transaction, error) {
	txType, err := transactionvalueobjects.NewTransactionType(transactionType)
	if err != nil {
		return nil, err
	}

	curr, err := sharedvalueobjects.NewCurrency(currency)
	if err != nil {
		return nil, err
	}

	amountCents := int64(amount * 100)
	money, err := sharedvalueobjects.NewMoney(amountCents, curr)
	if err != nil {
		return nil, err
	}

	desc, err := transactionvalueobjects.NewTransactionDescription(description)
	if err != nil {
		return nil, err
	}

	return entities.NewTransaction(userID, accountID, txType, money, desc, date)
}

// Extended mock repository with error support
type mockListTransactionRepository struct {
	*mockTransactionRepository
	findByUserIDErr             error
	findByAccountIDErr          error
	findByUserIDAndAccountIDErr error
	findByUserIDAndTypeErr      error
}

func newMockListTransactionRepository() *mockListTransactionRepository {
	return &mockListTransactionRepository{
		mockTransactionRepository: newMockTransactionRepository(),
	}
}

func (m *mockListTransactionRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	if m.findByUserIDErr != nil {
		return nil, m.findByUserIDErr
	}
	return m.mockTransactionRepository.FindByUserID(userID)
}

func (m *mockListTransactionRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	if m.findByAccountIDErr != nil {
		return nil, m.findByAccountIDErr
	}
	return m.mockTransactionRepository.FindByAccountID(accountID)
}

func (m *mockListTransactionRepository) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	if m.findByUserIDAndAccountIDErr != nil {
		return nil, m.findByUserIDAndAccountIDErr
	}
	return m.mockTransactionRepository.FindByUserIDAndAccountID(userID, accountID)
}

func (m *mockListTransactionRepository) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	if m.findByUserIDAndTypeErr != nil {
		return nil, m.findByUserIDAndTypeErr
	}
	return m.mockTransactionRepository.FindByUserIDAndType(userID, transactionType)
}
func (m *mockListTransactionRepository) FindByUserIDAndDateRange(userID identityvalueobjects.UserID, startDate, endDate time.Time) ([]*entities.Transaction, error) {
	return m.mockTransactionRepository.FindByUserIDAndDateRange(userID, startDate, endDate)
}
func (m *mockListTransactionRepository) FindByUserIDAndDateRangeWithCurrency(userID identityvalueobjects.UserID, startDate, endDate time.Time, currency string) ([]*entities.Transaction, error) {
	return m.mockTransactionRepository.FindByUserIDAndDateRangeWithCurrency(userID, startDate, endDate, currency)
}
func (m *mockListTransactionRepository) FindByUserIDWithPagination(userID identityvalueobjects.UserID, offset, limit int) ([]*entities.Transaction, int64, error) {
	return m.mockTransactionRepository.FindByUserIDWithPagination(userID, offset, limit)
}

func TestListTransactionsUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()
	accountID1 := accountvalueobjects.GenerateAccountID()
	accountID2 := accountvalueobjects.GenerateAccountID()
	date := time.Now()

	tests := []struct {
		name      string
		input     dtos.ListTransactionsInput
		setupMock func(*mockListTransactionRepository)
		wantError bool
		errorMsg  string
		wantCount int
	}{
		{
			name: "list all transactions for user",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: "",
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {
				// Create test transactions
				tx1, _ := createTestTransaction(userID, accountID1, "INCOME", 100.00, "BRL", "Salário", date)
				tx2, _ := createTestTransaction(userID, accountID1, "EXPENSE", 50.00, "BRL", "Compra", date)
				tx3, _ := createTestTransaction(userID, accountID2, "INCOME", 200.00, "BRL", "Venda", date)
				_ = m.Save(tx1)
				_ = m.Save(tx2)
				_ = m.Save(tx3)
			},
			wantError: false,
			wantCount: 3,
		},
		{
			name: "list transactions filtered by account",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: accountID1.Value(),
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {
				// Create test transactions
				tx1, _ := createTestTransaction(userID, accountID1, "INCOME", 100.00, "BRL", "Salário", date)
				tx2, _ := createTestTransaction(userID, accountID1, "EXPENSE", 50.00, "BRL", "Compra", date)
				tx3, _ := createTestTransaction(userID, accountID2, "INCOME", 200.00, "BRL", "Venda", date)
				_ = m.Save(tx1)
				_ = m.Save(tx2)
				_ = m.Save(tx3)
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name: "list transactions filtered by type",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: "",
				Type:      "INCOME",
			},
			setupMock: func(m *mockListTransactionRepository) {
				// Create test transactions
				tx1, _ := createTestTransaction(userID, accountID1, "INCOME", 100.00, "BRL", "Salário", date)
				tx2, _ := createTestTransaction(userID, accountID1, "EXPENSE", 50.00, "BRL", "Compra", date)
				tx3, _ := createTestTransaction(userID, accountID2, "INCOME", 200.00, "BRL", "Venda", date)
				_ = m.Save(tx1)
				_ = m.Save(tx2)
				_ = m.Save(tx3)
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name: "list transactions filtered by account and type",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: accountID1.Value(),
				Type:      "INCOME",
			},
			setupMock: func(m *mockListTransactionRepository) {
				// Create test transactions
				tx1, _ := createTestTransaction(userID, accountID1, "INCOME", 100.00, "BRL", "Salário", date)
				tx2, _ := createTestTransaction(userID, accountID1, "EXPENSE", 50.00, "BRL", "Compra", date)
				tx3, _ := createTestTransaction(userID, accountID2, "INCOME", 200.00, "BRL", "Venda", date)
				_ = m.Save(tx1)
				_ = m.Save(tx2)
				_ = m.Save(tx3)
			},
			wantError: false,
			wantCount: 1,
		},
		{
			name: "list transactions for user with no transactions",
			input: dtos.ListTransactionsInput{
				UserID:    userID2.Value(),
				AccountID: "",
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {
				// No transactions for this user
			},
			wantError: false,
			wantCount: 0,
		},
		{
			name: "invalid user ID",
			input: dtos.ListTransactionsInput{
				UserID:    "invalid-uuid",
				AccountID: "",
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid user ID",
		},
		{
			name: "invalid account ID",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: "invalid-uuid",
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid account ID",
		},
		{
			name: "invalid transaction type",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: "",
				Type:      "INVALID",
			},
			setupMock: func(m *mockListTransactionRepository) {},
			wantError: true,
			errorMsg:  "invalid transaction type",
		},
		{
			name: "repository error",
			input: dtos.ListTransactionsInput{
				UserID:    userID.Value(),
				AccountID: "",
				Type:      "",
			},
			setupMock: func(m *mockListTransactionRepository) {
				m.findByUserIDErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to find transactions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockListTransactionRepository()
			tt.setupMock(mockRepo)

			useCase := NewListTransactionsUseCase(mockRepo)
			output, err := useCase.Execute(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("ListTransactionsUseCase.Execute() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				if tt.errorMsg != "" && err != nil && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("ListTransactionsUseCase.Execute() error = %v, want error containing %v", err, tt.errorMsg)
				}
				if output != nil {
					t.Errorf("ListTransactionsUseCase.Execute() output = %v, want nil", output)
				}
			} else {
				if err != nil {
					t.Errorf("ListTransactionsUseCase.Execute() unexpected error = %v", err)
					return
				}
				if output == nil {
					t.Errorf("ListTransactionsUseCase.Execute() output = nil, want ListTransactionsOutput")
					return
				}
				if output.Count != tt.wantCount {
					t.Errorf("ListTransactionsUseCase.Execute() output.Count = %v, want %v", output.Count, tt.wantCount)
				}
				if len(output.Transactions) != tt.wantCount {
					t.Errorf("ListTransactionsUseCase.Execute() len(output.Transactions) = %v, want %v", len(output.Transactions), tt.wantCount)
				}
			}
		})
	}
}
