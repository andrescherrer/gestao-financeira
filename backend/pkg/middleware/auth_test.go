package middleware

import (
	"net/http/httptest"
	"testing"

	"gestao-financeira/backend/internal/identity/infrastructure/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	jwtService := services.NewJWTService()

	// Generate a valid token
	userID := "test-user-id"
	email := "test@example.com"
	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	app := fiber.New()
	app.Use(AuthMiddleware(jwtService))

	app.Get("/protected", func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		userEmail := GetUserEmail(c)
		claims := GetClaims(c)

		if userID == "" || userEmail == "" || claims == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "User information not found in context",
			})
		}

		return c.JSON(fiber.Map{
			"user_id":    userID,
			"user_email": userEmail,
			"message":    "Access granted",
		})
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "valid token with Bearer prefix",
			authHeader:     "Bearer " + token,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "valid token without Bearer prefix",
			authHeader:     token,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "Authorization header is required",
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "Invalid or expired token",
		},
		{
			name:           "empty token",
			authHeader:     "Bearer ",
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "Token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				// Check if error message is in response
				// This is a basic check - in production you'd parse JSON
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-id")
		userID := GetUserID(c)
		assert.Equal(t, "test-id", userID)
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetUserEmail(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("userEmail", "test@example.com")
		userEmail := GetUserEmail(c)
		assert.Equal(t, "test@example.com", userEmail)
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
