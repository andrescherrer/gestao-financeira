package usecases

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// mockTransactionRepository is a mock implementation of TransactionRepository for testing.
type mockTransactionRepository struct {
	transactions []*entities.Transaction
}

func (m *mockTransactionRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// Implement other required methods with empty implementations
func (m *mockTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepository) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepository) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepository) Save(transaction *entities.Transaction) error {
	return nil
}
func (m *mockTransactionRepository) Delete(id transactionvalueobjects.TransactionID) error {
	return nil
}
func (m *mockTransactionRepository) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	return false, nil
}
func (m *mockTransactionRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	return 0, nil
}
func (m *mockTransactionRepository) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	return 0, nil
}
func (m *mockTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepository) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	return nil, nil
}

func TestMonthlyReportUseCase_Execute(t *testing.T) {
	// Create test data
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	// Create transactions for January 2025
	incomeAmount, _ := sharedvalueobjects.NewMoney(100000, currency) // R$ 1000.00
	expenseAmount, _ := sharedvalueobjects.NewMoney(50000, currency) // R$ 500.00

	incomeTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount,
		transactionvalueobjects.MustTransactionDescription("Salary"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	expenseTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount,
		transactionvalueobjects.MustTransactionDescription("Groceries"),
		time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
	)

	// Create mock repository
	mockRepo := &mockTransactionRepository{
		transactions: []*entities.Transaction{incomeTx, expenseTx},
	}

	// Create use case
	useCase := NewMonthlyReportUseCase(mockRepo)

	// Execute
	input := dtos.MonthlyReportInput{
		UserID:   userID.Value(),
		Year:     2025,
		Month:    1,
		Currency: "BRL",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	// Validate output
	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if output.UserID != userID.Value() {
		t.Errorf("Execute() output.UserID = %v, want %v", output.UserID, userID.Value())
	}

	if output.Year != 2025 {
		t.Errorf("Execute() output.Year = %v, want 2025", output.Year)
	}

	if output.Month != 1 {
		t.Errorf("Execute() output.Month = %v, want 1", output.Month)
	}

	if output.TotalIncome != 1000.00 {
		t.Errorf("Execute() output.TotalIncome = %v, want 1000.00", output.TotalIncome)
	}

	if output.TotalExpense != 500.00 {
		t.Errorf("Execute() output.TotalExpense = %v, want 500.00", output.TotalExpense)
	}

	if output.Balance != 500.00 {
		t.Errorf("Execute() output.Balance = %v, want 500.00", output.Balance)
	}

	if output.IncomeCount != 1 {
		t.Errorf("Execute() output.IncomeCount = %v, want 1", output.IncomeCount)
	}

	if output.ExpenseCount != 1 {
		t.Errorf("Execute() output.ExpenseCount = %v, want 1", output.ExpenseCount)
	}

	if output.TotalCount != 2 {
		t.Errorf("Execute() output.TotalCount = %v, want 2", output.TotalCount)
	}
}

func TestMonthlyReportUseCase_Execute_InvalidInput(t *testing.T) {
	mockRepo := &mockTransactionRepository{transactions: []*entities.Transaction{}}
	useCase := NewMonthlyReportUseCase(mockRepo)

	tests := []struct {
		name  string
		input dtos.MonthlyReportInput
	}{
		{
			name: "invalid user ID",
			input: dtos.MonthlyReportInput{
				UserID: "invalid",
				Year:   2025,
				Month:  1,
			},
		},
		{
			name: "invalid year",
			input: dtos.MonthlyReportInput{
				UserID: "123e4567-e89b-12d3-a456-426614174000",
				Year:   1999,
				Month:  1,
			},
		},
		{
			name: "invalid month",
			input: dtos.MonthlyReportInput{
				UserID: "123e4567-e89b-12d3-a456-426614174000",
				Year:   2025,
				Month:  13,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase.Execute(tt.input)
			if err == nil {
				t.Errorf("Execute() error = nil, want error")
			}
		})
	}
}
