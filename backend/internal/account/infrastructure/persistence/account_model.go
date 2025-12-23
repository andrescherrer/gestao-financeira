package persistence

import (
	"time"

	"gorm.io/gorm"
)

// AccountModel represents the database model for Account entity.
// This is the persistence model, separate from the domain entity.
type AccountModel struct {
	ID        string         `gorm:"type:uuid;primary_key"`
	UserID    string         `gorm:"type:uuid;index;not null"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Type      string         `gorm:"type:varchar(50);not null"`              // BANK, WALLET, INVESTMENT, CREDIT_CARD
	Balance   int64          `gorm:"type:bigint;not null;default:0"`         // Amount in cents
	Currency  string         `gorm:"type:varchar(3);not null;default:'BRL'"` // Currency code
	Context   string         `gorm:"type:varchar(20);not null"`              // PERSONAL, BUSINESS
	IsActive  bool           `gorm:"default:true;not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (AccountModel) TableName() string {
	return "accounts"
}
