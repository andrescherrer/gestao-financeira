package persistence

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&TransactionModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestTransactionEntity creates a test transaction entity.
func createTestTransactionEntity(t *testing.T, userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) *entities.Transaction {
	transactionType := transactionvalueobjects.IncomeType()
	currency, err := sharedvalueobjects.NewCurrency("BRL")
	if err != nil {
		t.Fatalf("Failed to create currency: %v", err)
	}

	amount, err := sharedvalueobjects.NewMoney(100000, currency) // 1000.00 BRL
	if err != nil {
		t.Fatalf("Failed to create amount: %v", err)
	}

	description, err := transactionvalueobjects.NewTransactionDescription("Compra de supermercado")
	if err != nil {
		t.Fatalf("Failed to create description: %v", err)
	}

	date := time.Now()

	transaction, err := entities.NewTransaction(userID, accountID, transactionType, amount, description, date)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	return transaction
}

func TestGormTransactionRepository_FindByID(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transaction := createTestTransactionEntity(t, userID, accountID)

	// Save transaction first
	err := repo.Save(transaction)
	if err != nil {
		t.Fatalf("Failed to save transaction: %v", err)
	}

	tests := []struct {
		name          string
		transactionID transactionvalueobjects.TransactionID
		wantError     bool
		wantNil       bool
	}{
		{
			name:          "find existing transaction",
			transactionID: transaction.ID(),
			wantError:     false,
			wantNil:       false,
		},
		{
			name:          "find non-existent transaction",
			transactionID: transactionvalueobjects.GenerateTransactionID(),
			wantError:     false,
			wantNil:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByID(tt.transactionID)

			if (err != nil) != tt.wantError {
				t.Errorf("FindByID() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantNil && result != nil {
				t.Errorf("FindByID() = %v, want nil", result)
			}

			if !tt.wantNil && result == nil {
				t.Errorf("FindByID() = nil, want transaction")
			}

			if !tt.wantNil && result != nil {
				if !result.ID().Equals(transaction.ID()) {
					t.Errorf("FindByID() transaction ID = %v, want %v", result.ID(), transaction.ID())
				}
			}
		})
	}
}

func TestGormTransactionRepository_FindByUserID(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID1 := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create transactions for user1
	transaction1 := createTestTransactionEntity(t, userID1, accountID)
	transaction2 := createTestTransactionEntity(t, userID1, accountID)
	transaction3 := createTestTransactionEntity(t, userID2, accountID)

	// Save transactions
	_ = repo.Save(transaction1)
	_ = repo.Save(transaction2)
	_ = repo.Save(transaction3)

	// Find transactions for user1
	transactions, err := repo.FindByUserID(userID1)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("FindByUserID() returned %d transactions, want 2", len(transactions))
	}

	// Find transactions for user2
	transactions, err = repo.FindByUserID(userID2)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("FindByUserID() returned %d transactions, want 1", len(transactions))
	}
}

