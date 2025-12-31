package persistence

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"

	"gorm.io/gorm"
)

// GormAccountRepository implements AccountRepository using GORM.
type GormAccountRepository struct {
	db *gorm.DB
}

// NewGormAccountRepository creates a new GORM account repository.
func NewGormAccountRepository(db *gorm.DB) repositories.AccountRepository {
	return &GormAccountRepository{db: db}
}

// FindByID finds an account by its ID.
func (r *GormAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	var model AccountModel
	if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find account by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all accounts for a given user.
// Optimized to use index idx_accounts_user_active.
func (r *GormAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	var models []AccountModel
	// Use index idx_accounts_user_active for optimal performance
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by user ID: %w", err)
	}

	accounts := make([]*entities.Account, 0, len(models))
	for _, model := range models {
		account, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert account model to domain: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// FindByUserIDAndContext finds all accounts for a given user filtered by context.
func (r *GormAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	var models []AccountModel
	if err := r.db.Where("user_id = ? AND context = ?", userID.Value(), context.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by user ID and context: %w", err)
	}

	accounts := make([]*entities.Account, 0, len(models))
	for _, model := range models {
		account, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert account model to domain: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// Save saves or updates an account.
func (r *GormAccountRepository) Save(account *entities.Account) error {
	model := r.toModel(account)

	// Check if account exists
	var existing AccountModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new account
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check account existence: %w", err)
	} else {
		// Update existing account
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update account: %w", err)
		}
	}

	return nil
}

// Delete deletes an account by its ID (soft delete).
func (r *GormAccountRepository) Delete(id valueobjects.AccountID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&AccountModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// Restore restores a soft-deleted account by setting deleted_at to NULL.
func (r *GormAccountRepository) Restore(id valueobjects.AccountID) error {
	if err := r.db.Unscoped().
		Model(&AccountModel{}).
		Where("id = ?", id.Value()).
		Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("failed to restore account: %w", err)
	}

	return nil
}

// PermanentDelete permanently deletes an account (hard delete).
func (r *GormAccountRepository) PermanentDelete(id valueobjects.AccountID) error {
	if err := r.db.Unscoped().
		Where("id = ?", id.Value()).
		Delete(&AccountModel{}).Error; err != nil {
		return fmt.Errorf("failed to permanently delete account: %w", err)
	}

	return nil
}

// Exists checks if an account with the given ID already exists.
// Optimized to use LIMIT 1 for better performance.
func (r *GormAccountRepository) Exists(id valueobjects.AccountID) (bool, error) {
	var count int64
	if err := r.db.Model(&AccountModel{}).
		Where("id = ?", id.Value()).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of accounts for a given user.
func (r *GormAccountRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&AccountModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	return count, nil
}

// FindByUserIDWithPagination finds accounts for a given user with pagination.
func (r *GormAccountRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	offset, limit int,
) ([]*entities.Account, int64, error) {
	var models []AccountModel
	var total int64

	query := r.db.Model(&AccountModel{}).Where("user_id = ?", userID.Value())

	// Filter by context if provided
	if context != "" {
		query = query.Where("context = ?", context)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	// Get paginated results
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find accounts: %w", err)
	}

	accounts := make([]*entities.Account, 0, len(models))
	for _, model := range models {
		account, err := r.toDomain(&model)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert account model to domain: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, total, nil
}

// toDomain converts an AccountModel to an Account entity.
func (r *GormAccountRepository) toDomain(model *AccountModel) (*entities.Account, error) {
	// Convert value objects
	accountID, err := valueobjects.NewAccountID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	accountName, err := valueobjects.NewAccountName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid account name: %w", err)
	}

	accountType, err := valueobjects.NewAccountType(model.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid account type: %w", err)
	}

	currency, err := sharedvalueobjects.NewCurrency(model.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	balance, err := sharedvalueobjects.NewMoney(model.Balance, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid balance: %w", err)
	}

	context, err := sharedvalueobjects.NewAccountContext(model.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid account context: %w", err)
	}

	// Reconstruct account entity from persisted data
	return entities.AccountFromPersistence(
		accountID,
		userID,
		accountName,
		accountType,
		balance,
		context,
		model.CreatedAt,
		model.UpdatedAt,
		model.IsActive,
	)
}

// toModel converts an Account entity to an AccountModel.
func (r *GormAccountRepository) toModel(account *entities.Account) *AccountModel {
	balance := account.Balance()

	return &AccountModel{
		ID:        account.ID().Value(),
		UserID:    account.UserID().Value(),
		Name:      account.Name().Value(),
		Type:      account.AccountType().Value(),
		Balance:   balance.Amount(), // Amount in cents
		Currency:  balance.Currency().Code(),
		Context:   account.Context().Value(),
		IsActive:  account.IsActive(),
		CreatedAt: account.CreatedAt(),
		UpdatedAt: account.UpdatedAt(),
	}
}
