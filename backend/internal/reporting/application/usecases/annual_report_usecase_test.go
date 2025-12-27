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

func TestAnnualReportUseCase_Execute(t *testing.T) {
	// Create test data
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	// Create transactions for 2025 (January and February)
	incomeAmount1, _ := sharedvalueobjects.NewMoney(100000, currency)  // R$ 1000.00
	expenseAmount1, _ := sharedvalueobjects.NewMoney(50000, currency)  // R$ 500.00
	incomeAmount2, _ := sharedvalueobjects.NewMoney(200000, currency)  // R$ 2000.00
	expenseAmount2, _ := sharedvalueobjects.NewMoney(100000, currency) // R$ 1000.00

	incomeTx1, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount1,
		transactionvalueobjects.MustTransactionDescription("January Salary"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	expenseTx1, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount1,
		transactionvalueobjects.MustTransactionDescription("January Groceries"),
		time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
	)

	incomeTx2, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount2,
		transactionvalueobjects.MustTransactionDescription("February Salary"),
		time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
	)

	expenseTx2, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount2,
		transactionvalueobjects.MustTransactionDescription("February Groceries"),
		time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
	)

	// Create mock repository
	mockRepo := &mockTransactionRepository{
		transactions: []*entities.Transaction{incomeTx1, expenseTx1, incomeTx2, expenseTx2},
	}

	// Create use case
	useCase := NewAnnualReportUseCase(mockRepo)

	// Execute
	input := dtos.AnnualReportInput{
		UserID:   userID.Value(),
		Year:     2025,
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

	// Total income should be 3000.00 (1000 + 2000)
	if output.TotalIncome != 3000.00 {
		t.Errorf("Execute() output.TotalIncome = %v, want 3000.00", output.TotalIncome)
	}

	// Total expense should be 1500.00 (500 + 1000)
	if output.TotalExpense != 1500.00 {
		t.Errorf("Execute() output.TotalExpense = %v, want 1500.00", output.TotalExpense)
	}

	// Balance should be 1500.00 (3000 - 1500)
	if output.Balance != 1500.00 {
		t.Errorf("Execute() output.Balance = %v, want 1500.00", output.Balance)
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

	// Validate monthly breakdown
	if len(output.MonthlyBreakdown) != 12 {
		t.Errorf("Execute() output.MonthlyBreakdown length = %v, want 12", len(output.MonthlyBreakdown))
	}

	// Check January (month 1)
	january := output.MonthlyBreakdown[0]
	if january.Month != 1 {
		t.Errorf("Execute() january.Month = %v, want 1", january.Month)
	}
	if january.TotalIncome != 1000.00 {
		t.Errorf("Execute() january.TotalIncome = %v, want 1000.00", january.TotalIncome)
	}
	if january.TotalExpense != 500.00 {
		t.Errorf("Execute() january.TotalExpense = %v, want 500.00", january.TotalExpense)
	}

	// Check February (month 2)
	february := output.MonthlyBreakdown[1]
	if february.Month != 2 {
		t.Errorf("Execute() february.Month = %v, want 2", february.Month)
	}
	if february.TotalIncome != 2000.00 {
		t.Errorf("Execute() february.TotalIncome = %v, want 2000.00", february.TotalIncome)
	}
	if february.TotalExpense != 1000.00 {
		t.Errorf("Execute() february.TotalExpense = %v, want 1000.00", february.TotalExpense)
	}
}

func TestAnnualReportUseCase_Execute_InvalidInput(t *testing.T) {
	mockRepo := &mockTransactionRepository{transactions: []*entities.Transaction{}}
	useCase := NewAnnualReportUseCase(mockRepo)

	tests := []struct {
		name  string
		input dtos.AnnualReportInput
	}{
		{
			name: "invalid user ID",
			input: dtos.AnnualReportInput{
				UserID: "invalid",
				Year:   2025,
			},
		},
		{
			name: "invalid year",
			input: dtos.AnnualReportInput{
				UserID: "123e4567-e89b-12d3-a456-426614174000",
				Year:   1999,
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
