package repositories

import (
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/domain/entities"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
)

// InvestmentRepository defines the interface for investment persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type InvestmentRepository interface {
	// FindByID finds an investment by its ID.
	// Returns nil if the investment is not found.
	FindByID(id investmentvalueobjects.InvestmentID) (*entities.Investment, error)

	// FindByUserID finds all investments for a given user.
	// Returns an empty slice if no investments are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Investment, error)

	// FindByAccountID finds all investments for a given account.
	// Returns an empty slice if no investments are found.
	FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Investment, error)

	// FindByType finds all investments for a given user filtered by type.
	// Returns an empty slice if no investments are found.
	FindByType(userID identityvalueobjects.UserID, investmentType investmentvalueobjects.InvestmentType) ([]*entities.Investment, error)

	// Save saves or updates an investment.
	// If the investment already exists (by ID), it updates it.
	// If the investment doesn't exist, it creates a new one.
	Save(investment *entities.Investment) error

	// Delete deletes an investment by its ID.
	Delete(id investmentvalueobjects.InvestmentID) error

	// Exists checks if an investment with the given ID already exists.
	// Returns true if the investment exists, false otherwise.
	Exists(id investmentvalueobjects.InvestmentID) (bool, error)

	// Count returns the total number of investments for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)

	// FindByUserIDWithPagination finds investments for a given user with pagination.
	// If context is provided, filters by context.
	// If investmentType is provided, filters by type.
	// Returns investments, total count, and error.
	FindByUserIDWithPagination(
		userID identityvalueobjects.UserID,
		context string,
		investmentType string,
		offset, limit int,
	) ([]*entities.Investment, int64, error)
}
