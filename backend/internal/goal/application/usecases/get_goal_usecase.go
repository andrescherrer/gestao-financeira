package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// GetGoalUseCase handles retrieving a single goal by ID.
type GetGoalUseCase struct {
	goalRepository repositories.GoalRepository
}

// NewGetGoalUseCase creates a new GetGoalUseCase instance.
func NewGetGoalUseCase(
	goalRepository repositories.GoalRepository,
) *GetGoalUseCase {
	return &GetGoalUseCase{
		goalRepository: goalRepository,
	}
}

// Execute performs the goal retrieval.
// It validates the input, retrieves the goal from the repository,
// and returns it as a DTO.
func (uc *GetGoalUseCase) Execute(input dtos.GetGoalInput) (*dtos.GetGoalOutput, error) {
	// Create goal ID value object
	goalID, err := goalvalueobjects.NewGoalID(input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %w", err)
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
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	if !goal.UserID().Equals(userID) {
		return nil, errors.New("goal does not belong to user")
	}

	currentAmount := goal.CurrentAmount()

	// Convert to output DTO
	output := &dtos.GetGoalOutput{
		GoalID:        goal.ID().Value(),
		UserID:        goal.UserID().Value(),
		Name:          goal.Name().Name(),
		TargetAmount:  goal.TargetAmount().Float64(),
		CurrentAmount: currentAmount.Float64(),
		Currency:      currentAmount.Currency().Code(),
		Deadline:      goal.Deadline().Format("2006-01-02"),
		Context:       goal.Context().Value(),
		Status:        goal.Status().Value(),
		Progress:      goal.CalculateProgress(),
		RemainingDays: goal.CalculateRemainingDays(),
		CreatedAt:     goal.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     goal.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
