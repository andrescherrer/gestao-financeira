package valueobjects

import (
	"strings"
	"testing"
)

func TestNewInvestmentType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid STOCK", "STOCK", false},
		{"valid FUND", "FUND", false},
		{"valid CDB", "CDB", false},
		{"valid TREASURY", "TREASURY", false},
		{"valid CRYPTO", "CRYPTO", false},
		{"valid OTHER", "OTHER", false},
		{"lowercase stock", "stock", false},
		{"mixed case", "Stock", false},
		{"invalid type", "INVALID", true},
		{"empty type", "", true},
		{"only spaces", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investmentType, err := NewInvestmentType(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NewInvestmentType() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				// Verify it's uppercase
				if investmentType.Value() != strings.ToUpper(strings.TrimSpace(tt.input)) {
					t.Errorf("NewInvestmentType() value = %v, want %v", investmentType.Value(), strings.ToUpper(strings.TrimSpace(tt.input)))
				}
			}
		})
	}
}

func TestInvestmentType_Value(t *testing.T) {
	investmentType := MustInvestmentType("STOCK")
	value := investmentType.Value()

	if value != "STOCK" {
		t.Errorf("InvestmentType.Value() = %v, want %v", value, "STOCK")
	}
}

func TestInvestmentType_String(t *testing.T) {
	investmentType := MustInvestmentType("STOCK")
	str := investmentType.String()

	if str != "STOCK" {
		t.Errorf("InvestmentType.String() = %v, want %v", str, "STOCK")
	}
}

func TestInvestmentType_DisplayName(t *testing.T) {
	tests := []struct {
		name           string
		investmentType InvestmentType
		want           string
	}{
		{"STOCK", StockType(), "Stock"},
		{"FUND", FundType(), "Fund"},
		{"CDB", CDBType(), "CDB"},
		{"TREASURY", TreasuryType(), "Treasury"},
		{"CRYPTO", CryptoType(), "Crypto"},
		{"OTHER", OtherType(), "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentType.DisplayName(); got != tt.want {
				t.Errorf("InvestmentType.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentType_IsMethods(t *testing.T) {
	tests := []struct {
		name           string
		investmentType InvestmentType
		isStock        bool
		isFund         bool
		isCDB          bool
		isTreasury     bool
		isCrypto       bool
		isOther        bool
	}{
		{"STOCK", StockType(), true, false, false, false, false, false},
		{"FUND", FundType(), false, true, false, false, false, false},
		{"CDB", CDBType(), false, false, true, false, false, false},
		{"TREASURY", TreasuryType(), false, false, false, true, false, false},
		{"CRYPTO", CryptoType(), false, false, false, false, true, false},
		{"OTHER", OtherType(), false, false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentType.IsStock(); got != tt.isStock {
				t.Errorf("InvestmentType.IsStock() = %v, want %v", got, tt.isStock)
			}
			if got := tt.investmentType.IsFund(); got != tt.isFund {
				t.Errorf("InvestmentType.IsFund() = %v, want %v", got, tt.isFund)
			}
			if got := tt.investmentType.IsCDB(); got != tt.isCDB {
				t.Errorf("InvestmentType.IsCDB() = %v, want %v", got, tt.isCDB)
			}
			if got := tt.investmentType.IsTreasury(); got != tt.isTreasury {
				t.Errorf("InvestmentType.IsTreasury() = %v, want %v", got, tt.isTreasury)
			}
			if got := tt.investmentType.IsCrypto(); got != tt.isCrypto {
				t.Errorf("InvestmentType.IsCrypto() = %v, want %v", got, tt.isCrypto)
			}
			if got := tt.investmentType.IsOther(); got != tt.isOther {
				t.Errorf("InvestmentType.IsOther() = %v, want %v", got, tt.isOther)
			}
		})
	}
}

func TestInvestmentType_RequiresQuantity(t *testing.T) {
	tests := []struct {
		name           string
		investmentType InvestmentType
		want           bool
	}{
		{"STOCK", StockType(), true},
		{"FUND", FundType(), true},
		{"CDB", CDBType(), false},
		{"TREASURY", TreasuryType(), false},
		{"CRYPTO", CryptoType(), true},
		{"OTHER", OtherType(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentType.RequiresQuantity(); got != tt.want {
				t.Errorf("InvestmentType.RequiresQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentType_HasVariableValue(t *testing.T) {
	tests := []struct {
		name           string
		investmentType InvestmentType
		want           bool
	}{
		{"STOCK", StockType(), true},
		{"FUND", FundType(), true},
		{"CDB", CDBType(), false},
		{"TREASURY", TreasuryType(), false},
		{"CRYPTO", CryptoType(), true},
		{"OTHER", OtherType(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentType.HasVariableValue(); got != tt.want {
				t.Errorf("InvestmentType.HasVariableValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentType_Equals(t *testing.T) {
	investmentType1 := MustInvestmentType("STOCK")
	investmentType2 := MustInvestmentType("FUND")
	investmentType3 := MustInvestmentType("STOCK")

	if !investmentType1.Equals(investmentType3) {
		t.Error("InvestmentType.Equals() = false for equal types, want true")
	}

	if investmentType1.Equals(investmentType2) {
		t.Error("InvestmentType.Equals() = true for different types, want false")
	}
}

func TestMustInvestmentType(t *testing.T) {
	// Valid type should not panic
	validType := "STOCK"
	investmentType := MustInvestmentType(validType)
	if investmentType.Value() != "STOCK" {
		t.Error("MustInvestmentType() returned wrong value for valid type")
	}

	// Invalid type should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustInvestmentType() should panic for invalid type")
		}
	}()
	MustInvestmentType("INVALID")
}
