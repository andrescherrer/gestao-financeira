package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	transactionevents "gestao-financeira/backend/internal/transaction/domain/events"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// DeleteTransactionUseCase handles transaction deletion.
type DeleteTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
	eventBus              *eventbus.EventBus
}

// NewDeleteTransactionUseCase creates a new DeleteTransactionUseCase instance.
func NewDeleteTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
	eventBus *eventbus.EventBus,
) *DeleteTransactionUseCase {
	return &DeleteTransactionUseCase{
		transactionRepository: transactionRepository,
		eventBus:              eventBus,
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

	// Store transaction details before deletion (needed for TransactionDeleted event)
	accountID := transaction.AccountID().Value()
	transactionType := transaction.TransactionType().Value()
	amount := transaction.Amount()

	// Delete transaction (soft delete)
	if err := uc.transactionRepository.Delete(transactionID); err != nil {
		return nil, fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Publish TransactionDeleted event for balance update
	deleteEvent := transactionevents.NewTransactionDeleted(
		transactionID.Value(),
		accountID,
		transactionType,
		amount,
	)
	if err := uc.eventBus.Publish(deleteEvent); err != nil {
		// Log error but don't fail the transaction deletion
		// In production, you might want to handle this differently
		_ = err // Ignore for now, but should be logged
	}

	output := &dtos.DeleteTransactionOutput{
		Message:       "Transaction deleted successfully",
		TransactionID: transactionID.Value(),
	}

	return output, nil
}
