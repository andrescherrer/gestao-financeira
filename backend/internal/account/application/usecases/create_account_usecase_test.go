package usecases

import (
	"errors"
	"testing"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// mockAccountRepository is a mock implementation of AccountRepository for testing.
type mockAccountRepository struct {
	accounts map[string]*entities.Account
	saveErr  error
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{
		accounts: make(map[string]*entities.Account),
	}
}

func (m *mockAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	account, exists := m.accounts[id.Value()]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (m *mockAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	var result []*entities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	var result []*entities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() && account.Context().Value() == context.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepository) Save(account *entities.Account) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.accounts[account.ID().Value()] = account
	return nil
}

func (m *mockAccountRepository) Delete(id valueobjects.AccountID) error {
	delete(m.accounts, id.Value())
	return nil
}

func (m *mockAccountRepository) Exists(id valueobjects.AccountID) (bool, error) {
	_, exists := m.accounts[id.Value()]
	return exists, nil
}

func (m *mockAccountRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			count++
		}
	}
	return count, nil
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	eventBus := eventbus.NewEventBus()

	tests := []struct {
		name      string
		input     dtos.CreateAccountInput
		setupMock func(*mockAccountRepository)
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid account creation",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 100.50,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: false,
		},
		{
			name: "invalid user ID",
			input: dtos.CreateAccountInput{
				UserID:         "invalid-uuid",
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: true,
			errorMsg:  "invalid user ID",
		},
		{
			name: "invalid account name",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "AB", // Too short
				Type:           "BANK",
				InitialBalance: 0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: true,
			errorMsg:  "invalid account name",
		},
		{
			name: "invalid account type",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "Conta Corrente",
				Type:           "INVALID",
				InitialBalance: 0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: true,
			errorMsg:  "invalid account type",
		},
		{
			name: "invalid currency",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 0,
				Currency:       "INVALID",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: true,
			errorMsg:  "invalid currency",
		},
		{
			name: "invalid context",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 0,
				Currency:       "BRL",
				Context:        "INVALID",
			},
			setupMock: func(m *mockAccountRepository) {},
			wantError: true,
			errorMsg:  "invalid account context",
		},
		{
			name: "repository save error",
			input: dtos.CreateAccountInput{
				UserID:         userID.Value(),
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			setupMock: func(m *mockAccountRepository) {
				m.saveErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to save account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockAccountRepository()
			tt.setupMock(mockRepo)

			useCase := NewCreateAccountUseCase(mockRepo, eventBus)
			output, err := useCase.Execute(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateAccountUseCase.Execute() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				if err != nil && tt.errorMsg != "" {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("CreateAccountUseCase.Execute() error = %v, want error containing %v", err, tt.errorMsg)
					}
				}
				return
			}

			// Verify output
			if output == nil {
				t.Error("CreateAccountUseCase.Execute() output is nil")
				return
			}

			if output.AccountID == "" {
				t.Error("CreateAccountUseCase.Execute() output.AccountID is empty")
			}

			if output.UserID != tt.input.UserID {
				t.Errorf("CreateAccountUseCase.Execute() output.UserID = %v, want %v", output.UserID, tt.input.UserID)
			}

			if output.Name != tt.input.Name {
				t.Errorf("CreateAccountUseCase.Execute() output.Name = %v, want %v", output.Name, tt.input.Name)
			}

			if output.Type != tt.input.Type {
				t.Errorf("CreateAccountUseCase.Execute() output.Type = %v, want %v", output.Type, tt.input.Type)
			}

			// Verify account was saved
			accountID, _ := valueobjects.NewAccountID(output.AccountID)
			savedAccount, err := mockRepo.FindByID(accountID)
			if err != nil {
				t.Errorf("CreateAccountUseCase.Execute() failed to find saved account: %v", err)
			}
			if savedAccount == nil {
				t.Error("CreateAccountUseCase.Execute() account was not saved")
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
