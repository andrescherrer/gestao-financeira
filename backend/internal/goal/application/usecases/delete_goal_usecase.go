package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// DeleteGoalUseCase handles deleting a goal.
type DeleteGoalUseCase struct {
	goalRepository repositories.GoalRepository
}

// NewDeleteGoalUseCase creates a new DeleteGoalUseCase instance.
func NewDeleteGoalUseCase(
	goalRepository repositories.GoalRepository,
) *DeleteGoalUseCase {
	return &DeleteGoalUseCase{
		goalRepository: goalRepository,
	}
}

// Execute performs deleting a goal.
func (uc *DeleteGoalUseCase) Execute(input dtos.DeleteGoalInput) (*dtos.DeleteGoalOutput, error) {
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

	// Find goal by ID to verify ownership
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

	// Delete goal (soft delete)
	if err := uc.goalRepository.Delete(goalID); err != nil {
		return nil, fmt.Errorf("failed to delete goal: %w", err)
	}

	// Build output
	output := &dtos.DeleteGoalOutput{
		Success: true,
	}

	return output, nil
}
