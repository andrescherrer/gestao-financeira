package valueobjects

import (
	"fmt"
	"strings"
)

// AccountType represents the type of an account.
type AccountType struct {
	value string
}

// Valid account type values
const (
	Bank       = "BANK"
	Wallet     = "WALLET"
	Investment = "INVESTMENT"
	CreditCard = "CREDIT_CARD"
)

// ValidAccountTypes is a map of all supported account types.
var ValidAccountTypes = map[string]string{
	Bank:       "Bank",
	Wallet:     "Wallet",
	Investment: "Investment",
	CreditCard: "Credit Card",
}

// NewAccountType creates a new AccountType value object.
func NewAccountType(value string) (AccountType, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	if !IsValidAccountType(value) {
		return AccountType{}, fmt.Errorf("invalid account type: %s. Supported values: BANK, WALLET, INVESTMENT, CREDIT_CARD", value)
	}

	return AccountType{value: value}, nil
}

// MustAccountType creates a new AccountType and panics if invalid.
// Use this only when you are certain the type is valid (e.g., in tests).
func MustAccountType(value string) AccountType {
	accountType, err := NewAccountType(value)
	if err != nil {
		panic(err)
	}
	return accountType
}

// IsValidAccountType checks if an account type value is valid.
func IsValidAccountType(value string) bool {
	value = strings.ToUpper(strings.TrimSpace(value))
	_, exists := ValidAccountTypes[value]
	return exists
}

// Value returns the account type value.
func (at AccountType) Value() string {
	return at.value
}

// String returns the account type as a string (implements fmt.Stringer).
func (at AccountType) String() string {
	return at.value
}

// DisplayName returns the human-readable name of the account type.
func (at AccountType) DisplayName() string {
	if name, exists := ValidAccountTypes[at.value]; exists {
		return name
	}
	return at.value
}

// IsBank checks if the account type is Bank.
func (at AccountType) IsBank() bool {
	return at.value == Bank
}

// IsWallet checks if the account type is Wallet.
func (at AccountType) IsWallet() bool {
	return at.value == Wallet
}

// IsInvestment checks if the account type is Investment.
func (at AccountType) IsInvestment() bool {
	return at.value == Investment
}

// IsCreditCard checks if the account type is Credit Card.
func (at AccountType) IsCreditCard() bool {
	return at.value == CreditCard
}

// Equals checks if two AccountType values are equal.
func (at AccountType) Equals(other AccountType) bool {
	return at.value == other.value
}

// BankType returns an AccountType for Bank.
func BankType() AccountType {
	return AccountType{value: Bank}
}

// WalletType returns an AccountType for Wallet.
func WalletType() AccountType {
	return AccountType{value: Wallet}
}

// InvestmentType returns an AccountType for Investment.
func InvestmentType() AccountType {
	return AccountType{value: Investment}
}

// CreditCardType returns an AccountType for Credit Card.
func CreditCardType() AccountType {
	return AccountType{value: CreditCard}
}
