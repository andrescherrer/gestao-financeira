package dtos

// DeleteGoalInput represents the input data for deleting a goal.
type DeleteGoalInput struct {
	GoalID string `json:"goal_id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
}

// DeleteGoalOutput represents the output data after deleting a goal.
type DeleteGoalOutput struct {
	Success bool `json:"success"`
}
