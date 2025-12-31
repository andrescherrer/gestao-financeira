package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/pkg/cache"
)

// mockAccountRepository is a mock implementation for testing
type mockAccountRepository struct {
	accounts map[string]*entities.Account
	users    map[string][]*entities.Account
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{
		accounts: make(map[string]*entities.Account),
		users:    make(map[string][]*entities.Account),
	}
}

func (m *mockAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	return m.accounts[id.Value()], nil
}

func (m *mockAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	return m.users[userID.Value()], nil
}

func (m *mockAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	allAccounts := m.users[userID.Value()]
	var filtered []*entities.Account
	for _, acc := range allAccounts {
		if acc.Context().Equals(context) {
			filtered = append(filtered, acc)
		}
	}
	return filtered, nil
}

func (m *mockAccountRepository) Save(account *entities.Account) error {
	m.accounts[account.ID().Value()] = account
	userID := account.UserID().Value()
	m.users[userID] = append(m.users[userID], account)
	return nil
}

func (m *mockAccountRepository) Delete(id valueobjects.AccountID) error {
	delete(m.accounts, id.Value())
	return nil
}

func (m *mockAccountRepository) Exists(id valueobjects.AccountID) (bool, error) {
	_, exists := m.accounts[id.Value()]
	return exists, nil
}
func (m *mockAccountRepository) FindByUserIDWithPagination(userID identityvalueobjects.UserID, context string, offset, limit int) ([]*entities.Account, int64, error) {
	all, _ := m.FindByUserID(userID)
	var filtered []*entities.Account
	for _, acc := range all {
		if context == "" || acc.Context().Value() == context {
			filtered = append(filtered, acc)
		}
	}
	total := int64(len(filtered))
	start := offset
	end := offset + limit
	if start > len(filtered) {
		return []*entities.Account{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (m *mockAccountRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	return int64(len(m.users[userID.Value()])), nil
}

func TestCachedAccountRepository_FindByID(t *testing.T) {
	// Skip if Redis is not available
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	mockRepo := newMockAccountRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountName, _ := valueobjects.NewAccountName("Test Account")
	accountType := valueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency) // R$ 1000.00
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	account, err := entities.NewAccount(userID, accountName, accountType, balance, context)
	require.NoError(t, err)
	require.NotNil(t, account)
	accountID := account.ID()
	mockRepo.accounts[accountID.Value()] = account

	cachedRepo := NewCachedAccountRepository(mockRepo, cacheService, 1*time.Minute).(*CachedAccountRepository)

	// First call - should hit repository
	result1, err := cachedRepo.FindByID(accountID)
	require.NoError(t, err)
	require.NotNil(t, result1)
	assert.Equal(t, accountID.Value(), result1.ID().Value())

	// Second call - should hit cache
	result2, err := cachedRepo.FindByID(accountID)
	require.NoError(t, err)
	require.NotNil(t, result2)
	assert.Equal(t, accountID.Value(), result2.ID().Value())
}

func TestCachedAccountRepository_FindByUserID(t *testing.T) {
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	mockRepo := newMockAccountRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountName, _ := valueobjects.NewAccountName("Test Account")
	accountType := valueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	account, err := entities.NewAccount(userID, accountName, accountType, balance, context)
	require.NoError(t, err)
	require.NotNil(t, account)
	mockRepo.users[userID.Value()] = []*entities.Account{account}

	cachedRepo := NewCachedAccountRepository(mockRepo, cacheService, 1*time.Minute).(*CachedAccountRepository)

	// First call - should hit repository
	result1, err := cachedRepo.FindByUserID(userID)
	require.NoError(t, err)
	require.Len(t, result1, 1)

	// Second call - should hit cache
	result2, err := cachedRepo.FindByUserID(userID)
	require.NoError(t, err)
	require.Len(t, result2, 1)
}

func TestCachedAccountRepository_Save_InvalidatesCache(t *testing.T) {
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	mockRepo := newMockAccountRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountName, _ := valueobjects.NewAccountName("Test Account")
	accountType := valueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	account, err := entities.NewAccount(userID, accountName, accountType, balance, context)
	require.NoError(t, err)
	require.NotNil(t, account)

	cachedRepo := NewCachedAccountRepository(mockRepo, cacheService, 1*time.Minute).(*CachedAccountRepository)

	// Save account
	err = cachedRepo.Save(account)
	require.NoError(t, err)

	// Cache should be invalidated, so next FindByID should hit repository
	result, err := cachedRepo.FindByID(account.ID())
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestCachedAccountRepository_WithoutCache(t *testing.T) {
	// Test that repository works without cache
	mockRepo := newMockAccountRepository()
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountName, _ := valueobjects.NewAccountName("Test Account")
	accountType := valueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	account, err := entities.NewAccount(userID, accountName, accountType, balance, context)
	require.NoError(t, err)
	require.NotNil(t, account)
	mockRepo.accounts[account.ID().Value()] = account

	// Create cached repository with nil cache (should return original repository)
	cachedRepo := NewCachedAccountRepository(mockRepo, nil, 1*time.Minute)

	// Should work the same as original repository
	result, err := cachedRepo.FindByID(account.ID())
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, account.ID().Value(), result.ID().Value())
}
