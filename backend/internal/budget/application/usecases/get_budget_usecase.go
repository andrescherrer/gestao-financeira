package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// GetBudgetUseCase handles retrieving a single budget by ID.
type GetBudgetUseCase struct {
	budgetRepository repositories.BudgetRepository
}

// NewGetBudgetUseCase creates a new GetBudgetUseCase instance.
func NewGetBudgetUseCase(
	budgetRepository repositories.BudgetRepository,
) *GetBudgetUseCase {
	return &GetBudgetUseCase{
		budgetRepository: budgetRepository,
	}
}

// Execute performs the budget retrieval.
// It validates the input, retrieves the budget from the repository,
// and returns it as a DTO.
func (uc *GetBudgetUseCase) Execute(input dtos.GetBudgetInput) (*dtos.GetBudgetOutput, error) {
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
		return nil, fmt.Errorf("budget not found")
	}

	// Verify that the budget belongs to the user
	if !budget.UserID().Equals(userID) {
		return nil, fmt.Errorf("budget not found")
	}

	// Build output
	budgetAmount := budget.Amount()
	output := &dtos.GetBudgetOutput{
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
	}

	return output, nil
}
