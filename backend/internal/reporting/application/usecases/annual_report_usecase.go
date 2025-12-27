package usecases

import (
	"fmt"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
)

// AnnualReportUseCase handles generating annual financial reports.
type AnnualReportUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewAnnualReportUseCase creates a new AnnualReportUseCase instance.
func NewAnnualReportUseCase(
	transactionRepository repositories.TransactionRepository,
) *AnnualReportUseCase {
	return &AnnualReportUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute generates an annual report for the specified user and year.
func (uc *AnnualReportUseCase) Execute(input dtos.AnnualReportInput) (*dtos.AnnualReportOutput, error) {
	// Validate user ID
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Validate year
	if input.Year < 2000 || input.Year > 2100 {
		return nil, fmt.Errorf("invalid year: must be between 2000 and 2100")
	}

	// Create date range for the year
	startDate := time.Date(input.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(input.Year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)

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
		Month    int
	}

	var filteredTransactions []transactionData

	for _, tx := range allTransactions {
		txDate := tx.Date()

		// Check if transaction is within the year
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
			Month:    int(txDate.Month()),
		})
	}

	// Determine currency
	currency := input.Currency
	if currency == "" {
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

	// Calculate totals and monthly breakdown
	var totalIncomeCents int64 = 0
	var totalExpenseCents int64 = 0
	incomeCount := 0
	expenseCount := 0

	// Monthly totals (1-12)
	monthlyIncome := make(map[int]int64)
	monthlyExpense := make(map[int]int64)
	monthlyIncomeCount := make(map[int]int)
	monthlyExpenseCount := make(map[int]int)

	for _, tx := range filteredTransactions {
		// Only count transactions with matching currency
		if tx.Currency != currency {
			continue
		}

		if tx.Type == "INCOME" {
			totalIncomeCents += tx.Amount
			incomeCount++
			monthlyIncome[tx.Month] += tx.Amount
			monthlyIncomeCount[tx.Month]++
		} else if tx.Type == "EXPENSE" {
			totalExpenseCents += tx.Amount
			expenseCount++
			monthlyExpense[tx.Month] += tx.Amount
			monthlyExpenseCount[tx.Month]++
		}
	}

	// Convert totals to float64
	totalIncome, _ := sharedvalueobjects.NewMoney(totalIncomeCents, currencyVO)
	totalExpense, _ := sharedvalueobjects.NewMoney(totalExpenseCents, currencyVO)
	balance, _ := totalIncome.Subtract(totalExpense)

	// Build monthly breakdown
	monthlyBreakdown := make([]dtos.MonthlySummary, 0, 12)
	for month := 1; month <= 12; month++ {
		monthIncomeCents := monthlyIncome[month]
		monthExpenseCents := monthlyExpense[month]

		monthIncomeMoney, _ := sharedvalueobjects.NewMoney(monthIncomeCents, currencyVO)
		monthExpenseMoney, _ := sharedvalueobjects.NewMoney(monthExpenseCents, currencyVO)
		monthBalance, _ := monthIncomeMoney.Subtract(monthExpenseMoney)

		monthlyBreakdown = append(monthlyBreakdown, dtos.MonthlySummary{
			Month:        month,
			TotalIncome:  monthIncomeMoney.Float64(),
			TotalExpense: monthExpenseMoney.Float64(),
			Balance:      monthBalance.Float64(),
			IncomeCount:  monthlyIncomeCount[month],
			ExpenseCount: monthlyExpenseCount[month],
		})
	}

	// Build output
	output := &dtos.AnnualReportOutput{
		UserID:           input.UserID,
		Year:             input.Year,
		Currency:         currency,
		TotalIncome:      totalIncome.Float64(),
		TotalExpense:     totalExpense.Float64(),
		Balance:          balance.Float64(),
		IncomeCount:      incomeCount,
		ExpenseCount:     expenseCount,
		TotalCount:       incomeCount + expenseCount,
		MonthlyBreakdown: monthlyBreakdown,
	}

	return output, nil
}
