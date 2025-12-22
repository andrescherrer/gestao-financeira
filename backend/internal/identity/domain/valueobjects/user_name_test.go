package valueobjects

import (
	"strings"
	"testing"
)

func TestNewUserName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantError bool
	}{
		{"valid name", "John", "Doe", false},
		{"single character first name", "J", "Doe", true},
		{"empty first name", "", "Doe", true},
		{"empty last name", "John", "", true},
		{"name with hyphen", "Jean-Pierre", "Dupont", false},
		{"name with apostrophe", "O'Brien", "Smith", false},
		{"name with spaces", "Mary Jane", "Watson", false},
		{"too long first name", strings.Repeat("a", 101), "Doe", true},
		{"too long last name", "John", strings.Repeat("a", 101), true},
		{"name with numbers", "John123", "Doe", true},
		{"name with special chars", "John@", "Doe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userName, err := NewUserName(tt.firstName, tt.lastName)
			if (err != nil) != tt.wantError {
				t.Errorf("NewUserName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && userName.IsEmpty() {
				t.Error("NewUserName() returned empty name for valid input")
			}
		})
	}
}

func TestNewUserNameFromFullName(t *testing.T) {
	tests := []struct {
		name      string
		fullName  string
		wantFirst string
		wantLast  string
		wantError bool
	}{
		{"full name", "John Doe", "John", "Doe", false},
		{"three words", "Mary Jane Watson", "Mary", "Jane Watson", false},
		{"single name", "Madonna", "Madonna", "", false},
		{"empty name", "", "", "", true},
		{"name with extra spaces", "  John   Doe  ", "John", "Doe", false},
		{"name with hyphen", "Jean-Pierre Dupont", "Jean-Pierre", "Dupont", false},
		{"too long name", strings.Repeat("a", 201), "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userName, err := NewUserNameFromFullName(tt.fullName)
			if (err != nil) != tt.wantError {
				t.Errorf("NewUserNameFromFullName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if userName.FirstName() != tt.wantFirst {
					t.Errorf("NewUserNameFromFullName() firstName = %v, want %v", userName.FirstName(), tt.wantFirst)
				}
				if userName.LastName() != tt.wantLast {
					t.Errorf("NewUserNameFromFullName() lastName = %v, want %v", userName.LastName(), tt.wantLast)
				}
			}
		})
	}
}

func TestUserName_FirstName(t *testing.T) {
	userName, _ := NewUserName("John", "Doe")
	if userName.FirstName() != "John" {
		t.Errorf("UserName.FirstName() = %v, want John", userName.FirstName())
	}
}

func TestUserName_LastName(t *testing.T) {
	userName, _ := NewUserName("John", "Doe")
	if userName.LastName() != "Doe" {
		t.Errorf("UserName.LastName() = %v, want Doe", userName.LastName())
	}
}

func TestUserName_FullName(t *testing.T) {
	userName, _ := NewUserName("John", "Doe")
	if userName.FullName() != "John Doe" {
		t.Errorf("UserName.FullName() = %v, want John Doe", userName.FullName())
	}
}

func TestUserName_String(t *testing.T) {
	userName, _ := NewUserName("John", "Doe")
	if userName.String() != "John Doe" {
		t.Errorf("UserName.String() = %v, want John Doe", userName.String())
	}
}

func TestUserName_Initials(t *testing.T) {
	singleName, _ := NewUserNameFromFullName("Madonna")

	tests := []struct {
		name     string
		userName UserName
		want     string
	}{
		{"full name", MustUserName("John", "Doe"), "JD"},
		{"single name", singleName, "M"},
		{"lowercase", MustUserName("john", "doe"), "JD"},
		{"with hyphen", MustUserName("Jean-Pierre", "Dupont"), "JD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.userName.Initials(); got != tt.want {
				t.Errorf("UserName.Initials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserName_Equals(t *testing.T) {
	userName1, _ := NewUserName("John", "Doe")
	userName2, _ := NewUserName("John", "Doe")
	userName3, _ := NewUserName("Jane", "Doe")

	tests := []struct {
		name string
		u1   UserName
		u2   UserName
		want bool
	}{
		{"equal names", userName1, userName2, true},
		{"different names", userName1, userName3, false},
		{"same instance", userName1, userName1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u1.Equals(tt.u2); got != tt.want {
				t.Errorf("UserName.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserName_IsEmpty(t *testing.T) {
	userName, _ := NewUserName("John", "Doe")
	empty := UserName{}

	if userName.IsEmpty() {
		t.Error("UserName.IsEmpty() = true for valid name, want false")
	}
	if !empty.IsEmpty() {
		t.Error("UserName.IsEmpty() = false for empty name, want true")
	}
}

func TestUserName_HasLastName(t *testing.T) {
	userName1, _ := NewUserName("John", "Doe")
	userName2, _ := NewUserNameFromFullName("Madonna")

	if !userName1.HasLastName() {
		t.Error("UserName.HasLastName() = false for name with last name, want true")
	}
	if userName2.HasLastName() {
		t.Error("UserName.HasLastName() = true for name without last name, want false")
	}
}

func TestUserName_EdgeCases(t *testing.T) {
	// Test trimming whitespace
	userName, _ := NewUserName("  John  ", "  Doe  ")
	if userName.FirstName() != "John" {
		t.Errorf("UserName should trim whitespace, got firstName = %v", userName.FirstName())
	}
	if userName.LastName() != "Doe" {
		t.Errorf("UserName should trim whitespace, got lastName = %v", userName.LastName())
	}

	// Test minimum length
	_, err := NewUserName("Jo", "Do")
	if err != nil {
		t.Errorf("UserName should accept 2-character names, got error: %v", err)
	}

	// Test invalid starting/ending characters
	_, err = NewUserName("-John", "Doe")
	if err == nil {
		t.Error("UserName should reject names starting with hyphen")
	}

	_, err = NewUserName("John-", "Doe")
	if err == nil {
		t.Error("UserName should reject names ending with hyphen")
	}
}

func TestMustUserName(t *testing.T) {
	// Valid name should not panic
	userName := MustUserName("John", "Doe")
	if userName.IsEmpty() {
		t.Error("MustUserName() returned empty name for valid input")
	}

	// Invalid name should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustUserName() should panic for invalid name")
		}
	}()
	MustUserName("", "Doe")
}
