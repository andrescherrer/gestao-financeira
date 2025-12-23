package persistence

import (
	"time"

	"gorm.io/gorm"
)

// TransactionModel represents the database model for Transaction entity.
// This is the persistence model, separate from the domain entity.
type TransactionModel struct {
	ID          string         `gorm:"type:uuid;primary_key"`
	UserID      string         `gorm:"type:uuid;index;not null"`
	AccountID   string         `gorm:"type:uuid;index;not null"`
	Type        string         `gorm:"type:varchar(20);not null"`              // INCOME, EXPENSE
	Amount      int64          `gorm:"type:bigint;not null"`                   // Amount in cents
	Currency    string         `gorm:"type:varchar(3);not null;default:'BRL'"` // Currency code
	Description string         `gorm:"type:varchar(500);not null"`
	Date        time.Time      `gorm:"type:date;not null;index"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (TransactionModel) TableName() string {
	return "transactions"
}
