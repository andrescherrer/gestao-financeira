package dtos

// GetBudgetInput represents the input data for retrieving a budget.
type GetBudgetInput struct {
	BudgetID string `json:"budget_id" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}

// GetBudgetOutput represents the output data for a budget.
type GetBudgetOutput struct {
	BudgetID   string  `json:"budget_id"`
	UserID     string  `json:"user_id"`
	CategoryID string  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	PeriodType string  `json:"period_type"`
	Year       int     `json:"year"`
	Month      *int    `json:"month,omitempty"`
	Context    string  `json:"context"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
