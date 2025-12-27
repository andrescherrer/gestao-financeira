package services

import (
	"fmt"
	"time"

	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// RecurringTransactionProcessor processes recurring transactions and generates new instances.
type RecurringTransactionProcessor struct {
	transactionRepository repositories.TransactionRepository
	eventBus              *eventbus.EventBus
}

// NewRecurringTransactionProcessor creates a new RecurringTransactionProcessor instance.
func NewRecurringTransactionProcessor(
	transactionRepository repositories.TransactionRepository,
	eventBus *eventbus.EventBus,
) *RecurringTransactionProcessor {
	return &RecurringTransactionProcessor{
		transactionRepository: transactionRepository,
		eventBus:              eventBus,
	}
}

// ProcessRecurringTransactions processes all active recurring transactions and generates new instances if needed.
// Returns the number of transactions created and any errors encountered.
func (p *RecurringTransactionProcessor) ProcessRecurringTransactions() (int, error) {
	// Find all active recurring transactions
	recurringTransactions, err := p.transactionRepository.FindActiveRecurringTransactions()
	if err != nil {
		return 0, fmt.Errorf("failed to find active recurring transactions: %w", err)
	}

	createdCount := 0
	now := time.Now()

	for _, recurringTx := range recurringTransactions {
		// Check if we need to create a new instance
		shouldCreate, nextDate, err := p.shouldCreateNextInstance(recurringTx, now)
		if err != nil {
			// Log error but continue processing other transactions
			continue
		}

		if !shouldCreate {
			continue
		}

		// Create new transaction instance
		parentID := recurringTx.ID()
		newTransaction, err := p.createNextInstance(recurringTx, nextDate, &parentID)
		if err != nil {
			// Log error but continue processing other transactions
			continue
		}

		// Save the new transaction
		if err := p.transactionRepository.Save(newTransaction); err != nil {
			// Log error but continue processing other transactions
			continue
		}

		// Publish domain events
		domainEvents := newTransaction.GetEvents()
		for _, event := range domainEvents {
			if err := p.eventBus.Publish(event); err != nil {
				// Log error but don't fail
				_ = err
			}
		}
		newTransaction.ClearEvents()

		createdCount++
	}

	return createdCount, nil
}

// shouldCreateNextInstance determines if a new instance should be created for a recurring transaction.
// Returns (shouldCreate, nextDate, error).
func (p *RecurringTransactionProcessor) shouldCreateNextInstance(
	recurringTx *entities.Transaction,
	now time.Time,
) (bool, time.Time, error) {
	if !recurringTx.IsRecurring() {
		return false, time.Time{}, fmt.Errorf("transaction is not recurring")
	}

	frequency := recurringTx.RecurrenceFrequency()
	if frequency == nil {
		return false, time.Time{}, fmt.Errorf("recurrence frequency is nil")
	}

	// Calculate next execution date based on frequency
	// Start from the original transaction date and keep adding periods until we reach or pass today
	nextDate := p.calculateNextDate(recurringTx.Date(), frequency, now)

	// Check if next date is in the future (not yet time to create)
	// We compare dates only (not time) to allow creation on the same day
	nextDateOnly := time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, nextDate.Location())
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if nextDateOnly.After(nowOnly) {
		return false, time.Time{}, nil // Not yet time to create
	}

	// Check if we've passed the end date
	endDate := recurringTx.RecurrenceEndDate()
	if endDate != nil && !endDate.IsZero() && nextDate.After(*endDate) {
		return false, time.Time{}, nil // Past end date
	}

	// Check if an instance for this date already exists
	// We'll check by looking for transactions with the same parent and date
	// This is a simple check - in production, you might want a more robust approach
	exists, err := p.instanceExistsForDate(recurringTx.ID(), nextDate)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("failed to check if instance exists: %w", err)
	}

	if exists {
		return false, time.Time{}, nil // Instance already exists
	}

	return true, nextDate, nil
}

// calculateNextDate calculates the next execution date based on the frequency.
// It starts from the original transaction date and keeps adding periods until we reach or pass today.
func (p *RecurringTransactionProcessor) calculateNextDate(
	originalDate time.Time,
	frequency *transactionvalueobjects.RecurrenceFrequency,
	now time.Time,
) time.Time {
	// Start from the original transaction date
	nextDate := originalDate
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Keep adding periods until we reach or pass today
	for {
		nextDateOnly := time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, nextDate.Location())

		// If we've reached or passed today, this is the next date to create
		if !nextDateOnly.Before(nowOnly) {
			break
		}

		// Add one period
		switch {
		case frequency.IsDaily():
			nextDate = nextDate.AddDate(0, 0, 1)
		case frequency.IsWeekly():
			nextDate = nextDate.AddDate(0, 0, 7)
		case frequency.IsMonthly():
			nextDate = nextDate.AddDate(0, 1, 0)
		case frequency.IsYearly():
			nextDate = nextDate.AddDate(1, 0, 0)
		default:
			// Should not happen, but return originalDate + 1 day as fallback
			nextDate = nextDate.AddDate(0, 0, 1)
		}
	}

	return nextDate
}

// instanceExistsForDate checks if an instance already exists for a given parent and date.
func (p *RecurringTransactionProcessor) instanceExistsForDate(
	parentID transactionvalueobjects.TransactionID,
	date time.Time,
) (bool, error) {
	existing, err := p.transactionRepository.FindByParentIDAndDate(parentID, date)
	if err != nil {
		return false, fmt.Errorf("failed to check if instance exists: %w", err)
	}
	return existing != nil, nil
}

// createNextInstance creates a new transaction instance from a recurring transaction.
func (p *RecurringTransactionProcessor) createNextInstance(
	recurringTx *entities.Transaction,
	nextDate time.Time,
	parentID *transactionvalueobjects.TransactionID,
) (*entities.Transaction, error) {
	// Create new transaction with same properties but new date and parent reference
	newTransaction, err := entities.NewTransactionWithRecurrence(
		recurringTx.UserID(),
		recurringTx.AccountID(),
		recurringTx.TransactionType(),
		recurringTx.Amount(),
		recurringTx.Description(),
		nextDate,
		false, // This instance is not recurring
		nil,   // No recurrence frequency for instances
		nil,   // No end date for instances
		parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create next instance: %w", err)
	}

	return newTransaction, nil
}