func TestGormTransactionRepository_FindByAccountID(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID1 := accountvalueobjects.GenerateAccountID()
	accountID2 := accountvalueobjects.GenerateAccountID()

	// Create transactions for different accounts
	transaction1 := createTestTransactionEntity(t, userID, accountID1)
	transaction2 := createTestTransactionEntity(t, userID, accountID1)
	transaction3 := createTestTransactionEntity(t, userID, accountID2)

	// Save transactions
	_ = repo.Save(transaction1)
	_ = repo.Save(transaction2)
	_ = repo.Save(transaction3)

	// Find transactions for account1
	transactions, err := repo.FindByAccountID(accountID1)
	if err != nil {
		t.Fatalf("FindByAccountID() error = %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("FindByAccountID() returned %d transactions, want 2", len(transactions))
	}
}

func TestGormTransactionRepository_Save(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transaction := createTestTransactionEntity(t, userID, accountID)

	// Test create
	err := repo.Save(transaction)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify transaction was saved
	saved, err := repo.FindByID(transaction.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Fatal("Save() transaction was not saved")
	}

	// Test update
	newDescription, _ := transactionvalueobjects.NewTransactionDescription("Nova descrição")
	transaction.UpdateDescription(newDescription)

	err = repo.Save(transaction)
	if err != nil {
		t.Fatalf("Save() error on update = %v", err)
	}

	// Verify transaction was updated
	updated, err := repo.FindByID(transaction.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if updated.Description().Value() != "Nova descrição" {
		t.Errorf("Save() updated description = %v, want 'Nova descrição'", updated.Description().Value())
	}
}

func TestGormTransactionRepository_Delete(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transaction := createTestTransactionEntity(t, userID, accountID)

	// Save transaction
	err := repo.Save(transaction)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete transaction
	err = repo.Delete(transaction.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify transaction was deleted (soft delete)
	deleted, err := repo.FindByID(transaction.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	// With soft delete, FindByID should return nil
	if deleted != nil {
		t.Error("Delete() transaction still exists after deletion")
	}
}

func TestGormTransactionRepository_Exists(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transaction := createTestTransactionEntity(t, userID, accountID)

	// Transaction should not exist yet
	exists, err := repo.Exists(transaction.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if exists {
		t.Error("Exists() returned true for non-existent transaction")
	}

	// Save transaction
	err = repo.Save(transaction)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Transaction should exist now
	exists, err = repo.Exists(transaction.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if !exists {
		t.Error("Exists() returned false for existing transaction")
	}
}

func TestGormTransactionRepository_Count(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID1 := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create transactions for user1
	transaction1 := createTestTransactionEntity(t, userID1, accountID)
	transaction2 := createTestTransactionEntity(t, userID1, accountID)
	transaction3 := createTestTransactionEntity(t, userID2, accountID)

	// Save transactions
	_ = repo.Save(transaction1)
	_ = repo.Save(transaction2)
	_ = repo.Save(transaction3)

	// Count transactions for user1
	count, err := repo.Count(userID1)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 2 {
		t.Errorf("Count() = %d, want 2", count)
	}

	// Count transactions for user2
	count, err = repo.Count(userID2)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 1 {
		t.Errorf("Count() = %d, want 1", count)
	}
}

func TestGormTransactionRepository_CountByAccountID(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID1 := accountvalueobjects.GenerateAccountID()
	accountID2 := accountvalueobjects.GenerateAccountID()

	// Create transactions for different accounts
	transaction1 := createTestTransactionEntity(t, userID, accountID1)
	transaction2 := createTestTransactionEntity(t, userID, accountID1)
	transaction3 := createTestTransactionEntity(t, userID, accountID2)

	// Save transactions
	_ = repo.Save(transaction1)
	_ = repo.Save(transaction2)
	_ = repo.Save(transaction3)

	// Count transactions for account1
	count, err := repo.CountByAccountID(accountID1)
	if err != nil {
		t.Fatalf("CountByAccountID() error = %v", err)
	}

	if count != 2 {
		t.Errorf("CountByAccountID() = %d, want 2", count)
	}
}

func TestGormTransactionRepository_FindByUserIDAndType(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create transactions with different types
	incomeType := transactionvalueobjects.IncomeType()
	expenseType := transactionvalueobjects.ExpenseType()

	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	description, _ := transactionvalueobjects.NewTransactionDescription("Test")
	date := time.Now()

	transaction1, _ := entities.NewTransaction(userID, accountID, incomeType, amount, description, date)
	transaction2, _ := entities.NewTransaction(userID, accountID, expenseType, amount, description, date)
	transaction3, _ := entities.NewTransaction(userID, accountID, incomeType, amount, description, date)

	// Save transactions
	_ = repo.Save(transaction1)
	_ = repo.Save(transaction2)
	_ = repo.Save(transaction3)

	// Find income transactions
	transactions, err := repo.FindByUserIDAndType(userID, incomeType)
	if err != nil {
		t.Fatalf("FindByUserIDAndType() error = %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("FindByUserIDAndType() returned %d transactions, want 2", len(transactions))
	}

	// Find expense transactions
	transactions, err = repo.FindByUserIDAndType(userID, expenseType)
	if err != nil {
		t.Fatalf("FindByUserIDAndType() error = %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("FindByUserIDAndType() returned %d transactions, want 1", len(transactions))
	}
}
