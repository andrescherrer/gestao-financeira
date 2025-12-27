package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// CacheService provides caching functionality using Redis.
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService creates a new CacheService instance.
func NewCacheService(redisURL string) (*CacheService, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info().Msg("Redis cache service initialized")

	return &CacheService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache.
func (c *CacheService) Get(key string) ([]byte, error) {
	val, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Key doesn't exist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}
	return val, nil
}

// Set stores a value in cache with expiration.
func (c *CacheService) Set(key string, value []byte, expiration time.Duration) error {
	if err := c.client.Set(c.ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	return nil
}

// SetJSON stores a JSON-serializable value in cache.
func (c *CacheService) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.Set(key, jsonData, expiration)
}

// GetJSON retrieves and unmarshals a JSON value from cache.
func (c *CacheService) GetJSON(key string, dest interface{}) (bool, error) {
	data, err := c.Get(key)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil // Key doesn't exist
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return true, nil
}

// Delete removes a key from cache.
func (c *CacheService) Delete(key string) error {
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// DeletePattern removes all keys matching a pattern.
func (c *CacheService) DeletePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			log.Warn().Err(err).Str("key", iter.Val()).Msg("Failed to delete cache key")
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan cache keys: %w", err)
	}
	return nil
}

// Close closes the Redis connection.
func (c *CacheService) Close() error {
	return c.client.Close()
}
