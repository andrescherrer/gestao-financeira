package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// UpdateProgressUseCase handles updating goal progress.
type UpdateProgressUseCase struct {
	goalRepository repositories.GoalRepository
	eventBus       *eventbus.EventBus
}

// NewUpdateProgressUseCase creates a new UpdateProgressUseCase instance.
func NewUpdateProgressUseCase(
	goalRepository repositories.GoalRepository,
	eventBus *eventbus.EventBus,
) *UpdateProgressUseCase {
	return &UpdateProgressUseCase{
		goalRepository: goalRepository,
		eventBus:       eventBus,
	}
}

// Execute performs updating goal progress.
func (uc *UpdateProgressUseCase) Execute(input dtos.UpdateProgressInput) (*dtos.UpdateProgressOutput, error) {
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

	// Create progress amount (convert float to cents)
	progressCents := int64(input.Amount * 100)
	currency := goal.TargetAmount().Currency()
	progressAmount, err := sharedvalueobjects.NewMoney(progressCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid progress amount: %w", err)
	}

	// Update progress
	if err := goal.UpdateProgress(progressAmount); err != nil {
		return nil, fmt.Errorf("failed to update progress: %w", err)
	}

	// Save goal to repository
	if err := uc.goalRepository.Save(goal); err != nil {
		return nil, fmt.Errorf("failed to save goal: %w", err)
	}

	// Publish domain events
	domainEvents := goal.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			_ = err // Ignore for now, but should be logged
		}
	}
	goal.ClearEvents()

	// Build output
	currentAmount := goal.CurrentAmount()
	output := &dtos.UpdateProgressOutput{
		GoalID:        goal.ID().Value(),
		CurrentAmount: currentAmount.Float64(),
		TargetAmount:  goal.TargetAmount().Float64(),
		Currency:      currentAmount.Currency().Code(),
		Progress:      goal.CalculateProgress(),
		Status:        goal.Status().Value(),
		RemainingDays: goal.CalculateRemainingDays(),
		UpdatedAt:     goal.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
