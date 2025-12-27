package dtos

// GetBudgetProgressInput represents the input data for getting budget progress.
type GetBudgetProgressInput struct {
	BudgetID string `json:"budget_id" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}

// GetBudgetProgressOutput represents the output data for budget progress.
type GetBudgetProgressOutput struct {
	BudgetID       string  `json:"budget_id"`
	CategoryID     string  `json:"category_id"`
	Budgeted       float64 `json:"budgeted"`
	Spent          float64 `json:"spent"`
	Remaining      float64 `json:"remaining"`
	PercentageUsed float64 `json:"percentage_used"`
	Currency       string  `json:"currency"`
	IsExceeded     bool    `json:"is_exceeded"`
	PeriodType     string  `json:"period_type"`
	Year           int     `json:"year"`
	Month          *int    `json:"month,omitempty"`
}
