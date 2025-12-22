package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/application/usecases"
	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockUserRepository is a mock implementation of UserRepository for testing.
type mockUserRepository struct {
	users  map[string]*entities.User
	exists map[string]bool
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*entities.User),
		exists: make(map[string]bool),
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
	m.exists[user.Email().Value()] = true
	return nil
}

func (m *mockUserRepository) Delete(id valueobjects.UserID) error {
	delete(m.users, id.Value())
	return nil
}

func (m *mockUserRepository) Exists(email valueobjects.Email) (bool, error) {
	return m.exists[email.Value()], nil
}

func (m *mockUserRepository) Count() (int64, error) {
	return int64(len(m.users)), nil
}

func TestAuthHandler_Register(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockUserRepository()
	eventBus := eventbus.NewEventBus()
	registerUseCase := usecases.NewRegisterUserUseCase(mockRepo, eventBus)
	loginUseCase := usecases.NewLoginUseCase(mockRepo, services.NewJWTService())
	handler := NewAuthHandler(registerUseCase, loginUseCase)

	app.Post("/register", handler.Register)

	tests := []struct {
		name           string
		body           dtos.RegisterUserInput
		expectedStatus int
		expectedError  string
		setupMock      func()
	}{
		{
			name: "successful registration",
			body: dtos.RegisterUserInput{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: fiber.StatusCreated,
			setupMock:      func() {},
		},
		{
			name: "missing email",
			body: dtos.RegisterUserInput{
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Email is required",
		},
		{
			name: "missing password",
			body: dtos.RegisterUserInput{
				Email:     "user@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Password is required",
		},
		{
			name: "missing first name",
			body: dtos.RegisterUserInput{
				Email:    "user@example.com",
				Password: "password123",
				LastName: "Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "First name is required",
		},
		{
			name: "missing last name",
			body: dtos.RegisterUserInput{
				Email:     "user@example.com",
				Password:  "password123",
				FirstName: "John",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Last name is required",
		},
		{
			name: "email already exists",
			body: dtos.RegisterUserInput{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Smith",
			},
			expectedStatus: fiber.StatusConflict,
			expectedError:  "already exists",
			setupMock: func() {
				mockRepo.exists["existing@example.com"] = true
			},
		},
		{
			name: "invalid email format",
			body: dtos.RegisterUserInput{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "invalid password (too short)",
			body: dtos.RegisterUserInput{
				Email:     "user@example.com",
				Password:  "short",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Password must be at least 8 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.users = make(map[string]*entities.User)
			mockRepo.exists = make(map[string]bool)
			if tt.setupMock != nil {
				tt.setupMock()
			}

			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(bodyJSON))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				// Verify error message in response body
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				if errorMsg, ok := result["error"].(string); ok {
					assert.Contains(t, errorMsg, tt.expectedError)
				}
			} else if tt.expectedStatus == fiber.StatusCreated {
				// Verify success response
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "User registered successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockUserRepository()
	eventBus := eventbus.NewEventBus()
	jwtService := services.NewJWTService()

	// Create a test user
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")
	user, _ := entities.NewUser(email, passwordHash, name)
	mockRepo.users[user.ID().Value()] = user

	registerUseCase := usecases.NewRegisterUserUseCase(mockRepo, eventBus)
	loginUseCase := usecases.NewLoginUseCase(mockRepo, jwtService)
	handler := NewAuthHandler(registerUseCase, loginUseCase)

	app.Post("/login", handler.Login)

	tests := []struct {
		name           string
		body           dtos.LoginInput
		expectedStatus int
		expectedError  string
		setupMock      func()
	}{
		{
			name: "successful login",
			body: dtos.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "missing email",
			body: dtos.LoginInput{
				Password: "password123",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Email is required",
		},
		{
			name: "missing password",
			body: dtos.LoginInput{
				Email: "user@example.com",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Password is required",
		},
		{
			name: "invalid email format",
			body: dtos.LoginInput{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid email format",
		},
		{
			name: "user not found",
			body: dtos.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "Invalid email or password",
		},
		{
			name: "wrong password",
			body: dtos.LoginInput{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "Invalid email or password",
		},
		{
			name: "inactive user",
			body: dtos.LoginInput{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			expectedStatus: fiber.StatusForbidden,
			expectedError:  "inactive",
			setupMock: func() {
				// Create inactive user
				email, _ := valueobjects.NewEmail("inactive@example.com")
				passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
				name, _ := valueobjects.NewUserName("Inactive", "User")
				user, _ := entities.NewUser(email, passwordHash, name)
				user.Deactivate()
				mockRepo.users[user.ID().Value()] = user
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(bodyJSON))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				// Verify error message in response body
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				if errorMsg, ok := result["error"].(string); ok {
					assert.Contains(t, errorMsg, tt.expectedError)
				}
			} else if tt.expectedStatus == fiber.StatusOK {
				// Verify success response
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "Login successful", result["message"])
				assert.NotNil(t, result["data"])
				if data, ok := result["data"].(map[string]interface{}); ok {
					assert.NotEmpty(t, data["token"])
					assert.Equal(t, "user@example.com", data["email"])
				}
			}
		})
	}
}

func TestAuthHandler_InvalidJSON(t *testing.T) {
	app := fiber.New()
	mockRepo := newMockUserRepository()
	eventBus := eventbus.NewEventBus()
	registerUseCase := usecases.NewRegisterUserUseCase(mockRepo, eventBus)
	loginUseCase := usecases.NewLoginUseCase(mockRepo, services.NewJWTService())
	handler := NewAuthHandler(registerUseCase, loginUseCase)

	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	tests := []struct {
		name           string
		endpoint       string
		body           string
		expectedStatus int
	}{
		{
			name:           "invalid JSON for register",
			endpoint:       "/register",
			body:           "{invalid json}",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid JSON for login",
			endpoint:       "/login",
			body:           "{invalid json}",
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tt.endpoint, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
