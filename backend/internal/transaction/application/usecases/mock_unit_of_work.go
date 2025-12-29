package usecases

import (
	"fmt"

	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	transactionrepositories "gestao-financeira/backend/internal/transaction/domain/repositories"
)

// mockUnitOfWork is a mock implementation of UnitOfWork for testing.
type mockUnitOfWork struct {
	transactionRepository transactionrepositories.TransactionRepository
	accountRepository     accountrepositories.AccountRepository
	inTransaction         bool
	beginErr              error
	commitErr             error
	rollbackErr           error
}

// newMockUnitOfWork creates a new mock UnitOfWork instance.
func newMockUnitOfWork(
	transactionRepository transactionrepositories.TransactionRepository,
	accountRepository accountrepositories.AccountRepository,
) *mockUnitOfWork {
	return &mockUnitOfWork{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		inTransaction:         false,
	}
}

// Begin starts a new database transaction.
func (m *mockUnitOfWork) Begin() error {
	if m.inTransaction {
		return fmt.Errorf("transaction already in progress")
	}
	if m.beginErr != nil {
		return m.beginErr
	}
	m.inTransaction = true
	return nil
}

// Commit commits the current transaction.
func (m *mockUnitOfWork) Commit() error {
	if !m.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}
	if m.commitErr != nil {
		return m.commitErr
	}
	m.inTransaction = false
	return nil
}

// Rollback rolls back the current transaction.
func (m *mockUnitOfWork) Rollback() error {
	if !m.inTransaction {
		return fmt.Errorf("no transaction in progress")
	}
	if m.rollbackErr != nil {
		return m.rollbackErr
	}
	m.inTransaction = false
	return nil
}

// TransactionRepository returns a TransactionRepository.
func (m *mockUnitOfWork) TransactionRepository() transactionrepositories.TransactionRepository {
	return m.transactionRepository
}

// AccountRepository returns an AccountRepository.
func (m *mockUnitOfWork) AccountRepository() accountrepositories.AccountRepository {
	return m.accountRepository
}

// IsInTransaction returns true if a transaction is currently in progress.
func (m *mockUnitOfWork) IsInTransaction() bool {
	return m.inTransaction
}
