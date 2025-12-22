package persistence

import (
	"time"

	"gorm.io/gorm"
)

// UserModel represents the database model for User entity.
// This is the persistence model, separate from the domain entity.
type UserModel struct {
	ID           string         `gorm:"type:uuid;primary_key"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	FirstName    string         `gorm:"type:varchar(100);not null"`
	LastName     string         `gorm:"type:varchar(100);not null"`
	IsActive     bool           `gorm:"default:true;not null"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}
