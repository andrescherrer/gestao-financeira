package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeadersMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check security headers
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "off", resp.Header.Get("X-DNS-Prefetch-Control"))
	assert.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal(t, "geolocation=(), microphone=(), camera=()", resp.Header.Get("Permissions-Policy"))
	assert.Equal(t, "default-src 'self'; script-src 'none'; style-src 'none'; img-src 'none'; font-src 'none'; connect-src 'self'; frame-ancestors 'none';", resp.Header.Get("Content-Security-Policy"))
	assert.Equal(t, "none", resp.Header.Get("X-Permitted-Cross-Domain-Policies"))
}

func TestSecurityHeadersMiddleware_HSTS(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("env", "production")
		return c.Next()
	})
	app.Use(SecurityHeadersMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check HSTS header in production
	assert.Equal(t, "max-age=31536000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
}

func TestSecurityHeadersMiddleware_NoHSTSInDev(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("env", "development")
		return c.Next()
	})
	app.Use(SecurityHeadersMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// HSTS should not be set in development
	assert.Empty(t, resp.Header.Get("Strict-Transport-Security"))
}

