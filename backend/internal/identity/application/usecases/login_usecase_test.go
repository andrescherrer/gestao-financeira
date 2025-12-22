package usecases

import (
	"strings"
	"testing"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
)

// mockJWTService is a mock implementation of JWTService for testing.
type mockJWTService struct{}

func (m *mockJWTService) GenerateToken(userID, email string) (string, error) {
	return "mock-jwt-token-" + userID, nil
}

func (m *mockJWTService) ValidateToken(tokenString string) (*services.Claims, error) {
	return nil, nil
}

func (m *mockJWTService) GetTokenExpiration(tokenString string) (interface{}, error) {
	return nil, nil
}

func TestLoginUseCase_Execute(t *testing.T) {
	// Create a test user
	email, _ := valueobjects.NewEmail("user@example.com")
	passwordHash, _ := valueobjects.NewPasswordHashFromPlain("password123")
	name, _ := valueobjects.NewUserName("John", "Doe")
	user, _ := entities.NewUser(email, passwordHash, name)

	tests := []struct {
		name          string
		input         dtos.LoginInput
		setupMock     func(*mockUserRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful login",
			input: dtos.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				m.users[user.ID().Value()] = user
			},
			wantError: false,
		},
		{
			name: "user not found",
			input: dtos.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				// User doesn't exist
			},
			wantError:     true,
			errorContains: "invalid email or password",
		},
		{
			name: "invalid password",
			input: dtos.LoginInput{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *mockUserRepository) {
				m.users[user.ID().Value()] = user
			},
			wantError:     true,
			errorContains: "invalid email or password",
		},
		{
			name: "inactive user",
			input: dtos.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			setupMock: func(m *mockUserRepository) {
				user.Deactivate()
				m.users[user.ID().Value()] = user
			},
			wantError:     true,
			errorContains: "inactive",
		},
		{
			name: "invalid email format",
			input: dtos.LoginInput{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock:     func(m *mockUserRepository) {},
			wantError:     true,
			errorContains: "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)

			jwtService := &mockJWTService{}
			useCase := NewLoginUseCase(mockRepo, jwtService)

			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("LoginUseCase.Execute() expected error, got nil")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("LoginUseCase.Execute() error = %v, want error containing %v", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("LoginUseCase.Execute() unexpected error = %v", err)
				return
			}

			if output == nil {
				t.Errorf("LoginUseCase.Execute() expected output, got nil")
				return
			}

			if output.Token == "" {
				t.Errorf("LoginUseCase.Execute() output.Token is empty")
			}

			if output.Email != tt.input.Email {
				t.Errorf("LoginUseCase.Execute() output.Email = %v, want %v", output.Email, tt.input.Email)
			}

			if output.UserID == "" {
				t.Errorf("LoginUseCase.Execute() output.UserID is empty")
			}
		})
	}
}
