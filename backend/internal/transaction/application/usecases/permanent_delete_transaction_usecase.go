package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// PermanentDeleteTransactionUseCase handles permanent transaction deletion (hard delete).
// This should only be used by administrators.
type PermanentDeleteTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewPermanentDeleteTransactionUseCase creates a new PermanentDeleteTransactionUseCase instance.
func NewPermanentDeleteTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
) *PermanentDeleteTransactionUseCase {
	return &PermanentDeleteTransactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute performs the permanent transaction deletion.
func (uc *PermanentDeleteTransactionUseCase) Execute(input dtos.PermanentDeleteTransactionInput) (*dtos.PermanentDeleteTransactionOutput, error) {
	// Create transaction ID value object
	transactionID, err := transactionvalueobjects.NewTransactionID(input.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	// Check if transaction exists (including soft-deleted)
	transaction, err := uc.transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	// If transaction doesn't exist at all, return error
	if transaction == nil {
		// Check if it's soft-deleted by trying to restore first
		// If restore fails, it means it doesn't exist
		repo, ok := uc.transactionRepository.(interface {
			Restore(transactionvalueobjects.TransactionID) error
		})
		if !ok || repo.Restore(transactionID) != nil {
			return nil, errors.New("transaction not found")
		}
		// If restore succeeded, we need to delete it again first
		// Actually, let's just proceed with permanent delete
	}

	// Try to permanently delete
	repo, ok := uc.transactionRepository.(interface {
		PermanentDelete(transactionvalueobjects.TransactionID) error
	})
	if !ok {
		return nil, errors.New("repository does not support permanent delete operation")
	}

	if err := repo.PermanentDelete(transactionID); err != nil {
		return nil, fmt.Errorf("failed to permanently delete transaction: %w", err)
	}

	output := &dtos.PermanentDeleteTransactionOutput{
		Message:       "Transaction permanently deleted successfully",
		TransactionID: transactionID.Value(),
	}

	return output, nil
}
