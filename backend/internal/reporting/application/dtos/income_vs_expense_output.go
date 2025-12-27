package dtos

// IncomeVsExpenseOutput represents the output of an income vs expense report.
type IncomeVsExpenseOutput struct {
	UserID   string `json:"user_id"`
	Currency string `json:"currency"`

	// Overall summary
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	Difference   float64 `json:"difference"` // Income - Expense

	// Counts
	IncomeCount  int `json:"income_count"`
	ExpenseCount int `json:"expense_count"`
	TotalCount   int `json:"total_count"`

	// Period breakdown (if group_by is specified)
	PeriodBreakdown []PeriodSummary `json:"period_breakdown,omitempty"`
}

// PeriodSummary represents a summary for a specific period.
type PeriodSummary struct {
	Period       string  `json:"period"` // Date or period identifier
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	IncomeCount  int     `json:"income_count"`
	ExpenseCount int     `json:"expense_count"`
}
