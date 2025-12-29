package persistence

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func TestGormUnitOfWork_Begin(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Test initial state
	if uow.IsInTransaction() {
		t.Error("UnitOfWork should not be in transaction initially")
	}

	// Test Begin
	err := uow.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v, want nil", err)
	}

	// Test that transaction is in progress
	if !uow.IsInTransaction() {
		t.Error("UnitOfWork should be in transaction after Begin()")
	}

	// Test that repositories are created
	if uow.TransactionRepository() == nil {
		t.Error("TransactionRepository should not be nil after Begin()")
	}
	if uow.AccountRepository() == nil {
		t.Error("AccountRepository should not be nil after Begin()")
	}
}

func TestGormUnitOfWork_Begin_AlreadyInTransaction(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Begin first transaction
	err := uow.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v, want nil", err)
	}

	// Try to begin again (should fail)
	err = uow.Begin()
	if err == nil {
		t.Error("Begin() should fail when transaction already in progress")
	}
}

func TestGormUnitOfWork_Commit(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Begin transaction
	err := uow.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v, want nil", err)
	}

	// Commit transaction
	err = uow.Commit()
	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}

	// Test that transaction is no longer in progress
	if uow.IsInTransaction() {
		t.Error("UnitOfWork should not be in transaction after Commit()")
	}
}

func TestGormUnitOfWork_Commit_NoTransaction(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Try to commit without beginning (should fail)
	err := uow.Commit()
	if err == nil {
		t.Error("Commit() should fail when no transaction in progress")
	}
}

func TestGormUnitOfWork_Rollback(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Begin transaction
	err := uow.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v, want nil", err)
	}

	// Rollback transaction
	err = uow.Rollback()
	if err != nil {
		t.Fatalf("Rollback() error = %v, want nil", err)
	}

	// Test that transaction is no longer in progress
	if uow.IsInTransaction() {
		t.Error("UnitOfWork should not be in transaction after Rollback()")
	}
}

func TestGormUnitOfWork_Rollback_NoTransaction(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Try to rollback without beginning (should fail)
	err := uow.Rollback()
	if err == nil {
		t.Error("Rollback() should fail when no transaction in progress")
	}
}

func TestGormUnitOfWork_Repositories_WithoutTransaction(t *testing.T) {
	db := setupTestDB(t)
	uow := NewGormUnitOfWork(db).(*GormUnitOfWork)

	// Test that repositories are available even without transaction
	// They should use the main DB connection
	transactionRepo := uow.TransactionRepository()
	if transactionRepo == nil {
		t.Error("TransactionRepository should not be nil even without transaction")
	}

	accountRepo := uow.AccountRepository()
	if accountRepo == nil {
		t.Error("AccountRepository should not be nil even without transaction")
	}
}
