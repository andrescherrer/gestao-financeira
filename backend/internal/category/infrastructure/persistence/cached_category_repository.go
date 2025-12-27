package persistence

import (
	"fmt"
	"time"

	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/pkg/cache"
)

// CachedCategoryRepository wraps a CategoryRepository with caching.
type CachedCategoryRepository struct {
	repository repositories.CategoryRepository
	cache      *cache.CacheService
	ttl        time.Duration
}

// NewCachedCategoryRepository creates a new cached category repository.
func NewCachedCategoryRepository(repository repositories.CategoryRepository, cacheService *cache.CacheService, ttl time.Duration) repositories.CategoryRepository {
	if cacheService == nil {
		return repository // Return original repository if cache is not available
	}
	return &CachedCategoryRepository{
		repository: repository,
		cache:      cacheService,
		ttl:        ttl,
	}
}

// cacheKey generates a cache key for category operations.
func (r *CachedCategoryRepository) cacheKey(operation string, params ...string) string {
	key := fmt.Sprintf("category:%s", operation)
	for _, param := range params {
		key += fmt.Sprintf(":%s", param)
	}
	return key
}

// invalidateUserCache invalidates all cache entries for a user.
func (r *CachedCategoryRepository) invalidateUserCache(userID string) {
	if r.cache == nil {
		return
	}
	pattern := fmt.Sprintf("category:*:%s", userID)
	_ = r.cache.DeletePattern(pattern) // Ignore errors
}

// FindByID finds a category by its ID with caching.
func (r *CachedCategoryRepository) FindByID(id valueobjects.CategoryID) (*entities.Category, error) {
	cacheKey := r.cacheKey("find_by_id", id.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData *cachedCategoryData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found && cachedData != nil {
			category, err := cacheDataToCategory(cachedData)
			if err == nil && category != nil {
				return category, nil
			}
		}
	}

	// Get from repository
	category, err := r.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result (even if nil)
	if r.cache != nil && category != nil {
		cacheData := categoryToCacheData(category)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return category, nil
}

// FindByUserID finds all categories for a given user with caching.
func (r *CachedCategoryRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Category, error) {
	cacheKey := r.cacheKey("find_by_user_id", userID.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData []*cachedCategoryData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found {
			categories, err := cacheDataToCategories(cachedData)
			if err == nil {
				return categories, nil
			}
		}
	}

	// Get from repository
	categories, err := r.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		cacheData := categoriesToCacheData(categories)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return categories, nil
}

// FindByUserIDAndActive finds all active categories for a given user with caching.
func (r *CachedCategoryRepository) FindByUserIDAndActive(userID identityvalueobjects.UserID, isActive bool) ([]*entities.Category, error) {
	cacheKey := r.cacheKey("find_by_user_id_active", userID.Value(), fmt.Sprintf("%t", isActive))

	// Try to get from cache
	if r.cache != nil {
		var cachedData []*cachedCategoryData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found {
			categories, err := cacheDataToCategories(cachedData)
			if err == nil {
				return categories, nil
			}
		}
	}

	// Get from repository
	categories, err := r.repository.FindByUserIDAndActive(userID, isActive)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		cacheData := categoriesToCacheData(categories)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return categories, nil
}

// FindByUserIDAndSlug finds a category by user ID and slug with caching.
func (r *CachedCategoryRepository) FindByUserIDAndSlug(userID identityvalueobjects.UserID, slug valueobjects.CategorySlug) (*entities.Category, error) {
	cacheKey := r.cacheKey("find_by_user_id_slug", userID.Value(), slug.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData *cachedCategoryData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found && cachedData != nil {
			category, err := cacheDataToCategory(cachedData)
			if err == nil && category != nil {
				return category, nil
			}
		}
	}

	// Get from repository
	category, err := r.repository.FindByUserIDAndSlug(userID, slug)
	if err != nil {
		return nil, err
	}

	// Cache the result (even if nil)
	if r.cache != nil && category != nil {
		cacheData := categoryToCacheData(category)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return category, nil
}

// Save saves or updates a category and invalidates cache.
func (r *CachedCategoryRepository) Save(category *entities.Category) error {
	// Save to repository
	err := r.repository.Save(category)
	if err != nil {
		return err
	}

	// Invalidate cache for this category and user
	if category != nil {
		categoryIDKey := r.cacheKey("find_by_id", category.ID().Value())
		_ = r.cache.Delete(categoryIDKey) // Ignore errors

		slugKey := r.cacheKey("find_by_user_id_slug", category.UserID().Value(), category.Slug().Value())
		_ = r.cache.Delete(slugKey) // Ignore errors

		r.invalidateUserCache(category.UserID().Value())
	}

	return nil
}

// Delete deletes a category and invalidates cache.
func (r *CachedCategoryRepository) Delete(id valueobjects.CategoryID) error {
	// Get category first to know user ID for cache invalidation
	category, _ := r.repository.FindByID(id)

	// Delete from repository
	err := r.repository.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if category != nil {
		categoryIDKey := r.cacheKey("find_by_id", id.Value())
		_ = r.cache.Delete(categoryIDKey) // Ignore errors

		slugKey := r.cacheKey("find_by_user_id_slug", category.UserID().Value(), category.Slug().Value())
		_ = r.cache.Delete(slugKey) // Ignore errors

		r.invalidateUserCache(category.UserID().Value())
	}

	return nil
}

// Exists checks if a category exists (no caching for boolean checks).
func (r *CachedCategoryRepository) Exists(id valueobjects.CategoryID) (bool, error) {
	return r.repository.Exists(id)
}

// Count returns the total number of categories for a user (no caching for counts).
func (r *CachedCategoryRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	return r.repository.Count(userID)
}
