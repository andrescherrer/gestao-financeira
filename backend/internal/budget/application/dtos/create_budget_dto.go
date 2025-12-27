package dtos

// CreateBudgetInput represents the input data for budget creation.
type CreateBudgetInput struct {
	UserID     string  `json:"user_id" validate:"required,uuid"`
	CategoryID string  `json:"category_id" validate:"required,uuid"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Currency   string  `json:"currency" validate:"required,oneof=BRL USD EUR"`
	PeriodType string  `json:"period_type" validate:"required,oneof=MONTHLY YEARLY"`
	Year       int     `json:"year" validate:"required,min=1900,max=3000"`
	Month      *int    `json:"month" validate:"omitempty,min=1,max=12"`
	Context    string  `json:"context" validate:"required,oneof=PERSONAL BUSINESS"`
}

// CreateBudgetOutput represents the output data after budget creation.
type CreateBudgetOutput struct {
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
}
