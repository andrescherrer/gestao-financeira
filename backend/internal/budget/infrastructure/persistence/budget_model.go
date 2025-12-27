package persistence

import (
	"time"

	"gorm.io/gorm"
)

// BudgetModel represents the database model for Budget entity.
// This is the persistence model, separate from the domain entity.
type BudgetModel struct {
	ID         string         `gorm:"type:uuid;primary_key"`
	UserID     string         `gorm:"type:uuid;index;not null"`
	CategoryID string         `gorm:"type:uuid;index;not null"`
	Amount     int64          `gorm:"not null"` // Amount in cents
	Currency   string         `gorm:"type:varchar(3);not null"`
	PeriodType string         `gorm:"type:varchar(10);not null"` // MONTHLY or YEARLY
	Year       int            `gorm:"not null"`
	Month      *int           `gorm:"type:integer"`              // NULL for yearly periods
	Context    string         `gorm:"type:varchar(20);not null"` // PERSONAL or BUSINESS
	IsActive   bool           `gorm:"default:true;not null"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (BudgetModel) TableName() string {
	return "budgets"
}
