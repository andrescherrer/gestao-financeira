package usecases

import (
	"fmt"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// ListTransactionsUseCase handles listing transactions for a user.
type ListTransactionsUseCase struct {
	transactionRepository repositories.TransactionRepository
}

// NewListTransactionsUseCase creates a new ListTransactionsUseCase instance.
func NewListTransactionsUseCase(
	transactionRepository repositories.TransactionRepository,
) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{
		transactionRepository: transactionRepository,
	}
}

// Execute performs the transaction listing.
// It validates the input, retrieves transactions from the repository,
// and returns them as DTOs. Supports filtering by account ID and/or type.
func (uc *ListTransactionsUseCase) Execute(input dtos.ListTransactionsInput) (*dtos.ListTransactionsOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var domainTransactions []*entities.Transaction

	// If both account ID and type are provided
	if input.AccountID != "" && input.Type != "" {
		accountID, err := accountvalueobjects.NewAccountID(input.AccountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account ID: %w", err)
		}

		transactionType, err := transactionvalueobjects.NewTransactionType(input.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction type: %w", err)
		}

		// Find by user ID, account ID, and filter by type manually
		allTransactions, err := uc.transactionRepository.FindByUserIDAndAccountID(userID, accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}

		// Filter by type
		for _, tx := range allTransactions {
			if tx.TransactionType().Value() == transactionType.Value() {
				domainTransactions = append(domainTransactions, tx)
			}
		}
	} else if input.AccountID != "" {
		// If only account ID is provided
		accountID, err := accountvalueobjects.NewAccountID(input.AccountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account ID: %w", err)
		}

		domainTransactions, err = uc.transactionRepository.FindByUserIDAndAccountID(userID, accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}
	} else if input.Type != "" {
		// If only type is provided
		transactionType, err := transactionvalueobjects.NewTransactionType(input.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction type: %w", err)
		}

		domainTransactions, err = uc.transactionRepository.FindByUserIDAndType(userID, transactionType)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}
	} else {
		// Find all transactions for the user
		domainTransactions, err = uc.transactionRepository.FindByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find transactions: %w", err)
		}
	}

	transactions := uc.toTransactionOutputs(domainTransactions)

	output := &dtos.ListTransactionsOutput{
		Transactions: transactions,
		Count:        len(transactions),
	}

	return output, nil
}

// toTransactionOutputs converts domain transactions to DTOs.
func (uc *ListTransactionsUseCase) toTransactionOutputs(domainTransactions []*entities.Transaction) []*dtos.TransactionOutput {
	outputs := make([]*dtos.TransactionOutput, 0, len(domainTransactions))
	for _, transaction := range domainTransactions {
		amount := transaction.Amount()
		output := &dtos.TransactionOutput{
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
		outputs = append(outputs, output)
	}
	return outputs
}
