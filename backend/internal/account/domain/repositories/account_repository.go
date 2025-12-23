package repositories

import (
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// AccountRepository defines the interface for account persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type AccountRepository interface {
	// FindByID finds an account by its ID.
	// Returns nil if the account is not found.
	FindByID(id valueobjects.AccountID) (*entities.Account, error)

	// FindByUserID finds all accounts for a given user.
	// Returns an empty slice if no accounts are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error)

	// FindByUserIDAndContext finds all accounts for a given user filtered by context.
	// Returns an empty slice if no accounts are found.
	FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error)

	// Save saves or updates an account.
	// If the account already exists (by ID), it updates it.
	// If the account doesn't exist, it creates a new one.
	Save(account *entities.Account) error

	// Delete deletes an account by its ID.
	Delete(id valueobjects.AccountID) error

	// Exists checks if an account with the given ID already exists.
	// Returns true if the account exists, false otherwise.
	Exists(id valueobjects.AccountID) (bool, error)

	// Count returns the total number of accounts for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)
}
