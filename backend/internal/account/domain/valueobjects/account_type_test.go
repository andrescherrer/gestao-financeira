package valueobjects

import (
	"strings"
	"testing"
)

func TestNewAccountType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid BANK", "BANK", false},
		{"valid WALLET", "WALLET", false},
		{"valid INVESTMENT", "INVESTMENT", false},
		{"valid CREDIT_CARD", "CREDIT_CARD", false},
		{"lowercase bank", "bank", false},
		{"mixed case", "Bank", false},
		{"invalid type", "INVALID", true},
		{"empty type", "", true},
		{"only spaces", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountType, err := NewAccountType(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NewAccountType() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				// Verify it's uppercase
				if accountType.Value() != strings.ToUpper(strings.TrimSpace(tt.input)) {
					t.Errorf("NewAccountType() value = %v, want %v", accountType.Value(), strings.ToUpper(strings.TrimSpace(tt.input)))
				}
			}
		})
	}
}

func TestAccountType_Value(t *testing.T) {
	accountType := MustAccountType("BANK")
	value := accountType.Value()

	if value != "BANK" {
		t.Errorf("AccountType.Value() = %v, want %v", value, "BANK")
	}
}

func TestAccountType_String(t *testing.T) {
	accountType := MustAccountType("BANK")
	str := accountType.String()

	if str != "BANK" {
		t.Errorf("AccountType.String() = %v, want %v", str, "BANK")
	}
}

func TestAccountType_DisplayName(t *testing.T) {
	tests := []struct {
		name        string
		accountType AccountType
		want        string
	}{
		{"BANK", BankType(), "Bank"},
		{"WALLET", WalletType(), "Wallet"},
		{"INVESTMENT", InvestmentType(), "Investment"},
		{"CREDIT_CARD", CreditCardType(), "Credit Card"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.accountType.DisplayName(); got != tt.want {
				t.Errorf("AccountType.DisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountType_IsMethods(t *testing.T) {
	tests := []struct {
		name         string
		accountType  AccountType
		isBank       bool
		isWallet     bool
		isInvestment bool
		isCreditCard bool
	}{
		{"BANK", BankType(), true, false, false, false},
		{"WALLET", WalletType(), false, true, false, false},
		{"INVESTMENT", InvestmentType(), false, false, true, false},
		{"CREDIT_CARD", CreditCardType(), false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.accountType.IsBank(); got != tt.isBank {
				t.Errorf("AccountType.IsBank() = %v, want %v", got, tt.isBank)
			}
			if got := tt.accountType.IsWallet(); got != tt.isWallet {
				t.Errorf("AccountType.IsWallet() = %v, want %v", got, tt.isWallet)
			}
			if got := tt.accountType.IsInvestment(); got != tt.isInvestment {
				t.Errorf("AccountType.IsInvestment() = %v, want %v", got, tt.isInvestment)
			}
			if got := tt.accountType.IsCreditCard(); got != tt.isCreditCard {
				t.Errorf("AccountType.IsCreditCard() = %v, want %v", got, tt.isCreditCard)
			}
		})
	}
}

func TestAccountType_Equals(t *testing.T) {
	accountType1 := MustAccountType("BANK")
	accountType2 := MustAccountType("WALLET")
	accountType3 := MustAccountType("BANK")

	if !accountType1.Equals(accountType3) {
		t.Error("AccountType.Equals() = false for equal types, want true")
	}

	if accountType1.Equals(accountType2) {
		t.Error("AccountType.Equals() = true for different types, want false")
	}
}

func TestMustAccountType(t *testing.T) {
	// Valid type should not panic
	validType := "BANK"
	accountType := MustAccountType(validType)
	if accountType.Value() != "BANK" {
		t.Error("MustAccountType() returned wrong value for valid type")
	}

	// Invalid type should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAccountType() should panic for invalid type")
		}
	}()
	MustAccountType("INVALID")
}
