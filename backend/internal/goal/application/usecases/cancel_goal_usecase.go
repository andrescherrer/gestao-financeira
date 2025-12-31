package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// CancelGoalUseCase handles canceling a goal.
type CancelGoalUseCase struct {
	goalRepository repositories.GoalRepository
}

// NewCancelGoalUseCase creates a new CancelGoalUseCase instance.
func NewCancelGoalUseCase(
	goalRepository repositories.GoalRepository,
) *CancelGoalUseCase {
	return &CancelGoalUseCase{
		goalRepository: goalRepository,
	}
}

// Execute performs canceling a goal.
func (uc *CancelGoalUseCase) Execute(input dtos.CancelGoalInput) (*dtos.CancelGoalOutput, error) {
	// Create goal ID value object
	goalID, err := goalvalueobjects.NewGoalID(input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %w", err)
	}

	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find goal by ID
	goal, err := uc.goalRepository.FindByID(goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to find goal: %w", err)
	}

	// Check if goal exists
	if goal == nil {
		return nil, errors.New("goal not found")
	}

	// Verify goal belongs to user
	if !goal.UserID().Equals(userID) {
		return nil, errors.New("goal does not belong to user")
	}

	// Cancel goal
	if err := goal.Cancel(); err != nil {
		return nil, fmt.Errorf("failed to cancel goal: %w", err)
	}

	// Save goal to repository
	if err := uc.goalRepository.Save(goal); err != nil {
		return nil, fmt.Errorf("failed to save goal: %w", err)
	}

	// Build output
	output := &dtos.CancelGoalOutput{
		GoalID:    goal.ID().Value(),
		Status:    goal.Status().Value(),
		UpdatedAt: goal.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
