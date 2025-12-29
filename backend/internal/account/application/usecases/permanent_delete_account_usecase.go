package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
)

// PermanentDeleteAccountUseCase handles permanent account deletion (hard delete).
// This should only be used by administrators.
type PermanentDeleteAccountUseCase struct {
	accountRepository repositories.AccountRepository
}

// NewPermanentDeleteAccountUseCase creates a new PermanentDeleteAccountUseCase instance.
func NewPermanentDeleteAccountUseCase(
	accountRepository repositories.AccountRepository,
) *PermanentDeleteAccountUseCase {
	return &PermanentDeleteAccountUseCase{
		accountRepository: accountRepository,
	}
}

// Execute performs the permanent account deletion.
func (uc *PermanentDeleteAccountUseCase) Execute(input dtos.PermanentDeleteAccountInput) (*dtos.PermanentDeleteAccountOutput, error) {
	// Create account ID value object
	accountID, err := valueobjects.NewAccountID(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Try to permanently delete
	repo, ok := uc.accountRepository.(interface {
		PermanentDelete(valueobjects.AccountID) error
	})
	if !ok {
		return nil, errors.New("repository does not support permanent delete operation")
	}

	if err := repo.PermanentDelete(accountID); err != nil {
		return nil, fmt.Errorf("failed to permanently delete account: %w", err)
	}

	output := &dtos.PermanentDeleteAccountOutput{
		Message:   "Account permanently deleted successfully",
		AccountID: accountID.Value(),
	}

	return output, nil
}
