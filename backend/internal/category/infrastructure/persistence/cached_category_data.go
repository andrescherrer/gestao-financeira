package persistence

import (
	"time"

	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// cachedCategoryData is a serializable representation of Category for caching.
type cachedCategoryData struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// categoryToCacheData converts a Category entity to cachedCategoryData.
func categoryToCacheData(category *entities.Category) *cachedCategoryData {
	if category == nil {
		return nil
	}
	return &cachedCategoryData{
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

// cacheDataToCategory converts cachedCategoryData back to Category entity.
func cacheDataToCategory(data *cachedCategoryData) (*entities.Category, error) {
	if data == nil {
		return nil, nil
	}

	categoryID, err := valueobjects.NewCategoryID(data.ID)
	if err != nil {
		return nil, err
	}

	userID, err := identityvalueobjects.NewUserID(data.UserID)
	if err != nil {
		return nil, err
	}

	categoryName, err := valueobjects.NewCategoryName(data.Name)
	if err != nil {
		return nil, err
	}

	categorySlug, err := valueobjects.NewCategorySlug(data.Slug)
	if err != nil {
		return nil, err
	}

	return entities.CategoryFromPersistence(
		categoryID,
		userID,
		categoryName,
		categorySlug,
		data.Description,
		data.CreatedAt,
		data.UpdatedAt,
		data.IsActive,
	)
}

// categoriesToCacheData converts a slice of Category entities to cachedCategoryData.
func categoriesToCacheData(categories []*entities.Category) []*cachedCategoryData {
	if categories == nil {
		return nil
	}
	result := make([]*cachedCategoryData, 0, len(categories))
	for _, category := range categories {
		result = append(result, categoryToCacheData(category))
	}
	return result
}

// cacheDataToCategories converts cachedCategoryData slice back to Category entities.
func cacheDataToCategories(data []*cachedCategoryData) ([]*entities.Category, error) {
	if data == nil {
		return nil, nil
	}
	result := make([]*entities.Category, 0, len(data))
	for _, item := range data {
		category, err := cacheDataToCategory(item)
		if err != nil {
			return nil, err
		}
		if category != nil {
			result = append(result, category)
		}
	}
	return result, nil
}
