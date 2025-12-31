package usecases

import (
	"fmt"
	"time"

	"gestao-financeira/backend/internal/goal/application/dtos"
	"gestao-financeira/backend/internal/goal/domain/entities"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateGoalUseCase handles goal creation.
type CreateGoalUseCase struct {
	goalRepository repositories.GoalRepository
	eventBus       *eventbus.EventBus
}

// NewCreateGoalUseCase creates a new CreateGoalUseCase instance.
func NewCreateGoalUseCase(
	goalRepository repositories.GoalRepository,
	eventBus *eventbus.EventBus,
) *CreateGoalUseCase {
	return &CreateGoalUseCase{
		goalRepository: goalRepository,
		eventBus:       eventBus,
	}
}

// Execute performs the goal creation.
// It validates the input, creates value objects, creates a new goal entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateGoalUseCase) Execute(input dtos.CreateGoalInput) (*dtos.CreateGoalOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create goal name value object
	name, err := goalvalueobjects.NewGoalName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid goal name: %w", err)
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(input.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	// Create target amount (convert float to cents)
	targetAmountCents := int64(input.TargetAmount * 100)
	targetAmount, err := sharedvalueobjects.NewMoney(targetAmountCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid target amount: %w", err)
	}

	// Parse deadline
	deadline, err := time.Parse("2006-01-02", input.Deadline)
	if err != nil {
		return nil, fmt.Errorf("invalid deadline format: %w", err)
	}

	// Create account context value object
	context, err := sharedvalueobjects.NewAccountContext(input.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid account context: %w", err)
	}

	// Create goal entity
	goal, err := entities.NewGoal(
		userID,
		name,
		targetAmount,
		deadline,
		context,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	// Save goal to repository
	if err := uc.goalRepository.Save(goal); err != nil {
		return nil, fmt.Errorf("failed to save goal: %w", err)
	}

	// Publish domain events
	domainEvents := goal.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the goal creation
			_ = err // Ignore for now, but should be logged
		}
	}
	goal.ClearEvents()

	// Build output
	currentAmount := goal.CurrentAmount()
	output := &dtos.CreateGoalOutput{
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
	}

	return output, nil
}
