package dtos

// DeleteBudgetInput represents the input data for deleting a budget.
type DeleteBudgetInput struct {
	BudgetID string `json:"budget_id" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}
