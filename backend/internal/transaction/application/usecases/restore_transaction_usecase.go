package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// RestoreTransactionUseCase handles transaction restoration from soft delete.
type RestoreTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewRestoreTransactionUseCase creates a new RestoreTransactionUseCase instance.
func NewRestoreTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
) *RestoreTransactionUseCase {
	return &RestoreTransactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute performs the transaction restoration.
func (uc *RestoreTransactionUseCase) Execute(input dtos.RestoreTransactionInput) (*dtos.RestoreTransactionOutput, error) {
	// Create transaction ID value object
	transactionID, err := transactionvalueobjects.NewTransactionID(input.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	// Check if transaction exists (including soft-deleted)
	// We need to check if it's soft-deleted first
	transaction, err := uc.transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	// If transaction exists and is not deleted, return error
	if transaction != nil {
		return nil, errors.New("transaction is not deleted")
	}

	// Try to restore (repository should handle checking if it's soft-deleted)
	// We need to cast to the concrete type to access Restore method
	repo, ok := uc.transactionRepository.(interface {
		Restore(transactionvalueobjects.TransactionID) error
	})
	if !ok {
		return nil, errors.New("repository does not support restore operation")
	}

	if err := repo.Restore(transactionID); err != nil {
		return nil, fmt.Errorf("failed to restore transaction: %w", err)
	}

	output := &dtos.RestoreTransactionOutput{
		Message:       "Transaction restored successfully",
		TransactionID: transactionID.Value(),
	}

	return output, nil
}
