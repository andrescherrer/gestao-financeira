package usecases

import (
	"errors"
	"fmt"
	"time"

	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	transactionevents "gestao-financeira/backend/internal/transaction/domain/events"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// UpdateTransactionUseCase handles transaction updates.
type UpdateTransactionUseCase struct {
	transactionRepository repositories.TransactionRepository
	eventBus              *eventbus.EventBus
}

// NewUpdateTransactionUseCase creates a new UpdateTransactionUseCase instance.
func NewUpdateTransactionUseCase(
	transactionRepository repositories.TransactionRepository,
	eventBus *eventbus.EventBus,
) *UpdateTransactionUseCase {
	return &UpdateTransactionUseCase{
		transactionRepository: transactionRepository,
		eventBus:              eventBus,
	}
}

// Execute performs the transaction update.
// It validates the input, retrieves the transaction, updates the specified fields,
// saves it to the repository, and publishes domain events.
func (uc *UpdateTransactionUseCase) Execute(input dtos.UpdateTransactionInput) (*dtos.UpdateTransactionOutput, error) {
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

	// Store old values BEFORE any updates (needed for TransactionUpdated event)
	oldType := transaction.TransactionType().Value()
	oldAmount := transaction.Amount()

	// Update type if provided
	if input.Type != nil {
		transactionType, err := transactionvalueobjects.NewTransactionType(*input.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction type: %w", err)
		}
		if err := transaction.UpdateType(transactionType); err != nil {
			return nil, fmt.Errorf("failed to update transaction type: %w", err)
		}
	}

	// Update amount if provided
	if input.Amount != nil {
		// Get currency - use existing if not provided, otherwise use new currency
		var currency sharedvalueobjects.Currency
		if input.Currency != nil {
			var err error
			currency, err = sharedvalueobjects.NewCurrency(*input.Currency)
			if err != nil {
				return nil, fmt.Errorf("invalid currency: %w", err)
			}
		} else {
			currency = transaction.Amount().Currency()
		}

		// Convert float to cents
		amountCents := int64(*input.Amount * 100)
		amount, err := sharedvalueobjects.NewMoney(amountCents, currency)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}

		if err := transaction.UpdateAmount(amount); err != nil {
			return nil, fmt.Errorf("failed to update transaction amount: %w", err)
		}
	}

	// Update description if provided
	if input.Description != nil {
		description, err := transactionvalueobjects.NewTransactionDescription(*input.Description)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction description: %w", err)
		}
		if err := transaction.UpdateDescription(description); err != nil {
			return nil, fmt.Errorf("failed to update transaction description: %w", err)
		}
	}

	// Update date if provided
	if input.Date != nil {
		date, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: expected YYYY-MM-DD, got %s", *input.Date)
		}
		if err := transaction.UpdateDate(date); err != nil {
			return nil, fmt.Errorf("failed to update transaction date: %w", err)
		}
	}

	// Check if at least one field was provided for update
	if input.Type == nil && input.Amount == nil && input.Description == nil && input.Date == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Save transaction to repository
	if err := uc.transactionRepository.Save(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Get new values after update
	newType := transaction.TransactionType().Value()
	newAmount := transaction.Amount()

	// Check if amount or type changed (these affect account balance)
	amountChanged := !oldAmount.Equals(newAmount)
	typeChanged := oldType != newType

	// Publish domain events from entity
	domainEvents := transaction.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the transaction update
			// In production, you might want to handle this differently
			// (e.g., store events in an outbox pattern)
			_ = err // Ignore for now, but should be logged
		}
	}
	transaction.ClearEvents()

	// If amount or type changed, publish TransactionUpdated event for balance update
	if amountChanged || typeChanged {
		updateEvent := transactionevents.NewTransactionUpdated(
			transaction.ID().Value(),
			transaction.AccountID().Value(),
			oldType,
			oldAmount,
			newType,
			newAmount,
		)
		if err := uc.eventBus.Publish(updateEvent); err != nil {
			// Log error but don't fail the transaction update
			_ = err // Ignore for now, but should be logged
		}
	}

	// Build output
	amount := transaction.Amount()
	output := &dtos.UpdateTransactionOutput{
		TransactionID: transaction.ID().Value(),
		UserID:        transaction.UserID().Value(),
		AccountID:     transaction.AccountID().Value(),
		Type:          transaction.TransactionType().Value(),
		Amount:        amount.Float64(),
		Currency:      amount.Currency().Code(),
		Description:   transaction.Description().Value(),
		Date:          transaction.Date().Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
