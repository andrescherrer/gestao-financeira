package entities

import (
	"testing"
	"time"

	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

func TestNewCategory(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	description := "Gastos com alimentação"

	category, err := NewCategory(userID, categoryName, description)
	if err != nil {
		t.Fatalf("NewCategory() error = %v, want nil", err)
	}

	if category == nil {
		t.Fatal("NewCategory() returned nil category")
	}

	if category.ID().IsEmpty() {
		t.Error("NewCategory() returned category with empty ID")
	}

	if !category.UserID().Equals(userID) {
		t.Error("NewCategory() returned category with wrong user ID")
	}

	if !category.Name().Equals(categoryName) {
		t.Error("NewCategory() returned category with wrong name")
	}

	if category.Description() != description {
		t.Errorf("NewCategory() description = %v, want %v", category.Description(), description)
	}

	if !category.IsActive() {
		t.Error("NewCategory() returned inactive category")
	}

	if category.CreatedAt().IsZero() {
		t.Error("NewCategory() returned category with zero created_at")
	}

	if category.UpdatedAt().IsZero() {
		t.Error("NewCategory() returned category with zero updated_at")
	}

	// Check domain events
	events := category.GetEvents()
	if len(events) == 0 {
		t.Error("NewCategory() should have domain events")
	}
}

func TestNewCategory_InvalidInput(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")

	tests := []struct {
		name         string
		userID       identityvalueobjects.UserID
		categoryName valueobjects.CategoryName
		description  string
		wantError    bool
	}{
		{
			name:         "empty user ID",
			userID:       identityvalueobjects.UserID{},
			categoryName: categoryName,
			description:  "Description",
			wantError:    true,
		},
		{
			name:         "empty category name",
			userID:       userID,
			categoryName: valueobjects.CategoryName{},
			description:  "Description",
			wantError:    true,
		},
		{
			name:         "description too long",
			userID:       userID,
			categoryName: categoryName,
			description:  string(make([]byte, 501)), // 501 characters
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCategory(tt.userID, tt.categoryName, tt.description)
			if (err != nil) != tt.wantError {
				t.Errorf("NewCategory() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCategory_UpdateName(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	category, _ := NewCategory(userID, categoryName, "Description")

	newName, _ := valueobjects.NewCategoryName("Transporte")
	if err := category.UpdateName(newName); err != nil {
		t.Fatalf("UpdateName() error = %v, want nil", err)
	}

	if !category.Name().Equals(newName) {
		t.Error("UpdateName() did not update the name")
	}

	// Check domain events
	events := category.GetEvents()
	if len(events) == 0 {
		t.Error("UpdateName() should have domain events")
	}
}

func TestCategory_UpdateDescription(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	category, _ := NewCategory(userID, categoryName, "Old description")

	newDescription := "New description"
	if err := category.UpdateDescription(newDescription); err != nil {
		t.Fatalf("UpdateDescription() error = %v, want nil", err)
	}

	if category.Description() != newDescription {
		t.Errorf("UpdateDescription() description = %v, want %v", category.Description(), newDescription)
	}

	// Check domain events
	events := category.GetEvents()
	if len(events) == 0 {
		t.Error("UpdateDescription() should have domain events")
	}
}

func TestCategory_Deactivate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	category, _ := NewCategory(userID, categoryName, "Description")

	if err := category.Deactivate(); err != nil {
		t.Fatalf("Deactivate() error = %v, want nil", err)
	}

	if category.IsActive() {
		t.Error("Deactivate() did not deactivate the category")
	}

	// Check domain events
	events := category.GetEvents()
	if len(events) == 0 {
		t.Error("Deactivate() should have domain events")
	}
}

func TestCategory_Activate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	category, _ := NewCategory(userID, categoryName, "Description")

	// Deactivate first
	category.Deactivate()

	// Then activate
	if err := category.Activate(); err != nil {
		t.Fatalf("Activate() error = %v, want nil", err)
	}

	if !category.IsActive() {
		t.Error("Activate() did not activate the category")
	}

	// Check domain events
	events := category.GetEvents()
	if len(events) == 0 {
		t.Error("Activate() should have domain events")
	}
}

func TestCategoryFromPersistence(t *testing.T) {
	categoryID := valueobjects.GenerateCategoryID()
	userID := identityvalueobjects.GenerateUserID()
	categoryName, _ := valueobjects.NewCategoryName("Alimentação")
	description := "Description"
	createdAt := time.Now()
	updatedAt := time.Now()

	categorySlug := valueobjects.GenerateSlugFromName(categoryName.Value())
	category, err := CategoryFromPersistence(
		categoryID,
		userID,
		categoryName,
		categorySlug,
		description,
		createdAt,
		updatedAt,
		true,
	)

	if err != nil {
		t.Fatalf("CategoryFromPersistence() error = %v, want nil", err)
	}

	if category == nil {
		t.Fatal("CategoryFromPersistence() returned nil category")
	}

	if !category.ID().Equals(categoryID) {
		t.Error("CategoryFromPersistence() returned category with wrong ID")
	}

	if !category.UserID().Equals(userID) {
		t.Error("CategoryFromPersistence() returned category with wrong user ID")
	}

	if !category.Name().Equals(categoryName) {
		t.Error("CategoryFromPersistence() returned category with wrong name")
	}

	if category.Description() != description {
		t.Errorf("CategoryFromPersistence() description = %v, want %v", category.Description(), description)
	}

	// CategoryFromPersistence should not have domain events
	events := category.GetEvents()
	if len(events) != 0 {
		t.Error("CategoryFromPersistence() should not have domain events")
	}
}
