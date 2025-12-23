package persistence

import (
	"testing"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&AccountModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestAccount creates a test account entity.
func createTestAccountEntity(t *testing.T, userID identityvalueobjects.UserID) *entities.Account {
	accountName, err := valueobjects.NewAccountName("Conta Corrente")
	if err != nil {
		t.Fatalf("Failed to create account name: %v", err)
	}

	accountType := valueobjects.BankType()
	currency, err := sharedvalueobjects.NewCurrency("BRL")
	if err != nil {
		t.Fatalf("Failed to create currency: %v", err)
	}

	balance, err := sharedvalueobjects.NewMoney(100000, currency) // 1000.00 BRL
	if err != nil {
		t.Fatalf("Failed to create balance: %v", err)
	}

	context := sharedvalueobjects.PersonalContext()

	account, err := entities.NewAccount(userID, accountName, accountType, balance, context)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	return account
}

func TestGormAccountRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Save account first
	err := repo.Save(account)
	if err != nil {
		t.Fatalf("Failed to save account: %v", err)
	}

	tests := []struct {
		name      string
		accountID valueobjects.AccountID
		wantError bool
		wantNil   bool
	}{
		{
			name:      "find existing account",
			accountID: account.ID(),
			wantError: false,
			wantNil:   false,
		},
		{
			name:      "find non-existent account",
			accountID: valueobjects.GenerateAccountID(),
			wantError: false,
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByID(tt.accountID)

			if (err != nil) != tt.wantError {
				t.Errorf("FindByID() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantNil && result != nil {
				t.Errorf("FindByID() = %v, want nil", result)
			}

			if !tt.wantNil && result == nil {
				t.Errorf("FindByID() = nil, want account")
			}

			if !tt.wantNil && result != nil {
				if !result.ID().Equals(account.ID()) {
					t.Errorf("FindByID() account ID = %v, want %v", result.ID(), account.ID())
				}
			}
		})
	}
}

func TestGormAccountRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID1 := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()

	// Create accounts for user1
	account1 := createTestAccountEntity(t, userID1)
	account2 := createTestAccountEntity(t, userID1)
	account3 := createTestAccountEntity(t, userID2)

	// Save accounts
	_ = repo.Save(account1)
	_ = repo.Save(account2)
	_ = repo.Save(account3)

	// Find accounts for user1
	accounts, err := repo.FindByUserID(userID1)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("FindByUserID() returned %d accounts, want 2", len(accounts))
	}

	// Find accounts for user2
	accounts, err = repo.FindByUserID(userID2)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("FindByUserID() returned %d accounts, want 1", len(accounts))
	}

	// Find accounts for non-existent user
	userID3 := identityvalueobjects.GenerateUserID()
	accounts, err = repo.FindByUserID(userID3)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(accounts) != 0 {
		t.Errorf("FindByUserID() returned %d accounts, want 0", len(accounts))
	}
}

