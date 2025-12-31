package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
)

// GetInvestmentUseCase handles retrieving a single investment by ID.
type GetInvestmentUseCase struct {
	investmentRepository repositories.InvestmentRepository
}

// NewGetInvestmentUseCase creates a new GetInvestmentUseCase instance.
func NewGetInvestmentUseCase(
	investmentRepository repositories.InvestmentRepository,
) *GetInvestmentUseCase {
	return &GetInvestmentUseCase{
		investmentRepository: investmentRepository,
	}
}

// Execute performs the investment retrieval.
// It validates the input, retrieves the investment from the repository,
// and returns it as a DTO.
func (uc *GetInvestmentUseCase) Execute(input dtos.GetInvestmentInput) (*dtos.GetInvestmentOutput, error) {
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

	// Check if investment exists
	if investment == nil {
		return nil, errors.New("investment not found")
	}

	// Calculate return
	returnObj := investment.CalculateReturn()
	currentValue := investment.CurrentValue()

	// Convert to output DTO
	output := &dtos.GetInvestmentOutput{
		InvestmentID:     investment.ID().Value(),
		UserID:           investment.UserID().Value(),
		AccountID:        investment.AccountID().Value(),
		Type:             investment.InvestmentType().Value(),
		Name:             investment.Name().Name(),
		Ticker:           investment.Name().Ticker(),
		PurchaseDate:     investment.PurchaseDate().Format("2006-01-02"),
		PurchaseAmount:   investment.PurchaseAmount().Float64(),
		CurrentValue:     currentValue.Float64(),
		Currency:         currentValue.Currency().Code(),
		Quantity:         investment.Quantity(),
		Context:          investment.Context().Value(),
		ReturnAbsolute:   returnObj.Absolute().Float64(),
		ReturnPercentage: returnObj.Percentage(),
		CreatedAt:        investment.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        investment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
