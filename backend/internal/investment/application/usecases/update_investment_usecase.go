package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// UpdateInvestmentUseCase handles updating an investment.
type UpdateInvestmentUseCase struct {
	investmentRepository repositories.InvestmentRepository
	eventBus            *eventbus.EventBus
}

// NewUpdateInvestmentUseCase creates a new UpdateInvestmentUseCase instance.
func NewUpdateInvestmentUseCase(
	investmentRepository repositories.InvestmentRepository,
	eventBus *eventbus.EventBus,
) *UpdateInvestmentUseCase {
	return &UpdateInvestmentUseCase{
		investmentRepository: investmentRepository,
		eventBus:              eventBus,
	}
}

// Execute performs the investment update.
// It validates the input, updates the investment entity,
// saves it to the repository, and publishes domain events.
func (uc *UpdateInvestmentUseCase) Execute(input dtos.UpdateInvestmentInput) (*dtos.UpdateInvestmentOutput, error) {
	// Create investment ID value object
	investmentID, err := investmentvalueobjects.NewInvestmentID(input.InvestmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid investment ID: %w", err)
	}

	// Find investment by ID
	investment, err := uc.investmentRepository.FindByID(investmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find investment: %w", err)
	}

	if investment == nil {
		return nil, errors.New("investment not found")
	}

	// Update current value if provided
	if input.CurrentValue != nil {
		currency := investment.CurrentValue().Currency()
		currentValueCents := int64(*input.CurrentValue * 100)
		currentValue, err := sharedvalueobjects.NewMoney(currentValueCents, currency)
		if err != nil {
			return nil, fmt.Errorf("invalid current value: %w", err)
		}

		if err := investment.UpdateCurrentValue(currentValue); err != nil {
			return nil, fmt.Errorf("failed to update current value: %w", err)
		}
	}

	// Update quantity if provided
	if input.Quantity != nil {
		investmentType := investment.InvestmentType()
		if investmentType.RequiresQuantity() {
			// Calculate difference
			currentQuantity := investment.Quantity()
			if currentQuantity == nil {
				// Add quantity
				if err := investment.AddQuantity(*input.Quantity); err != nil {
					return nil, fmt.Errorf("failed to add quantity: %w", err)
				}
			} else {
				// Update quantity based on difference
				diff := *input.Quantity - *currentQuantity
				if diff > 0 {
					if err := investment.AddQuantity(diff); err != nil {
						return nil, fmt.Errorf("failed to add quantity: %w", err)
					}
				} else if diff < 0 {
					if err := investment.RemoveQuantity(-diff); err != nil {
						return nil, fmt.Errorf("failed to remove quantity: %w", err)
					}
				}
			}
		} else {
			return nil, fmt.Errorf("this investment type does not support quantity")
		}
	}

	// Save investment to repository
	if err := uc.investmentRepository.Save(investment); err != nil {
		return nil, fmt.Errorf("failed to save investment: %w", err)
	}

	// Publish domain events
	domainEvents := investment.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			_ = err // Ignore for now, but should be logged
		}
	}
	investment.ClearEvents()

	// Calculate return
	returnObj := investment.CalculateReturn()
	currentValue := investment.CurrentValue()

	// Build output
	output := &dtos.UpdateInvestmentOutput{
		InvestmentID:    investment.ID().Value(),
		UserID:          investment.UserID().Value(),
		AccountID:       investment.AccountID().Value(),
		Type:            investment.InvestmentType().Value(),
		Name:            investment.Name().Name(),
		Ticker:           investment.Name().Ticker(),
		PurchaseDate:     investment.PurchaseDate().Format("2006-01-02"),
		PurchaseAmount:   investment.PurchaseAmount().Float64(),
		CurrentValue:     currentValue.Float64(),
		Currency:         currentValue.Currency().Code(),
		Quantity:         investment.Quantity(),
		Context:          investment.Context().Value(),
		ReturnAbsolute:   returnObj.Absolute().Float64(),
		ReturnPercentage: returnObj.Percentage(),
		UpdatedAt:        investment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}

