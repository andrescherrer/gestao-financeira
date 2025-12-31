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

// AddContributionUseCase handles adding a contribution to a goal.
type AddContributionUseCase struct {
	goalRepository repositories.GoalRepository
	eventBus       *eventbus.EventBus
}

// NewAddContributionUseCase creates a new AddContributionUseCase instance.
func NewAddContributionUseCase(
	goalRepository repositories.GoalRepository,
	eventBus *eventbus.EventBus,
) *AddContributionUseCase {
	return &AddContributionUseCase{
		goalRepository: goalRepository,
		eventBus:       eventBus,
	}
}

// Execute performs adding a contribution to a goal.
func (uc *AddContributionUseCase) Execute(input dtos.AddContributionInput) (*dtos.AddContributionOutput, error) {
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

	// Create contribution amount (convert float to cents)
	contributionCents := int64(input.Amount * 100)
	currency := goal.TargetAmount().Currency()
	contributionAmount, err := sharedvalueobjects.NewMoney(contributionCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid contribution amount: %w", err)
	}

	// Add contribution
	if err := goal.AddContribution(contributionAmount); err != nil {
		return nil, fmt.Errorf("failed to add contribution: %w", err)
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
	output := &dtos.AddContributionOutput{
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
