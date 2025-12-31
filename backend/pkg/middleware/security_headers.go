package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersMiddleware creates a middleware that adds security headers to responses
// Similar to Helmet.js for Express, this middleware adds various security headers
func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-Content-Type-Options: Prevents MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking attacks
		c.Set("X-Frame-Options", "DENY")

		// X-XSS-Protection: Enables XSS filter in older browsers (deprecated but still used)
		c.Set("X-XSS-Protection", "1; mode=block")

		// X-DNS-Prefetch-Control: Controls DNS prefetching
		c.Set("X-DNS-Prefetch-Control", "off")

		// Referrer-Policy: Controls referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Controls browser features (formerly Feature-Policy)
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Strict-Transport-Security (HSTS): Force HTTPS (only in production)
		// Note: This should only be set if HTTPS is enabled
		// We'll check environment or config to enable this
		env := c.Locals("env")
		if env == "production" {
			// HSTS: max-age=31536000 (1 year), includeSubDomains, preload
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Content-Security-Policy: Controls resource loading
		// For API, we can be restrictive since we only serve JSON
		csp := "default-src 'self'; script-src 'none'; style-src 'none'; img-src 'none'; font-src 'none'; connect-src 'self'; frame-ancestors 'none';"
		c.Set("Content-Security-Policy", csp)

		// X-Permitted-Cross-Domain-Policies: Controls cross-domain policy
		c.Set("X-Permitted-Cross-Domain-Policies", "none")

		// Continue with the request
		return c.Next()
	}
}

