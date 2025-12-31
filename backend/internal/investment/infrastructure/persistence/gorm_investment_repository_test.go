package persistence

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/domain/entities"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
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
	err = db.AutoMigrate(&InvestmentModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestInvestmentEntity creates a test investment entity.
func createTestInvestmentEntity(t *testing.T, userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) *entities.Investment {
	investmentType := investmentvalueobjects.StockType()
	ticker := "PETR4"
	name, err := investmentvalueobjects.NewInvestmentName("Petrobras", &ticker)
	if err != nil {
		t.Fatalf("Failed to create investment name: %v", err)
	}

	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	currency, err := sharedvalueobjects.NewCurrency("BRL")
	if err != nil {
		t.Fatalf("Failed to create currency: %v", err)
	}

	purchaseAmount, err := sharedvalueobjects.NewMoneyFromFloat(1000.0, currency)
	if err != nil {
		t.Fatalf("Failed to create purchase amount: %v", err)
	}

	quantity := floatPtr(100.0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	investment, err := entities.NewInvestment(
		userID,
		accountID,
		investmentType,
		name,
		purchaseDate,
		purchaseAmount,
		quantity,
		context,
	)
	if err != nil {
		t.Fatalf("Failed to create investment: %v", err)
	}

	return investment
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestGormInvestmentRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investment := createTestInvestmentEntity(t, userID, accountID)

	// Save investment first
	err := repo.Save(investment)
	if err != nil {
		t.Fatalf("Failed to save investment: %v", err)
	}

	tests := []struct {
		name         string
		investmentID investmentvalueobjects.InvestmentID
		wantError    bool
		wantNil      bool
	}{
		{
			name:         "find existing investment",
			investmentID: investment.ID(),
			wantError:    false,
			wantNil:      false,
		},
		{
			name:         "find non-existent investment",
			investmentID: investmentvalueobjects.GenerateInvestmentID(),
			wantError:    false,
			wantNil:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.FindByID(tt.investmentID)

			if (err != nil) != tt.wantError {
				t.Errorf("FindByID() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if (result == nil) != tt.wantNil {
				t.Errorf("FindByID() result = %v, wantNil %v", result == nil, tt.wantNil)
			}

			if !tt.wantNil && result != nil {
				if result.ID() != tt.investmentID {
					t.Errorf("FindByID() ID = %v, want %v", result.ID(), tt.investmentID)
				}
			}
		})
	}
}

func TestGormInvestmentRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create and save multiple investments
	investment1 := createTestInvestmentEntity(t, userID, accountID)
	err := repo.Save(investment1)
	if err != nil {
		t.Fatalf("Failed to save investment: %v", err)
	}

	investment2 := createTestInvestmentEntity(t, userID, accountID)
	err = repo.Save(investment2)
	if err != nil {
		t.Fatalf("Failed to save investment: %v", err)
	}

	// Find by user ID
	investments, err := repo.FindByUserID(userID)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(investments) != 2 {
		t.Errorf("FindByUserID() returned %d investments, want 2", len(investments))
	}
}

func TestGormInvestmentRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investment := createTestInvestmentEntity(t, userID, accountID)

	// Save investment
	err := repo.Save(investment)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify it was saved
	saved, err := repo.FindByID(investment.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Error("Save() failed - investment not found after save")
	}

	if saved.ID() != investment.ID() {
		t.Errorf("Save() ID mismatch: got %v, want %v", saved.ID(), investment.ID())
	}
}

func TestGormInvestmentRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investment := createTestInvestmentEntity(t, userID, accountID)

	// Save investment
	err := repo.Save(investment)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete investment
	err = repo.Delete(investment.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it was deleted (soft delete)
	deleted, err := repo.FindByID(investment.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if deleted != nil {
		t.Error("Delete() failed - investment still found after delete")
	}
}

func TestGormInvestmentRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investment := createTestInvestmentEntity(t, userID, accountID)

	// Check non-existent investment
	exists, err := repo.Exists(investment.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() returned true for non-existent investment")
	}

	// Save investment
	err = repo.Save(investment)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check existing investment
	exists, err = repo.Exists(investment.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() returned false for existing investment")
	}
}

func TestGormInvestmentRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormInvestmentRepository(db).(*GormInvestmentRepository)

	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()

	// Create and save multiple investments
	investment1 := createTestInvestmentEntity(t, userID, accountID)
	err := repo.Save(investment1)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	investment2 := createTestInvestmentEntity(t, userID, accountID)
	err = repo.Save(investment2)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Count investments
	count, err := repo.Count(userID)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}

	if count != 2 {
		t.Errorf("Count() = %d, want 2", count)
	}
}
