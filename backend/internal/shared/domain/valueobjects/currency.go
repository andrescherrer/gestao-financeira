package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

// Currency represents a currency code (ISO 4217).
type Currency struct {
	code string
}

// Supported currency codes
const (
	BRL = "BRL" // Brazilian Real
	USD = "USD" // US Dollar
	EUR = "EUR" // Euro
)

// ValidCurrencies is a map of all supported currency codes.
var ValidCurrencies = map[string]string{
	BRL: "Brazilian Real",
	USD: "US Dollar",
	EUR: "Euro",
}

// NewCurrency creates a new Currency value object.
func NewCurrency(code string) (Currency, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	if !IsValidCurrency(code) {
		return Currency{}, fmt.Errorf("invalid currency code: %s. Supported codes: BRL, USD, EUR", code)
	}

	return Currency{code: code}, nil
}

// MustCurrency creates a new Currency value object and panics if the code is invalid.
// Use this only when you are certain the currency code is valid.
func MustCurrency(code string) Currency {
	currency, err := NewCurrency(code)
	if err != nil {
		panic(err)
	}
	return currency
}

// IsValidCurrency checks if a currency code is valid.
func IsValidCurrency(code string) bool {
	code = strings.ToUpper(strings.TrimSpace(code))
	_, exists := ValidCurrencies[code]
	return exists
}

// Code returns the currency code (e.g., "BRL", "USD", "EUR").
func (c Currency) Code() string {
	return c.code
}

// Name returns the full name of the currency.
func (c Currency) Name() string {
	if name, exists := ValidCurrencies[c.code]; exists {
		return name
	}
	return c.code
}

// Symbol returns the currency symbol.
func (c Currency) Symbol() string {
	switch c.code {
	case BRL:
		return "R$"
	case USD:
		return "$"
	case EUR:
		return "€"
	default:
		return c.code
	}
}

// Equals checks if two Currency values are equal.
func (c Currency) Equals(other Currency) bool {
	return c.code == other.code
}

// String returns the currency code as a string.
func (c Currency) String() string {
	return c.code
}

// BRL returns a Currency for Brazilian Real.
func BRLCurrency() Currency {
	return Currency{code: BRL}
}

// USDCurrency returns a Currency for US Dollar.
func USDCurrency() Currency {
	return Currency{code: USD}
}

// EURCurrency returns a Currency for Euro.
func EURCurrency() Currency {
	return Currency{code: EUR}
}

// ParseCurrency attempts to parse a currency code from a string.
// It accepts both uppercase and lowercase codes.
func ParseCurrency(s string) (Currency, error) {
	if s == "" {
		return Currency{}, errors.New("currency code cannot be empty")
	}

	return NewCurrency(s)
}
