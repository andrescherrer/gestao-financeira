package repositories

import (
	"gestao-financeira/backend/internal/budget/domain/entities"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// BudgetRepository defines the interface for budget persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type BudgetRepository interface {
	// FindByID finds a budget by its ID.
	// Returns nil if the budget is not found.
	FindByID(id valueobjects.BudgetID) (*entities.Budget, error)

	// FindByUserID finds all budgets for a given user.
	// Returns an empty slice if no budgets are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Budget, error)

	// FindByCategoryAndPeriod finds a budget by category ID and period.
	// Returns nil if the budget is not found.
	FindByCategoryAndPeriod(categoryID categoryvalueobjects.CategoryID, period valueobjects.BudgetPeriod) (*entities.Budget, error)

	// FindByPeriod finds all budgets for a given user and period.
	// Returns an empty slice if no budgets are found.
	FindByPeriod(userID identityvalueobjects.UserID, period valueobjects.BudgetPeriod) ([]*entities.Budget, error)

	// Save saves or updates a budget.
	// If the budget already exists (by ID), it updates it.
	// If the budget doesn't exist, it creates a new one.
	Save(budget *entities.Budget) error

	// Delete deletes a budget by its ID.
	Delete(id valueobjects.BudgetID) error

	// Exists checks if a budget with the given ID already exists.
	// Returns true if the budget exists, false otherwise.
	Exists(id valueobjects.BudgetID) (bool, error)

	// Count returns the total number of budgets for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)
}
