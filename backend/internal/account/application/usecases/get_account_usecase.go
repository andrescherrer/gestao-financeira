package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
)

// GetAccountUseCase handles retrieving a single account by ID.
type GetAccountUseCase struct {
	accountRepository repositories.AccountRepository
}

// NewGetAccountUseCase creates a new GetAccountUseCase instance.
func NewGetAccountUseCase(
	accountRepository repositories.AccountRepository,
) *GetAccountUseCase {
	return &GetAccountUseCase{
		accountRepository: accountRepository,
	}
}

// Execute performs the account retrieval.
// It validates the input, retrieves the account from the repository,
// and returns it as a DTO.
func (uc *GetAccountUseCase) Execute(input dtos.GetAccountInput) (*dtos.GetAccountOutput, error) {
	// Create account ID value object
	accountID, err := valueobjects.NewAccountID(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Find account by ID
	account, err := uc.accountRepository.FindByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	// Check if account exists
	if account == nil {
		return nil, errors.New("account not found")
	}

	// Convert to output DTO
	balance := account.Balance()
	output := &dtos.GetAccountOutput{
		AccountID: account.ID().Value(),
		UserID:    account.UserID().Value(),
		Name:      account.Name().Value(),
		Type:      account.AccountType().Value(),
		Balance:   balance.Float64(),
		Currency:  balance.Currency().Code(),
		Context:   account.Context().Value(),
		IsActive:  account.IsActive(),
		CreatedAt: account.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: account.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
