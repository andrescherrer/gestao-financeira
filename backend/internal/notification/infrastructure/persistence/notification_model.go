package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// NotificationModel represents the database model for Notification entity.
// This is the persistence model, separate from the domain entity.
type NotificationModel struct {
	ID        string         `gorm:"type:uuid;primary_key"`
	UserID    string         `gorm:"type:uuid;index;not null"`
	Title     string         `gorm:"type:varchar(200);not null"`
	Message   string         `gorm:"type:varchar(1000);not null"`
	Type      string         `gorm:"type:varchar(20);not null"`                  // INFO, WARNING, SUCCESS, ERROR
	Status    string         `gorm:"type:varchar(20);not null;default:'UNREAD'"` // UNREAD, READ, ARCHIVED
	ReadAt    *time.Time     `gorm:"type:timestamp"`
	Metadata  JSONB          `gorm:"type:jsonb"` // Additional data as JSON
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (NotificationModel) TableName() string {
	return "notifications"
}

// JSONB is a custom type for handling JSONB columns in PostgreSQL
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
