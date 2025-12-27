package dtos

// MonthlyReportOutput represents the output of a monthly report.
type MonthlyReportOutput struct {
	UserID   string `json:"user_id"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Currency string `json:"currency"`

	// Summary
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`

	// Transaction counts
	IncomeCount  int `json:"income_count"`
	ExpenseCount int `json:"expense_count"`
	TotalCount   int `json:"total_count"`

	// Category breakdown (optional, can be added later)
	CategoryBreakdown []MonthlyCategorySummary `json:"category_breakdown,omitempty"`
}

// MonthlyCategorySummary represents a summary of transactions by category for monthly reports.
type MonthlyCategorySummary struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name,omitempty"`
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
	Type         string  `json:"type"` // INCOME or EXPENSE
}
