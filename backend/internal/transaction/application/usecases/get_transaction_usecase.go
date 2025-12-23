package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// GetTransactionUseCase handles retrieving a single transaction by ID.
type GetTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewGetTransactionUseCase creates a new GetTransactionUseCase instance.
func NewGetTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
) *GetTransactionUseCase {
	return &GetTransactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute performs the transaction retrieval.
// It validates the input, retrieves the transaction from the repository,
// and returns it as a DTO.
func (uc *GetTransactionUseCase) Execute(input dtos.GetTransactionInput) (*dtos.GetTransactionOutput, error) {
	// Create transaction ID value object
	transactionID, err := transactionvalueobjects.NewTransactionID(input.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	// Find transaction by ID
	transaction, err := uc.transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	// Check if transaction exists
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Convert to output DTO
	amount := transaction.Amount()
	output := &dtos.GetTransactionOutput{
		TransactionID: transaction.ID().Value(),
		UserID:        transaction.UserID().Value(),
		AccountID:     transaction.AccountID().Value(),
		Type:          transaction.TransactionType().Value(),
		Amount:        amount.Float64(),
		Currency:      amount.Currency().Code(),
		Description:   transaction.Description().Value(),
		Date:          transaction.Date().Format("2006-01-02"),
		CreatedAt:     transaction.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     transaction.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
