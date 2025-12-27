package dtos

import "time"

// CategoryReportInput represents the input for generating a category report.
type CategoryReportInput struct {
	UserID     string     `json:"user_id" validate:"required,uuid"`
	CategoryID string     `json:"category_id,omitempty" validate:"omitempty,uuid"` // Optional: filter by specific category
	StartDate  *time.Time `json:"start_date,omitempty"`                            // Optional: start date filter
	EndDate    *time.Time `json:"end_date,omitempty"`                              // Optional: end date filter
	Currency   string     `json:"currency,omitempty" validate:"omitempty,oneof=BRL USD EUR"`
}
