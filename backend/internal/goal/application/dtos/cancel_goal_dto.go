package dtos

// CancelGoalInput represents the input data for canceling a goal.
type CancelGoalInput struct {
	GoalID string `json:"goal_id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
}

// CancelGoalOutput represents the output data after canceling a goal.
type CancelGoalOutput struct {
	GoalID    string `json:"goal_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}
