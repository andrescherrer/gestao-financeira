package persistence

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/budget/domain/entities"
	"gestao-financeira/backend/internal/budget/domain/repositories"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"

	"gorm.io/gorm"
)

// GormBudgetRepository implements BudgetRepository using GORM.
type GormBudgetRepository struct {
	db *gorm.DB
}

// NewGormBudgetRepository creates a new GORM budget repository.
func NewGormBudgetRepository(db *gorm.DB) repositories.BudgetRepository {
	return &GormBudgetRepository{db: db}
}

// FindByID finds a budget by its ID.
func (r *GormBudgetRepository) FindByID(id valueobjects.BudgetID) (*entities.Budget, error) {
	var model BudgetModel
	if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find budget by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all budgets for a given user.
// Optimized to exclude soft-deleted records.
func (r *GormBudgetRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Budget, error) {
	var models []BudgetModel
	// Exclude soft-deleted records for better performance
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find budgets by user ID: %w", err)
	}

	budgets := make([]*entities.Budget, 0, len(models))
	for _, model := range models {
		budget, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert budget model to domain: %w", err)
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// FindByCategoryAndPeriod finds a budget by category ID and period.
func (r *GormBudgetRepository) FindByCategoryAndPeriod(categoryID categoryvalueobjects.CategoryID, period valueobjects.BudgetPeriod) (*entities.Budget, error) {
	var model BudgetModel
	query := r.db.Where("category_id = ? AND period_type = ? AND year = ?", categoryID.Value(), string(period.PeriodType()), period.Year())

	if period.IsMonthly() {
		query = query.Where("month = ?", *period.Month())
	} else {
		query = query.Where("month IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find budget by category and period: %w", err)
	}

	return r.toDomain(&model)
}

// FindByPeriod finds all budgets for a given user and period.
func (r *GormBudgetRepository) FindByPeriod(userID identityvalueobjects.UserID, period valueobjects.BudgetPeriod) ([]*entities.Budget, error) {
	var models []BudgetModel
	query := r.db.Where("user_id = ? AND period_type = ? AND year = ?", userID.Value(), string(period.PeriodType()), period.Year())

	if period.IsMonthly() {
		query = query.Where("month = ?", *period.Month())
	} else {
		query = query.Where("month IS NULL")
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find budgets by period: %w", err)
	}

	budgets := make([]*entities.Budget, 0, len(models))
	for _, model := range models {
		budget, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert budget model to domain: %w", err)
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// Save saves or updates a budget.
func (r *GormBudgetRepository) Save(budget *entities.Budget) error {
	model := r.toModel(budget)

	// Check if budget exists
	var existing BudgetModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new budget
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create budget: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check budget existence: %w", err)
	} else {
		// Update existing budget
		if err := r.db.Model(&BudgetModel{}).Where("id = ?", model.ID).
			Select("amount", "currency", "period_type", "year", "month", "context", "is_active", "updated_at").
			Updates(map[string]interface{}{
				"amount":      model.Amount,
				"currency":    model.Currency,
				"period_type": model.PeriodType,
				"year":        model.Year,
				"month":       model.Month,
				"context":     model.Context,
				"is_active":   model.IsActive,
				"updated_at":  model.UpdatedAt,
			}).Error; err != nil {
			return fmt.Errorf("failed to update budget: %w", err)
		}
	}

	return nil
}

// Delete deletes a budget by its ID (soft delete).
func (r *GormBudgetRepository) Delete(id valueobjects.BudgetID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&BudgetModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	return nil
}

// Restore restores a soft-deleted budget by setting deleted_at to NULL.
func (r *GormBudgetRepository) Restore(id valueobjects.BudgetID) error {
	if err := r.db.Unscoped().
		Model(&BudgetModel{}).
		Where("id = ?", id.Value()).
		Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("failed to restore budget: %w", err)
	}

	return nil
}

// PermanentDelete permanently deletes a budget (hard delete).
func (r *GormBudgetRepository) PermanentDelete(id valueobjects.BudgetID) error {
	if err := r.db.Unscoped().
		Where("id = ?", id.Value()).
		Delete(&BudgetModel{}).Error; err != nil {
		return fmt.Errorf("failed to permanently delete budget: %w", err)
	}

	return nil
}

// Exists checks if a budget with the given ID already exists.
// Optimized to use LIMIT 1 for better performance.
func (r *GormBudgetRepository) Exists(id valueobjects.BudgetID) (bool, error) {
	var count int64
	if err := r.db.Model(&BudgetModel{}).
		Where("id = ?", id.Value()).
		Limit(1).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check budget existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of budgets for a given user.
func (r *GormBudgetRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&BudgetModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count budgets: %w", err)
	}

	return count, nil
}

// toDomain converts a BudgetModel to a Budget domain entity.
func (r *GormBudgetRepository) toDomain(model *BudgetModel) (*entities.Budget, error) {
	budgetID, err := valueobjects.NewBudgetID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid budget ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	categoryID, err := categoryvalueobjects.NewCategoryID(model.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	currency, err := sharedvalueobjects.NewCurrency(model.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid currency: %w", err)
	}

	amount, err := sharedvalueobjects.NewMoney(model.Amount, currency)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	var period valueobjects.BudgetPeriod
	if model.PeriodType == "MONTHLY" {
		if model.Month == nil {
			return nil, fmt.Errorf("month is required for monthly periods")
		}
		period, err = valueobjects.NewMonthlyBudgetPeriod(model.Year, *model.Month)
	} else {
		period, err = valueobjects.NewYearlyBudgetPeriod(model.Year)
	}
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	context, err := sharedvalueobjects.NewAccountContext(model.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid context: %w", err)
	}

	return entities.BudgetFromPersistence(
		budgetID,
		userID,
		categoryID,
		amount,
		period,
		context,
		model.CreatedAt,
		model.UpdatedAt,
		model.IsActive,
	)
}

// toModel converts a Budget domain entity to a BudgetModel.
func (r *GormBudgetRepository) toModel(budget *entities.Budget) *BudgetModel {
	amount := budget.Amount()
	period := budget.Period()

	model := &BudgetModel{
		ID:         budget.ID().Value(),
		UserID:     budget.UserID().Value(),
		CategoryID: budget.CategoryID().Value(),
		Amount:     amount.Amount(),
		Currency:   amount.Currency().Code(),
		PeriodType: string(period.PeriodType()),
		Year:       period.Year(),
		Month:      period.Month(),
		Context:    budget.Context().Value(),
		IsActive:   budget.IsActive(),
		CreatedAt:  budget.CreatedAt(),
		UpdatedAt:  budget.UpdatedAt(),
	}

	return model
}
