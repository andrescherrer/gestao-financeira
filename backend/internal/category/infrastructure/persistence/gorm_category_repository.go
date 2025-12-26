package persistence

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"

	"gorm.io/gorm"
)

// GormCategoryRepository implements CategoryRepository using GORM.
type GormCategoryRepository struct {
	db *gorm.DB
}

// NewGormCategoryRepository creates a new GORM category repository.
func NewGormCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &GormCategoryRepository{db: db}
}

// FindByID finds a category by its ID.
func (r *GormCategoryRepository) FindByID(id valueobjects.CategoryID) (*entities.Category, error) {
	var model CategoryModel
	if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUserID finds all categories for a given user.
func (r *GormCategoryRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Category, error) {
	var models []CategoryModel
	if err := r.db.Where("user_id = ?", userID.Value()).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find categories by user ID: %w", err)
	}

	categories := make([]*entities.Category, 0, len(models))
	for _, model := range models {
		category, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert category model to domain: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindByUserIDAndActive finds all active categories for a given user.
func (r *GormCategoryRepository) FindByUserIDAndActive(userID identityvalueobjects.UserID, isActive bool) ([]*entities.Category, error) {
	var models []CategoryModel
	if err := r.db.Where("user_id = ? AND is_active = ?", userID.Value(), isActive).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find categories by user ID and active status: %w", err)
	}

	categories := make([]*entities.Category, 0, len(models))
	for _, model := range models {
		category, err := r.toDomain(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert category model to domain: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Save saves or updates a category.
func (r *GormCategoryRepository) Save(category *entities.Category) error {
	model := r.toModel(category)

	// Check if category exists
	var existing CategoryModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new category
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create category: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check category existence: %w", err)
	} else {
		// Update existing category
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}
	}

	return nil
}

// Delete deletes a category by its ID.
func (r *GormCategoryRepository) Delete(id valueobjects.CategoryID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&CategoryModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// Exists checks if a category with the given ID already exists.
func (r *GormCategoryRepository) Exists(id valueobjects.CategoryID) (bool, error) {
	var count int64
	if err := r.db.Model(&CategoryModel{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of categories for a given user.
func (r *GormCategoryRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	var count int64
	if err := r.db.Model(&CategoryModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}

// toDomain converts a CategoryModel to a Category domain entity.
func (r *GormCategoryRepository) toDomain(model *CategoryModel) (*entities.Category, error) {
	categoryID, err := valueobjects.NewCategoryID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	userID, err := identityvalueobjects.NewUserID(model.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	categoryName, err := valueobjects.NewCategoryName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid category name: %w", err)
	}

	// If slug is empty in database, generate it from name (for backward compatibility)
	var categorySlug valueobjects.CategorySlug
	if model.Slug == "" {
		categorySlug = valueobjects.GenerateSlugFromName(model.Name)
	} else {
		categorySlug, err = valueobjects.NewCategorySlug(model.Slug)
		if err != nil {
			// If slug is invalid, generate a new one
			categorySlug = valueobjects.GenerateSlugFromName(model.Name)
		}
	}

	return entities.CategoryFromPersistence(
		categoryID,
		userID,
		categoryName,
		categorySlug,
		model.Description,
		model.CreatedAt,
		model.UpdatedAt,
		model.IsActive,
	)
}

// toModel converts a Category domain entity to a CategoryModel.
func (r *GormCategoryRepository) toModel(category *entities.Category) *CategoryModel {
	return &CategoryModel{
		ID:          category.ID().Value(),
		UserID:      category.UserID().Value(),
		Name:        category.Name().Value(),
		Slug:        category.Slug().Value(),
		Description: category.Description(),
		IsActive:    category.IsActive(),
		CreatedAt:   category.CreatedAt(),
		UpdatedAt:   category.UpdatedAt(),
	}
}
