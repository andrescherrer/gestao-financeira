package valueobjects

import (
	"testing"
)

func TestNewAccountID(t *testing.T) {
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
			accountID, err := NewAccountID(tt.id)
			if (err != nil) != tt.wantError {
				t.Errorf("NewAccountID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && accountID.IsEmpty() {
				t.Error("NewAccountID() returned empty ID for valid UUID")
			}
		})
	}
}

func TestGenerateAccountID(t *testing.T) {
	accountID := GenerateAccountID()

	if accountID.IsEmpty() {
		t.Error("GenerateAccountID() returned empty ID")
	}

	// Generate another one and verify they're different
	accountID2 := GenerateAccountID()
	if accountID.Equals(accountID2) {
		t.Error("GenerateAccountID() should generate unique IDs")
	}
}

func TestAccountID_Value(t *testing.T) {
	accountID := GenerateAccountID()
	value := accountID.Value()

	if value == "" {
		t.Error("AccountID.Value() returned empty string")
	}
}

func TestAccountID_String(t *testing.T) {
	accountID := GenerateAccountID()
	str := accountID.String()

	if str == "" {
		t.Error("AccountID.String() returned empty string")
	}

	if str != accountID.Value() {
		t.Errorf("AccountID.String() = %v, want %v", str, accountID.Value())
	}
}

func TestAccountID_Equals(t *testing.T) {
	accountID1 := GenerateAccountID()
	accountID2 := GenerateAccountID()
	accountID3 := MustAccountID(accountID1.Value())

	if !accountID1.Equals(accountID3) {
		t.Error("AccountID.Equals() = false for equal IDs, want true")
	}

	if accountID1.Equals(accountID2) {
		t.Error("AccountID.Equals() = true for different IDs, want false")
	}
}

func TestAccountID_IsEmpty(t *testing.T) {
	accountID := GenerateAccountID()
	empty := AccountID{}

	if accountID.IsEmpty() {
		t.Error("AccountID.IsEmpty() = true for valid ID, want false")
	}
	if !empty.IsEmpty() {
		t.Error("AccountID.IsEmpty() = false for empty ID, want true")
	}
}

func TestMustAccountID(t *testing.T) {
	// Valid UUID should not panic
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := MustAccountID(validUUID)
	if accountID.IsEmpty() {
		t.Error("MustAccountID() returned empty ID for valid UUID")
	}

	// Invalid UUID should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAccountID() should panic for invalid UUID")
		}
	}()
	MustAccountID("invalid-uuid")
}
