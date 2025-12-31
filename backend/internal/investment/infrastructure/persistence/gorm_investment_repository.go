package persistence

import (
	"errors"
	"fmt"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/domain/entities"
	"gestao-financeira/backend/internal/investment/domain/repositories"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"

	"gorm.io/gorm"
)

// GormInvestmentRepository implements InvestmentRepository using GORM.
type GormInvestmentRepository struct {
	db *gorm.DB
}

// NewGormInvestmentRepository creates a new GORM investment repository.
func NewGormInvestmentRepository(db *gorm.DB) repositories.InvestmentRepository {
	return &GormInvestmentRepository{db: db}
}

// FindByID finds an investment by its ID.
func (r *GormInvestmentRepository) FindByID(id investmentvalueobjects.InvestmentID) (*entities.Investment, error) {
	var model InvestmentModel
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find investment by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all investments for a given user.
func (r *GormInvestmentRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Investment, error) {
	var models []InvestmentModel
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find investments by user ID: %w", err)
	}

	investments := make([]*entities.Investment, 0, len(models))
	for _, model := range models {
		investment, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert investment model to domain: %w", err)
		}
		investments = append(investments, investment)
	}

	return investments, nil
}

// FindByAccountID finds all investments for a given account.
func (r *GormInvestmentRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Investment, error) {
	var models []InvestmentModel
	if err := r.db.Where("account_id = ? AND deleted_at IS NULL", accountID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find investments by account ID: %w", err)
	}

	investments := make([]*entities.Investment, 0, len(models))
	for _, model := range models {
		investment, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert investment model to domain: %w", err)
		}
		investments = append(investments, investment)
	}

	return investments, nil
}

// FindByType finds all investments for a given user filtered by type.
func (r *GormInvestmentRepository) FindByType(userID identityvalueobjects.UserID, investmentType investmentvalueobjects.InvestmentType) ([]*entities.Investment, error) {
	var models []InvestmentModel
	if err := r.db.Where("user_id = ? AND type = ? AND deleted_at IS NULL", userID.Value(), investmentType.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find investments by type: %w", err)
	}

	investments := make([]*entities.Investment, 0, len(models))
	for _, model := range models {
		investment, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert investment model to domain: %w", err)
		}
		investments = append(investments, investment)
	}

	return investments, nil
}

// Save saves or updates an investment.
func (r *GormInvestmentRepository) Save(investment *entities.Investment) error {
	model := r.toModel(investment)

	// Check if investment exists
	var existing InvestmentModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new investment
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create investment: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check investment existence: %w", err)
	} else {
		// Update existing investment
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update investment: %w", err)
		}
	}

	return nil
}

// Delete deletes an investment by its ID (soft delete).
func (r *GormInvestmentRepository) Delete(id investmentvalueobjects.InvestmentID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&InvestmentModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete investment: %w", err)
	}
	return nil
}

// Exists checks if an investment with the given ID already exists.
func (r *GormInvestmentRepository) Exists(id investmentvalueobjects.InvestmentID) (bool, error) {
	var count int64
	if err := r.db.Model(&InvestmentModel{}).Where("id = ? AND deleted_at IS NULL", id.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check investment existence: %w", err)
	}
	return count > 0, nil
}

// Count returns the total number of investments for a given user.
func (r *GormInvestmentRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&InvestmentModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count investments: %w", err)
	}
	return count, nil
}

// FindByUserIDWithPagination finds investments for a given user with pagination.
func (r *GormInvestmentRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	investmentType string,
	offset, limit int,
) ([]*entities.Investment, int64, error) {
	var models []InvestmentModel
	var total int64

	query := r.db.Model(&InvestmentModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value())

	// Apply context filter if provided
	if context != "" {
		query = query.Where("context = ?", context)
	}

	// Apply type filter if provided
	if investmentType != "" {
		query = query.Where("type = ?", investmentType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count investments: %w", err)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find investments: %w", err)
	}

	investments := make([]*entities.Investment, 0, len(models))
	for _, model := range models {
		investment, err := r.toDomain(&model)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert investment model to domain: %w", err)
		}
		investments = append(investments, investment)
	}

	return investments, total, nil
}

// toDomain converts a persistence model to a domain entity.
func (r *GormInvestmentRepository) toDomain(model *InvestmentModel) (*entities.Investment, error) {
	// Create value objects
	investmentID, err := investmentvalueobjects.NewInvestmentID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid investment ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	accountID, err := accountvalueobjects.NewAccountID(model.AccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	investmentType, err := investmentvalueobjects.NewInvestmentType(model.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid investment type: %w", err)
	}

	name, err := investmentvalueobjects.NewInvestmentName(model.Name, model.Ticker)
	if err != nil {
		return nil, fmt.Errorf("invalid investment name: %w", err)
	}

	purchaseCurrency, err := sharedvalueobjects.NewCurrency(model.PurchaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase currency: %w", err)
	}

	purchaseAmount, err := sharedvalueobjects.NewMoney(model.PurchaseAmount, purchaseCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase amount: %w", err)
	}

	currentCurrency, err := sharedvalueobjects.NewCurrency(model.CurrentCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid current currency: %w", err)
	}

	currentValue, err := sharedvalueobjects.NewMoney(model.CurrentValue, currentCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid current value: %w", err)
	}

	context, err := sharedvalueobjects.NewAccountContext(model.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid context: %w", err)
	}

	return entities.InvestmentFromPersistence(
		investmentID,
		userID,
		accountID,
		investmentType,
		name,
		model.PurchaseDate,
		purchaseAmount,
		currentValue,
		model.Quantity,
		context,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

// toModel converts a domain entity to a persistence model.
func (r *GormInvestmentRepository) toModel(investment *entities.Investment) *InvestmentModel {
	purchaseAmount := investment.PurchaseAmount()
	currentValue := investment.CurrentValue()

	model := &InvestmentModel{
		ID:               investment.ID().Value(),
		UserID:           investment.UserID().Value(),
		AccountID:        investment.AccountID().Value(),
		Type:             investment.InvestmentType().Value(),
		Name:             investment.Name().Name(),
		Ticker:           investment.Name().Ticker(),
		PurchaseDate:     investment.PurchaseDate(),
		PurchaseAmount:   purchaseAmount.Amount(),
		PurchaseCurrency: purchaseAmount.Currency().Code(),
		CurrentValue:     currentValue.Amount(),
		CurrentCurrency:  currentValue.Currency().Code(),
		Quantity:         investment.Quantity(),
		Context:          investment.Context().Value(),
		CreatedAt:        investment.CreatedAt(),
		UpdatedAt:        investment.UpdatedAt(),
	}

	return model
}
