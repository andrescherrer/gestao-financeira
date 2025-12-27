package dtos

// AnnualReportInput represents the input for generating an annual report.
type AnnualReportInput struct {
	UserID   string `json:"user_id" validate:"required,uuid"`
	Year     int    `json:"year" validate:"required,min=2000,max=2100"`
	Currency string `json:"currency,omitempty" validate:"omitempty,oneof=BRL USD EUR"`
}
