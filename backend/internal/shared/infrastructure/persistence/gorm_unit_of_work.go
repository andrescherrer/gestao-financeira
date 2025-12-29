package persistence

import (
	"fmt"

	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	sharedrepositories "gestao-financeira/backend/internal/shared/domain/repositories"
	transactionrepositories "gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"

	"gorm.io/gorm"
)

// GormUnitOfWork implements the UnitOfWork interface using GORM.
type GormUnitOfWork struct {
	db                    *gorm.DB
	tx                    *gorm.DB
	transactionRepository transactionrepositories.TransactionRepository
	accountRepository     accountrepositories.AccountRepository
	inTransaction         bool
}

// NewGormUnitOfWork creates a new GormUnitOfWork instance.
func NewGormUnitOfWork(db *gorm.DB) sharedrepositories.UnitOfWork {
	return &GormUnitOfWork{
		db:            db,
		inTransaction: false,
	}
}

// Begin starts a new database transaction.
func (uow *GormUnitOfWork) Begin() error {
	if uow.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}

	uow.tx = uow.db.Begin()
	if uow.tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", uow.tx.Error)
	}

	uow.inTransaction = true

	// Create repositories that use the transaction
	uow.transactionRepository = transactionpersistence.NewGormTransactionRepository(uow.tx)
	uow.accountRepository = accountpersistence.NewGormAccountRepository(uow.tx)

	return nil
}

// Commit commits the current transaction.
func (uow *GormUnitOfWork) Commit() error {
	if !uow.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if err := uow.tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	uow.inTransaction = false
	uow.tx = nil
	uow.transactionRepository = nil
	uow.accountRepository = nil

	return nil
}

// Rollback rolls back the current transaction.
func (uow *GormUnitOfWork) Rollback() error {
	if !uow.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}

	if err := uow.tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	uow.inTransaction = false
	uow.tx = nil
	uow.transactionRepository = nil
	uow.accountRepository = nil

	return nil
}

// TransactionRepository returns a TransactionRepository that operates within the current transaction.
func (uow *GormUnitOfWork) TransactionRepository() transactionrepositories.TransactionRepository {
	if uow.inTransaction && uow.transactionRepository != nil {
		return uow.transactionRepository
	}
	// If no transaction, return a repository that uses the main DB connection
	return transactionpersistence.NewGormTransactionRepository(uow.db)
}

// AccountRepository returns an AccountRepository that operates within the current transaction.
func (uow *GormUnitOfWork) AccountRepository() accountrepositories.AccountRepository {
	if uow.inTransaction && uow.accountRepository != nil {
		return uow.accountRepository
	}
	// If no transaction, return a repository that uses the main DB connection
	return accountpersistence.NewGormAccountRepository(uow.db)
}

// IsInTransaction returns true if a transaction is currently in progress.
func (uow *GormUnitOfWork) IsInTransaction() bool {
	return uow.inTransaction
}
