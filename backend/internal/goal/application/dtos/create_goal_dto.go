package dtos

// CreateGoalInput represents the input data for goal creation.
type CreateGoalInput struct {
	UserID       string  `json:"user_id" validate:"required,uuid"`
	Name         string  `json:"name" validate:"required,min=3,max=200,no_sql_injection,no_xss,utf8"`
	TargetAmount float64 `json:"target_amount" validate:"required,gt=0"`
	Currency     string  `json:"currency" validate:"required,oneof=BRL USD EUR"`
	Deadline     string  `json:"deadline" validate:"required,datetime=2006-01-02"`
	Context      string  `json:"context" validate:"required,oneof=PERSONAL BUSINESS"`
}

// CreateGoalOutput represents the output data after goal creation.
type CreateGoalOutput struct {
	GoalID        string  `json:"goal_id"`
	UserID        string  `json:"user_id"`
	Name          string  `json:"name"`
	TargetAmount  float64 `json:"target_amount"`
	CurrentAmount float64 `json:"current_amount"`
	Currency      string  `json:"currency"`
	Deadline      string  `json:"deadline"`
	Context       string  `json:"context"`
	Status        string  `json:"status"`
	Progress      float64 `json:"progress"`
	RemainingDays int     `json:"remaining_days"`
	CreatedAt     string  `json:"created_at"`
}
