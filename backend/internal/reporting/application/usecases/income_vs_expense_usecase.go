package usecases

import (
	"fmt"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
)

// IncomeVsExpenseUseCase handles generating income vs expense comparison reports.
type IncomeVsExpenseUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewIncomeVsExpenseUseCase creates a new IncomeVsExpenseUseCase instance.
func NewIncomeVsExpenseUseCase(
	transactionRepository repositories.TransactionRepository,
) *IncomeVsExpenseUseCase {
	return &IncomeVsExpenseUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute generates an income vs expense comparison report.
func (uc *IncomeVsExpenseUseCase) Execute(input dtos.IncomeVsExpenseInput) (*dtos.IncomeVsExpenseOutput, error) {
	// Validate user ID
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Validate group_by if specified
	if input.GroupBy != "" && input.GroupBy != "day" && input.GroupBy != "week" && input.GroupBy != "month" && input.GroupBy != "year" {
		return nil, fmt.Errorf("invalid group_by: must be one of: day, week, month, year")
	}

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
		Date     time.Time
	}

	var filteredTransactions []transactionData

	for _, tx := range allTransactions {
		// Filter by date range if specified
		if input.StartDate != nil && tx.Date().Before(*input.StartDate) {
			continue
		}
		if input.EndDate != nil && tx.Date().After(*input.EndDate) {
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
			Date:     tx.Date(),
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

	// Calculate totals
	var totalIncomeCents int64 = 0
	var totalExpenseCents int64 = 0
	incomeCount := 0
	expenseCount := 0

	// Period breakdown map (if group_by is specified)
	periodMap := make(map[string]struct {
		IncomeCents  int64
		ExpenseCents int64
		IncomeCount  int
		ExpenseCount int
	})

	for _, tx := range filteredTransactions {
		// Only count transactions with matching currency
		if tx.Currency != currency {
			continue
		}

		if tx.Type == "INCOME" {
			totalIncomeCents += tx.Amount
			incomeCount++
		} else if tx.Type == "EXPENSE" {
			totalExpenseCents += tx.Amount
			expenseCount++
		}

		// Add to period breakdown if group_by is specified
		if input.GroupBy != "" {
			periodKey := uc.getPeriodKey(tx.Date, input.GroupBy)
			period := periodMap[periodKey]
			if tx.Type == "INCOME" {
				period.IncomeCents += tx.Amount
				period.IncomeCount++
			} else if tx.Type == "EXPENSE" {
				period.ExpenseCents += tx.Amount
				period.ExpenseCount++
			}
			periodMap[periodKey] = period
		}
	}

	// Convert totals to Money objects
	totalIncome, _ := sharedvalueobjects.NewMoney(totalIncomeCents, currencyVO)
	totalExpense, _ := sharedvalueobjects.NewMoney(totalExpenseCents, currencyVO)
	balance, _ := totalIncome.Subtract(totalExpense)

	// Build period breakdown if group_by is specified
	var periodBreakdown []dtos.PeriodSummary
	if input.GroupBy != "" {
		for periodKey, period := range periodMap {
			periodIncome, _ := sharedvalueobjects.NewMoney(period.IncomeCents, currencyVO)
			periodExpense, _ := sharedvalueobjects.NewMoney(period.ExpenseCents, currencyVO)
			periodBalance, _ := periodIncome.Subtract(periodExpense)

			periodBreakdown = append(periodBreakdown, dtos.PeriodSummary{
				Period:       periodKey,
				TotalIncome:  periodIncome.Float64(),
				TotalExpense: periodExpense.Float64(),
				Balance:      periodBalance.Float64(),
				IncomeCount:  period.IncomeCount,
				ExpenseCount: period.ExpenseCount,
			})
		}
	}

	// Build output
	output := &dtos.IncomeVsExpenseOutput{
		UserID:          input.UserID,
		Currency:        currency,
		TotalIncome:     totalIncome.Float64(),
		TotalExpense:    totalExpense.Float64(),
		Balance:         balance.Float64(),
		Difference:      balance.Float64(), // Same as balance (income - expense)
		IncomeCount:     incomeCount,
		ExpenseCount:    expenseCount,
		TotalCount:      incomeCount + expenseCount,
		PeriodBreakdown: periodBreakdown,
	}

	return output, nil
}

// getPeriodKey generates a period key based on the date and group_by option.
func (uc *IncomeVsExpenseUseCase) getPeriodKey(date time.Time, groupBy string) string {
	switch groupBy {
	case "day":
		return date.Format("2006-01-02")
	case "week":
		// Get the start of the week (Monday)
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday becomes 7
		}
		weekStart := date.AddDate(0, 0, -(weekday - 1))
		return weekStart.Format("2006-01-02")
	case "month":
		return date.Format("2006-01")
	case "year":
		return date.Format("2006")
	default:
		return date.Format("2006-01-02")
	}
}
