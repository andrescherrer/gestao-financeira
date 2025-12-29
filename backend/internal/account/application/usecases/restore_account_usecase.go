package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
)

// RestoreAccountUseCase handles account restoration from soft delete.
type RestoreAccountUseCase struct {
	accountRepository repositories.AccountRepository
}

// NewRestoreAccountUseCase creates a new RestoreAccountUseCase instance.
func NewRestoreAccountUseCase(
	accountRepository repositories.AccountRepository,
) *RestoreAccountUseCase {
	return &RestoreAccountUseCase{
		accountRepository: accountRepository,
	}
}

// Execute performs the account restoration.
func (uc *RestoreAccountUseCase) Execute(input dtos.RestoreAccountInput) (*dtos.RestoreAccountOutput, error) {
	// Create account ID value object
	accountID, err := valueobjects.NewAccountID(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Check if account exists
	account, err := uc.accountRepository.FindByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	// If account exists and is not deleted, return error
	if account != nil {
		return nil, errors.New("account is not deleted")
	}

	// Try to restore
	repo, ok := uc.accountRepository.(interface {
		Restore(valueobjects.AccountID) error
	})
	if !ok {
		return nil, errors.New("repository does not support restore operation")
	}

	if err := repo.Restore(accountID); err != nil {
		return nil, fmt.Errorf("failed to restore account: %w", err)
	}

	output := &dtos.RestoreAccountOutput{
		Message:   "Account restored successfully",
		AccountID: accountID.Value(),
	}

	return output, nil
}
