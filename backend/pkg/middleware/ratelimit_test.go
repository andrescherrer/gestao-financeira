package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gestao-financeira/backend/pkg/cache"
)

func TestRateLimitMiddleware_WithoutCache(t *testing.T) {
	app := fiber.New()
	config := DefaultRateLimitConfig()
	config.CacheService = nil // No cache = no rate limiting

	app.Use(RateLimitMiddleware(config))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRateLimitMiddleware_WithCache(t *testing.T) {
	// Skip if Redis is not available
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	app := fiber.New()
	config := DefaultRateLimitConfig()
	config.CacheService = cacheService
	config.MaxRequests = 2 // Allow only 2 requests
	config.Window = 1 * time.Minute

	app.Use(RateLimitMiddleware(config))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// First request - should succeed
	req1 := httptest.NewRequest("GET", "/test", nil)
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
	assert.Equal(t, "2", resp1.Header.Get("X-RateLimit-Limit"))
	assert.Equal(t, "1", resp1.Header.Get("X-RateLimit-Remaining"))

	// Second request - should succeed
	req2 := httptest.NewRequest("GET", "/test", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	assert.Equal(t, "0", resp2.Header.Get("X-RateLimit-Remaining"))

	// Third request - should be rate limited
	req3 := httptest.NewRequest("GET", "/test", nil)
	resp3, err := app.Test(req3)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp3.StatusCode)
}

func TestUserRateLimitMiddleware(t *testing.T) {
	// Skip if Redis is not available
	redisURL := "redis://localhost:6379"
	cacheService, err := cache.NewCacheService(redisURL)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
		return
	}
	defer cacheService.Close()

	app := fiber.New()
	config := DefaultRateLimitConfig()
	config.CacheService = cacheService
	config.MaxRequests = 2
	config.Window = 1 * time.Minute

	app.Use(UserRateLimitMiddleware(config))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Request without user ID - should use IP
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
