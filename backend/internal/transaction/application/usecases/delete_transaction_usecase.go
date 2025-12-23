package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// DeleteTransactionUseCase handles transaction deletion.
type DeleteTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewDeleteTransactionUseCase creates a new DeleteTransactionUseCase instance.
func NewDeleteTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
) *DeleteTransactionUseCase {
	return &DeleteTransactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute performs the transaction deletion.
// It validates the input, checks if the transaction exists,
// deletes it from the repository (soft delete).
func (uc *DeleteTransactionUseCase) Execute(input dtos.DeleteTransactionInput) (*dtos.DeleteTransactionOutput, error) {
	// Create transaction ID value object
	transactionID, err := transactionvalueobjects.NewTransactionID(input.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	// Check if transaction exists
	transaction, err := uc.transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Delete transaction (soft delete)
	if err := uc.transactionRepository.Delete(transactionID); err != nil {
		return nil, fmt.Errorf("failed to delete transaction: %w", err)
	}

	output := &dtos.DeleteTransactionOutput{
		Message:       "Transaction deleted successfully",
		TransactionID: transactionID.Value(),
	}

	return output, nil
}
