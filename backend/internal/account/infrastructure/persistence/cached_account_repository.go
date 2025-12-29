package persistence

import (
	"fmt"
	"time"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/repositories"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/pkg/cache"
)

// CachedAccountRepository wraps an AccountRepository with caching.
type CachedAccountRepository struct {
	repository repositories.AccountRepository
	cache      *cache.CacheService
	ttl        time.Duration
}

// NewCachedAccountRepository creates a new cached account repository.
func NewCachedAccountRepository(repository repositories.AccountRepository, cacheService *cache.CacheService, ttl time.Duration) repositories.AccountRepository {
	if cacheService == nil {
		return repository // Return original repository if cache is not available
	}
	return &CachedAccountRepository{
		repository: repository,
		cache:      cacheService,
		ttl:        ttl,
	}
}

// cacheKey generates a cache key for account operations.
func (r *CachedAccountRepository) cacheKey(operation string, params ...string) string {
	key := fmt.Sprintf("account:%s", operation)
	for _, param := range params {
		key += fmt.Sprintf(":%s", param)
	}
	return key
}

// invalidateUserCache invalidates all cache entries for a user.
func (r *CachedAccountRepository) invalidateUserCache(userID string) {
	if r.cache == nil {
		return
	}
	pattern := fmt.Sprintf("account:*:%s", userID)
	_ = r.cache.DeletePattern(pattern) // Ignore errors
}

// FindByID finds an account by its ID with caching.
func (r *CachedAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	cacheKey := r.cacheKey("find_by_id", id.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData *cachedAccountData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found && cachedData != nil {
			account, err := cacheDataToAccount(cachedData)
			if err == nil && account != nil {
				return account, nil
			}
		}
	}

	// Get from repository
	account, err := r.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result (even if nil)
	if r.cache != nil && account != nil {
		cacheData := accountToCacheData(account)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return account, nil
}

// FindByUserID finds all accounts for a given user with caching.
func (r *CachedAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	cacheKey := r.cacheKey("find_by_user_id", userID.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData []*cachedAccountData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found {
			accounts, err := cacheDataToAccounts(cachedData)
			if err == nil {
				return accounts, nil
			}
		}
	}

	// Get from repository
	accounts, err := r.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		cacheData := accountsToCacheData(accounts)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return accounts, nil
}

// FindByUserIDAndContext finds all accounts for a given user filtered by context with caching.
func (r *CachedAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	cacheKey := r.cacheKey("find_by_user_id_context", userID.Value(), context.Value())

	// Try to get from cache
	if r.cache != nil {
		var cachedData []*cachedAccountData
		found, err := r.cache.GetJSON(cacheKey, &cachedData)
		if err == nil && found {
			accounts, err := cacheDataToAccounts(cachedData)
			if err == nil {
				return accounts, nil
			}
		}
	}

	// Get from repository
	accounts, err := r.repository.FindByUserIDAndContext(userID, context)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		cacheData := accountsToCacheData(accounts)
		_ = r.cache.SetJSON(cacheKey, cacheData, r.ttl) // Ignore cache errors
	}

	return accounts, nil
}

// Save saves or updates an account and invalidates cache.
func (r *CachedAccountRepository) Save(account *entities.Account) error {
	// Save to repository
	err := r.repository.Save(account)
	if err != nil {
		return err
	}

	// Invalidate cache for this account and user
	if account != nil {
		accountIDKey := r.cacheKey("find_by_id", account.ID().Value())
		_ = r.cache.Delete(accountIDKey) // Ignore errors
		r.invalidateUserCache(account.UserID().Value())
	}

	return nil
}

// Delete deletes an account and invalidates cache.
func (r *CachedAccountRepository) Delete(id valueobjects.AccountID) error {
	// Get account first to know user ID for cache invalidation
	account, _ := r.repository.FindByID(id)

	// Delete from repository
	err := r.repository.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if account != nil {
		accountIDKey := r.cacheKey("find_by_id", id.Value())
		_ = r.cache.Delete(accountIDKey) // Ignore errors
		r.invalidateUserCache(account.UserID().Value())
	}

	return nil
}

// Exists checks if an account exists (no caching for boolean checks).
func (r *CachedAccountRepository) Exists(id valueobjects.AccountID) (bool, error) {
	return r.repository.Exists(id)
}

// Count returns the total number of accounts for a user (no caching for counts).
func (r *CachedAccountRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	return r.repository.Count(userID)
}

// FindByUserIDWithPagination finds accounts with pagination (no caching for paginated results).
func (r *CachedAccountRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	offset, limit int,
) ([]*entities.Account, int64, error) {
	// Paginated results are not cached to avoid cache complexity
	return r.repository.FindByUserIDWithPagination(userID, context, offset, limit)
}
