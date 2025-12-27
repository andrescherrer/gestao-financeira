package persistence

import (
	"testing"

	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"

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
	err = db.AutoMigrate(&CategoryModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestCategory creates a test category entity.
func createTestCategory(t *testing.T, userID identityvalueobjects.UserID) *entities.Category {
	categoryName, err := valueobjects.NewCategoryName("Alimentação")
	if err != nil {
		t.Fatalf("Failed to create category name: %v", err)
	}

	category, err := entities.NewCategory(userID, categoryName, "Categoria de alimentação")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	return category
}

func TestGormCategoryRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Save category first
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Failed to save category: %v", err)
	}

	tests := []struct {
		name       string
		categoryID valueobjects.CategoryID
		wantError  bool
		wantNil    bool
	}{
		{
			name:       "find existing category",
			categoryID: category.ID(),
			wantError:  false,
			wantNil:    false,
		},
		{
			name:       "find non-existing category",
			categoryID: valueobjects.GenerateCategoryID(),
			wantError:  false,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByID(tt.categoryID)
			if (err != nil) != tt.wantError {
				t.Errorf("FindByID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("FindByID() got = %v, wantNil %v", got, tt.wantNil)
			}
			if !tt.wantNil && got != nil {
				if !got.ID().Equals(tt.categoryID) {
					t.Errorf("FindByID() category ID = %v, want %v", got.ID(), tt.categoryID)
				}
			}
		})
	}
}

func TestGormCategoryRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID1 := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()

	// Create categories for user1
	category1 := createTestCategory(t, userID1)
	category2 := createTestCategory(t, userID1)
	category2Name, _ := valueobjects.NewCategoryName("Transporte")
	category2.UpdateName(category2Name)

	// Create category for user2 (different name to avoid slug conflict)
	category3Name, _ := valueobjects.NewCategoryName("Saúde")
	category3, _ := entities.NewCategory(userID2, category3Name, "Categoria de saúde")

	// Save all categories
	err := repo.Save(category1)
	if err != nil {
		t.Fatalf("Failed to save category1: %v", err)
	}
	err = repo.Save(category2)
	if err != nil {
		t.Fatalf("Failed to save category2: %v", err)
	}
	err = repo.Save(category3)
	if err != nil {
		t.Fatalf("Failed to save category3: %v", err)
	}

	// Find categories for user1
	categories, err := repo.FindByUserID(userID1)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("FindByUserID() found %d categories, want 2", len(categories))
	}

	// Find categories for user2
	categories, err = repo.FindByUserID(userID2)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}

	if len(categories) != 1 {
		t.Errorf("FindByUserID() found %d categories, want 1", len(categories))
	}
}

func TestGormCategoryRepository_FindByUserIDAndSlug(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Save category
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Failed to save category: %v", err)
	}

	// Find by slug
	slug := category.Slug()
	found, err := repo.FindByUserIDAndSlug(userID, slug)
	if err != nil {
		t.Fatalf("FindByUserIDAndSlug() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByUserIDAndSlug() returned nil for existing category")
	}

	if !found.ID().Equals(category.ID()) {
		t.Errorf("FindByUserIDAndSlug() category ID = %v, want %v", found.ID(), category.ID())
	}

	// Find non-existing slug
	nonExistentSlug := valueobjects.GenerateSlugFromName("NonExistent")
	found, err = repo.FindByUserIDAndSlug(userID, nonExistentSlug)
	if err != nil {
		t.Fatalf("FindByUserIDAndSlug() error = %v", err)
	}

	if found != nil {
		t.Error("FindByUserIDAndSlug() should return nil for non-existing slug")
	}
}

func TestGormCategoryRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Test create
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify category was saved
	saved, err := repo.FindByID(category.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Fatal("Save() category was not saved")
	}

	if !saved.ID().Equals(category.ID()) {
		t.Errorf("Save() category ID = %v, want %v", saved.ID(), category.ID())
	}

	// Test update
	newName, _ := valueobjects.NewCategoryName("Alimentação Atualizada")
	category.UpdateName(newName)

	err = repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error on update = %v", err)
	}

	// Verify category was updated
	updated, err := repo.FindByID(category.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if updated.Name().Value() != "Alimentação Atualizada" {
		t.Errorf("Save() updated name = %v, want 'Alimentação Atualizada'", updated.Name().Value())
	}
}

func TestGormCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Save category
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete category
	err = repo.Delete(category.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify category was deleted (soft delete)
	deleted, err := repo.FindByID(category.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	// With soft delete, FindByID should return nil
	if deleted != nil {
		t.Error("Delete() category still exists after deletion")
	}
}

func TestGormCategoryRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Check non-existing category
	exists, err := repo.Exists(category.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() should return false for non-existing category")
	}

	// Save category
	err = repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check existing category
	exists, err = repo.Exists(category.ID())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing category")
	}
}

func TestGormCategoryRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()

	// Initial count should be 0
	count, err := repo.Count(userID)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Count() = %v, want 0", count)
	}

	// Create and save categories
	category1 := createTestCategory(t, userID)
	err = repo.Save(category1)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	count, err = repo.Count(userID)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 1 {
		t.Errorf("Count() = %v, want 1", count)
	}

	// Create second category
	category2 := createTestCategory(t, userID)
	category2Name, _ := valueobjects.NewCategoryName("Transporte")
	category2.UpdateName(category2Name)
	err = repo.Save(category2)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	count, err = repo.Count(userID)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 2 {
		t.Errorf("Count() = %v, want 2", count)
	}
}

func TestGormCategoryRepository_FindByUserIDAndActive(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()

	// Create and save active category
	category1 := createTestCategory(t, userID)
	err := repo.Save(category1)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Create, save, then deactivate category2
	category2Name, _ := valueobjects.NewCategoryName("Transporte")
	category2, _ := entities.NewCategory(userID, category2Name, "Categoria de transporte")
	err = repo.Save(category2)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Now deactivate and save again
	category2.Deactivate()
	err = repo.Save(category2)
	if err != nil {
		t.Fatalf("Save() error after deactivate = %v", err)
	}

	// Reload category2 to ensure state is persisted
	reloadedCategory2, err := repo.FindByID(category2.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if reloadedCategory2 == nil {
		t.Fatal("FindByID() returned nil for category2")
	}
	if reloadedCategory2.IsActive() {
		t.Error("Category2 should be inactive after Deactivate() and Save()")
	}

	// Find active categories
	activeCategories, err := repo.FindByUserIDAndActive(userID, true)
	if err != nil {
		t.Fatalf("FindByUserIDAndActive() error = %v", err)
	}

	if len(activeCategories) != 1 {
		t.Errorf("FindByUserIDAndActive() found %d active categories, want 1", len(activeCategories))
	}

	// Find inactive categories
	inactiveCategories, err := repo.FindByUserIDAndActive(userID, false)
	if err != nil {
		t.Fatalf("FindByUserIDAndActive() error = %v", err)
	}

	if len(inactiveCategories) != 1 {
		t.Errorf("FindByUserIDAndActive() found %d inactive categories, want 1", len(inactiveCategories))
	}
}

func TestGormCategoryRepository_toDomain(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Find category to trigger toDomain
	found, err := repo.FindByID(category.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByID() returned nil")
	}

	// Verify all fields are correctly converted
	if !found.ID().Equals(category.ID()) {
		t.Errorf("toDomain() ID = %v, want %v", found.ID(), category.ID())
	}

	if !found.Name().Equals(category.Name()) {
		t.Errorf("toDomain() Name = %v, want %v", found.Name(), category.Name())
	}

	if found.Description() != category.Description() {
		t.Errorf("toDomain() Description = %v, want %v", found.Description(), category.Description())
	}

	if found.IsActive() != category.IsActive() {
		t.Errorf("toDomain() IsActive = %v, want %v", found.IsActive(), category.IsActive())
	}
}

func TestGormCategoryRepository_toModel(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormCategoryRepository(db).(*GormCategoryRepository)

	userID := identityvalueobjects.GenerateUserID()
	category := createTestCategory(t, userID)

	// Save category to trigger toModel
	err := repo.Save(category)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify category was saved correctly by finding it
	found, err := repo.FindByID(category.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("toModel() category was not saved correctly")
	}

	// Verify all fields match
	if !found.ID().Equals(category.ID()) {
		t.Errorf("toModel() ID mismatch")
	}

	if !found.Name().Equals(category.Name()) {
		t.Errorf("toModel() Name mismatch")
	}

	if found.Description() != category.Description() {
		t.Errorf("toModel() Description mismatch")
	}
}
