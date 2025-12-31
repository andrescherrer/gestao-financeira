package dtos

// UpdateProgressInput represents the input data for updating goal progress.
type UpdateProgressInput struct {
	GoalID string  `json:"goal_id" validate:"required,uuid"`
	UserID string  `json:"user_id" validate:"required,uuid"`
	Amount float64 `json:"amount" validate:"required,gte=0"`
}

// UpdateProgressOutput represents the output data after updating progress.
type UpdateProgressOutput struct {
	GoalID        string  `json:"goal_id"`
	CurrentAmount float64 `json:"current_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Currency      string  `json:"currency"`
	Progress      float64 `json:"progress"`
	Status        string  `json:"status"`
	RemainingDays int     `json:"remaining_days"`
	UpdatedAt     string  `json:"updated_at"`
}
