package middleware

import (
	"fmt"
	"strings"
	"time"

	"gestao-financeira/backend/pkg/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// EndpointRateLimitConfig defines rate limit configuration for a specific endpoint pattern.
type EndpointRateLimitConfig struct {
	Pattern     string        // Endpoint pattern (e.g., "/api/v1/auth/login", "/api/v1/transactions")
	MaxRequests int           // Maximum requests per window
	Window      time.Duration // Time window
	Method      string        // HTTP method (empty for all methods)
}

// GranularRateLimitConfig holds configuration for granular rate limiting.
type GranularRateLimitConfig struct {
	CacheService *cache.CacheService       // Redis cache service
	Default      *RateLimitConfig          // Default rate limit for unmatched endpoints (nil uses defaults)
	Endpoints    []EndpointRateLimitConfig // Specific endpoint configurations
}

// GranularRateLimitMiddleware creates a rate limiting middleware with different limits per endpoint.
func GranularRateLimitMiddleware(config GranularRateLimitConfig) fiber.Handler {
	if config.CacheService == nil {
		log.Warn().Msg("Granular rate limiting disabled: cache service not available")
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Ensure default config exists and is valid
	defaultConfig := config.Default
	if defaultConfig == nil {
		defaultConfigValue := DefaultRateLimitConfig()
		defaultConfig = &defaultConfigValue
		config.Default = defaultConfig
	}

	if defaultConfig.MaxRequests <= 0 {
		defaultConfig.MaxRequests = 100
	}
	if defaultConfig.Window <= 0 {
		defaultConfig.Window = 1 * time.Minute
	}
	if defaultConfig.KeyGenerator == nil {
		defaultConfig.KeyGenerator = func(c *fiber.Ctx) string {
			return fmt.Sprintf("ratelimit:ip:%s", c.IP())
		}
	}
	if defaultConfig.Message == "" {
		defaultConfig.Message = "Too many requests, please try again later"
	}

	return func(c *fiber.Ctx) error {
		// Find matching endpoint configuration
		endpointConfig := findEndpointConfig(c.Path(), c.Method(), config.Endpoints)

		var maxRequests int
		var window time.Duration
		var message string

		if endpointConfig != nil {
			// Use endpoint-specific configuration
			maxRequests = endpointConfig.MaxRequests
			window = endpointConfig.Window
			message = fmt.Sprintf("Too many requests to %s, please try again later", c.Path())
		} else {
			// Use default configuration
			defaultConfig := config.Default
			maxRequests = defaultConfig.MaxRequests
			window = defaultConfig.Window
			message = defaultConfig.Message
		}

		// Generate cache key
		keyGenerator := config.Default.KeyGenerator
		if keyGenerator == nil {
			keyGenerator = func(c *fiber.Ctx) string {
				return fmt.Sprintf("ratelimit:ip:%s", c.IP())
			}
		}

		key := keyGenerator(c)
		cacheKey := fmt.Sprintf("ratelimit:%s:%s:%d",
			strings.ReplaceAll(c.Path(), "/", "_"),
			key,
			time.Now().Unix()/int64(window.Seconds()))

		// Get current count
		countBytes, err := config.CacheService.Get(cacheKey)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get rate limit count, allowing request")
			return c.Next()
		}

		var count int
		if countBytes != nil {
			count = parseCount(string(countBytes))
		}

		// Check if limit exceeded
		if count >= maxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       message,
				"code":        fiber.StatusTooManyRequests,
				"retry_after": int(window.Seconds()),
				"endpoint":    c.Path(),
			})
		}

		// Increment count
		newCount := count + 1
		err = config.CacheService.Set(cacheKey, []byte(fmt.Sprintf("%d", newCount)), window)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to set rate limit count, allowing request")
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests-newCount))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		return c.Next()
	}
}

// findEndpointConfig finds the first matching endpoint configuration for the given path and method.
func findEndpointConfig(path, method string, configs []EndpointRateLimitConfig) *EndpointRateLimitConfig {
	for _, config := range configs {
		// Check if method matches (empty method means all methods)
		if config.Method != "" && config.Method != method {
			continue
		}

		// Check if path matches
		if matchesPattern(path, config.Pattern) {
			return &config
		}
	}
	return nil
}

// matchesPattern checks if a path matches a pattern.
// Supports exact match and prefix match (if pattern ends with "*").
func matchesPattern(path, pattern string) bool {
	if pattern == "" {
		return false
	}

	// Exact match
	if path == pattern {
		return true
	}

	// Prefix match (if pattern ends with "*")
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}

	return false
}

// parseCount safely parses a count string to int.
func parseCount(s string) int {
	if s == "" {
		return 0
	}

	var count int
	_, err := fmt.Sscanf(s, "%d", &count)
	if err != nil {
		return 0
	}
	return count
}

// DefaultGranularRateLimitConfig returns a default granular rate limit configuration.
func DefaultGranularRateLimitConfig(cacheService *cache.CacheService) GranularRateLimitConfig {
	defaultConfig := RateLimitConfig{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("ratelimit:ip:%s", c.IP())
		},
		Message: "Too many requests, please try again later",
	}

	return GranularRateLimitConfig{
		CacheService: cacheService,
		Default:      &defaultConfig,
		Endpoints: []EndpointRateLimitConfig{
			// Auth endpoints - more restrictive
			{
				Pattern:     "/api/v1/auth/login",
				MaxRequests: 10,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
			{
				Pattern:     "/api/v1/auth/register",
				MaxRequests: 5,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
			// Write operations - moderate restriction
			{
				Pattern:     "/api/v1/transactions",
				MaxRequests: 30,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
			{
				Pattern:     "/api/v1/transactions/*",
				MaxRequests: 30,
				Window:      1 * time.Minute,
				Method:      "PUT",
			},
			{
				Pattern:     "/api/v1/transactions/*",
				MaxRequests: 30,
				Window:      1 * time.Minute,
				Method:      "DELETE",
			},
			{
				Pattern:     "/api/v1/accounts",
				MaxRequests: 20,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
			{
				Pattern:     "/api/v1/categories",
				MaxRequests: 20,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
			{
				Pattern:     "/api/v1/budgets",
				MaxRequests: 20,
				Window:      1 * time.Minute,
				Method:      "POST",
			},
		},
	}
}
