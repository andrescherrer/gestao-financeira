package dtos

// CategoryReportOutput represents the output of a category report.
type CategoryReportOutput struct {
	UserID   string `json:"user_id"`
	Currency string `json:"currency"`

	// Summary by category
	CategoryBreakdown []CategorySummary `json:"category_breakdown"`

	// Totals
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`

	// Counts
	TotalCount int `json:"total_count"`
}

// CategorySummary represents a summary of transactions for a category.
type CategorySummary struct {
	CategoryID   string  `json:"category_id,omitempty"`   // Will be populated when category_id is added to transactions
	CategoryName string  `json:"category_name,omitempty"` // Will be populated when category_id is added to transactions
	Type         string  `json:"type"`                    // INCOME or EXPENSE
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"` // Percentage of total (income or expense)
}
