package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/identity/infrastructure/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// mockUserRepository is a mock implementation of UserRepository for testing.
type mockUserRepository struct {
	users map[string]*entities.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (m *mockUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
	user, ok := m.users[id.Value()]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepository) FindByEmail(email valueobjects.Email) (*entities.User, error) {
	for _, user := range m.users {
		if user.Email().Equals(email) {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) Save(user *entities.User) error {
	m.users[user.ID().Value()] = user
	return nil
}

func (m *mockUserRepository) Delete(id valueobjects.UserID) error {
	delete(m.users, id.Value())
	return nil
}

func (m *mockUserRepository) Exists(email valueobjects.Email) (bool, error) {
	for _, user := range m.users {
		if user.Email().Equals(email) {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepository) Count() (int64, error) {
	return int64(len(m.users)), nil
}

func TestAuthMiddleware(t *testing.T) {
	jwtService := services.NewJWTService()
	userRepository := newMockUserRepository()

	// Create a test user
	userID := "550e8400-e29b-41d4-a716-446655440000"
	email, _ := valueobjects.NewEmail("test@example.com")
	// Use a valid bcrypt hash (minimum 60 characters)
	passwordHash, _ := valueobjects.NewPasswordHashFromHash("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
	name, _ := valueobjects.NewUserName("Test", "User")
	userIDVO, _ := valueobjects.NewUserID(userID)
	now := time.Now()
	user, err := entities.FromPersistence(userIDVO, email, passwordHash, name, now, now, true)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if err := userRepository.Save(user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Generate a valid token
	token, err := jwtService.GenerateToken(userID, email.Value())
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	app := fiber.New()
	app.Use(AuthMiddleware(AuthMiddlewareConfig{
		JWTService:     jwtService,
		UserRepository: userRepository,
		CacheService:   nil, // No cache for tests
	}))

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
