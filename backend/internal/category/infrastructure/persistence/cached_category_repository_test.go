package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/pkg/cache"
)

// mockCategoryRepository is a mock implementation for testing
type mockCategoryRepository struct {
	categories map[string]*entities.Category
	users      map[string][]*entities.Category
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[string]*entities.Category),
		users:      make(map[string][]*entities.Category),
	}
}

func (m *mockCategoryRepository) FindByID(id valueobjects.CategoryID) (*entities.Category, error) {
	return m.categories[id.Value()], nil
}

func (m *mockCategoryRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Category, error) {
	return m.users[userID.Value()], nil
}

func (m *mockCategoryRepository) FindByUserIDAndActive(userID identityvalueobjects.UserID, isActive bool) ([]*entities.Category, error) {
	allCategories := m.users[userID.Value()]
	var filtered []*entities.Category
	for _, cat := range allCategories {
		if cat.IsActive() == isActive {
			filtered = append(filtered, cat)
		}
	}
	return filtered, nil
}

func (m *mockCategoryRepository) FindByUserIDAndSlug(userID identityvalueobjects.UserID, slug valueobjects.CategorySlug) (*entities.Category, error) {
	allCategories := m.users[userID.Value()]
	for _, cat := range allCategories {
		if cat.Slug().Equals(slug) {
			return cat, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryRepository) Save(category *entities.Category) error {
	m.categories[category.ID().Value()] = category
	userID := category.UserID().Value()
	m.users[userID] = append(m.users[userID], category)
	return nil
}

func (m *mockCategoryRepository) Delete(id valueobjects.CategoryID) error {
	delete(m.categories, id.Value())
	return nil
}

func (m *mockCategoryRepository) Exists(id valueobjects.CategoryID) (bool, error) {
	_, exists := m.categories[id.Value()]
	return exists, nil
}
func (m *mockCategoryRepository) FindByUserIDWithPagination(userID identityvalueobjects.UserID, isActive *bool, offset, limit int) ([]*entities.Category, int64, error) {
	all, _ := m.FindByUserID(userID)
	var filtered []*entities.Category
	for _, cat := range all {
		if isActive == nil || cat.IsActive() == *isActive {
			filtered = append(filtered, cat)
		}
	}
	total := int64(len(filtered))
	start := offset
	end := offset + limit
	if start > len(filtered) {
		return []*entities.Category{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (m *mockCategoryRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	return int64(len(m.users[userID.Value()])), nil
}

func TestCachedCategoryRepository_FindByUserID(t *testing.T) {
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	mockRepo := newMockCategoryRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	categoryName, _ := valueobjects.NewCategoryName("Test Category")

	category, err := entities.NewCategory(userID, categoryName, "Test description")
	require.NoError(t, err)
	require.NotNil(t, category)
	mockRepo.users[userID.Value()] = []*entities.Category{category}

	cachedRepo := NewCachedCategoryRepository(mockRepo, cacheService, 1*time.Minute).(*CachedCategoryRepository)

	// First call - should hit repository
	result1, err := cachedRepo.FindByUserID(userID)
	require.NoError(t, err)
	require.Len(t, result1, 1)

	// Second call - should hit cache
	result2, err := cachedRepo.FindByUserID(userID)
	require.NoError(t, err)
	require.Len(t, result2, 1)
}

func TestCachedCategoryRepository_WithoutCache(t *testing.T) {
	// Test that repository works without cache
	mockRepo := newMockCategoryRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	categoryName, _ := valueobjects.NewCategoryName("Test Category")

	category, err := entities.NewCategory(userID, categoryName, "Test description")
	require.NoError(t, err)
	require.NotNil(t, category)
	mockRepo.categories[category.ID().Value()] = category

	// Create cached repository with nil cache (should return original repository)
	cachedRepo := NewCachedCategoryRepository(mockRepo, nil, 1*time.Minute)

	// Should work the same as original repository
	result, err := cachedRepo.FindByID(category.ID())
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, category.ID().Value(), result.ID().Value())
}
