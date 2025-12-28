package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/pkg/cache"
)

// AuthMiddlewareConfig holds configuration for the auth middleware.
type AuthMiddlewareConfig struct {
	JWTService     *services.JWTService
	UserRepository repositories.UserRepository
	CacheService   *cache.CacheService // Optional: for caching user existence checks
	CacheTTL       time.Duration       // TTL for user existence cache (default: 30 seconds)
}

// AuthMiddleware creates a middleware that validates JWT tokens and verifies user existence.
// It extracts the token from the Authorization header, validates it,
// verifies that the user exists in the database, and adds user information to the request context.
func AuthMiddleware(config AuthMiddlewareConfig) fiber.Handler {
	// Set default cache TTL if not provided
	if config.CacheTTL == 0 {
		config.CacheTTL = 30 * time.Second
	}

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
		claims, err := config.JWTService.ValidateToken(tokenString)
		if err != nil {
			log.Warn().Err(err).Str("path", c.Path()).Msg("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"code":  fiber.StatusUnauthorized,
			})
		}

		// Verify user exists in database (with caching)
		userID, err := valueobjects.NewUserID(claims.UserID)
		if err != nil {
			log.Warn().Err(err).Str("path", c.Path()).Str("user_id", claims.UserID).Msg("Invalid user ID format in token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
				"code":  fiber.StatusUnauthorized,
			})
		}

		// Check cache first (if cache service is available)
		cacheKey := fmt.Sprintf("auth:user_exists:%s", claims.UserID)
		needDatabaseCheck := true

		if config.CacheService != nil {
			cached, err := config.CacheService.Get(cacheKey)
			if err == nil && cached != nil {
				// User existence is cached, skip database check
				needDatabaseCheck = false
				log.Debug().Str("user_id", claims.UserID).Msg("User existence verified from cache")
			}
		}

		// Check database if not in cache or cache not available
		if needDatabaseCheck {
			user, err := config.UserRepository.FindByID(userID)
			if err != nil {
				log.Error().Err(err).Str("user_id", claims.UserID).Str("path", c.Path()).Msg("Failed to check user existence")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
					"code":  fiber.StatusInternalServerError,
				})
			}

			if user == nil {
				log.Warn().Str("user_id", claims.UserID).Str("path", c.Path()).Msg("User not found in database")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User not found",
					"code":  fiber.StatusUnauthorized,
				})
			}

			// Cache user existence (only if user exists and cache is available)
			if config.CacheService != nil {
				_ = config.CacheService.Set(cacheKey, []byte("1"), config.CacheTTL)
			}
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
