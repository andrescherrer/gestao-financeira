package valueobjects

import (
	"fmt"
	"strings"
)

// InvestmentType represents the type of an investment.
type InvestmentType struct {
	value string
}

// Valid investment type values
const (
	Stock    = "STOCK"
	Fund     = "FUND"
	CDB      = "CDB"
	Treasury = "TREASURY"
	Crypto   = "CRYPTO"
	Other    = "OTHER"
)

// ValidInvestmentTypes is a map of all supported investment types.
var ValidInvestmentTypes = map[string]string{
	Stock:    "Stock",
	Fund:     "Fund",
	CDB:      "CDB",
	Treasury: "Treasury",
	Crypto:   "Crypto",
	Other:    "Other",
}

// NewInvestmentType creates a new InvestmentType value object.
func NewInvestmentType(value string) (InvestmentType, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	if !IsValidInvestmentType(value) {
		return InvestmentType{}, fmt.Errorf("invalid investment type: %s. Supported values: STOCK, FUND, CDB, TREASURY, CRYPTO, OTHER", value)
	}

	return InvestmentType{value: value}, nil
}

// MustInvestmentType creates a new InvestmentType and panics if invalid.
// Use this only when you are certain the type is valid (e.g., in tests).
func MustInvestmentType(value string) InvestmentType {
	investmentType, err := NewInvestmentType(value)
	if err != nil {
		panic(err)
	}
	return investmentType
}

// IsValidInvestmentType checks if an investment type value is valid.
func IsValidInvestmentType(value string) bool {
	value = strings.ToUpper(strings.TrimSpace(value))
	_, exists := ValidInvestmentTypes[value]
	return exists
}

// Value returns the investment type value.
func (it InvestmentType) Value() string {
	return it.value
}

// String returns the investment type as a string (implements fmt.Stringer).
func (it InvestmentType) String() string {
	return it.value
}

// DisplayName returns the human-readable name of the investment type.
func (it InvestmentType) DisplayName() string {
	if name, exists := ValidInvestmentTypes[it.value]; exists {
		return name
	}
	return it.value
}

// IsStock checks if the investment type is Stock.
func (it InvestmentType) IsStock() bool {
	return it.value == Stock
}

// IsFund checks if the investment type is Fund.
func (it InvestmentType) IsFund() bool {
	return it.value == Fund
}

// IsCDB checks if the investment type is CDB.
func (it InvestmentType) IsCDB() bool {
	return it.value == CDB
}

// IsTreasury checks if the investment type is Treasury.
func (it InvestmentType) IsTreasury() bool {
	return it.value == Treasury
}

// IsCrypto checks if the investment type is Crypto.
func (it InvestmentType) IsCrypto() bool {
	return it.value == Crypto
}

// IsOther checks if the investment type is Other.
func (it InvestmentType) IsOther() bool {
	return it.value == Other
}

// Equals checks if two InvestmentType values are equal.
func (it InvestmentType) Equals(other InvestmentType) bool {
	return it.value == other.value
}

// RequiresQuantity checks if this investment type requires quantity tracking.
// Stocks, Funds, and Crypto typically require quantity.
func (it InvestmentType) RequiresQuantity() bool {
	return it.value == Stock || it.value == Fund || it.value == Crypto
}

// HasVariableValue checks if this investment type has variable value.
// All investment types have variable value except fixed income like CDB and Treasury.
func (it InvestmentType) HasVariableValue() bool {
	return it.value == Stock || it.value == Fund || it.value == Crypto || it.value == Other
}

// StockType returns an InvestmentType for Stock.
func StockType() InvestmentType {
	return InvestmentType{value: Stock}
}

// FundType returns an InvestmentType for Fund.
func FundType() InvestmentType {
	return InvestmentType{value: Fund}
}

// CDBType returns an InvestmentType for CDB.
func CDBType() InvestmentType {
	return InvestmentType{value: CDB}
}

// TreasuryType returns an InvestmentType for Treasury.
func TreasuryType() InvestmentType {
	return InvestmentType{value: Treasury}
}

// CryptoType returns an InvestmentType for Crypto.
func CryptoType() InvestmentType {
	return InvestmentType{value: Crypto}
}

// OtherType returns an InvestmentType for Other.
func OtherType() InvestmentType {
	return InvestmentType{value: Other}
}
