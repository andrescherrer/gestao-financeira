package dtos

import (
	"gestao-financeira/backend/pkg/pagination"
)

// ListGoalsInput represents the input data for listing goals.
type ListGoalsInput struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	Context string `json:"context,omitempty" validate:"omitempty,oneof=PERSONAL BUSINESS"`
	Status  string `json:"status,omitempty" validate:"omitempty,oneof=IN_PROGRESS COMPLETED OVERDUE CANCELLED"`
	Page    string `json:"page,omitempty" validate:"omitempty,numeric"`
	Limit   string `json:"limit,omitempty" validate:"omitempty,numeric"`
}

// GoalListItem represents a goal in the list.
type GoalListItem struct {
	GoalID        string  `json:"goal_id"`
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

// ListGoalsOutput represents the output data for listing goals.
type ListGoalsOutput struct {
	Goals      []GoalListItem               `json:"goals"`
	Count      int                          `json:"count"`
	Pagination *pagination.PaginationResult `json:"pagination,omitempty"`
}
