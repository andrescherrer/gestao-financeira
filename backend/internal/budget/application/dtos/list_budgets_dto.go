package dtos

import "gestao-financeira/backend/pkg/pagination"

// ListBudgetsInput represents the input data for listing budgets.
type ListBudgetsInput struct {
	UserID     string `json:"user_id" validate:"required,uuid"`
	CategoryID string `json:"category_id,omitempty" validate:"omitempty,uuid"`
	PeriodType string `json:"period_type,omitempty" validate:"omitempty,oneof=MONTHLY YEARLY"`
	Year       *int   `json:"year,omitempty" validate:"omitempty,min=1900,max=3000"`
	Month      *int   `json:"month,omitempty" validate:"omitempty,min=1,max=12"`
	Context    string `json:"context,omitempty" validate:"omitempty,oneof=PERSONAL BUSINESS"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Page       string `json:"page,omitempty"`  // Query parameter
	Limit      string `json:"limit,omitempty"` // Query parameter
}

// ListBudgetsOutput represents the output data for listing budgets.
type ListBudgetsOutput struct {
	Budgets    []BudgetOutput               `json:"budgets"`
	Total      int64                        `json:"total"`
	Pagination *pagination.PaginationResult `json:"pagination,omitempty"`
}

// BudgetOutput represents a budget in the list output.
type BudgetOutput struct {
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
