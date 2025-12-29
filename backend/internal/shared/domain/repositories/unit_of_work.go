package repositories

import (
	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	transactionrepositories "gestao-financeira/backend/internal/transaction/domain/repositories"
)

// UnitOfWork represents a unit of work pattern for managing database transactions.
// It ensures that all operations within a unit of work are executed atomically.
// If any operation fails, all changes are rolled back.
type UnitOfWork interface {
	// Begin starts a new database transaction.
	// Returns an error if a transaction is already in progress.
	Begin() error

	// Commit commits the current transaction.
	// All changes made within the transaction are persisted.
	// Returns an error if no transaction is in progress or if commit fails.
	Commit() error

	// Rollback rolls back the current transaction.
	// All changes made within the transaction are discarded.
	// Returns an error if no transaction is in progress or if rollback fails.
	Rollback() error

	// TransactionRepository returns a TransactionRepository that operates within the current transaction.
	// If no transaction is in progress, it returns a repository that operates outside of a transaction.
	TransactionRepository() transactionrepositories.TransactionRepository

	// AccountRepository returns an AccountRepository that operates within the current transaction.
	// If no transaction is in progress, it returns a repository that operates outside of a transaction.
	AccountRepository() accountrepositories.AccountRepository

	// IsInTransaction returns true if a transaction is currently in progress.
	IsInTransaction() bool
}
