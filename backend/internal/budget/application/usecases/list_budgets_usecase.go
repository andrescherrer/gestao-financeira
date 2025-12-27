package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/entities"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// ListBudgetsUseCase handles listing budgets for a user.
type ListBudgetsUseCase struct {
	budgetRepository repositories.BudgetRepository
}

// NewListBudgetsUseCase creates a new ListBudgetsUseCase instance.
func NewListBudgetsUseCase(
	budgetRepository repositories.BudgetRepository,
) *ListBudgetsUseCase {
	return &ListBudgetsUseCase{
		budgetRepository: budgetRepository,
	}
}

// Execute performs the budget listing.
// It validates the input, retrieves budgets from the repository,
// and returns them as DTOs.
func (uc *ListBudgetsUseCase) Execute(input dtos.ListBudgetsInput) (*dtos.ListBudgetsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get all budgets for user
	budgets, err := uc.budgetRepository.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list budgets: %w", err)
	}

	// Apply filters
	filteredBudgets := make([]*entities.Budget, 0)
	for _, budget := range budgets {
		// Filter by category
		if input.CategoryID != "" {
			categoryID, err := categoryvalueobjects.NewCategoryID(input.CategoryID)
			if err == nil && !budget.CategoryID().Equals(categoryID) {
				continue
			}
		}

		// Filter by period type
		if input.PeriodType != "" {
			if string(budget.Period().PeriodType()) != input.PeriodType {
				continue
			}
		}

		// Filter by year
		if input.Year != nil {
			if budget.Period().Year() != *input.Year {
				continue
			}
		}

		// Filter by month
		if input.Month != nil {
			budgetMonth := budget.Period().Month()
			if budgetMonth == nil || *budgetMonth != *input.Month {
				continue
			}
		}

		// Filter by context
		if input.Context != "" {
			if budget.Context().Value() != input.Context {
				continue
			}
		}

		// Filter by is_active
		if input.IsActive != nil {
			if budget.IsActive() != *input.IsActive {
				continue
			}
		}

		filteredBudgets = append(filteredBudgets, budget)
	}

	// Build output
	budgetOutputs := make([]dtos.BudgetOutput, 0, len(filteredBudgets))
	for _, budget := range filteredBudgets {
		budgetAmount := budget.Amount()
		budgetOutputs = append(budgetOutputs, dtos.BudgetOutput{
			BudgetID:   budget.ID().Value(),
			UserID:     budget.UserID().Value(),
			CategoryID: budget.CategoryID().Value(),
			Amount:     budgetAmount.Float64(),
			Currency:   budgetAmount.Currency().Code(),
			PeriodType: string(budget.Period().PeriodType()),
			Year:       budget.Period().Year(),
			Month:      budget.Period().Month(),
			Context:    budget.Context().Value(),
			IsActive:   budget.IsActive(),
			CreatedAt:  budget.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  budget.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	output := &dtos.ListBudgetsOutput{
		Budgets: budgetOutputs,
		Total:   int64(len(budgetOutputs)),
	}

	return output, nil
}
