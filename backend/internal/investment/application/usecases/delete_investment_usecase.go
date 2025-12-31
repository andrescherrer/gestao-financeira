package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
)

// DeleteInvestmentUseCase handles deleting an investment.
type DeleteInvestmentUseCase struct {
	investmentRepository repositories.InvestmentRepository
}

// NewDeleteInvestmentUseCase creates a new DeleteInvestmentUseCase instance.
func NewDeleteInvestmentUseCase(
	investmentRepository repositories.InvestmentRepository,
) *DeleteInvestmentUseCase {
	return &DeleteInvestmentUseCase{
		investmentRepository: investmentRepository,
	}
}

// Execute performs the investment deletion.
// It validates the input, checks if the investment exists,
// and deletes it from the repository (soft delete).
func (uc *DeleteInvestmentUseCase) Execute(input dtos.DeleteInvestmentInput) (*dtos.DeleteInvestmentOutput, error) {
	// Create investment ID value object
	investmentID, err := investmentvalueobjects.NewInvestmentID(input.InvestmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid investment ID: %w", err)
	}

	// Check if investment exists
	exists, err := uc.investmentRepository.Exists(investmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check investment existence: %w", err)
	}

	if !exists {
		return nil, errors.New("investment not found")
	}

	// Delete investment (soft delete)
	if err := uc.investmentRepository.Delete(investmentID); err != nil {
		return nil, fmt.Errorf("failed to delete investment: %w", err)
	}

	output := &dtos.DeleteInvestmentOutput{
		Success: true,
	}

	return output, nil
}
