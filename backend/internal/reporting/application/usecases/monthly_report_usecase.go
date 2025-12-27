package usecases

import (
	"fmt"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// MonthlyReportUseCase handles generating monthly financial reports.
type MonthlyReportUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewMonthlyReportUseCase creates a new MonthlyReportUseCase instance.
func NewMonthlyReportUseCase(
	transactionRepository repositories.TransactionRepository,
) *MonthlyReportUseCase {
	return &MonthlyReportUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute generates a monthly report for the specified user, year, and month.
func (uc *MonthlyReportUseCase) Execute(input dtos.MonthlyReportInput) (*dtos.MonthlyReportOutput, error) {
	// Validate user ID
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Validate year and month
	if input.Year < 2000 || input.Year > 2100 {
		return nil, fmt.Errorf("invalid year: must be between 2000 and 2100")
	}
	if input.Month < 1 || input.Month > 12 {
		return nil, fmt.Errorf("invalid month: must be between 1 and 12")
	}

	// Create date range for the month
	startDate := time.Date(input.Year, time.Month(input.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond) // Last moment of the month

	// Get all transactions for the user
	allTransactions, err := uc.transactionRepository.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}

	// Filter transactions by date range and currency (if specified)
	type transactionData struct {
		Type     string
		Amount   int64
		Currency string
	}

	var filteredTransactions []transactionData

	for _, tx := range allTransactions {
		txDate := tx.Date()

		// Check if transaction is within the month
		if txDate.Before(startDate) || txDate.After(endDate) {
			continue
		}

		// Filter by currency if specified
		if input.Currency != "" {
			currency, err := sharedvalueobjects.NewCurrency(input.Currency)
			if err == nil && !tx.Amount().Currency().Equals(currency) {
				continue
			}
		}

		txType := tx.TransactionType()
		amount := tx.Amount()

		filteredTransactions = append(filteredTransactions, transactionData{
			Type:     txType.Value(),
			Amount:   amount.Amount(),
			Currency: amount.Currency().Code(),
		})
	}

	// Calculate totals
	var totalIncomeCents int64 = 0
	var totalExpenseCents int64 = 0
	incomeCount := 0
	expenseCount := 0

	currency := input.Currency
	if currency == "" {
		// Use the first transaction's currency if not specified
		if len(filteredTransactions) > 0 {
			currency = filteredTransactions[0].Currency
		} else {
			currency = "BRL" // Default
		}
	}

	currencyVO, err := sharedvalueobjects.NewCurrency(currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	for _, tx := range filteredTransactions {
		// Only count transactions with matching currency
		if tx.Currency != currency {
			continue
		}

		if tx.Type == transactionvalueobjects.Income {
			totalIncomeCents += tx.Amount
			incomeCount++
		} else if tx.Type == transactionvalueobjects.Expense {
			totalExpenseCents += tx.Amount
			expenseCount++
		}
	}

	// Convert cents to float64
	totalIncome, _ := sharedvalueobjects.NewMoney(totalIncomeCents, currencyVO)
	totalExpense, _ := sharedvalueobjects.NewMoney(totalExpenseCents, currencyVO)
	balance, _ := totalIncome.Subtract(totalExpense)

	// Build output
	output := &dtos.MonthlyReportOutput{
		UserID:       input.UserID,
		Year:         input.Year,
		Month:        input.Month,
		Currency:     currency,
		TotalIncome:  totalIncome.Float64(),
		TotalExpense: totalExpense.Float64(),
		Balance:      balance.Float64(),
		IncomeCount:  incomeCount,
		ExpenseCount: expenseCount,
		TotalCount:   incomeCount + expenseCount,
	}

	return output, nil
}
