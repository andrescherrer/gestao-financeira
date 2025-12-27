package usecases

import (
	"fmt"
	"strings"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/entities"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateBudgetUseCase handles budget creation.
type CreateBudgetUseCase struct {
	budgetRepository repositories.BudgetRepository
	eventBus         *eventbus.EventBus
}

// NewCreateBudgetUseCase creates a new CreateBudgetUseCase instance.
func NewCreateBudgetUseCase(
	budgetRepository repositories.BudgetRepository,
	eventBus *eventbus.EventBus,
) *CreateBudgetUseCase {
	return &CreateBudgetUseCase{
		budgetRepository: budgetRepository,
		eventBus:         eventBus,
	}
}

// Execute performs the budget creation.
// It validates the input, creates value objects, creates a new budget entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateBudgetUseCase) Execute(input dtos.CreateBudgetInput) (*dtos.CreateBudgetOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create category ID value object
	categoryID, err := categoryvalueobjects.NewCategoryID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(input.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	// Create amount (convert float to cents)
	amountCents := int64(input.Amount * 100)
	amount, err := sharedvalueobjects.NewMoney(amountCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// Create period value object
	var period valueobjects.BudgetPeriod
	if input.PeriodType == "MONTHLY" {
		if input.Month == nil {
			return nil, fmt.Errorf("month is required for monthly periods")
		}
		period, err = valueobjects.NewMonthlyBudgetPeriod(input.Year, *input.Month)
	} else {
		period, err = valueobjects.NewYearlyBudgetPeriod(input.Year)
	}
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	// Create context value object
	context, err := sharedvalueobjects.NewAccountContext(input.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid context: %w", err)
	}

	// Check if budget already exists for this category and period
	existingBudget, err := uc.budgetRepository.FindByCategoryAndPeriod(categoryID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to check if budget exists: %w", err)
	}
	if existingBudget != nil {
		return nil, fmt.Errorf("budget already exists for this category and period")
	}

	// Create budget entity
	budget, err := entities.NewBudget(userID, categoryID, amount, period, context)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	// Save budget to repository
	if err := uc.budgetRepository.Save(budget); err != nil {
		// Check if error is about duplicate (fallback in case check above didn't catch it)
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "unique constraint") {
			return nil, fmt.Errorf("budget already exists for this category and period")
		}
		return nil, fmt.Errorf("failed to save budget: %w", err)
	}

	// Publish domain events
	domainEvents := budget.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the budget creation
			_ = err // Ignore for now, but should be logged
		}
	}
	budget.ClearEvents()

	// Build output
	budgetAmount := budget.Amount()
	output := &dtos.CreateBudgetOutput{
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
	}

	return output, nil
}
