package persistence

import (
	"errors"
	"fmt"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	"gestao-financeira/backend/internal/transaction/domain/repositories"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"

	"gorm.io/gorm"
)

// GormTransactionRepository implements TransactionRepository using GORM.
type GormTransactionRepository struct {
	db *gorm.DB
}

// NewGormTransactionRepository creates a new GORM transaction repository.
func NewGormTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &GormTransactionRepository{db: db}
}

// FindByID finds a transaction by its ID.
func (r *GormTransactionRepository) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	var model TransactionModel
	if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all transactions for a given user.
func (r *GormTransactionRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var models []TransactionModel
	if err := r.db.Where("user_id = ?", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find transactions by user ID: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))
	for _, model := range models {
		transaction, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction model to domain: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindByAccountID finds all transactions for a given account.
func (r *GormTransactionRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var models []TransactionModel
	if err := r.db.Where("account_id = ?", accountID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find transactions by account ID: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))
	for _, model := range models {
		transaction, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction model to domain: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindByUserIDAndAccountID finds all transactions for a given user and account.
func (r *GormTransactionRepository) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var models []TransactionModel
	if err := r.db.Where("user_id = ? AND account_id = ?", userID.Value(), accountID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find transactions by user ID and account ID: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))
	for _, model := range models {
		transaction, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction model to domain: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindByUserIDAndType finds all transactions for a given user filtered by type.
func (r *GormTransactionRepository) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	var models []TransactionModel
	if err := r.db.Where("user_id = ? AND type = ?", userID.Value(), transactionType.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find transactions by user ID and type: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))
	for _, model := range models {
		transaction, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction model to domain: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// Save saves or updates a transaction.
func (r *GormTransactionRepository) Save(transaction *entities.Transaction) error {
	model := r.toModel(transaction)

	// Check if transaction exists
	var existing TransactionModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new transaction
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check transaction existence: %w", err)
	} else {
		// Update existing transaction
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}
	}

	return nil
}

// Delete deletes a transaction by its ID.
func (r *GormTransactionRepository) Delete(id transactionvalueobjects.TransactionID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&TransactionModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}

// Exists checks if a transaction with the given ID already exists.
func (r *GormTransactionRepository) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	var count int64
	if err := r.db.Model(&TransactionModel{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check transaction existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of transactions for a given user.
func (r *GormTransactionRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&TransactionModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}

// CountByAccountID returns the total number of transactions for a given account.
func (r *GormTransactionRepository) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	var count int64
	if err := r.db.Model(&TransactionModel{}).Where("account_id = ?", accountID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count transactions by account ID: %w", err)
	}

	return count, nil
}

// FindActiveRecurringTransactions finds all active recurring transactions that need to be processed.
func (r *GormTransactionRepository) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	var models []TransactionModel
	today := time.Now().Format("2006-01-02")

	// Find transactions where:
	// - is_recurring = true
	// - (recurrence_end_date IS NULL OR recurrence_end_date >= today)
	// - parent_transaction_id IS NULL (only parent transactions, not generated instances)
	query := r.db.Where("is_recurring = ? AND parent_transaction_id IS NULL", true).
		Where("recurrence_end_date IS NULL OR recurrence_end_date >= ?", today)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find active recurring transactions: %w", err)
	}

	transactions := make([]*entities.Transaction, 0, len(models))
	for _, model := range models {
		transaction, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction model to domain: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// FindByParentIDAndDate finds a transaction instance by parent transaction ID and date.
func (r *GormTransactionRepository) FindByParentIDAndDate(
	parentID transactionvalueobjects.TransactionID,
	date time.Time,
) (*entities.Transaction, error) {
	var model TransactionModel
	dateStr := date.Format("2006-01-02")

	if err := r.db.Where("parent_transaction_id = ? AND date = ?", parentID.Value(), dateStr).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction by parent ID and date: %w", err)
	}

	return r.toDomain(&model)
}

// toDomain converts a TransactionModel to a Transaction entity.
func (r *GormTransactionRepository) toDomain(model *TransactionModel) (*entities.Transaction, error) {
	// Convert value objects
	transactionID, err := transactionvalueobjects.NewTransactionID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	accountID, err := accountvalueobjects.NewAccountID(model.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	transactionType, err := transactionvalueobjects.NewTransactionType(model.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction type: %w", err)
	}

	currency, err := valueobjects.NewCurrency(model.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	amount, err := valueobjects.NewMoney(model.Amount, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	description, err := transactionvalueobjects.NewTransactionDescription(model.Description)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction description: %w", err)
	}

	// Convert recurrence fields
	var recurrenceFrequency *transactionvalueobjects.RecurrenceFrequency
	if model.IsRecurring && model.RecurrenceFrequency != nil {
		rf, err := transactionvalueobjects.NewRecurrenceFrequency(*model.RecurrenceFrequency)
		if err != nil {
			return nil, fmt.Errorf("invalid recurrence frequency: %w", err)
		}
		recurrenceFrequency = &rf
	}

	var parentTransactionID *transactionvalueobjects.TransactionID
	if model.ParentTransactionID != nil {
		pid, err := transactionvalueobjects.NewTransactionID(*model.ParentTransactionID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent transaction ID: %w", err)
		}
		parentTransactionID = &pid
	}

	// Reconstruct transaction entity from persisted data
	return entities.TransactionFromPersistenceWithRecurrence(
		transactionID,
		userID,
		accountID,
		transactionType,
		amount,
		description,
		model.Date,
		model.CreatedAt,
		model.UpdatedAt,
		model.IsRecurring,
		recurrenceFrequency,
		model.RecurrenceEndDate,
		parentTransactionID,
	)
}

// toModel converts a Transaction entity to a TransactionModel.
func (r *GormTransactionRepository) toModel(transaction *entities.Transaction) *TransactionModel {
	amount := transaction.Amount()

	var recurrenceFrequency *string
	if transaction.RecurrenceFrequency() != nil {
		freq := transaction.RecurrenceFrequency().Value()
		recurrenceFrequency = &freq
	}

	var parentTransactionID *string
	if transaction.ParentTransactionID() != nil {
		pid := transaction.ParentTransactionID().Value()
		parentTransactionID = &pid
	}

	return &TransactionModel{
		ID:                  transaction.ID().Value(),
		UserID:              transaction.UserID().Value(),
		AccountID:           transaction.AccountID().Value(),
		Type:                transaction.TransactionType().Value(),
		Amount:              amount.Amount(), // Amount in cents
		Currency:            amount.Currency().Code(),
		Description:         transaction.Description().Value(),
		Date:                transaction.Date(),
		IsRecurring:         transaction.IsRecurring(),
		RecurrenceFrequency: recurrenceFrequency,
		RecurrenceEndDate:   transaction.RecurrenceEndDate(),
		ParentTransactionID: parentTransactionID,
		CreatedAt:           transaction.CreatedAt(),
		UpdatedAt:           transaction.UpdatedAt(),
	}
}
