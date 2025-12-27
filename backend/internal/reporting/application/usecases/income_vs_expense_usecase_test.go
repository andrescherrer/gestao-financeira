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

func TestIncomeVsExpenseUseCase_Execute(t *testing.T) {
	// Create test data
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	// Create transactions
	incomeAmount1, _ := sharedvalueobjects.NewMoney(100000, currency)  // R$ 1000.00
	incomeAmount2, _ := sharedvalueobjects.NewMoney(200000, currency)  // R$ 2000.00
	expenseAmount1, _ := sharedvalueobjects.NewMoney(50000, currency)  // R$ 500.00
	expenseAmount2, _ := sharedvalueobjects.NewMoney(150000, currency) // R$ 1500.00

	incomeTx1, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount1,
		transactionvalueobjects.MustTransactionDescription("Salary 1"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	incomeTx2, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount2,
		transactionvalueobjects.MustTransactionDescription("Salary 2"),
		time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
	)

	expenseTx1, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount1,
		transactionvalueobjects.MustTransactionDescription("Groceries"),
		time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
	)

	expenseTx2, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount2,
		transactionvalueobjects.MustTransactionDescription("Rent"),
		time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
	)

	// Create mock repository
	mockRepo := &mockTransactionRepository{
		transactions: []*entities.Transaction{incomeTx1, incomeTx2, expenseTx1, expenseTx2},
	}

	// Create use case
	useCase := NewIncomeVsExpenseUseCase(mockRepo)

	// Execute
	input := dtos.IncomeVsExpenseInput{
		UserID:   userID.Value(),
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

	// Total income should be 3000.00 (1000 + 2000)
	if output.TotalIncome != 3000.00 {
		t.Errorf("Execute() output.TotalIncome = %v, want 3000.00", output.TotalIncome)
	}

	// Total expense should be 2000.00 (500 + 1500)
	if output.TotalExpense != 2000.00 {
		t.Errorf("Execute() output.TotalExpense = %v, want 2000.00", output.TotalExpense)
	}

	// Balance should be 1000.00 (3000 - 2000)
	if output.Balance != 1000.00 {
		t.Errorf("Execute() output.Balance = %v, want 1000.00", output.Balance)
	}

	// Difference should be same as balance
	if output.Difference != 1000.00 {
		t.Errorf("Execute() output.Difference = %v, want 1000.00", output.Difference)
	}

	if output.IncomeCount != 2 {
		t.Errorf("Execute() output.IncomeCount = %v, want 2", output.IncomeCount)
	}

	if output.ExpenseCount != 2 {
		t.Errorf("Execute() output.ExpenseCount = %v, want 2", output.ExpenseCount)
	}

	if output.TotalCount != 4 {
		t.Errorf("Execute() output.TotalCount = %v, want 4", output.TotalCount)
	}
}

func TestIncomeVsExpenseUseCase_Execute_WithGroupBy(t *testing.T) {
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	incomeAmount, _ := sharedvalueobjects.NewMoney(100000, currency)
	expenseAmount, _ := sharedvalueobjects.NewMoney(50000, currency)

	// January transaction
	incomeTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount,
		transactionvalueobjects.MustTransactionDescription("January Salary"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	// February transaction
	expenseTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount,
		transactionvalueobjects.MustTransactionDescription("February Groceries"),
		time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
	)

	mockRepo := &mockTransactionRepository{
		transactions: []*entities.Transaction{incomeTx, expenseTx},
	}

	useCase := NewIncomeVsExpenseUseCase(mockRepo)

	// Execute with group_by month
	input := dtos.IncomeVsExpenseInput{
		UserID:   userID.Value(),
		Currency: "BRL",
		GroupBy:  "month",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	// Should have period breakdown
	if len(output.PeriodBreakdown) == 0 {
		t.Error("Execute() output.PeriodBreakdown should not be empty")
	}

	// Find January and February in breakdown
	var january, february *dtos.PeriodSummary
	for i := range output.PeriodBreakdown {
		if output.PeriodBreakdown[i].Period == "2025-01" {
			january = &output.PeriodBreakdown[i]
		} else if output.PeriodBreakdown[i].Period == "2025-02" {
			february = &output.PeriodBreakdown[i]
		}
	}

	if january == nil {
		t.Fatal("Execute() january period not found in breakdown")
	}
	if january.TotalIncome != 1000.00 {
		t.Errorf("Execute() january.TotalIncome = %v, want 1000.00", january.TotalIncome)
	}

	if february == nil {
		t.Fatal("Execute() february period not found in breakdown")
	}
	if february.TotalExpense != 500.00 {
		t.Errorf("Execute() february.TotalExpense = %v, want 500.00", february.TotalExpense)
	}
}

func TestIncomeVsExpenseUseCase_Execute_InvalidInput(t *testing.T) {
	mockRepo := &mockTransactionRepository{transactions: []*entities.Transaction{}}
	useCase := NewIncomeVsExpenseUseCase(mockRepo)

	tests := []struct {
		name  string
		input dtos.IncomeVsExpenseInput
	}{
		{
			name: "invalid user ID",
			input: dtos.IncomeVsExpenseInput{
				UserID: "invalid",
			},
		},
		{
			name: "invalid group_by",
			input: dtos.IncomeVsExpenseInput{
				UserID:  "123e4567-e89b-12d3-a456-426614174000",
				GroupBy: "invalid",
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
