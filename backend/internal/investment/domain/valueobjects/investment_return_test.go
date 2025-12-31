package valueobjects

import (
	"testing"

	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

func TestNewInvestmentReturn(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	positiveMoney, _ := sharedvalueobjects.NewMoney(10000, currency) // R$ 100.00
	negativeMoney, _ := sharedvalueobjects.NewMoney(-5000, currency) // -R$ 50.00
	zeroMoney, _ := sharedvalueobjects.NewMoney(0, currency)

	tests := []struct {
		name       string
		absolute   sharedvalueobjects.Money
		percentage float64
	}{
		{"positive return", positiveMoney, 10.5},
		{"negative return", negativeMoney, -5.2},
		{"zero return", zeroMoney, 0.0},
		{"high percentage", positiveMoney, 150.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investmentReturn := NewInvestmentReturn(tt.absolute, tt.percentage)

			if !investmentReturn.Absolute().Equals(tt.absolute) {
				t.Errorf("InvestmentReturn.Absolute() = %v, want %v", investmentReturn.Absolute(), tt.absolute)
			}

			if investmentReturn.Percentage() != tt.percentage {
				t.Errorf("InvestmentReturn.Percentage() = %v, want %v", investmentReturn.Percentage(), tt.percentage)
			}
		})
	}
}

func TestInvestmentReturn_IsPositive(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	positiveMoney, _ := sharedvalueobjects.NewMoney(10000, currency)
	negativeMoney, _ := sharedvalueobjects.NewMoney(-5000, currency)
	zeroMoney, _ := sharedvalueobjects.NewMoney(0, currency)

	tests := []struct {
		name             string
		investmentReturn InvestmentReturn
		want             bool
	}{
		{"positive return", NewInvestmentReturn(positiveMoney, 10.5), true},
		{"negative return", NewInvestmentReturn(negativeMoney, -5.2), false},
		{"zero return", NewInvestmentReturn(zeroMoney, 0.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentReturn.IsPositive(); got != tt.want {
				t.Errorf("InvestmentReturn.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentReturn_IsNegative(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	positiveMoney, _ := sharedvalueobjects.NewMoney(10000, currency)
	negativeMoney, _ := sharedvalueobjects.NewMoney(-5000, currency)
	zeroMoney, _ := sharedvalueobjects.NewMoney(0, currency)

	tests := []struct {
		name             string
		investmentReturn InvestmentReturn
		want             bool
	}{
		{"positive return", NewInvestmentReturn(positiveMoney, 10.5), false},
		{"negative return", NewInvestmentReturn(negativeMoney, -5.2), true},
		{"zero return", NewInvestmentReturn(zeroMoney, 0.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentReturn.IsNegative(); got != tt.want {
				t.Errorf("InvestmentReturn.IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentReturn_IsZero(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	positiveMoney, _ := sharedvalueobjects.NewMoney(10000, currency)
	zeroMoney, _ := sharedvalueobjects.NewMoney(0, currency)

	tests := []struct {
		name             string
		investmentReturn InvestmentReturn
		want             bool
	}{
		{"positive return", NewInvestmentReturn(positiveMoney, 10.5), false},
		{"zero return", NewInvestmentReturn(zeroMoney, 0.0), true},
		{"negative return", NewInvestmentReturn(positiveMoney, -5.2), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.investmentReturn.IsZero(); got != tt.want {
				t.Errorf("InvestmentReturn.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvestmentReturn_Equals(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	money1, _ := sharedvalueobjects.NewMoney(10000, currency)
	money2, _ := sharedvalueobjects.NewMoney(20000, currency)

	return1 := NewInvestmentReturn(money1, 10.5)
	return2 := NewInvestmentReturn(money1, 10.5)
	return3 := NewInvestmentReturn(money2, 10.5)
	return4 := NewInvestmentReturn(money1, 20.0)

	if !return1.Equals(return2) {
		t.Error("InvestmentReturn.Equals() = false for equal returns, want true")
	}

	if return1.Equals(return3) {
		t.Error("InvestmentReturn.Equals() = true for different absolute values, want false")
	}

	if return1.Equals(return4) {
		t.Error("InvestmentReturn.Equals() = true for different percentages, want false")
	}
}

func TestInvestmentReturn_String(t *testing.T) {
	currency := sharedvalueobjects.MustCurrency("BRL")
	positiveMoney, _ := sharedvalueobjects.NewMoney(10000, currency) // R$ 100.00
	negativeMoney, _ := sharedvalueobjects.NewMoney(-5000, currency) // -R$ 50.00

	return1 := NewInvestmentReturn(positiveMoney, 10.5)
	return2 := NewInvestmentReturn(negativeMoney, -5.2)

	str1 := return1.String()
	if str1 == "" {
		t.Error("InvestmentReturn.String() returned empty string")
	}

	str2 := return2.String()
	if str2 == "" {
		t.Error("InvestmentReturn.String() returned empty string")
	}
}
