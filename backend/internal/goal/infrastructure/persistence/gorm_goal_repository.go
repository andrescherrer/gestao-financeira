package persistence

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/goal/domain/entities"
	"gestao-financeira/backend/internal/goal/domain/repositories"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"

	"gorm.io/gorm"
)

// GormGoalRepository implements GoalRepository using GORM.
type GormGoalRepository struct {
	db *gorm.DB
}

// NewGormGoalRepository creates a new GORM goal repository.
func NewGormGoalRepository(db *gorm.DB) repositories.GoalRepository {
	return &GormGoalRepository{db: db}
}

// FindByID finds a goal by its ID.
func (r *GormGoalRepository) FindByID(id goalvalueobjects.GoalID) (*entities.Goal, error) {
	var model GoalModel
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find goal by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all goals for a given user.
func (r *GormGoalRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Goal, error) {
	var models []GoalModel
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find goals by user ID: %w", err)
	}

	goals := make([]*entities.Goal, 0, len(models))
	for _, model := range models {
		goal, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert goal model to domain: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

// FindByStatus finds all goals for a given user filtered by status.
func (r *GormGoalRepository) FindByStatus(userID identityvalueobjects.UserID, status goalvalueobjects.GoalStatus) ([]*entities.Goal, error) {
	var models []GoalModel
	if err := r.db.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID.Value(), status.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find goals by status: %w", err)
	}

	goals := make([]*entities.Goal, 0, len(models))
	for _, model := range models {
		goal, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert goal model to domain: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

// Save saves or updates a goal.
func (r *GormGoalRepository) Save(goal *entities.Goal) error {
	model := r.toModel(goal)

	// Check if goal exists
	var existing GoalModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new goal
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create goal: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check goal existence: %w", err)
	} else {
		// Update existing goal
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update goal: %w", err)
		}
	}

	return nil
}

// Delete deletes a goal by its ID (soft delete).
func (r *GormGoalRepository) Delete(id goalvalueobjects.GoalID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&GoalModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}
	return nil
}

// Exists checks if a goal with the given ID already exists.
func (r *GormGoalRepository) Exists(id goalvalueobjects.GoalID) (bool, error) {
	var count int64
	if err := r.db.Model(&GoalModel{}).Where("id = ? AND deleted_at IS NULL", id.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check goal existence: %w", err)
	}
	return count > 0, nil
}

// Count returns the total number of goals for a given user.
func (r *GormGoalRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&GoalModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count goals: %w", err)
	}
	return count, nil
}

// FindByUserIDWithPagination finds goals for a given user with pagination.
func (r *GormGoalRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	status string,
	offset, limit int,
) ([]*entities.Goal, int64, error) {
	var models []GoalModel
	var total int64

	query := r.db.Model(&GoalModel{}).Where("user_id = ? AND deleted_at IS NULL", userID.Value())

	// Apply context filter if provided
	if context != "" {
		query = query.Where("context = ?", context)
	}

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count goals: %w", err)
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find goals: %w", err)
	}

	goals := make([]*entities.Goal, 0, len(models))
	for _, model := range models {
		goal, err := r.toDomain(&model)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert goal model to domain: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, total, nil
}

// toDomain converts a persistence model to a domain entity.
func (r *GormGoalRepository) toDomain(model *GoalModel) (*entities.Goal, error) {
	// Create value objects
	goalID, err := goalvalueobjects.NewGoalID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid goal ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	name, err := goalvalueobjects.NewGoalName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid goal name: %w", err)
	}

	targetCurrency, err := sharedvalueobjects.NewCurrency(model.TargetCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid target currency: %w", err)
	}

	targetAmount, err := sharedvalueobjects.NewMoney(model.TargetAmount, targetCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid target amount: %w", err)
	}

	currentCurrency, err := sharedvalueobjects.NewCurrency(model.CurrentCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid current currency: %w", err)
	}

	currentAmount, err := sharedvalueobjects.NewMoney(model.CurrentAmount, currentCurrency)
	if err != nil {
		return nil, fmt.Errorf("invalid current amount: %w", err)
	}

	context, err := sharedvalueobjects.NewAccountContext(model.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid context: %w", err)
	}

	status, err := goalvalueobjects.NewGoalStatus(model.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	return entities.GoalFromPersistence(
		goalID,
		userID,
		name,
		targetAmount,
		currentAmount,
		model.Deadline,
		context,
		status,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

// toModel converts a domain entity to a persistence model.
func (r *GormGoalRepository) toModel(goal *entities.Goal) *GoalModel {
	targetAmount := goal.TargetAmount()
	currentAmount := goal.CurrentAmount()

	model := &GoalModel{
		ID:              goal.ID().Value(),
		UserID:          goal.UserID().Value(),
		Name:            goal.Name().Name(),
		TargetAmount:    targetAmount.Amount(),
		TargetCurrency:  targetAmount.Currency().Code(),
		CurrentAmount:   currentAmount.Amount(),
		CurrentCurrency: currentAmount.Currency().Code(),
		Deadline:        goal.Deadline(),
		Context:         goal.Context().Value(),
		Status:          goal.Status().Value(),
		CreatedAt:       goal.CreatedAt(),
		UpdatedAt:       goal.UpdatedAt(),
	}

	return model
}
