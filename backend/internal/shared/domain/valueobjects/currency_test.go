package valueobjects

import (
	"testing"
)

func TestNewCurrency(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantCode  string
		wantError bool
	}{
		{"valid BRL", "BRL", "BRL", false},
		{"valid USD", "USD", "USD", false},
		{"valid EUR", "EUR", "EUR", false},
		{"lowercase", "brl", "BRL", false},
		{"mixed case", "UsD", "USD", false},
		{"invalid code", "INVALID", "", true},
		{"empty code", "", "", true},
		{"whitespace", "  BRL  ", "BRL", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currency, err := NewCurrency(tt.code)
			if (err != nil) != tt.wantError {
				t.Errorf("NewCurrency() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && currency.Code() != tt.wantCode {
				t.Errorf("NewCurrency() code = %v, want %v", currency.Code(), tt.wantCode)
			}
		})
	}
}

func TestCurrency_Code(t *testing.T) {
	brl, _ := NewCurrency("BRL")
	if brl.Code() != "BRL" {
		t.Errorf("Currency.Code() = %v, want BRL", brl.Code())
	}
}

func TestCurrency_Name(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		want     string
	}{
		{"BRL name", BRLCurrency(), "Brazilian Real"},
		{"USD name", USDCurrency(), "US Dollar"},
		{"EUR name", EURCurrency(), "Euro"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.currency.Name(); got != tt.want {
				t.Errorf("Currency.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrency_Symbol(t *testing.T) {
	tests := []struct {
		name     string
		currency Currency
		want     string
	}{
		{"BRL symbol", BRLCurrency(), "R$"},
		{"USD symbol", USDCurrency(), "$"},
		{"EUR symbol", EURCurrency(), "€"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.currency.Symbol(); got != tt.want {
				t.Errorf("Currency.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrency_Equals(t *testing.T) {
	brl1, _ := NewCurrency("BRL")
	brl2, _ := NewCurrency("BRL")
	usd, _ := NewCurrency("USD")

	tests := []struct {
		name string
		c1   Currency
		c2   Currency
		want bool
	}{
		{"equal currencies", brl1, brl2, true},
		{"different currencies", brl1, usd, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c1.Equals(tt.c2); got != tt.want {
				t.Errorf("Currency.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"valid BRL", "BRL", true},
		{"valid USD", "USD", true},
		{"valid EUR", "EUR", true},
		{"lowercase", "brl", true},
		{"invalid", "INVALID", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCurrency(tt.code); got != tt.want {
				t.Errorf("IsValidCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCurrency(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		wantCode  string
		wantError bool
	}{
		{"valid BRL", "BRL", "BRL", false},
		{"lowercase", "brl", "BRL", false},
		{"empty", "", "", true},
		{"invalid", "INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currency, err := ParseCurrency(tt.s)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseCurrency() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && currency.Code() != tt.wantCode {
				t.Errorf("ParseCurrency() code = %v, want %v", currency.Code(), tt.wantCode)
			}
		})
	}
}

func TestAllCurrencyHelpers(t *testing.T) {
	brl := BRLCurrency()
	usd := USDCurrency()
	eur := EURCurrency()

	if brl.Code() != "BRL" {
		t.Errorf("BRLCurrency() code = %v, want BRL", brl.Code())
	}
	if usd.Code() != "USD" {
		t.Errorf("USDCurrency() code = %v, want USD", usd.Code())
	}
	if eur.Code() != "EUR" {
		t.Errorf("EURCurrency() code = %v, want EUR", eur.Code())
	}
}
