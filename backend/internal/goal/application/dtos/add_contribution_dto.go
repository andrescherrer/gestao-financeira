package dtos

// AddContributionInput represents the input data for adding a contribution to a goal.
type AddContributionInput struct {
	GoalID string  `json:"goal_id" validate:"required,uuid"`
	UserID string  `json:"user_id" validate:"required,uuid"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// AddContributionOutput represents the output data after adding a contribution.
type AddContributionOutput struct {
	GoalID        string  `json:"goal_id"`
	CurrentAmount float64 `json:"current_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Currency      string  `json:"currency"`
	Progress      float64 `json:"progress"`
	Status        string  `json:"status"`
	RemainingDays int     `json:"remaining_days"`
	UpdatedAt     string  `json:"updated_at"`
}
