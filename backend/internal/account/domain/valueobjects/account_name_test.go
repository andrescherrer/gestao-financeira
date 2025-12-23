package valueobjects

import (
	"strings"
	"testing"
)

func TestNewAccountName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid name", "Conta Corrente", false},
		{"valid name with spaces", "  Conta Poupança  ", false},
		{"empty name", "", true},
		{"only spaces", "   ", true},
		{"too short", "AB", true},
		{"minimum length", "ABC", false},
		{"too long", string(make([]byte, 101)), true},
		{"maximum length", string(make([]byte, 100)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountName, err := NewAccountName(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NewAccountName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && accountName.IsEmpty() {
				t.Error("NewAccountName() returned empty name for valid input")
			}
			if !tt.wantError {
				// Verify trimming
				trimmed := accountName.Value()
				if trimmed != strings.TrimSpace(tt.input) {
					t.Errorf("NewAccountName() trimmed value = %v, want %v", trimmed, strings.TrimSpace(tt.input))
				}
			}
		})
	}
}

func TestAccountName_Value(t *testing.T) {
	accountName := MustAccountName("Conta Corrente")
	value := accountName.Value()

	if value != "Conta Corrente" {
		t.Errorf("AccountName.Value() = %v, want %v", value, "Conta Corrente")
	}
}

func TestAccountName_String(t *testing.T) {
	accountName := MustAccountName("Conta Corrente")
	str := accountName.String()

	if str != "Conta Corrente" {
		t.Errorf("AccountName.String() = %v, want %v", str, "Conta Corrente")
	}

	if str != accountName.Value() {
		t.Errorf("AccountName.String() = %v, want %v", str, accountName.Value())
	}
}

func TestAccountName_Equals(t *testing.T) {
	accountName1 := MustAccountName("Conta Corrente")
	accountName2 := MustAccountName("Conta Poupança")
	accountName3 := MustAccountName("Conta Corrente")

	if !accountName1.Equals(accountName3) {
		t.Error("AccountName.Equals() = false for equal names, want true")
	}

	if accountName1.Equals(accountName2) {
		t.Error("AccountName.Equals() = true for different names, want false")
	}
}

func TestAccountName_IsEmpty(t *testing.T) {
	accountName := MustAccountName("Conta Corrente")
	empty := AccountName{}

	if accountName.IsEmpty() {
		t.Error("AccountName.IsEmpty() = true for valid name, want false")
	}
	if !empty.IsEmpty() {
		t.Error("AccountName.IsEmpty() = false for empty name, want true")
	}
}

func TestMustAccountName(t *testing.T) {
	// Valid name should not panic
	validName := "Conta Corrente"
	accountName := MustAccountName(validName)
	if accountName.IsEmpty() {
		t.Error("MustAccountName() returned empty name for valid input")
	}

	// Invalid name should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAccountName() should panic for invalid name")
		}
	}()
	MustAccountName("AB")
}
