package valueobjects

import (
	"testing"
)

func TestNewUserID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	invalidUUID := "invalid-uuid"
	emptyUUID := ""

	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"valid UUID", validUUID, false},
		{"invalid UUID", invalidUUID, true},
		{"empty UUID", emptyUUID, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := NewUserID(tt.id)
			if (err != nil) != tt.wantError {
				t.Errorf("NewUserID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && userID.IsEmpty() {
				t.Error("NewUserID() returned empty ID for valid UUID")
			}
		})
	}
}

func TestGenerateUserID(t *testing.T) {
	userID := GenerateUserID()

	if userID.IsEmpty() {
		t.Error("GenerateUserID() returned empty ID")
	}

	// Generate another one and verify they're different
	userID2 := GenerateUserID()
	if userID.Equals(userID2) {
		t.Error("GenerateUserID() should generate unique IDs")
	}
}

func TestUserID_Value(t *testing.T) {
	userID := GenerateUserID()
	value := userID.Value()

	if value == "" {
		t.Error("UserID.Value() returned empty string")
	}
}

func TestUserID_String(t *testing.T) {
	userID := GenerateUserID()
	str := userID.String()

	if str == "" {
		t.Error("UserID.String() returned empty string")
	}

	if str != userID.Value() {
		t.Errorf("UserID.String() = %v, want %v", str, userID.Value())
	}
}

func TestUserID_Equals(t *testing.T) {
	userID1 := GenerateUserID()
	userID2 := GenerateUserID()
	userID3 := MustUserID(userID1.Value())

	if !userID1.Equals(userID3) {
		t.Error("UserID.Equals() = false for equal IDs, want true")
	}

	if userID1.Equals(userID2) {
		t.Error("UserID.Equals() = true for different IDs, want false")
	}
}

func TestUserID_IsEmpty(t *testing.T) {
	userID := GenerateUserID()
	empty := UserID{}

	if userID.IsEmpty() {
		t.Error("UserID.IsEmpty() = true for valid ID, want false")
	}
	if !empty.IsEmpty() {
		t.Error("UserID.IsEmpty() = false for empty ID, want true")
	}
}

func TestMustUserID(t *testing.T) {
	// Valid UUID should not panic
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	userID := MustUserID(validUUID)
	if userID.IsEmpty() {
		t.Error("MustUserID() returned empty ID for valid UUID")
	}

	// Invalid UUID should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustUserID() should panic for invalid UUID")
		}
	}()
	MustUserID("invalid-uuid")
}
