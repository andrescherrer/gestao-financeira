package dtos

import "time"

// IncomeVsExpenseInput represents the input for generating an income vs expense report.
type IncomeVsExpenseInput struct {
	UserID    string     `json:"user_id" validate:"required,uuid"`
	StartDate *time.Time `json:"start_date,omitempty"` // Optional: start date filter
	EndDate   *time.Time `json:"end_date,omitempty"`   // Optional: end date filter
	Currency  string     `json:"currency,omitempty" validate:"omitempty,oneof=BRL USD EUR"`
	GroupBy   string     `json:"group_by,omitempty" validate:"omitempty,oneof=day week month year"` // Optional: grouping period
}