func TestGormAccountRepository_FindByUserIDAndContext(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()

	// Create accounts with different contexts
	accountName1, _ := valueobjects.NewAccountName("Conta Pessoal")
	accountName2, _ := valueobjects.NewAccountName("Conta Empresarial")
	accountType := valueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency)
	personalContext := sharedvalueobjects.PersonalContext()
	businessContext := sharedvalueobjects.BusinessContext()

	account1, _ := entities.NewAccount(userID, accountName1, accountType, balance, personalContext)
	account2, _ := entities.NewAccount(userID, accountName2, accountType, balance, businessContext)
	account3, _ := entities.NewAccount(userID, accountName1, accountType, balance, personalContext)

	// Save accounts
	_ = repo.Save(account1)
	_ = repo.Save(account2)
	_ = repo.Save(account3)

	// Find personal accounts
	accounts, err := repo.FindByUserIDAndContext(userID, personalContext)
	if err != nil {
		t.Fatalf("FindByUserIDAndContext() error = %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("FindByUserIDAndContext() returned %d accounts, want 2", len(accounts))
	}

	// Find business accounts
	accounts, err = repo.FindByUserIDAndContext(userID, businessContext)
	if err != nil {
		t.Fatalf("FindByUserIDAndContext() error = %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("FindByUserIDAndContext() returned %d accounts, want 1", len(accounts))
	}
}

func TestGormAccountRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Test create
	err := repo.Save(account)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify account was saved
	saved, err := repo.FindByID(account.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Fatal("Save() account was not saved")
	}

	if !saved.ID().Equals(account.ID()) {
		t.Errorf("Save() account ID = %v, want %v", saved.ID(), account.ID())
	}

	// Test update
	newName, _ := valueobjects.NewAccountName("Conta Atualizada")
	account.UpdateName(newName)

	err = repo.Save(account)
	if err != nil {
		t.Fatalf("Save() error on update = %v", err)
	}

	// Verify account was updated
	updated, err := repo.FindByID(account.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if updated.Name().Value() != "Conta Atualizada" {
		t.Errorf("Save() updated name = %v, want 'Conta Atualizada'", updated.Name().Value())
	}
}

func TestGormAccountRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Save account
	err := repo.Save(account)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete account
	err = repo.Delete(account.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify account was deleted (soft delete)
	deleted, err := repo.FindByID(account.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	// With soft delete, FindByID should return nil
	if deleted != nil {
		t.Error("Delete() account still exists after deletion")
	}
}

func TestGormAccountRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Account should not exist yet
	exists, err := repo.Exists(account.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if exists {
		t.Error("Exists() returned true for non-existent account")
	}

	// Save account
	err = repo.Save(account)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Account should exist now
	exists, err = repo.Exists(account.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if !exists {
		t.Error("Exists() returned false for existing account")
	}
}

func TestGormAccountRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID1 := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()

	// Create accounts for user1
	account1 := createTestAccountEntity(t, userID1)
	account2 := createTestAccountEntity(t, userID1)
	account3 := createTestAccountEntity(t, userID2)

	// Save accounts
	_ = repo.Save(account1)
	_ = repo.Save(account2)
	_ = repo.Save(account3)

	// Count accounts for user1
	count, err := repo.Count(userID1)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 2 {
		t.Errorf("Count() = %d, want 2", count)
	}

	// Count accounts for user2
	count, err = repo.Count(userID2)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 1 {
		t.Errorf("Count() = %d, want 1", count)
	}

	// Count accounts for non-existent user
	userID3 := identityvalueobjects.GenerateUserID()
	count, err = repo.Count(userID3)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 0 {
		t.Errorf("Count() = %d, want 0", count)
	}
}

func TestGormAccountRepository_toDomain(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Convert to model
	model := repo.toModel(account)

	// Convert back to domain
	domainAccount, err := repo.toDomain(model)
	if err != nil {
		t.Fatalf("toDomain() error = %v", err)
	}

	// Verify conversion
	if !domainAccount.ID().Equals(account.ID()) {
		t.Errorf("toDomain() account ID = %v, want %v", domainAccount.ID(), account.ID())
	}

	if !domainAccount.UserID().Equals(account.UserID()) {
		t.Errorf("toDomain() user ID = %v, want %v", domainAccount.UserID(), account.UserID())
	}

	if domainAccount.Name().Value() != account.Name().Value() {
		t.Errorf("toDomain() name = %v, want %v", domainAccount.Name().Value(), account.Name().Value())
	}

	if !domainAccount.Balance().Equals(account.Balance()) {
		t.Errorf("toDomain() balance = %v, want %v", domainAccount.Balance(), account.Balance())
	}
}

func TestGormAccountRepository_toModel(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormAccountRepository(db).(*GormAccountRepository)

	userID := identityvalueobjects.GenerateUserID()
	account := createTestAccountEntity(t, userID)

	// Convert to model
	model := repo.toModel(account)

	// Verify model fields
	if model.ID != account.ID().Value() {
		t.Errorf("toModel() ID = %v, want %v", model.ID, account.ID().Value())
	}

	if model.UserID != account.UserID().Value() {
		t.Errorf("toModel() UserID = %v, want %v", model.UserID, account.UserID().Value())
	}

	if model.Name != account.Name().Value() {
		t.Errorf("toModel() Name = %v, want %v", model.Name, account.Name().Value())
	}

	if model.Balance != account.Balance().Amount() {
		t.Errorf("toModel() Balance = %v, want %v", model.Balance, account.Balance().Amount())
	}

	if model.Currency != account.Balance().Currency().Code() {
		t.Errorf("toModel() Currency = %v, want %v", model.Currency, account.Balance().Currency().Code())
	}
}
