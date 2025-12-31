package dtos

// GetGoalInput represents the input data for getting a goal.
type GetGoalInput struct {
	GoalID string `json:"goal_id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
}

// GetGoalOutput represents the output data for getting a goal.
type GetGoalOutput struct {
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
	UpdatedAt     string  `json:"updated_at"`
}
