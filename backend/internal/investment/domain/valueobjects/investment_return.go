package valueobjects

import (
	"fmt"

	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// InvestmentReturn represents the return of an investment (absolute and percentage).
type InvestmentReturn struct {
	absolute   sharedvalueobjects.Money // Absolute return (gain/loss)
	percentage float64                   // Percentage return
}

// NewInvestmentReturn creates a new InvestmentReturn value object.
func NewInvestmentReturn(absolute sharedvalueobjects.Money, percentage float64) InvestmentReturn {
	return InvestmentReturn{
		absolute:   absolute,
		percentage: percentage,
	}
}

// Absolute returns the absolute return as Money.
func (ir InvestmentReturn) Absolute() sharedvalueobjects.Money {
	return ir.absolute
}

// Percentage returns the percentage return as a float64.
func (ir InvestmentReturn) Percentage() float64 {
	return ir.percentage
}

// IsPositive checks if the return is positive (gain).
func (ir InvestmentReturn) IsPositive() bool {
	return ir.percentage > 0
}

// IsNegative checks if the return is negative (loss).
func (ir InvestmentReturn) IsNegative() bool {
	return ir.percentage < 0
}

// IsZero checks if the return is zero (no gain, no loss).
func (ir InvestmentReturn) IsZero() bool {
	return ir.percentage == 0
}

// Equals checks if two InvestmentReturn values are equal.
func (ir InvestmentReturn) Equals(other InvestmentReturn) bool {
	return ir.absolute.Equals(other.absolute) && ir.percentage == other.percentage
}

// String returns a string representation of the investment return.
func (ir InvestmentReturn) String() string {
	sign := "+"
	if ir.percentage < 0 {
		sign = ""
	}
	return fmt.Sprintf("%s%s (%.2f%%)", sign, ir.absolute.String(), ir.percentage)
}

