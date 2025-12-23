package usecases

import (
	"fmt"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// CreateTransactionUseCase handles transaction creation.
type CreateTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
	eventBus              *eventbus.EventBus
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase instance.
func NewCreateTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
	eventBus *eventbus.EventBus,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		transactionRepository: transactionRepository,
		eventBus:              eventBus,
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

	// Save transaction to repository
	if err := uc.transactionRepository.Save(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Publish domain events
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
