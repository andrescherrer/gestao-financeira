package valueobjects

import (
	"errors"
	"fmt"
)

// Money represents a monetary value with amount and currency.
// The amount is stored in cents (int64) to avoid floating-point precision issues.
type Money struct {
	amount   int64  // amount in cents
	currency string // currency code (e.g., "BRL", "USD", "EUR")
}

// NewMoney creates a new Money value object from cents.
func NewMoney(cents int64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	return Money{
		amount:   cents,
		currency: currency,
	}, nil
}

// NewMoneyFromFloat creates a new Money value object from a float64 amount.
// The amount is converted to cents (multiplied by 100).
func NewMoneyFromFloat(amount float64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}

	cents := int64(amount * 100)
	return Money{
		amount:   cents,
		currency: currency,
	}, nil
}

// Zero creates a Money with zero amount for the given currency.
func Zero(currency string) (Money, error) {
	return NewMoney(0, currency)
}

// Amount returns the amount in cents.
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code.
func (m Money) Currency() string {
	return m.currency
}

// Float64 returns the amount as a float64 (amount / 100).
func (m Money) Float64() float64 {
	return float64(m.amount) / 100.0
}

// Add adds two Money values. Both must have the same currency.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add money with different currencies: %s and %s", m.currency, other.currency)
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts other Money from this Money. Both must have the same currency.
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract money with different currencies: %s and %s", m.currency, other.currency)
	}

	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply multiplies the Money by a scalar factor.
func (m Money) Multiply(factor float64) Money {
	newAmount := int64(float64(m.amount) * factor)
	return Money{
		amount:   newAmount,
		currency: m.currency,
	}
}

// Divide divides the Money by a scalar divisor.
func (m Money) Divide(divisor float64) (Money, error) {
	if divisor == 0 {
		return Money{}, errors.New("cannot divide by zero")
	}

	newAmount := int64(float64(m.amount) / divisor)
	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// Negate returns the negative of this Money.
func (m Money) Negate() Money {
	return Money{
		amount:   -m.amount,
		currency: m.currency,
	}
}

// IsZero checks if the Money amount is zero.
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive checks if the Money amount is positive.
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative checks if the Money amount is negative.
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// Equals checks if two Money values are equal (same amount and currency).
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// GreaterThan checks if this Money is greater than other Money.
// Both must have the same currency.
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare money with different currencies: %s and %s", m.currency, other.currency)
	}

	return m.amount > other.amount, nil
}

// LessThan checks if this Money is less than other Money.
// Both must have the same currency.
func (m Money) LessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare money with different currencies: %s and %s", m.currency, other.currency)
	}

	return m.amount < other.amount, nil
}

// GreaterThanOrEqual checks if this Money is greater than or equal to other Money.
// Both must have the same currency.
func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare money with different currencies: %s and %s", m.currency, other.currency)
	}

	return m.amount >= other.amount, nil
}

// LessThanOrEqual checks if this Money is less than or equal to other Money.
// Both must have the same currency.
func (m Money) LessThanOrEqual(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare money with different currencies: %s and %s", m.currency, other.currency)
	}

	return m.amount <= other.amount, nil
}

// String returns a string representation of Money.
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Float64(), m.currency)
}

// Format returns a formatted string representation of Money.
// Example: "R$ 1.234,56" for BRL or "$1,234.56" for USD.
func (m Money) Format() string {
	amount := m.Float64()
	switch m.currency {
	case "BRL":
		return fmt.Sprintf("R$ %.2f", amount)
	case "USD":
		return fmt.Sprintf("$%.2f", amount)
	case "EUR":
		return fmt.Sprintf("€%.2f", amount)
	default:
		return fmt.Sprintf("%.2f %s", amount, m.currency)
	}
}
