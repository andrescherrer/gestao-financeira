package valueobjects

import (
	"errors"
	"strings"
)

const (
	// MinInvestmentNameLength is the minimum length for an investment name.
	MinInvestmentNameLength = 2
	// MaxInvestmentNameLength is the maximum length for an investment name.
	MaxInvestmentNameLength = 200
	// MaxTickerLength is the maximum length for a ticker symbol.
	MaxTickerLength = 20
)

// InvestmentName represents an investment name value object with optional ticker.
type InvestmentName struct {
	name   string
	ticker *string // Optional ticker symbol (e.g., "PETR4", "AAPL")
}

// NewInvestmentName creates a new InvestmentName value object.
func NewInvestmentName(name string, ticker *string) (InvestmentName, error) {
	trimmedName := strings.TrimSpace(name)

	if trimmedName == "" {
		return InvestmentName{}, errors.New("investment name cannot be empty")
	}

	if len(trimmedName) < MinInvestmentNameLength {
		return InvestmentName{}, errors.New("investment name must have at least 2 characters")
	}

	if len(trimmedName) > MaxInvestmentNameLength {
		return InvestmentName{}, errors.New("investment name must have at most 200 characters")
	}

	var trimmedTicker *string
	if ticker != nil {
		trimmed := strings.TrimSpace(*ticker)
		if trimmed != "" {
			if len(trimmed) > MaxTickerLength {
				return InvestmentName{}, errors.New("ticker must have at most 20 characters")
			}
			trimmedTicker = &trimmed
		}
	}

	return InvestmentName{
		name:   trimmedName,
		ticker: trimmedTicker,
	}, nil
}

// MustInvestmentName creates a new InvestmentName and panics if invalid.
// Use this only when you are certain the name is valid (e.g., in tests).
func MustInvestmentName(name string, ticker *string) InvestmentName {
	investmentName, err := NewInvestmentName(name, ticker)
	if err != nil {
		panic(err)
	}
	return investmentName
}

// Name returns the investment name as a string.
func (in InvestmentName) Name() string {
	return in.name
}

// Ticker returns the ticker symbol if available, otherwise returns nil.
func (in InvestmentName) Ticker() *string {
	return in.ticker
}

// HasTicker checks if the investment has a ticker symbol.
func (in InvestmentName) HasTicker() bool {
	return in.ticker != nil && *in.ticker != ""
}

// String returns the investment name as a string (implements fmt.Stringer).
// If ticker is available, returns "Name (TICKER)", otherwise returns "Name".
func (in InvestmentName) String() string {
	if in.HasTicker() {
		return in.name + " (" + *in.ticker + ")"
	}
	return in.name
}

// Equals checks if two InvestmentName values are equal.
func (in InvestmentName) Equals(other InvestmentName) bool {
	if in.name != other.name {
		return false
	}

	// Compare tickers
	if in.ticker == nil && other.ticker == nil {
		return true
	}
	if in.ticker == nil || other.ticker == nil {
		return false
	}
	return *in.ticker == *other.ticker
}

// IsEmpty checks if the investment name is empty.
func (in InvestmentName) IsEmpty() bool {
	return in.name == ""
}

