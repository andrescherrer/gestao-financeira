package dtos

// UpdateBudgetInput represents the input data for updating a budget.
type UpdateBudgetInput struct {
	BudgetID   string   `json:"budget_id" validate:"required,uuid"`
	UserID     string   `json:"user_id" validate:"required,uuid"`
	Amount     *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	PeriodType *string  `json:"period_type,omitempty" validate:"omitempty,oneof=MONTHLY YEARLY"`
	Year       *int     `json:"year,omitempty" validate:"omitempty,min=1900,max=3000"`
	Month      *int     `json:"month,omitempty" validate:"omitempty,min=1,max=12"`
	IsActive   *bool    `json:"is_active,omitempty"`
}

// UpdateBudgetOutput represents the output data after budget update.
type UpdateBudgetOutput struct {
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
	UpdatedAt  string  `json:"updated_at"`
}
