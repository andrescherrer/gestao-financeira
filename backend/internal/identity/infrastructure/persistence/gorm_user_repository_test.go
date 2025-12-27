package persistence

import (
	"testing"

	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"

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
	err = db.AutoMigrate(&UserModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestUser creates a test user entity.
func createTestUser(t *testing.T) *entities.User {
	email, err := valueobjects.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	passwordHash, err := valueobjects.NewPasswordHashFromPlain("password123")
	if err != nil {
		t.Fatalf("Failed to create password hash: %v", err)
	}

	name, err := valueobjects.NewUserName("John", "Doe")
	if err != nil {
		t.Fatalf("Failed to create name: %v", err)
	}

	user, err := entities.NewUser(email, passwordHash, name)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	return user
}

func TestGormUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Save user first
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	tests := []struct {
		name      string
		userID    valueobjects.UserID
		wantError bool
		wantNil   bool
	}{
		{
			name:      "find existing user",
			userID:    user.ID(),
			wantError: false,
			wantNil:   false,
		},
		{
			name:      "find non-existing user",
			userID:    valueobjects.GenerateUserID(),
			wantError: false,
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByID(tt.userID)
			if (err != nil) != tt.wantError {
				t.Errorf("FindByID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("FindByID() got = %v, wantNil %v", got, tt.wantNil)
			}
			if !tt.wantNil && got != nil {
				if !got.ID().Equals(tt.userID) {
					t.Errorf("FindByID() user ID = %v, want %v", got.ID(), tt.userID)
				}
			}
		})
	}
}

func TestGormUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Save user first
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	tests := []struct {
		name      string
		email     valueobjects.Email
		wantError bool
		wantNil   bool
	}{
		{
			name:      "find existing user by email",
			email:     user.Email(),
			wantError: false,
			wantNil:   false,
		},
		{
			name: "find non-existing user by email",
			email: func() valueobjects.Email {
				email, _ := valueobjects.NewEmail("nonexistent@example.com")
				return email
			}(),
			wantError: false,
			wantNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByEmail(tt.email)
			if (err != nil) != tt.wantError {
				t.Errorf("FindByEmail() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("FindByEmail() got = %v, wantNil %v", got, tt.wantNil)
			}
			if !tt.wantNil && got != nil {
				if !got.Email().Equals(tt.email) {
					t.Errorf("FindByEmail() user email = %v, want %v", got.Email(), tt.email)
				}
			}
		})
	}
}

func TestGormUserRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Test create
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify user was saved
	saved, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if saved == nil {
		t.Fatal("Save() user was not saved")
	}

	if !saved.ID().Equals(user.ID()) {
		t.Errorf("Save() user ID = %v, want %v", saved.ID(), user.ID())
	}

	// Test update
	newName, _ := valueobjects.NewUserName("Jane", "Smith")
	user.UpdateName(newName)

	err = repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error on update = %v", err)
	}

	// Verify user was updated
	updated, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if updated.Name().FirstName() != "Jane" {
		t.Errorf("Save() updated first name = %v, want 'Jane'", updated.Name().FirstName())
	}

	if updated.Name().LastName() != "Smith" {
		t.Errorf("Save() updated last name = %v, want 'Smith'", updated.Name().LastName())
	}
}

func TestGormUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Save user
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete user
	err = repo.Delete(user.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify user was deleted
	deleted, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if deleted != nil {
		t.Error("Delete() user still exists after deletion")
	}
}

func TestGormUserRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Check non-existing user
	exists, err := repo.Exists(user.Email())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() should return false for non-existing user")
	}

	// Save user
	err = repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check existing user
	exists, err = repo.Exists(user.Email())
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing user")
	}
}

func TestGormUserRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	// Initial count should be 0
	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Count() = %v, want 0", count)
	}

	// Create and save users
	user1 := createTestUser(t)
	err = repo.Save(user1)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	count, err = repo.Count()
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 1 {
		t.Errorf("Count() = %v, want 1", count)
	}

	// Create second user
	user2 := createTestUser(t)
	email2, _ := valueobjects.NewEmail("user2@example.com")
	user2.UpdateEmail(email2)
	err = repo.Save(user2)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	count, err = repo.Count()
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 2 {
		t.Errorf("Count() = %v, want 2", count)
	}
}

func TestGormUserRepository_toDomain(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Find user to trigger toDomain
	found, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByID() returned nil")
	}

	// Verify all fields are correctly converted
	if !found.ID().Equals(user.ID()) {
		t.Errorf("toDomain() ID = %v, want %v", found.ID(), user.ID())
	}

	if !found.Email().Equals(user.Email()) {
		t.Errorf("toDomain() Email = %v, want %v", found.Email(), user.Email())
	}

	if found.Name().FirstName() != user.Name().FirstName() {
		t.Errorf("toDomain() FirstName = %v, want %v", found.Name().FirstName(), user.Name().FirstName())
	}

	if found.Name().LastName() != user.Name().LastName() {
		t.Errorf("toDomain() LastName = %v, want %v", found.Name().LastName(), user.Name().LastName())
	}

	if found.IsActive() != user.IsActive() {
		t.Errorf("toDomain() IsActive = %v, want %v", found.IsActive(), user.IsActive())
	}
}

func TestGormUserRepository_toModel(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormUserRepository(db).(*GormUserRepository)

	user := createTestUser(t)

	// Save user to trigger toModel
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify user was saved correctly by finding it
	found, err := repo.FindByID(user.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("toModel() user was not saved correctly")
	}

	// Verify all fields match
	if !found.ID().Equals(user.ID()) {
		t.Errorf("toModel() ID mismatch")
	}

	if !found.Email().Equals(user.Email()) {
		t.Errorf("toModel() Email mismatch")
	}

	if found.Name().FirstName() != user.Name().FirstName() {
		t.Errorf("toModel() FirstName mismatch")
	}

	if found.Name().LastName() != user.Name().LastName() {
		t.Errorf("toModel() LastName mismatch")
	}
}
