package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/entities"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/pkg/pagination"
)

// ListGoalsUseCase handles listing goals for a user.
type ListGoalsUseCase struct {
	goalRepository repositories.GoalRepository
}

// NewListGoalsUseCase creates a new ListGoalsUseCase instance.
func NewListGoalsUseCase(
	goalRepository repositories.GoalRepository,
) *ListGoalsUseCase {
	return &ListGoalsUseCase{
		goalRepository: goalRepository,
	}
}

// Execute performs the goal listing.
// It validates the input, retrieves goals from the repository,
// and returns them as DTOs. Supports pagination.
func (uc *ListGoalsUseCase) Execute(input dtos.ListGoalsInput) (*dtos.ListGoalsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse pagination parameters
	paginationParams := pagination.ParsePaginationParams(input.Page, input.Limit)
	usePagination := input.Page != "" || input.Limit != ""

	var domainGoals []*entities.Goal
	var total int64

	// Check if we should use pagination
	if usePagination {
		// Use paginated query
		domainGoals, total, err = uc.goalRepository.FindByUserIDWithPagination(
			userID,
			input.Context,
			input.Status,
			paginationParams.CalculateOffset(),
			paginationParams.Limit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find goals: %w", err)
		}
	} else {
		// Use non-paginated query (backward compatibility)
		if input.Status != "" {
			// Create status value object
			status, err := goalvalueobjects.NewGoalStatus(input.Status)
			if err != nil {
				return nil, fmt.Errorf("invalid goal status: %w", err)
			}

			// Find goals by user ID and status
			domainGoals, err = uc.goalRepository.FindByStatus(userID, status)
			if err != nil {
				return nil, fmt.Errorf("failed to find goals: %w", err)
			}
			total = int64(len(domainGoals))
		} else {
			// Find all goals for the user
			domainGoals, err = uc.goalRepository.FindByUserID(userID)
			if err != nil {
				return nil, fmt.Errorf("failed to find goals: %w", err)
			}
			total = int64(len(domainGoals))
		}
	}

	goals := uc.toGoalOutputs(domainGoals)

	output := &dtos.ListGoalsOutput{
		Goals: goals,
		Count: len(goals),
	}

	// Add pagination metadata if pagination was used
	if usePagination {
		paginationResult := pagination.BuildPaginationResult(paginationParams, total)
		output.Pagination = &paginationResult
	}

	return output, nil
}

// toGoalOutputs converts domain goals to DTOs.
func (uc *ListGoalsUseCase) toGoalOutputs(domainGoals []*entities.Goal) []dtos.GoalListItem {
	outputs := make([]dtos.GoalListItem, 0, len(domainGoals))
	for _, goal := range domainGoals {
		currentAmount := goal.CurrentAmount()
		output := dtos.GoalListItem{
			GoalID:        goal.ID().Value(),
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
		}
		outputs = append(outputs, output)
	}
	return outputs
}
