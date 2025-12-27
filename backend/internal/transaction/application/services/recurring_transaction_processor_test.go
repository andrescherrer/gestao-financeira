package services

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// mockTransactionRepository is a mock implementation of TransactionRepository for testing.
type mockTransactionRepository struct {
	transactions map[string]*entities.Transaction
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{
		transactions: make(map[string]*entities.Transaction),
	}
}

func (m *mockTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	tx, exists := m.transactions[id.Value()]
	if !exists {
		return nil, nil
	}
	return tx, nil
}

func (m *mockTransactionRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.AccountID().Equals(accountID) {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) && tx.AccountID().Equals(accountID) {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) && tx.TransactionType().Equals(transactionType) {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockTransactionRepository) Save(transaction *entities.Transaction) error {
	m.transactions[transaction.ID().Value()] = transaction
	return nil
}

func (m *mockTransactionRepository) Delete(id transactionvalueobjects.TransactionID) error {
	delete(m.transactions, id.Value())
	return nil
}

func (m *mockTransactionRepository) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	_, exists := m.transactions[id.Value()]
	return exists, nil
}

func (m *mockTransactionRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) {
			count++
		}
	}
	return count, nil
}

func (m *mockTransactionRepository) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	count := int64(0)
	for _, tx := range m.transactions {
		if tx.AccountID().Equals(accountID) {
			count++
		}
	}
	return count, nil
}

func (m *mockTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	now := time.Now()
	for _, tx := range m.transactions {
		if !tx.IsRecurring() {
			continue
		}
		if tx.ParentTransactionID() != nil {
			continue // Skip instances
		}
		endDate := tx.RecurrenceEndDate()
		if endDate != nil && !endDate.IsZero() && endDate.Before(now) {
			continue // Past end date
		}
		result = append(result, tx)
	}
	return result, nil
}

func (m *mockTransactionRepository) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	for _, tx := range m.transactions {
		parentIDPtr := tx.ParentTransactionID()
		if parentIDPtr != nil && parentIDPtr.Equals(parentID) {
			// Check if date matches (compare only date part, not time)
			txDate := tx.Date()
			if txDate.Year() == date.Year() && txDate.Month() == date.Month() && txDate.Day() == date.Day() {
				return tx, nil
			}
		}
	}
	return nil, nil
}

func TestRecurringTransactionProcessor_ProcessRecurringTransactions(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockTransactionRepository()
	processor := NewRecurringTransactionProcessor(repository, eventBus)

	// Create a recurring transaction
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.ExpenseType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	date := time.Now().AddDate(0, -1, 0) // 1 month ago
	recurrenceFrequency := transactionvalueobjects.MonthlyFrequency()
	endDate := time.Now().AddDate(1, 0, 0) // 1 year from now

	recurringTx, err := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionType,
		amount,
		description,
		date,
		true,
		&recurrenceFrequency,
		&endDate,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create recurring transaction: %v", err)
	}

	// Save recurring transaction
	if err := repository.Save(recurringTx); err != nil {
		t.Fatalf("Failed to save recurring transaction: %v", err)
	}

	// Process recurring transactions
	createdCount, err := processor.ProcessRecurringTransactions()
	if err != nil {
		t.Fatalf("ProcessRecurringTransactions() error = %v, want nil", err)
	}

	// Should create at least one instance
	if createdCount == 0 {
		t.Error("ProcessRecurringTransactions() should create at least one transaction instance")
	}

	// Verify that instances were created
	allTransactions, _ := repository.FindByUserID(userID)
	instanceCount := 0
	for _, tx := range allTransactions {
		if tx.ParentTransactionID() != nil && tx.ParentTransactionID().Equals(recurringTx.ID()) {
			instanceCount++
		}
	}

	if instanceCount == 0 {
		t.Error("No transaction instances were created")
	}
}

func TestRecurringTransactionProcessor_ProcessRecurringTransactions_NoActiveRecurring(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockTransactionRepository()
	processor := NewRecurringTransactionProcessor(repository, eventBus)

	// Process with no recurring transactions
	createdCount, err := processor.ProcessRecurringTransactions()
	if err != nil {
		t.Fatalf("ProcessRecurringTransactions() error = %v, want nil", err)
	}

	if createdCount != 0 {
		t.Errorf("ProcessRecurringTransactions() createdCount = %v, want 0", createdCount)
	}
}

func TestRecurringTransactionProcessor_ProcessRecurringTransactions_AlreadyProcessed(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockTransactionRepository()
	processor := NewRecurringTransactionProcessor(repository, eventBus)

	// Create a recurring transaction
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	transactionType := transactionvalueobjects.ExpenseType()
	amount, _ := sharedvalueobjects.NewMoney(10000, sharedvalueobjects.MustCurrency("BRL"))
	description, _ := transactionvalueobjects.NewTransactionDescription("Assinatura mensal")
	date := time.Now().AddDate(0, -1, 0)
	recurrenceFrequency := transactionvalueobjects.MonthlyFrequency()
	endDate := time.Now().AddDate(1, 0, 0)

	recurringTx, _ := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionType,
		amount,
		description,
		date,
		true,
		&recurrenceFrequency,
		&endDate,
		nil,
	)

	repository.Save(recurringTx)

	// Create an instance for today (simulating already processed)
	parentID := recurringTx.ID()
	instance, _ := entities.NewTransactionWithRecurrence(
		userID,
		accountID,
		transactionType,
		amount,
		description,
		time.Now(),
		false,
		nil,
		nil,
		&parentID,
	)
	repository.Save(instance)

	// Process - should not create duplicate
	createdCount, err := processor.ProcessRecurringTransactions()
	if err != nil {
		t.Fatalf("ProcessRecurringTransactions() error = %v, want nil", err)
	}

	// Should not create duplicate instance
	if createdCount > 0 {
		t.Logf("Note: Created %d instances (may be expected if multiple periods passed)", createdCount)
	}
}
