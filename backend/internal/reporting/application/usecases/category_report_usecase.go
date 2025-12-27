package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/dtos"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
)

// CategoryReportUseCase handles generating category-based financial reports.
// Note: Currently, transactions don't have category_id. This use case groups by transaction type.
// When category_id is added to transactions, this will be enhanced to group by category.
type CategoryReportUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewCategoryReportUseCase creates a new CategoryReportUseCase instance.
func NewCategoryReportUseCase(
	transactionRepository repositories.TransactionRepository,
) *CategoryReportUseCase {
	return &CategoryReportUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute generates a category report for the specified user.
// Currently groups by transaction type (INCOME/EXPENSE) since transactions don't have category_id yet.
func (uc *CategoryReportUseCase) Execute(input dtos.CategoryReportInput) (*dtos.CategoryReportOutput, error) {
	// Validate user ID
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
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

	// Group by type (INCOME/EXPENSE)
	// TODO: When category_id is added to transactions, group by category instead
	typeSummary := make(map[string]struct {
		TotalCents int64
		Count      int
	})

	for _, tx := range filteredTransactions {
		// Only count transactions with matching currency
		if tx.Currency != currency {
			continue
		}

		summary := typeSummary[tx.Type]
		summary.TotalCents += tx.Amount
		summary.Count++
		typeSummary[tx.Type] = summary
	}

	// Calculate totals
	var totalIncomeCents int64 = 0
	var totalExpenseCents int64 = 0
	totalCount := 0

	incomeSummary := typeSummary["INCOME"]
	expenseSummary := typeSummary["EXPENSE"]

	totalIncomeCents = incomeSummary.TotalCents
	totalExpenseCents = expenseSummary.TotalCents
	totalCount = incomeSummary.Count + expenseSummary.Count

	// Convert to Money objects
	totalIncome, _ := sharedvalueobjects.NewMoney(totalIncomeCents, currencyVO)
	totalExpense, _ := sharedvalueobjects.NewMoney(totalExpenseCents, currencyVO)
	balance, _ := totalIncome.Subtract(totalExpense)

	// Build category breakdown
	// Note: Currently grouping by type since transactions don't have category_id
	// When category_id is added, this will be enhanced to show actual categories
	categoryBreakdown := make([]dtos.CategorySummary, 0)

	// Income summary
	if incomeSummary.Count > 0 {
		incomeAmount, _ := sharedvalueobjects.NewMoney(incomeSummary.TotalCents, currencyVO)
		percentage := 0.0
		if totalIncomeCents > 0 {
			percentage = float64(incomeSummary.TotalCents) / float64(totalIncomeCents) * 100.0
		}

		categoryBreakdown = append(categoryBreakdown, dtos.CategorySummary{
			CategoryID:   "",       // Will be populated when category_id is added
			CategoryName: "Income", // Placeholder
			Type:         "INCOME",
			TotalAmount:  incomeAmount.Float64(),
			Count:        incomeSummary.Count,
			Percentage:   percentage,
		})
	}

	// Expense summary
	if expenseSummary.Count > 0 {
		expenseAmount, _ := sharedvalueobjects.NewMoney(expenseSummary.TotalCents, currencyVO)
		percentage := 0.0
		if totalExpenseCents > 0 {
			percentage = float64(expenseSummary.TotalCents) / float64(totalExpenseCents) * 100.0
		}

		categoryBreakdown = append(categoryBreakdown, dtos.CategorySummary{
			CategoryID:   "",        // Will be populated when category_id is added
			CategoryName: "Expense", // Placeholder
			Type:         "EXPENSE",
			TotalAmount:  expenseAmount.Float64(),
			Count:        expenseSummary.Count,
			Percentage:   percentage,
		})
	}

	// Build output
	output := &dtos.CategoryReportOutput{
		UserID:            input.UserID,
		Currency:          currency,
		CategoryBreakdown: categoryBreakdown,
		TotalIncome:       totalIncome.Float64(),
		TotalExpense:      totalExpense.Float64(),
		Balance:           balance.Float64(),
		TotalCount:        totalCount,
	}

	return output, nil
}
