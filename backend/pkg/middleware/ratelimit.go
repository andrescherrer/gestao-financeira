package middleware

import (
	"fmt"
	"strconv"
	"time"

	"gestao-financeira/backend/pkg/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RateLimitConfig holds configuration for rate limiting.
type RateLimitConfig struct {
	MaxRequests  int                     // Maximum number of requests
	Window       time.Duration           // Time window (e.g., 1 minute)
	CacheService *cache.CacheService     // Redis cache service (optional)
	KeyGenerator func(*fiber.Ctx) string // Function to generate cache key
	SkipSuccess  bool                    // Skip rate limiting on successful requests
	Message      string                  // Custom error message
}

// DefaultRateLimitConfig returns default rate limit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as default key
			return fmt.Sprintf("ratelimit:ip:%s", c.IP())
		},
		Message: "Too many requests, please try again later",
	}
}

// RateLimitMiddleware creates a rate limiting middleware.
// If cache service is nil, rate limiting is disabled (graceful degradation).
func RateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	if config.CacheService == nil {
		log.Warn().Msg("Rate limiting disabled: cache service not available")
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	if config.MaxRequests <= 0 {
		config.MaxRequests = 100
	}

	if config.Window <= 0 {
		config.Window = 1 * time.Minute
	}

	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return fmt.Sprintf("ratelimit:ip:%s", c.IP())
		}
	}

	if config.Message == "" {
		config.Message = "Too many requests, please try again later"
	}

	return func(c *fiber.Ctx) error {
		key := config.KeyGenerator(c)
		cacheKey := fmt.Sprintf("%s:%d", key, time.Now().Unix()/int64(config.Window.Seconds()))

		// Get current count
		countBytes, err := config.CacheService.Get(cacheKey)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get rate limit count, allowing request")
			return c.Next()
		}

		var count int
		if countBytes != nil {
			count, err = strconv.Atoi(string(countBytes))
			if err != nil {
				log.Warn().Err(err).Msg("Failed to parse rate limit count, resetting")
				count = 0
			}
		}

		// Check if limit exceeded
		if count >= config.MaxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       config.Message,
				"code":        fiber.StatusTooManyRequests,
				"retry_after": int(config.Window.Seconds()),
			})
		}

		// Increment count
		newCount := count + 1
		err = config.CacheService.Set(cacheKey, []byte(strconv.Itoa(newCount)), config.Window)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to set rate limit count, allowing request")
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(config.MaxRequests))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(config.MaxRequests-newCount))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))

		return c.Next()
	}
}

// UserRateLimitMiddleware creates a rate limiting middleware that uses user ID as key.
// Falls back to IP if user is not authenticated.
func UserRateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	originalKeyGenerator := config.KeyGenerator
	config.KeyGenerator = func(c *fiber.Ctx) string {
		// Try to get user ID from context
		userID := GetUserID(c)
		if userID != "" {
			return fmt.Sprintf("ratelimit:user:%s", userID)
		}
		// Fallback to IP
		if originalKeyGenerator != nil {
			return originalKeyGenerator(c)
		}
		return fmt.Sprintf("ratelimit:ip:%s", c.IP())
	}
	return RateLimitMiddleware(config)
}
