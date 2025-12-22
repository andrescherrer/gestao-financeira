package entities

import (
	"testing"
	"time"

	"gestao-financeira/backend/internal/identity/domain/valueobjects"
)

func TestNewUser(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	tests := []struct {
		name         string
		email        valueobjects.Email
		passwordHash valueobjects.PasswordHash
		nameVO       valueobjects.UserName
		wantError    bool
	}{
		{"valid user", email, passwordHash, name, false},
		{"empty email", valueobjects.Email{}, passwordHash, name, true},
		{"empty password hash", email, valueobjects.PasswordHash{}, name, true},
		{"empty name", email, passwordHash, valueobjects.UserName{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.passwordHash, tt.nameVO)
			if (err != nil) != tt.wantError {
				t.Errorf("NewUser() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if user.ID().IsEmpty() {
					t.Error("NewUser() returned user with empty ID")
				}
				if !user.IsActive() {
					t.Error("NewUser() returned inactive user")
				}
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("oldpassword123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)

	tests := []struct {
		name        string
		oldPassword string
		newPassword string
		wantError   bool
	}{
		{"correct old password", "oldpassword123", "newpassword123", false},
		{"incorrect old password", "wrongpassword", "newpassword123", true},
		{"invalid new password", "oldpassword123", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ChangePassword(tt.oldPassword, tt.newPassword)
			if (err != nil) != tt.wantError {
				t.Errorf("User.ChangePassword() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				// Verify new password works
				if !user.VerifyPassword(tt.newPassword) {
					t.Error("User.ChangePassword() new password does not verify")
				}
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)
	newName, _ := valueobjects.NewUserName("Jane", "Smith")

	err := user.UpdateName(newName)
	if err != nil {
		t.Errorf("User.UpdateName() error = %v, want nil", err)
	}

	if !user.Name().Equals(newName) {
		t.Error("User.UpdateName() name was not updated")
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)
	newEmail, _ := valueobjects.NewEmail("newuser@example.com")

	err := user.UpdateEmail(newEmail)
	if err != nil {
		t.Errorf("User.UpdateEmail() error = %v, want nil", err)
	}

	if !user.Email().Equals(newEmail) {
		t.Error("User.UpdateEmail() email was not updated")
	}
}

func TestUser_Deactivate(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)

	if !user.IsActive() {
		t.Error("User should be active after creation")
	}

	err := user.Deactivate()
	if err != nil {
		t.Errorf("User.Deactivate() error = %v, want nil", err)
	}

	if user.IsActive() {
		t.Error("User should be inactive after Deactivate()")
	}

	// Try to deactivate again
	err = user.Deactivate()
	if err == nil {
		t.Error("User.Deactivate() should return error when already inactive")
	}
}

func TestUser_Activate(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)
	user.Deactivate()

	err := user.Activate()
	if err != nil {
		t.Errorf("User.Activate() error = %v, want nil", err)
	}

	if !user.IsActive() {
		t.Error("User should be active after Activate()")
	}

	// Try to activate again
	err = user.Activate()
	if err == nil {
		t.Error("User.Activate() should return error when already active")
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)

	if !user.VerifyPassword("password123") {
		t.Error("User.VerifyPassword() = false for correct password, want true")
	}

	if user.VerifyPassword("wrongpassword") {
		t.Error("User.VerifyPassword() = true for incorrect password, want false")
	}
}

func TestUser_DomainEvents(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)

	// Check that UserRegistered event was created
	events := user.GetEvents()
	if len(events) == 0 {
		t.Error("NewUser() should create UserRegistered event")
	}

	if events[0].EventType() != "UserRegistered" {
		t.Errorf("NewUser() event type = %v, want UserRegistered", events[0].EventType())
	}

	// Clear events
	user.ClearEvents()
	if len(user.GetEvents()) != 0 {
		t.Error("User.ClearEvents() did not clear events")
	}

	// Change password should create event
	user.ChangePassword("password123", "newpassword123")
	events = user.GetEvents()
	if len(events) == 0 {
		t.Error("User.ChangePassword() should create UserPasswordChanged event")
	}
}

func TestUser_Timestamps(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	beforeCreation := time.Now()
	user, _ := NewUser(email, passwordHash, name)
	afterCreation := time.Now()

	if user.CreatedAt().Before(beforeCreation) || user.CreatedAt().After(afterCreation) {
		t.Error("User.CreatedAt() should be set to current time")
	}

	if user.UpdatedAt().Before(beforeCreation) || user.UpdatedAt().After(afterCreation) {
		t.Error("User.UpdatedAt() should be set to current time on creation")
	}

	// Update name and check updatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	newName, _ := valueobjects.NewUserName("Jane", "Smith")
	user.UpdateName(newName)

	if !user.UpdatedAt().After(user.CreatedAt()) {
		t.Error("User.UpdatedAt() should be updated after UpdateName()")
	}
}

func TestUser_InactiveUserOperations(t *testing.T) {
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")

	user, _ := NewUser(email, passwordHash, name)
	user.Deactivate()

	// Try to change password
	err := user.ChangePassword("password123", "newpassword123")
	if err == nil {
		t.Error("User.ChangePassword() should return error for inactive user")
	}

	// Try to update name
	newName, _ := valueobjects.NewUserName("Jane", "Smith")
	err = user.UpdateName(newName)
	if err == nil {
		t.Error("User.UpdateName() should return error for inactive user")
	}

	// Try to update email
	newEmail, _ := valueobjects.NewEmail("newuser@example.com")
	err = user.UpdateEmail(newEmail)
	if err == nil {
		t.Error("User.UpdateEmail() should return error for inactive user")
	}
}
