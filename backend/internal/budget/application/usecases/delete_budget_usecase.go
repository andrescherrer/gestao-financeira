package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// DeleteBudgetUseCase handles budget deletion.
type DeleteBudgetUseCase struct {
	budgetRepository repositories.BudgetRepository
	eventBus         *eventbus.EventBus
}

// NewDeleteBudgetUseCase creates a new DeleteBudgetUseCase instance.
func NewDeleteBudgetUseCase(
	budgetRepository repositories.BudgetRepository,
	eventBus *eventbus.EventBus,
) *DeleteBudgetUseCase {
	return &DeleteBudgetUseCase{
		budgetRepository: budgetRepository,
		eventBus:         eventBus,
	}
}

// Execute performs the budget deletion.
func (uc *DeleteBudgetUseCase) Execute(input dtos.DeleteBudgetInput) error {
	// Create budget ID value object
	budgetID, err := valueobjects.NewBudgetID(input.BudgetID)
	if err != nil {
		return fmt.Errorf("invalid budget ID: %w", err)
	}

	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Find budget by ID
	budget, err := uc.budgetRepository.FindByID(budgetID)
	if err != nil {
		return fmt.Errorf("failed to find budget: %w", err)
	}

	if budget == nil {
		return errors.New("budget not found")
	}

	// Verify that the budget belongs to the user
	if !budget.UserID().Equals(userID) {
		return errors.New("budget not found")
	}

	// Delete budget
	if err := uc.budgetRepository.Delete(budgetID); err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	// Publish domain event
	event := events.NewBaseDomainEvent(
		"BudgetDeleted",
		budgetID.Value(),
		"Budget",
	)
	if err := uc.eventBus.Publish(&event); err != nil {
		_ = err // Ignore for now, but should be logged
	}

	return nil
}
