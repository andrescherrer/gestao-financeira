package valueobjects

import (
	"strings"
	"testing"
)

func TestNewInvestmentName(t *testing.T) {
	ticker := "PETR4"
	emptyTicker := ""
	longTicker := strings.Repeat("A", 21) // 21 characters
	longName := strings.Repeat("A", 201)  // 201 characters

	tests := []struct {
		name      string
		invName   string
		ticker    *string
		wantError bool
	}{
		{"valid name without ticker", "Petrobras", nil, false},
		{"valid name with ticker", "Petrobras", &ticker, false},
		{"valid name with empty ticker", "Petrobras", &emptyTicker, false},
		{"empty name", "", nil, true},
		{"name too short", "A", nil, true},
		{"name too long", longName, nil, true},
		{"ticker too long", "Petrobras", &longTicker, true},
		{"name with spaces", "  Petrobras  ", nil, false},
		{"ticker with spaces", "  PETR4  ", &ticker, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investmentName, err := NewInvestmentName(tt.invName, tt.ticker)
			if (err != nil) != tt.wantError {
				t.Errorf("NewInvestmentName() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if investmentName.IsEmpty() {
					t.Error("NewInvestmentName() returned empty name for valid input")
				}
				if investmentName.Name() != strings.TrimSpace(tt.invName) {
					t.Errorf("NewInvestmentName() name = %v, want %v", investmentName.Name(), strings.TrimSpace(tt.invName))
				}
			}
		})
	}
}

func TestInvestmentName_Name(t *testing.T) {
	investmentName := MustInvestmentName("Petrobras", nil)
	name := investmentName.Name()

	if name != "Petrobras" {
		t.Errorf("InvestmentName.Name() = %v, want %v", name, "Petrobras")
	}
}

func TestInvestmentName_Ticker(t *testing.T) {
	ticker := "PETR4"
	investmentNameWithTicker := MustInvestmentName("Petrobras", &ticker)
	investmentNameWithoutTicker := MustInvestmentName("Petrobras", nil)

	if !investmentNameWithTicker.HasTicker() {
		t.Error("InvestmentName.HasTicker() = false for name with ticker, want true")
	}

	if investmentNameWithoutTicker.HasTicker() {
		t.Error("InvestmentName.HasTicker() = true for name without ticker, want false")
	}

	if investmentNameWithTicker.Ticker() == nil {
		t.Error("InvestmentName.Ticker() returned nil for name with ticker")
	} else if *investmentNameWithTicker.Ticker() != "PETR4" {
		t.Errorf("InvestmentName.Ticker() = %v, want %v", *investmentNameWithTicker.Ticker(), "PETR4")
	}
}

func TestInvestmentName_String(t *testing.T) {
	ticker := "PETR4"
	investmentNameWithTicker := MustInvestmentName("Petrobras", &ticker)
	investmentNameWithoutTicker := MustInvestmentName("Petrobras", nil)

	expectedWithTicker := "Petrobras (PETR4)"
	if investmentNameWithTicker.String() != expectedWithTicker {
		t.Errorf("InvestmentName.String() = %v, want %v", investmentNameWithTicker.String(), expectedWithTicker)
	}

	expectedWithoutTicker := "Petrobras"
	if investmentNameWithoutTicker.String() != expectedWithoutTicker {
		t.Errorf("InvestmentName.String() = %v, want %v", investmentNameWithoutTicker.String(), expectedWithoutTicker)
	}
}

func TestInvestmentName_Equals(t *testing.T) {
	ticker1 := "PETR4"
	ticker2 := "VALE3"
	investmentName1 := MustInvestmentName("Petrobras", &ticker1)
	investmentName2 := MustInvestmentName("Petrobras", &ticker1)
	investmentName3 := MustInvestmentName("Petrobras", &ticker2)
	investmentName4 := MustInvestmentName("Petrobras", nil)
	investmentName5 := MustInvestmentName("Vale", &ticker1)

	if !investmentName1.Equals(investmentName2) {
		t.Error("InvestmentName.Equals() = false for equal names, want true")
	}

	if investmentName1.Equals(investmentName3) {
		t.Error("InvestmentName.Equals() = true for different tickers, want false")
	}

	if investmentName1.Equals(investmentName4) {
		t.Error("InvestmentName.Equals() = true for one with ticker and one without, want false")
	}

	if investmentName1.Equals(investmentName5) {
		t.Error("InvestmentName.Equals() = true for different names, want false")
	}
}

func TestInvestmentName_IsEmpty(t *testing.T) {
	investmentName := MustInvestmentName("Petrobras", nil)
	empty := InvestmentName{}

	if investmentName.IsEmpty() {
		t.Error("InvestmentName.IsEmpty() = true for valid name, want false")
	}
	if !empty.IsEmpty() {
		t.Error("InvestmentName.IsEmpty() = false for empty name, want true")
	}
}

func TestMustInvestmentName(t *testing.T) {
	// Valid name should not panic
	validName := "Petrobras"
	ticker := "PETR4"
	investmentName := MustInvestmentName(validName, &ticker)
	if investmentName.IsEmpty() {
		t.Error("MustInvestmentName() returned empty name for valid input")
	}

	// Invalid name should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustInvestmentName() should panic for invalid name")
		}
	}()
	MustInvestmentName("", nil)
}

