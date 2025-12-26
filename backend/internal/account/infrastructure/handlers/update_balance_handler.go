package handlers

import (
	"fmt"

	"gestao-financeira/backend/internal/account/domain/repositories"
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	transactionevents "gestao-financeira/backend/internal/transaction/domain/events"

	"github.com/rs/zerolog/log"
)

// UpdateBalanceHandler handles domain events related to transactions
// and updates the account balance accordingly.
type UpdateBalanceHandler struct {
	accountRepository repositories.AccountRepository
}

// NewUpdateBalanceHandler creates a new UpdateBalanceHandler instance.
func NewUpdateBalanceHandler(accountRepository repositories.AccountRepository) *UpdateBalanceHandler {
	return &UpdateBalanceHandler{
		accountRepository: accountRepository,
	}
}

// HandleTransactionCreated handles TransactionCreated events and updates account balance.
func (h *UpdateBalanceHandler) HandleTransactionCreated(event events.DomainEvent) error {
	transactionCreated, ok := event.(*transactionevents.TransactionCreated)
	if !ok {
		return fmt.Errorf("expected TransactionCreated event, got %T", event)
	}

	// Get account ID
	accountID, err := accountvalueobjects.NewAccountID(transactionCreated.AccountID())
	if err != nil {
		return fmt.Errorf("invalid account ID in event: %w", err)
	}

	// Find account
	account, err := h.accountRepository.FindByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("account not found: %s", transactionCreated.AccountID())
	}

	// Create money value object from event
	currency, err := sharedvalueobjects.NewCurrency(transactionCreated.Currency())
	if err != nil {
		return fmt.Errorf("invalid currency in event: %w", err)
	}

	amount, err := sharedvalueobjects.NewMoney(transactionCreated.Amount(), currency)
	if err != nil {
		return fmt.Errorf("invalid amount in event: %w", err)
	}

	// Update balance based on transaction type
	// INCOME: credit (add to balance)
	// EXPENSE: debit (subtract from balance)
	if transactionCreated.TransactionType() == "INCOME" {
		if err := account.Credit(amount); err != nil {
			return fmt.Errorf("failed to credit account: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionCreated.AggregateID()).
			Str("type", "INCOME").
			Int64("amount", transactionCreated.Amount()).
			Msg("Account credited due to transaction creation")
	} else if transactionCreated.TransactionType() == "EXPENSE" {
		if err := account.Debit(amount); err != nil {
			return fmt.Errorf("failed to debit account: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionCreated.AggregateID()).
			Str("type", "EXPENSE").
			Int64("amount", transactionCreated.Amount()).
			Msg("Account debited due to transaction creation")
	}

	// Save updated account
	if err := h.accountRepository.Save(account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}

// HandleTransactionUpdated handles TransactionUpdated events and updates account balance.
// It reverses the old transaction effect and applies the new transaction effect.
func (h *UpdateBalanceHandler) HandleTransactionUpdated(event events.DomainEvent) error {
	transactionUpdated, ok := event.(*transactionevents.TransactionUpdated)
	if !ok {
		return fmt.Errorf("expected TransactionUpdated event, got %T", event)
	}

	// Get account ID
	accountID, err := accountvalueobjects.NewAccountID(transactionUpdated.AccountID())
	if err != nil {
		return fmt.Errorf("invalid account ID in event: %w", err)
	}

	// Find account
	account, err := h.accountRepository.FindByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("account not found: %s", transactionUpdated.AccountID())
	}

	// Create currency value object
	currency, err := sharedvalueobjects.NewCurrency(transactionUpdated.Currency())
	if err != nil {
		return fmt.Errorf("invalid currency in event: %w", err)
	}

	// Reverse old transaction effect
	oldAmount, err := sharedvalueobjects.NewMoney(transactionUpdated.OldAmount(), currency)
	if err != nil {
		return fmt.Errorf("invalid old amount in event: %w", err)
	}

	if transactionUpdated.OldType() == "INCOME" {
		// Reverse: debit (subtract) what was previously credited
		if err := account.Debit(oldAmount); err != nil {
			return fmt.Errorf("failed to reverse old income transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionUpdated.AggregateID()).
			Str("action", "reverse_income").
			Int64("amount", transactionUpdated.OldAmount()).
			Msg("Reversing old income transaction")
	} else if transactionUpdated.OldType() == "EXPENSE" {
		// Reverse: credit (add) what was previously debited
		if err := account.Credit(oldAmount); err != nil {
			return fmt.Errorf("failed to reverse old expense transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionUpdated.AggregateID()).
			Str("action", "reverse_expense").
			Int64("amount", transactionUpdated.OldAmount()).
			Msg("Reversing old expense transaction")
	}

	// Apply new transaction effect
	newAmount, err := sharedvalueobjects.NewMoney(transactionUpdated.NewAmount(), currency)
	if err != nil {
		return fmt.Errorf("invalid new amount in event: %w", err)
	}

	if transactionUpdated.NewType() == "INCOME" {
		// Apply: credit (add) new income
		if err := account.Credit(newAmount); err != nil {
			return fmt.Errorf("failed to apply new income transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionUpdated.AggregateID()).
			Str("type", "INCOME").
			Int64("amount", transactionUpdated.NewAmount()).
			Msg("Account credited due to transaction update")
	} else if transactionUpdated.NewType() == "EXPENSE" {
		// Apply: debit (subtract) new expense
		if err := account.Debit(newAmount); err != nil {
			return fmt.Errorf("failed to apply new expense transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionUpdated.AggregateID()).
			Str("type", "EXPENSE").
			Int64("amount", transactionUpdated.NewAmount()).
			Msg("Account debited due to transaction update")
	}

	// Save updated account
	if err := h.accountRepository.Save(account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}

// HandleTransactionDeleted handles TransactionDeleted events and reverses the balance update.
func (h *UpdateBalanceHandler) HandleTransactionDeleted(event events.DomainEvent) error {
	transactionDeleted, ok := event.(*transactionevents.TransactionDeleted)
	if !ok {
		return fmt.Errorf("expected TransactionDeleted event, got %T", event)
	}

	// Get account ID
	accountID, err := accountvalueobjects.NewAccountID(transactionDeleted.AccountID())
	if err != nil {
		return fmt.Errorf("invalid account ID in event: %w", err)
	}

	// Find account
	account, err := h.accountRepository.FindByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("account not found: %s", transactionDeleted.AccountID())
	}

	// Create money value object from event
	currency, err := sharedvalueobjects.NewCurrency(transactionDeleted.Currency())
	if err != nil {
		return fmt.Errorf("invalid currency in event: %w", err)
	}

	amount, err := sharedvalueobjects.NewMoney(transactionDeleted.Amount(), currency)
	if err != nil {
		return fmt.Errorf("invalid amount in event: %w", err)
	}

	// Reverse the transaction effect
	// INCOME: reverse credit (debit)
	// EXPENSE: reverse debit (credit)
	if transactionDeleted.TransactionType() == "INCOME" {
		// Reverse: debit (subtract) what was previously credited
		if err := account.Debit(amount); err != nil {
			return fmt.Errorf("failed to reverse income transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionDeleted.AggregateID()).
			Str("action", "reverse_income").
			Int64("amount", transactionDeleted.Amount()).
			Msg("Reversing income transaction due to deletion")
	} else if transactionDeleted.TransactionType() == "EXPENSE" {
		// Reverse: credit (add) what was previously debited
		if err := account.Credit(amount); err != nil {
			return fmt.Errorf("failed to reverse expense transaction: %w", err)
		}
		log.Info().
			Str("account_id", accountID.Value()).
			Str("transaction_id", transactionDeleted.AggregateID()).
			Str("action", "reverse_expense").
			Int64("amount", transactionDeleted.Amount()).
			Msg("Reversing expense transaction due to deletion")
	}

	// Save updated account
	if err := h.accountRepository.Save(account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}
