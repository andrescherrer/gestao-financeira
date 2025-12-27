package dtos

// AnnualReportOutput represents the output of an annual report.
type AnnualReportOutput struct {
	UserID   string `json:"user_id"`
	Year     int    `json:"year"`
	Currency string `json:"currency"`

	// Summary
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`

	// Transaction counts
	IncomeCount  int `json:"income_count"`
	ExpenseCount int `json:"expense_count"`
	TotalCount   int `json:"total_count"`

	// Monthly breakdown
	MonthlyBreakdown []MonthlySummary `json:"monthly_breakdown"`
}

// MonthlySummary represents a summary for a specific month.
type MonthlySummary struct {
	Month        int     `json:"month"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	IncomeCount  int     `json:"income_count"`
	ExpenseCount int     `json:"expense_count"`
}
