package persistence

import (
	"time"

	"gorm.io/gorm"
)

// InvestmentModel represents the database model for Investment entity.
// This is the persistence model, separate from the domain entity.
type InvestmentModel struct {
	ID             string         `gorm:"type:uuid;primary_key"`
	UserID         string         `gorm:"type:uuid;index;not null"`
	AccountID      string         `gorm:"type:uuid;index;not null"`
	Type           string         `gorm:"type:varchar(50);not null"`              // STOCK, FUND, CDB, TREASURY, CRYPTO, OTHER
	Name           string         `gorm:"type:varchar(200);not null"`
	Ticker         *string        `gorm:"type:varchar(20)"`                         // Optional ticker symbol
	PurchaseDate   time.Time      `gorm:"type:date;not null"`
	PurchaseAmount int64          `gorm:"type:bigint;not null"`                    // Amount in cents
	PurchaseCurrency string       `gorm:"type:varchar(3);not null;default:'BRL'"`   // Currency code
	CurrentValue   int64          `gorm:"type:bigint;not null"`                    // Amount in cents
	CurrentCurrency string        `gorm:"type:varchar(3);not null;default:'BRL'"` // Currency code
	Quantity       *float64       `gorm:"type:decimal(15,4)"`                      // Optional quantity
	Context        string         `gorm:"type:varchar(20);not null"`              // PERSONAL, BUSINESS
	CreatedAt      time.Time      `gorm:"not null"`
	UpdatedAt      time.Time      `gorm:"not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (InvestmentModel) TableName() string {
	return "investments"
}

