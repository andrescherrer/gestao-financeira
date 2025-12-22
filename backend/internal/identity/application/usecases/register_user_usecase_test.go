package usecases

import (
	"strings"
	"testing"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
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

func TestRegisterUserUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         dtos.RegisterUserInput
		setupMock     func(*mockUserRepository)
		wantError     bool
		errorContains string
	}{
		{
			name: "successful registration",
			input: dtos.RegisterUserInput{
				Email:     "user@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock: func(m *mockUserRepository) {
				// User doesn't exist
			},
			wantError: false,
		},
		{
			name: "email already exists",
			input: dtos.RegisterUserInput{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Smith",
			},
			setupMock: func(m *mockUserRepository) {
				m.exists["existing@example.com"] = true
			},
			wantError:     true,
			errorContains: "already exists",
		},
		{
			name: "invalid email",
			input: dtos.RegisterUserInput{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock:     func(m *mockUserRepository) {},
			wantError:     true,
			errorContains: "invalid email",
		},
		{
			name: "invalid password",
			input: dtos.RegisterUserInput{
				Email:     "user@example.com",
				Password:  "short",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock:     func(m *mockUserRepository) {},
			wantError:     true,
			errorContains: "invalid password",
		},
		{
			name: "invalid name",
			input: dtos.RegisterUserInput{
				Email:     "user@example.com",
				Password:  "password123",
				FirstName: "J",
				LastName:  "D",
			},
			setupMock:     func(m *mockUserRepository) {},
			wantError:     true,
			errorContains: "invalid name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockUserRepository()
			tt.setupMock(mockRepo)

			eventBus := eventbus.NewEventBus()
			useCase := NewRegisterUserUseCase(mockRepo, eventBus)

			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("RegisterUserUseCase.Execute() expected error, got nil")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("RegisterUserUseCase.Execute() error = %v, want error containing %v", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("RegisterUserUseCase.Execute() unexpected error = %v", err)
				return
			}

			if output == nil {
				t.Errorf("RegisterUserUseCase.Execute() expected output, got nil")
				return
			}

			if output.Email != tt.input.Email {
				t.Errorf("RegisterUserUseCase.Execute() output.Email = %v, want %v", output.Email, tt.input.Email)
			}

			if output.FirstName != tt.input.FirstName {
				t.Errorf("RegisterUserUseCase.Execute() output.FirstName = %v, want %v", output.FirstName, tt.input.FirstName)
			}

			if output.LastName != tt.input.LastName {
				t.Errorf("RegisterUserUseCase.Execute() output.LastName = %v, want %v", output.LastName, tt.input.LastName)
			}

			if output.UserID == "" {
				t.Errorf("RegisterUserUseCase.Execute() output.UserID is empty")
			}
		})
	}
}
