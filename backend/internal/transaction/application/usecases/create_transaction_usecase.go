package usecases

import (
	"fmt"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedrepositories "gestao-financeira/backend/internal/shared/domain/repositories"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// CreateTransactionUseCase handles transaction creation.
// It uses UnitOfWork to ensure atomicity when creating a transaction and updating account balance.
type CreateTransactionUseCase struct {
	unitOfWork sharedrepositories.UnitOfWork
	eventBus   *eventbus.EventBus
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase instance.
func NewCreateTransactionUseCase(
	unitOfWork sharedrepositories.UnitOfWork,
	eventBus *eventbus.EventBus,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		unitOfWork: unitOfWork,
		eventBus:   eventBus,
	}
}

// Execute performs the transaction creation.
// It validates the input, creates value objects, creates a new transaction entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateTransactionUseCase) Execute(input dtos.CreateTransactionInput) (*dtos.CreateTransactionOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create account ID value object
	accountID, err := accountvalueobjects.NewAccountID(input.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	// Create transaction type value object
	transactionType, err := transactionvalueobjects.NewTransactionType(input.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction type: %w", err)
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(input.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	// Create amount (convert float to cents)
	amountCents := int64(input.Amount * 100)
	amount, err := sharedvalueobjects.NewMoney(amountCents, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// Create transaction description value object
	description, err := transactionvalueobjects.NewTransactionDescription(input.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction description: %w", err)
	}

	// Parse transaction date
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: expected YYYY-MM-DD, got %s", input.Date)
	}

	// Create transaction entity
	transaction, err := entities.NewTransaction(userID, accountID, transactionType, amount, description, date)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Get repositories from UnitOfWork
	transactionRepository := uc.unitOfWork.TransactionRepository()
	accountRepository := uc.unitOfWork.AccountRepository()

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

	// Save transaction to repository (within transaction)
	if err := transactionRepository.Save(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Find account and update balance (within transaction)
	account, err := accountRepository.FindByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found: %s", accountID.Value())
	}

	// Update account balance based on transaction type
	if transactionType.Value() == "INCOME" {
		if err := account.Credit(amount); err != nil {
			return nil, fmt.Errorf("failed to credit account: %w", err)
		}
	} else if transactionType.Value() == "EXPENSE" {
		if err := account.Debit(amount); err != nil {
			return nil, fmt.Errorf("failed to debit account: %w", err)
		}
	}

	// Save updated account (within transaction)
	if err := accountRepository.Save(account); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	// Commit transaction (all operations succeed)
	if err := uc.unitOfWork.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish domain events (after successful commit)
	// Events are published outside the transaction to avoid blocking the commit
	domainEvents := transaction.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the transaction creation
			// In production, you might want to handle this differently
			// (e.g., store events in an outbox pattern)
			_ = err // Ignore for now, but should be logged
		}
	}
	transaction.ClearEvents()

	// Build output
	transactionAmount := transaction.Amount()
	output := &dtos.CreateTransactionOutput{
		TransactionID: transaction.ID().Value(),
		UserID:        transaction.UserID().Value(),
		AccountID:     transaction.AccountID().Value(),
		Type:          transaction.TransactionType().Value(),
		Amount:        transactionAmount.Float64(),
		Currency:      transactionAmount.Currency().Code(),
		Description:   transaction.Description().Value(),
		Date:          transaction.Date().Format("2006-01-02"),
		CreatedAt:     transaction.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
