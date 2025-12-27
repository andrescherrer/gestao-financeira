package dtos

// MonthlyReportInput represents the input for generating a monthly report.
type MonthlyReportInput struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	Year     int    `json:"year" validate:"required,min=2000,max=2100"`
	Month    int    `json:"month" validate:"required,min=1,max=12"`
	Currency string `json:"currency,omitempty" validate:"omitempty,oneof=BRL USD EUR"`
}
