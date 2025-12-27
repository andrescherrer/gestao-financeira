package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// UpdateBudgetUseCase handles budget updates.
type UpdateBudgetUseCase struct {
	budgetRepository repositories.BudgetRepository
	eventBus         *eventbus.EventBus
}

// NewUpdateBudgetUseCase creates a new UpdateBudgetUseCase instance.
func NewUpdateBudgetUseCase(
	budgetRepository repositories.BudgetRepository,
	eventBus *eventbus.EventBus,
) *UpdateBudgetUseCase {
	return &UpdateBudgetUseCase{
		budgetRepository: budgetRepository,
		eventBus:         eventBus,
	}
}

// Execute performs the budget update.
func (uc *UpdateBudgetUseCase) Execute(input dtos.UpdateBudgetInput) (*dtos.UpdateBudgetOutput, error) {
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

	// Find budget by ID
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

	// Update amount if provided
	if input.Amount != nil {
		currency := budget.Amount().Currency()
		amountCents := int64(*input.Amount * 100)
		newAmount, err := sharedvalueobjects.NewMoney(amountCents, currency)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}
		if err := budget.UpdateAmount(newAmount); err != nil {
			return nil, fmt.Errorf("failed to update budget amount: %w", err)
		}
	}

	// Update period if provided
	if input.PeriodType != nil || input.Year != nil || input.Month != nil {
		var period valueobjects.BudgetPeriod
		periodType := input.PeriodType
		if periodType == nil {
			periodTypeStr := string(budget.Period().PeriodType())
			periodType = &periodTypeStr
		}

		year := input.Year
		if year == nil {
			yearVal := budget.Period().Year()
			year = &yearVal
		}

		month := input.Month
		if *periodType == "MONTHLY" {
			if month == nil {
				budgetMonth := budget.Period().Month()
				if budgetMonth == nil {
					return nil, errors.New("month is required for monthly periods")
				}
				month = budgetMonth
			}
			var err error
			period, err = valueobjects.NewMonthlyBudgetPeriod(*year, *month)
			if err != nil {
				return nil, fmt.Errorf("invalid period: %w", err)
			}
		} else {
			var err error
			period, err = valueobjects.NewYearlyBudgetPeriod(*year)
			if err != nil {
				return nil, fmt.Errorf("invalid period: %w", err)
			}
		}

		if err := budget.UpdatePeriod(period); err != nil {
			return nil, fmt.Errorf("failed to update budget period: %w", err)
		}
	}

	// Update is_active if provided
	if input.IsActive != nil {
		if *input.IsActive {
			if err := budget.Activate(); err != nil {
				return nil, fmt.Errorf("failed to activate budget: %w", err)
			}
		} else {
			if err := budget.Deactivate(); err != nil {
				return nil, fmt.Errorf("failed to deactivate budget: %w", err)
			}
		}
	}

	// Check if at least one field was provided for update
	if input.Amount == nil && input.PeriodType == nil && input.Year == nil && input.Month == nil && input.IsActive == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Save budget to repository
	if err := uc.budgetRepository.Save(budget); err != nil {
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	// Publish domain events
	domainEvents := budget.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			_ = err // Ignore for now, but should be logged
		}
	}
	budget.ClearEvents()

	// Build output
	budgetAmount := budget.Amount()
	output := &dtos.UpdateBudgetOutput{
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
		UpdatedAt:  budget.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
