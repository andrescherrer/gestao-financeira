package usecases

import (
	"fmt"

	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/entities"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	"gestao-financeira/backend/pkg/pagination"
)

// ListInvestmentsUseCase handles listing investments for a user.
type ListInvestmentsUseCase struct {
	investmentRepository repositories.InvestmentRepository
}

// NewListInvestmentsUseCase creates a new ListInvestmentsUseCase instance.
func NewListInvestmentsUseCase(
	investmentRepository repositories.InvestmentRepository,
) *ListInvestmentsUseCase {
	return &ListInvestmentsUseCase{
		investmentRepository: investmentRepository,
	}
}

// Execute performs the investment listing.
// It validates the input, retrieves investments from the repository,
// and returns them as DTOs. Supports pagination.
func (uc *ListInvestmentsUseCase) Execute(input dtos.ListInvestmentsInput) (*dtos.ListInvestmentsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse pagination parameters
	paginationParams := pagination.ParsePaginationParams(input.Page, input.Limit)
	usePagination := input.Page != "" || input.Limit != ""

	var domainInvestments []*entities.Investment
	var total int64

	// Check if we should use pagination
	if usePagination {
		// Use paginated query
		domainInvestments, total, err = uc.investmentRepository.FindByUserIDWithPagination(
			userID,
			input.Context,
			input.Type,
			paginationParams.CalculateOffset(),
			paginationParams.Limit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find investments: %w", err)
		}
	} else {
		// Use non-paginated query (backward compatibility)
		if input.Type != "" {
			// Create investment type value object
			investmentType, err := investmentvalueobjects.NewInvestmentType(input.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid investment type: %w", err)
			}

			// Find investments by user ID and type
			domainInvestments, err = uc.investmentRepository.FindByType(userID, investmentType)
			if err != nil {
				return nil, fmt.Errorf("failed to find investments: %w", err)
			}
			total = int64(len(domainInvestments))
		} else {
			// Find all investments for the user
			domainInvestments, err = uc.investmentRepository.FindByUserID(userID)
			if err != nil {
				return nil, fmt.Errorf("failed to find investments: %w", err)
			}
			total = int64(len(domainInvestments))
		}
	}

	investments := uc.toInvestmentOutputs(domainInvestments)

	output := &dtos.ListInvestmentsOutput{
		Investments: investments,
		Count:       len(investments),
	}

	// Add pagination metadata if pagination was used
	if usePagination {
		paginationResult := pagination.BuildPaginationResult(paginationParams, total)
		output.Pagination = &paginationResult
	}

	return output, nil
}

// toInvestmentOutputs converts domain investments to DTOs.
func (uc *ListInvestmentsUseCase) toInvestmentOutputs(domainInvestments []*entities.Investment) []dtos.InvestmentOutput {
	outputs := make([]dtos.InvestmentOutput, 0, len(domainInvestments))
	for _, investment := range domainInvestments {
		returnObj := investment.CalculateReturn()
		currentValue := investment.CurrentValue()
		output := dtos.InvestmentOutput{
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
		outputs = append(outputs, output)
	}
	return outputs
}
