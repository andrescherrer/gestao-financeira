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

// createTestRecurringTransactionEntity creates a test recurring transaction entity.
func createTestRecurringTransactionEntity(
	t *testing.T,
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	frequency transactionvalueobjects.RecurrenceFrequency,
	date time.Time,
	endDate *time.Time,
) *entities.Transaction {
	transactionType := transactionvalueobjects.ExpenseType()
	currency, err := sharedvalueobjects.NewCurrency("BRL")
	if err != nil {
		t.Fatalf("Failed to create currency: %v", err)
	}

	amount, err := sharedvalueobjects.NewMoney(5000, currency) // 50.00 BRL
	if err != nil {
		t.Fatalf("Failed to create amount: %v", err)
	}

	description, err := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	if err != nil {
		t.Fatalf("Failed to create description: %v", err)
	}

	transaction, err := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionType,
		amount,
		description,
		date,
		true,
		&frequency,
		endDate,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create recurring transaction: %v", err)
	}

	return transaction
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

func TestGormTransactionRepository_FindActiveRecurringTransactions(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create active recurring transaction (no end date)
	frequency1 := transactionvalueobjects.MonthlyFrequency()
	recurringTx1 := createTestRecurringTransactionEntity(t, userID, accountID, frequency1, time.Now().AddDate(0, -1, 0), nil)
	_ = repo.Save(recurringTx1)

	// Create active recurring transaction with future end date
	frequency2 := transactionvalueobjects.WeeklyFrequency()
	endDate2 := time.Now().AddDate(0, 6, 0)
	recurringTx2 := createTestRecurringTransactionEntity(t, userID, accountID, frequency2, time.Now().AddDate(0, -1, 0), &endDate2)
	_ = repo.Save(recurringTx2)

	// Create expired recurring transaction (past end date)
	frequency3 := transactionvalueobjects.DailyFrequency()
	endDate3 := time.Now().AddDate(0, -1, 0) // Past date
	recurringTx3 := createTestRecurringTransactionEntity(t, userID, accountID, frequency3, time.Now().AddDate(0, -2, 0), &endDate3)
	_ = repo.Save(recurringTx3)

	// Create non-recurring transaction
	nonRecurring := createTestTransactionEntity(t, userID, accountID)
	_ = repo.Save(nonRecurring)

	// Create instance (child transaction) - should not be returned
	parentID := recurringTx1.ID()
	instance, _ := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionvalueobjects.ExpenseType(),
		recurringTx1.Amount(),
		recurringTx1.Description(),
		time.Now(),
		false,
		nil,
		nil,
		&parentID,
	)
	_ = repo.Save(instance)

	// Find active recurring transactions
	activeRecurring, err := repo.FindActiveRecurringTransactions()
	if err != nil {
		t.Fatalf("FindActiveRecurringTransactions() error = %v", err)
	}

	// Should return only active recurring transactions (2)
	if len(activeRecurring) != 2 {
		t.Errorf("FindActiveRecurringTransactions() returned %d transactions, want 2", len(activeRecurring))
	}

	// Verify all returned transactions are recurring
	for _, tx := range activeRecurring {
		if !tx.IsRecurring() {
			t.Error("FindActiveRecurringTransactions() returned non-recurring transaction")
		}
		if tx.ParentTransactionID() != nil {
			t.Error("FindActiveRecurringTransactions() returned instance transaction")
		}
	}
}

func TestGormTransactionRepository_FindByParentIDAndDate(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create parent recurring transaction
	frequency := transactionvalueobjects.MonthlyFrequency()
	parentTx := createTestRecurringTransactionEntity(t, userID, accountID, frequency, time.Now().AddDate(0, -1, 0), nil)
	_ = repo.Save(parentTx)

	// Create instance for specific date
	parentID := parentTx.ID()
	instanceDate := time.Now()
	instance, _ := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionvalueobjects.ExpenseType(),
		parentTx.Amount(),
		parentTx.Description(),
		instanceDate,
		false,
		nil,
		nil,
		&parentID,
	)
	_ = repo.Save(instance)

	// Find instance by parent ID and date
	found, err := repo.FindByParentIDAndDate(parentID, instanceDate)
	if err != nil {
		t.Fatalf("FindByParentIDAndDate() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByParentIDAndDate() should find the instance")
	}

	if !found.ID().Equals(instance.ID()) {
		t.Errorf("FindByParentIDAndDate() returned wrong transaction")
	}

	// Try to find non-existent instance
	nonExistentDate := time.Now().AddDate(0, 1, 0)
	notFound, err := repo.FindByParentIDAndDate(parentID, nonExistentDate)
	if err != nil {
		t.Fatalf("FindByParentIDAndDate() error = %v", err)
	}

	if notFound != nil {
		t.Error("FindByParentIDAndDate() should return nil for non-existent instance")
	}
}

func TestGormTransactionRepository_SaveRecurringTransaction(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewGormTransactionRepository(db).(*GormTransactionRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.AccountID{}
	accountID, _ = accountvalueobjects.NewAccountID(accountvalueobjects.GenerateAccountID().Value())

	frequency := transactionvalueobjects.MonthlyFrequency()
	endDate := time.Now().AddDate(1, 0, 0)
	recurringTx := createTestRecurringTransactionEntity(t, userID, accountID, frequency, time.Now(), &endDate)

	// Save recurring transaction
	err := repo.Save(recurringTx)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify transaction was saved with recurrence fields
	saved, err := repo.FindByID(recurringTx.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Fatal("Save() transaction was not saved")
	}

	if !saved.IsRecurring() {
		t.Error("Save() isRecurring should be true")
	}

	if saved.RecurrenceFrequency() == nil {
		t.Error("Save() recurrenceFrequency should not be nil")
	} else if !saved.RecurrenceFrequency().Equals(frequency) {
		t.Errorf("Save() recurrenceFrequency = %v, want %v", saved.RecurrenceFrequency(), frequency)
	}

	if saved.RecurrenceEndDate() == nil {
		t.Error("Save() recurrenceEndDate should not be nil")
	} else if !saved.RecurrenceEndDate().Equal(endDate) {
		t.Errorf("Save() recurrenceEndDate = %v, want %v", saved.RecurrenceEndDate(), endDate)
	}
}
