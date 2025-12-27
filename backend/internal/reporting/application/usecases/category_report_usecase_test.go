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

func TestCategoryReportUseCase_Execute(t *testing.T) {
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
	useCase := NewCategoryReportUseCase(mockRepo)

	// Execute
	input := dtos.CategoryReportInput{
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

	if output.TotalCount != 4 {
		t.Errorf("Execute() output.TotalCount = %v, want 4", output.TotalCount)
	}

	// Validate category breakdown
	if len(output.CategoryBreakdown) != 2 {
		t.Errorf("Execute() output.CategoryBreakdown length = %v, want 2", len(output.CategoryBreakdown))
	}

	// Find income and expense summaries
	var incomeSummary, expenseSummary *dtos.CategorySummary
	for i := range output.CategoryBreakdown {
		if output.CategoryBreakdown[i].Type == "INCOME" {
			incomeSummary = &output.CategoryBreakdown[i]
		} else if output.CategoryBreakdown[i].Type == "EXPENSE" {
			expenseSummary = &output.CategoryBreakdown[i]
		}
	}

	if incomeSummary == nil {
		t.Fatal("Execute() income summary not found")
	}
	if incomeSummary.TotalAmount != 3000.00 {
		t.Errorf("Execute() incomeSummary.TotalAmount = %v, want 3000.00", incomeSummary.TotalAmount)
	}
	if incomeSummary.Count != 2 {
		t.Errorf("Execute() incomeSummary.Count = %v, want 2", incomeSummary.Count)
	}

	if expenseSummary == nil {
		t.Fatal("Execute() expense summary not found")
	}
	if expenseSummary.TotalAmount != 2000.00 {
		t.Errorf("Execute() expenseSummary.TotalAmount = %v, want 2000.00", expenseSummary.TotalAmount)
	}
	if expenseSummary.Count != 2 {
		t.Errorf("Execute() expenseSummary.Count = %v, want 2", expenseSummary.Count)
	}
}

func TestCategoryReportUseCase_Execute_WithDateFilter(t *testing.T) {
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	incomeAmount, _ := sharedvalueobjects.NewMoney(100000, currency)
	expenseAmount, _ := sharedvalueobjects.NewMoney(50000, currency)

	// Transaction in January
	incomeTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount,
		transactionvalueobjects.MustTransactionDescription("January Salary"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	// Transaction in February
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

	useCase := NewCategoryReportUseCase(mockRepo)

	// Filter only January
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	input := dtos.CategoryReportInput{
		UserID:    userID.Value(),
		StartDate: &startDate,
		EndDate:   &endDate,
		Currency:  "BRL",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	// Should only have income (January transaction)
	if output.TotalIncome != 1000.00 {
		t.Errorf("Execute() output.TotalIncome = %v, want 1000.00", output.TotalIncome)
	}

	if output.TotalExpense != 0.00 {
		t.Errorf("Execute() output.TotalExpense = %v, want 0.00", output.TotalExpense)
	}
}

func TestCategoryReportUseCase_Execute_InvalidInput(t *testing.T) {
	mockRepo := &mockTransactionRepository{transactions: []*entities.Transaction{}}
	useCase := NewCategoryReportUseCase(mockRepo)

	input := dtos.CategoryReportInput{
		UserID: "invalid",
	}

	_, err := useCase.Execute(input)
	if err == nil {
		t.Errorf("Execute() error = nil, want error")
	}
}
