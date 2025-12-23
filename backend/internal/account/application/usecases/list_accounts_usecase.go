package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/repositories"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// ListAccountsUseCase handles listing accounts for a user.
type ListAccountsUseCase struct {
	accountRepository repositories.AccountRepository
}

// NewListAccountsUseCase creates a new ListAccountsUseCase instance.
func NewListAccountsUseCase(
	accountRepository repositories.AccountRepository,
) *ListAccountsUseCase {
	return &ListAccountsUseCase{
		accountRepository: accountRepository,
	}
}

// Execute performs the account listing.
// It validates the input, retrieves accounts from the repository,
// and returns them as DTOs.
func (uc *ListAccountsUseCase) Execute(input dtos.ListAccountsInput) (*dtos.ListAccountsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var domainAccounts []*entities.Account

	// If context is provided, filter by context
	if input.Context != "" {
		// Create account context value object
		context, err := sharedvalueobjects.NewAccountContext(input.Context)
		if err != nil {
			return nil, fmt.Errorf("invalid account context: %w", err)
		}

		// Find accounts by user ID and context
		domainAccounts, err = uc.accountRepository.FindByUserIDAndContext(userID, context)
		if err != nil {
			return nil, fmt.Errorf("failed to find accounts: %w", err)
		}
	} else {
		// Find all accounts for the user
		domainAccounts, err = uc.accountRepository.FindByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find accounts: %w", err)
		}
	}

	accounts := uc.toAccountOutputs(domainAccounts)

	output := &dtos.ListAccountsOutput{
		Accounts: accounts,
		Count:    len(accounts),
	}

	return output, nil
}

// toAccountOutputs converts domain accounts to DTOs.
func (uc *ListAccountsUseCase) toAccountOutputs(domainAccounts []*entities.Account) []*dtos.AccountOutput {
	outputs := make([]*dtos.AccountOutput, 0, len(domainAccounts))
	for _, account := range domainAccounts {
		balance := account.Balance()
		output := &dtos.AccountOutput{
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
		outputs = append(outputs, output)
	}
	return outputs
}
