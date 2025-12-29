package usecases

import (
	"errors"
	"fmt"

	sharedrepositories "gestao-financeira/backend/internal/shared/domain/repositories"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	transactionevents "gestao-financeira/backend/internal/transaction/domain/events"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// DeleteTransactionUseCase handles transaction deletion.
// It uses UnitOfWork to ensure atomicity when deleting a transaction and updating account balance.
type DeleteTransactionUseCase struct {
	unitOfWork sharedrepositories.UnitOfWork
	eventBus   *eventbus.EventBus
}

// NewDeleteTransactionUseCase creates a new DeleteTransactionUseCase instance.
func NewDeleteTransactionUseCase(
	unitOfWork sharedrepositories.UnitOfWork,
	eventBus *eventbus.EventBus,
) *DeleteTransactionUseCase {
	return &DeleteTransactionUseCase{
		unitOfWork: unitOfWork,
		eventBus:   eventBus,
	}
}

// Execute performs the transaction deletion.
// It validates the input, checks if the transaction exists,
// deletes it from the repository (soft delete), and updates account balance atomically.
func (uc *DeleteTransactionUseCase) Execute(input dtos.DeleteTransactionInput) (*dtos.DeleteTransactionOutput, error) {
	// Get repositories from UnitOfWork
	transactionRepository := uc.unitOfWork.TransactionRepository()
	accountRepository := uc.unitOfWork.AccountRepository()

	// Create transaction ID value object
	transactionID, err := transactionvalueobjects.NewTransactionID(input.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	// Begin transaction to ensure atomicity
	if err := uc.unitOfWork.Begin(); err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on error
	defer func() {
		if uc.unitOfWork.IsInTransaction() {
			if rollbackErr := uc.unitOfWork.Rollback(); rollbackErr != nil {
				// Log rollback error but don't fail the function
				_ = rollbackErr
			}
		}
	}()

	// Check if transaction exists (within transaction)
	transaction, err := transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Store transaction details before deletion (needed for balance reversal and TransactionDeleted event)
	accountID := transaction.AccountID()
	transactionType := transaction.TransactionType().Value()
	amount := transaction.Amount()

	// Find account and reverse transaction effect (within transaction)
	account, err := accountRepository.FindByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found: %s", accountID.Value())
	}

	// Reverse the transaction effect
	// INCOME: reverse credit (debit)
	// EXPENSE: reverse debit (credit)
	if transactionType == "INCOME" {
		// Reverse: debit (subtract) what was previously credited
		if err := account.Debit(amount); err != nil {
			return nil, fmt.Errorf("failed to reverse income transaction: %w", err)
		}
	} else if transactionType == "EXPENSE" {
		// Reverse: credit (add) what was previously debited
		if err := account.Credit(amount); err != nil {
			return nil, fmt.Errorf("failed to reverse expense transaction: %w", err)
		}
	}

	// Save updated account (within transaction)
	if err := accountRepository.Save(account); err != nil {
		return nil, fmt.Errorf("failed to save updated account: %w", err)
	}

	// Delete transaction (soft delete, within transaction)
	if err := transactionRepository.Delete(transactionID); err != nil {
		return nil, fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Commit transaction (all operations succeed)
	if err := uc.unitOfWork.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish TransactionDeleted event for other subscribers (after successful commit)
	deleteEvent := transactionevents.NewTransactionDeleted(
		transactionID.Value(),
		accountID.Value(),
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
