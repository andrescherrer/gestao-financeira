package valueobjects

import (
	"testing"
)

func TestNewInvestmentID(t *testing.T) {
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
			investmentID, err := NewInvestmentID(tt.id)
			if (err != nil) != tt.wantError {
				t.Errorf("NewInvestmentID() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && investmentID.IsEmpty() {
				t.Error("NewInvestmentID() returned empty ID for valid UUID")
			}
		})
	}
}

func TestGenerateInvestmentID(t *testing.T) {
	investmentID := GenerateInvestmentID()

	if investmentID.IsEmpty() {
		t.Error("GenerateInvestmentID() returned empty ID")
	}

	// Generate another one and verify they're different
	investmentID2 := GenerateInvestmentID()
	if investmentID.Equals(investmentID2) {
		t.Error("GenerateInvestmentID() should generate unique IDs")
	}
}

func TestInvestmentID_Value(t *testing.T) {
	investmentID := GenerateInvestmentID()
	value := investmentID.Value()

	if value == "" {
		t.Error("InvestmentID.Value() returned empty string")
	}
}

func TestInvestmentID_String(t *testing.T) {
	investmentID := GenerateInvestmentID()
	str := investmentID.String()

	if str == "" {
		t.Error("InvestmentID.String() returned empty string")
	}

	if str != investmentID.Value() {
		t.Errorf("InvestmentID.String() = %v, want %v", str, investmentID.Value())
	}
}

func TestInvestmentID_Equals(t *testing.T) {
	investmentID1 := GenerateInvestmentID()
	investmentID2 := GenerateInvestmentID()
	investmentID3 := MustInvestmentID(investmentID1.Value())

	if !investmentID1.Equals(investmentID3) {
		t.Error("InvestmentID.Equals() = false for equal IDs, want true")
	}

	if investmentID1.Equals(investmentID2) {
		t.Error("InvestmentID.Equals() = true for different IDs, want false")
	}
}

func TestInvestmentID_IsEmpty(t *testing.T) {
	investmentID := GenerateInvestmentID()
	empty := InvestmentID{}

	if investmentID.IsEmpty() {
		t.Error("InvestmentID.IsEmpty() = true for valid ID, want false")
	}
	if !empty.IsEmpty() {
		t.Error("InvestmentID.IsEmpty() = false for empty ID, want true")
	}
}

func TestMustInvestmentID(t *testing.T) {
	// Valid UUID should not panic
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	investmentID := MustInvestmentID(validUUID)
	if investmentID.IsEmpty() {
		t.Error("MustInvestmentID() returned empty ID for valid UUID")
	}

	// Invalid UUID should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustInvestmentID() should panic for invalid UUID")
		}
	}()
	MustInvestmentID("invalid-uuid")
}
