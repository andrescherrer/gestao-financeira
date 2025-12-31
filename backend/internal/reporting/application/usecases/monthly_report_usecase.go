package usecases

import (
	"fmt"
	"strconv"
	"time"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	reportingservices "gestao-financeira/backend/internal/reporting/infrastructure/services"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// MonthlyReportUseCase handles generating monthly financial reports.
type MonthlyReportUseCase struct {
	transactionRepository repositories.TransactionRepository
	cacheService          *reportingservices.ReportCacheService
}

// NewMonthlyReportUseCase creates a new MonthlyReportUseCase instance.
func NewMonthlyReportUseCase(
	transactionRepository repositories.TransactionRepository,
	cacheService *reportingservices.ReportCacheService,
) *MonthlyReportUseCase {
	return &MonthlyReportUseCase{
		transactionRepository: transactionRepository,
		cacheService:          cacheService,
	}
}

// Execute generates a monthly report for the specified user, year, and month.
func (uc *MonthlyReportUseCase) Execute(input dtos.MonthlyReportInput) (*dtos.MonthlyReportOutput, error) {
	// Try to get from cache first
	if uc.cacheService != nil {
		cacheKey := uc.cacheService.GenerateKey("monthly", map[string]string{
			"user_id":  input.UserID,
			"year":     strconv.Itoa(input.Year),
			"month":    strconv.Itoa(input.Month),
			"currency": input.Currency,
		})

		var cachedOutput dtos.MonthlyReportOutput
		found, err := uc.cacheService.Get(cacheKey, &cachedOutput)
		if err == nil && found {
			return &cachedOutput, nil
		}
	}

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
	endDate := startDate.AddDate(0, 1, 0).AddDate(0, 0, -1) // Last day of the month

	// Get transactions for the user within date range (optimized query)
	var transactions []*entities.Transaction
	if input.Currency != "" {
		// Filter by currency at database level
		transactions, err = uc.transactionRepository.FindByUserIDAndDateRangeWithCurrency(userID, startDate, endDate, input.Currency)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}
	} else {
		// Get all transactions in date range
		transactions, err = uc.transactionRepository.FindByUserIDAndDateRange(userID, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}
	}

	// Process transactions
	type transactionData struct {
		Type     string
		Amount   int64
		Currency string
	}

	filteredTransactions := make([]transactionData, 0, len(transactions))
	for _, tx := range transactions {
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

	// Cache the result
	if uc.cacheService != nil {
		cacheKey := uc.cacheService.GenerateKey("monthly", map[string]string{
			"user_id":  input.UserID,
			"year":     strconv.Itoa(input.Year),
			"month":    strconv.Itoa(input.Month),
			"currency": input.Currency,
		})
		_ = uc.cacheService.Set(cacheKey, output) // Ignore cache errors
	}

	return output, nil
}
