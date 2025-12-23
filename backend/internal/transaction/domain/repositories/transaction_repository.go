package repositories

import (
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// TransactionRepository defines the interface for transaction persistence operations.
type TransactionRepository interface {
	// FindByID finds a transaction by its ID.
	FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error)

	// FindByUserID finds all transactions for a given user.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error)

	// FindByAccountID finds all transactions for a given account.
	FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error)

	// FindByUserIDAndAccountID finds all transactions for a given user and account.
	FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error)

	// FindByUserIDAndType finds all transactions for a given user filtered by type.
	FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error)

	// Save saves or updates a transaction.
	Save(transaction *entities.Transaction) error

	// Delete deletes a transaction by its ID.
	Delete(id transactionvalueobjects.TransactionID) error

	// Exists checks if a transaction with the given ID already exists.
	Exists(id transactionvalueobjects.TransactionID) (bool, error)

	// Count returns the total number of transactions for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)

	// CountByAccountID returns the total number of transactions for a given account.
	CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error)
}
