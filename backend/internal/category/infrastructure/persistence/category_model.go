package persistence

import (
	"time"

	"gorm.io/gorm"
)

// CategoryModel represents the database model for Category entity.
// This is the persistence model, separate from the domain entity.
type CategoryModel struct {
	ID          string         `gorm:"type:uuid;primary_key"`
	UserID      string         `gorm:"type:uuid;index;not null"`
	Name        string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	IsActive    bool           `gorm:"default:true;not null"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (CategoryModel) TableName() string {
	return "categories"
}
