package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/identity/infrastructure/services"
)

// AuthMiddleware creates a middleware that validates JWT tokens.
// It extracts the token from the Authorization header, validates it,
// and adds user information to the request context.
func AuthMiddleware(jwtService *services.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn().Str("path", c.Path()).Msg("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
				"code":  fiber.StatusUnauthorized,
			})
		}

		// Remove "Bearer " prefix if present
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Try without prefix (some clients might not use Bearer)
			tokenString = authHeader
		}

		// Remove any whitespace
		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			log.Warn().Str("path", c.Path()).Msg("Empty token in Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is required",
				"code":  fiber.StatusUnauthorized,
			})
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Warn().Err(err).Str("path", c.Path()).Msg("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"code":  fiber.StatusUnauthorized,
			})
		}

		// Add user information to context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("claims", claims)

		// Continue to next handler
		return c.Next()
	}
}

// GetUserID extracts the user ID from the request context.
// Returns empty string if not found.
func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("userID").(string); ok {
		return userID
	}
	return ""
}

// GetUserEmail extracts the user email from the request context.
// Returns empty string if not found.
func GetUserEmail(c *fiber.Ctx) string {
	if userEmail, ok := c.Locals("userEmail").(string); ok {
		return userEmail
	}
	return ""
}

// GetClaims extracts the JWT claims from the request context.
// Returns nil if not found.
func GetClaims(c *fiber.Ctx) *services.Claims {
	if claims, ok := c.Locals("claims").(*services.Claims); ok {
		return claims
	}
	return nil
}
