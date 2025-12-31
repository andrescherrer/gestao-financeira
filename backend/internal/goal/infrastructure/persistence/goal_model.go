package persistence

import (
	"time"

	"gorm.io/gorm"
)

// GoalModel represents the database model for Goal entity.
// This is the persistence model, separate from the domain entity.
type GoalModel struct {
	ID              string         `gorm:"type:uuid;primary_key"`
	UserID          string         `gorm:"type:uuid;index;not null"`
	Name            string         `gorm:"type:varchar(200);not null"`
	TargetAmount    int64          `gorm:"type:bigint;not null"`                   // Amount in cents
	TargetCurrency  string         `gorm:"type:varchar(3);not null;default:'BRL'"` // Currency code
	CurrentAmount   int64          `gorm:"type:bigint;not null;default:0"`         // Amount in cents
	CurrentCurrency string         `gorm:"type:varchar(3);not null;default:'BRL'"` // Currency code
	Deadline        time.Time      `gorm:"type:date;not null"`
	Context         string         `gorm:"type:varchar(20);not null"`                       // PERSONAL, BUSINESS
	Status          string         `gorm:"type:varchar(20);not null;default:'IN_PROGRESS'"` // IN_PROGRESS, COMPLETED, OVERDUE, CANCELLED
	CreatedAt       time.Time      `gorm:"not null"`
	UpdatedAt       time.Time      `gorm:"not null"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (GoalModel) TableName() string {
	return "goals"
}
