package valueobjects

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	brl := BRLCurrency()

	tests := []struct {
		name      string
		cents     int64
		currency  Currency
		wantError bool
	}{
		{"valid money", 10000, brl, false},
		{"zero money", 0, brl, false},
		{"negative money", -1000, brl, false}, // Negative is allowed for debts
		{"large amount", 999999999999, brl, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.cents, tt.currency)
			if (err != nil) != tt.wantError {
				t.Errorf("NewMoney() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && money.Amount() != tt.cents {
				t.Errorf("NewMoney() amount = %d, want %d", money.Amount(), tt.cents)
			}
		})
	}
}

func TestNewMoneyFromFloat(t *testing.T) {
	brl := BRLCurrency()

	tests := []struct {
		name      string
		amount    float64
		currency  Currency
		wantCents int64
		wantError bool
	}{
		{"valid float", 100.50, brl, 10050, false},
		{"zero", 0.0, brl, 0, false},
		{"small amount", 0.01, brl, 1, false},
		{"large amount", 999999.99, brl, 99999999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoneyFromFloat(tt.amount, tt.currency)
			if (err != nil) != tt.wantError {
				t.Errorf("NewMoneyFromFloat() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && money.Amount() != tt.wantCents {
				t.Errorf("NewMoneyFromFloat() amount = %d, want %d", money.Amount(), tt.wantCents)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	brl := BRLCurrency()
	usd := USDCurrency()

	money1, _ := NewMoney(10000, brl) // R$ 100.00
	money2, _ := NewMoney(5000, brl)  // R$ 50.00
	money3, _ := NewMoney(10000, usd) // $100.00

	tests := []struct {
		name      string
		m1        Money
		m2        Money
		wantCents int64
		wantError bool
	}{
		{"same currency", money1, money2, 15000, false},
		{"different currencies", money1, money3, 0, true},
		{"zero addition", money1, Zero(brl), 10000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m1.Add(tt.m2)
			if (err != nil) != tt.wantError {
				t.Errorf("Money.Add() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && result.Amount() != tt.wantCents {
				t.Errorf("Money.Add() amount = %d, want %d", result.Amount(), tt.wantCents)
			}
		})
	}
}

func TestMoney_Subtract(t *testing.T) {
	brl := BRLCurrency()
	usd := USDCurrency()

	money1, _ := NewMoney(10000, brl) // R$ 100.00
	money2, _ := NewMoney(3000, brl)  // R$ 30.00
	money3, _ := NewMoney(10000, usd) // $100.00

	tests := []struct {
		name      string
		m1        Money
		m2        Money
		wantCents int64
		wantError bool
	}{
		{"same currency", money1, money2, 7000, false},
		{"different currencies", money1, money3, 0, true},
		{"zero subtraction", money1, Zero(brl), 10000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.m1.Subtract(tt.m2)
			if (err != nil) != tt.wantError {
				t.Errorf("Money.Subtract() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && result.Amount() != tt.wantCents {
				t.Errorf("Money.Subtract() amount = %d, want %d", result.Amount(), tt.wantCents)
			}
		})
	}
}

func TestMoney_Multiply(t *testing.T) {
	brl := BRLCurrency()
	money, _ := NewMoney(10000, brl) // R$ 100.00

	tests := []struct {
		name      string
		money     Money
		factor    float64
		wantCents int64
	}{
		{"multiply by 2", money, 2.0, 20000},
		{"multiply by 0.5", money, 0.5, 5000},
		{"multiply by 0", money, 0.0, 0},
		{"multiply by 1.5", money, 1.5, 15000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.money.Multiply(tt.factor)
			if result.Amount() != tt.wantCents {
				t.Errorf("Money.Multiply() amount = %d, want %d", result.Amount(), tt.wantCents)
			}
		})
	}
}

func TestMoney_Divide(t *testing.T) {
	brl := BRLCurrency()
	money, _ := NewMoney(10000, brl) // R$ 100.00

	tests := []struct {
		name      string
		money     Money
		divisor   float64
		wantError bool
	}{
		{"divide by 2", money, 2.0, false},
		{"divide by 0", money, 0.0, true},
		{"divide by 0.5", money, 0.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.money.Divide(tt.divisor)
			if (err != nil) != tt.wantError {
				t.Errorf("Money.Divide() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	brl := BRLCurrency()

	tests := []struct {
		name  string
		money Money
		want  bool
	}{
		{"zero money", Zero(brl), true},
		{"non-zero money", MustMoneyFromFloat(100.0, brl), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.money.IsZero(); got != tt.want {
				t.Errorf("Money.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_IsPositive(t *testing.T) {
	brl := BRLCurrency()

	tests := []struct {
		name  string
		money Money
		want  bool
	}{
		{"positive money", MustMoneyFromFloat(100.0, brl), true},
		{"zero money", Zero(brl), false},
		{"negative money", MustMoneyFromFloat(-100.0, brl), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.money.IsPositive(); got != tt.want {
				t.Errorf("Money.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_IsNegative(t *testing.T) {
	brl := BRLCurrency()

	tests := []struct {
		name  string
		money Money
		want  bool
	}{
		{"negative money", MustMoneyFromFloat(-100.0, brl), true},
		{"zero money", Zero(brl), false},
		{"positive money", MustMoneyFromFloat(100.0, brl), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.money.IsNegative(); got != tt.want {
				t.Errorf("Money.IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Equals(t *testing.T) {
	brl := BRLCurrency()
	usd := USDCurrency()

	money1, _ := NewMoney(10000, brl)
	money2, _ := NewMoney(10000, brl)
	money3, _ := NewMoney(5000, brl)
	money4, _ := NewMoney(10000, usd)

	tests := []struct {
		name string
		m1   Money
		m2   Money
		want bool
	}{
		{"equal amounts same currency", money1, money2, true},
		{"different amounts same currency", money1, money3, false},
		{"same amounts different currencies", money1, money4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Equals(tt.m2); got != tt.want {
				t.Errorf("Money.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	brl := BRLCurrency()
	usd := USDCurrency()

	money1, _ := NewMoney(10000, brl)
	money2, _ := NewMoney(5000, brl)
	money3, _ := NewMoney(10000, usd)

	tests := []struct {
		name      string
		m1        Money
		m2        Money
		want      bool
		wantError bool
	}{
		{"greater than", money1, money2, true, false},
		{"not greater than", money2, money1, false, false},
		{"different currencies", money1, money3, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m1.GreaterThan(tt.m2)
			if (err != nil) != tt.wantError {
				t.Errorf("Money.GreaterThan() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("Money.GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function for tests
func MustMoneyFromFloat(amount float64, currency Currency) Money {
	money, _ := NewMoneyFromFloat(amount, currency)
	return money
}
