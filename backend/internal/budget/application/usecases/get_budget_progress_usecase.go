package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	transactionrepositories "gestao-financeira/backend/internal/transaction/domain/repositories"
)

// GetBudgetProgressUseCase handles budget progress calculation.
type GetBudgetProgressUseCase struct {
	budgetRepository      repositories.BudgetRepository
	transactionRepository transactionrepositories.TransactionRepository
}

// NewGetBudgetProgressUseCase creates a new GetBudgetProgressUseCase instance.
func NewGetBudgetProgressUseCase(
	budgetRepository repositories.BudgetRepository,
	transactionRepository transactionrepositories.TransactionRepository,
) *GetBudgetProgressUseCase {
	return &GetBudgetProgressUseCase{
		budgetRepository:      budgetRepository,
		transactionRepository: transactionRepository,
	}
}

// Execute performs the budget progress calculation.
func (uc *GetBudgetProgressUseCase) Execute(input dtos.GetBudgetProgressInput) (*dtos.GetBudgetProgressOutput, error) {
	// Create budget ID value object
	budgetID, err := valueobjects.NewBudgetID(input.BudgetID)
	if err != nil {
		return nil, fmt.Errorf("invalid budget ID: %w", err)
	}

	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find budget
	budget, err := uc.budgetRepository.FindByID(budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to find budget: %w", err)
	}

	if budget == nil {
		return nil, errors.New("budget not found")
	}

	// Verify that the budget belongs to the user
	if !budget.UserID().Equals(userID) {
		return nil, errors.New("budget not found")
	}

	// Get period
	period := budget.Period()

	// Find all transactions for this user
	allTransactions, err := uc.transactionRepository.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}

	// Filter transactions by date range and type
	// Note: Currently, transactions don't have category_id, so we calculate based on all EXPENSE transactions
	// TODO: Add category_id to transactions to properly filter by category
	budgetAmount := budget.Amount()
	currency := budgetAmount.Currency()
	spentCents := int64(0)

	for _, transaction := range allTransactions {
		// Only count EXPENSE transactions
		if transaction.TransactionType().Value() != "EXPENSE" {
			continue
		}

		// Only count transactions within the period
		transactionDate := transaction.Date()
		if !period.Includes(transactionDate) {
			continue
		}

		// Only count transactions with the same currency
		if !transaction.Amount().Currency().Equals(currency) {
			continue
		}

		spentCents += transaction.Amount().Amount()
	}

	spent, err := sharedvalueobjects.NewMoney(spentCents, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to create spent amount: %w", err)
	}

	// Calculate remaining amount
	remaining, err := budgetAmount.Subtract(spent)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate remaining amount: %w", err)
	}

	// Calculate percentage used
	percentageUsed := float64(spentCents) / float64(budgetAmount.Amount()) * 100.0
	if percentageUsed > 100.0 {
		percentageUsed = 100.0
	}

	// Check if exceeded
	isExceeded := spentCents > budgetAmount.Amount()

	// Build output
	output := &dtos.GetBudgetProgressOutput{
		BudgetID:       budget.ID().Value(),
		CategoryID:     budget.CategoryID().Value(),
		Budgeted:       budgetAmount.Float64(),
		Spent:          spent.Float64(),
		Remaining:      remaining.Float64(),
		PercentageUsed: percentageUsed,
		Currency:       currency.Code(),
		IsExceeded:     isExceeded,
		PeriodType:     string(period.PeriodType()),
		Year:           period.Year(),
		Month:          period.Month(),
	}

	return output, nil
}
