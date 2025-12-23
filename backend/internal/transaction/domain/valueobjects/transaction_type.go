package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

// TransactionType represents a transaction type value object.
type TransactionType struct {
	value string
}

// Valid transaction type values
const (
	Income  = "INCOME"  // Income transaction
	Expense = "EXPENSE" // Expense transaction
)

// ValidTransactionTypes is a map of all supported transaction types.
var ValidTransactionTypes = map[string]string{
	Income:  "Income",
	Expense: "Expense",
}

// NewTransactionType creates a new TransactionType value object.
func NewTransactionType(value string) (TransactionType, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	if !IsValidTransactionType(value) {
		return TransactionType{}, fmt.Errorf("invalid transaction type: %s. Supported values: INCOME, EXPENSE", value)
	}

	return TransactionType{value: value}, nil
}

// MustTransactionType creates a new TransactionType value object and panics if the value is invalid.
// Use this only when you are certain the transaction type value is valid.
func MustTransactionType(value string) TransactionType {
	tt, err := NewTransactionType(value)
	if err != nil {
		panic(err)
	}
	return tt
}

// IsValidTransactionType checks if a transaction type value is valid.
func IsValidTransactionType(value string) bool {
	value = strings.ToUpper(strings.TrimSpace(value))
	_, exists := ValidTransactionTypes[value]
	return exists
}

// Value returns the transaction type value (INCOME or EXPENSE).
func (tt TransactionType) Value() string {
	return tt.value
}

// String returns the transaction type value as a string.
func (tt TransactionType) String() string {
	return tt.value
}

// DisplayName returns the human-readable name of the transaction type.
func (tt TransactionType) DisplayName() string {
	if name, exists := ValidTransactionTypes[tt.value]; exists {
		return name
	}
	return tt.value
}

// IsIncome checks if the transaction type is Income.
func (tt TransactionType) IsIncome() bool {
	return tt.value == Income
}

// IsExpense checks if the transaction type is Expense.
func (tt TransactionType) IsExpense() bool {
	return tt.value == Expense
}

// Equals checks if two TransactionType values are equal.
func (tt TransactionType) Equals(other TransactionType) bool {
	return tt.value == other.value
}

// IncomeType returns a TransactionType for Income.
func IncomeType() TransactionType {
	return TransactionType{value: Income}
}

// ExpenseType returns a TransactionType for Expense.
func ExpenseType() TransactionType {
	return TransactionType{value: Expense}
}

// ParseTransactionType attempts to parse a transaction type from a string.
// It accepts both uppercase and lowercase values.
func ParseTransactionType(s string) (TransactionType, error) {
	if s == "" {
		return TransactionType{}, errors.New("transaction type cannot be empty")
	}

	return NewTransactionType(s)
}

// AllTransactionTypes returns all valid transaction types.
func AllTransactionTypes() []TransactionType {
	return []TransactionType{
		IncomeType(),
		ExpenseType(),
	}
}
