package usecases

import (
	"errors"
	"testing"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// Extended mock repository with error support
type mockListAccountRepository struct {
	*mockAccountRepository
	findByUserIDErr           error
	findByUserIDAndContextErr error
}

func newMockListAccountRepository() *mockListAccountRepository {
	return &mockListAccountRepository{
		mockAccountRepository: newMockAccountRepository(),
	}
}

func (m *mockListAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	if m.findByUserIDErr != nil {
		return nil, m.findByUserIDErr
	}
	return m.mockAccountRepository.FindByUserID(userID)
}

func (m *mockListAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	if m.findByUserIDAndContextErr != nil {
		return nil, m.findByUserIDAndContextErr
	}
	return m.mockAccountRepository.FindByUserIDAndContext(userID, context)
}

func TestListAccountsUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	userID2 := identityvalueobjects.GenerateUserID()

	tests := []struct {
		name      string
		input     dtos.ListAccountsInput
		setupMock func(*mockListAccountRepository)
		wantError bool
		errorMsg  string
		wantCount int
	}{
		{
			name: "list all accounts for user",
			input: dtos.ListAccountsInput{
				UserID:  userID.Value(),
				Context: "",
			},
			setupMock: func(m *mockListAccountRepository) {
				// Create test accounts
				account1, _ := createTestAccount(userID, "Conta Corrente", "BANK", 1000.00, "BRL", "PERSONAL")
				account2, _ := createTestAccount(userID, "Conta Poupança", "BANK", 500.00, "BRL", "PERSONAL")
				account3, _ := createTestAccount(userID, "Conta Empresarial", "BANK", 2000.00, "BRL", "BUSINESS")
				_ = m.Save(account1)
				_ = m.Save(account2)
				_ = m.Save(account3)
			},
			wantError: false,
			wantCount: 3,
		},
		{
			name: "list accounts filtered by context",
			input: dtos.ListAccountsInput{
				UserID:  userID.Value(),
				Context: "PERSONAL",
			},
			setupMock: func(m *mockListAccountRepository) {
				// Create test accounts
				account1, _ := createTestAccount(userID, "Conta Corrente", "BANK", 1000.00, "BRL", "PERSONAL")
				account2, _ := createTestAccount(userID, "Conta Poupança", "BANK", 500.00, "BRL", "PERSONAL")
				account3, _ := createTestAccount(userID, "Conta Empresarial", "BANK", 2000.00, "BRL", "BUSINESS")
				_ = m.Save(account1)
				_ = m.Save(account2)
				_ = m.Save(account3)
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name: "list accounts for user with no accounts",
			input: dtos.ListAccountsInput{
				UserID:  userID2.Value(),
				Context: "",
			},
			setupMock: func(m *mockListAccountRepository) {
				// No accounts for this user
			},
			wantError: false,
			wantCount: 0,
		},
		{
			name: "invalid user ID",
			input: dtos.ListAccountsInput{
				UserID:  "invalid-uuid",
				Context: "",
			},
			setupMock: func(m *mockListAccountRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid user ID",
		},
		{
			name: "invalid context",
			input: dtos.ListAccountsInput{
				UserID:  userID.Value(),
				Context: "INVALID",
			},
			setupMock: func(m *mockListAccountRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid account context",
		},
		{
			name: "repository error",
			input: dtos.ListAccountsInput{
				UserID:  userID.Value(),
				Context: "",
			},
			setupMock: func(m *mockListAccountRepository) {
				m.findByUserIDErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to find accounts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockListAccountRepository()
			tt.setupMock(mockRepo)

			useCase := NewListAccountsUseCase(mockRepo)
			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output == nil {
				t.Errorf("expected output but got nil")
				return
			}

			if output.Count != tt.wantCount {
				t.Errorf("expected count %d, got %d", tt.wantCount, output.Count)
			}

			if len(output.Accounts) != tt.wantCount {
				t.Errorf("expected %d accounts, got %d", tt.wantCount, len(output.Accounts))
			}

			// Validate account structure
			for _, account := range output.Accounts {
				if account.AccountID == "" {
					t.Errorf("expected account ID to be set")
				}
				if account.UserID != tt.input.UserID {
					t.Errorf("expected user ID %s, got %s", tt.input.UserID, account.UserID)
				}
				if account.Name == "" {
					t.Errorf("expected account name to be set")
				}
				if account.Type == "" {
					t.Errorf("expected account type to be set")
				}
				if account.Currency == "" {
					t.Errorf("expected currency to be set")
				}
				if account.Context == "" {
					t.Errorf("expected context to be set")
				}
				if tt.input.Context != "" && account.Context != tt.input.Context {
					t.Errorf("expected context %s, got %s", tt.input.Context, account.Context)
				}
			}
		})
	}
}

// Helper function to create a test account
func createTestAccount(
	userID identityvalueobjects.UserID,
	name, accountType string,
	balance float64,
	currency, context string,
) (*entities.Account, error) {
	accountName, err := valueobjects.NewAccountName(name)
	if err != nil {
		return nil, err
	}

	accType, err := valueobjects.NewAccountType(accountType)
	if err != nil {
		return nil, err
	}

	curr, err := sharedvalueobjects.NewCurrency(currency)
	if err != nil {
		return nil, err
	}

	balanceCents := int64(balance * 100)
	money, err := sharedvalueobjects.NewMoney(balanceCents, curr)
	if err != nil {
		return nil, err
	}

	ctx, err := sharedvalueobjects.NewAccountContext(context)
	if err != nil {
		return nil, err
	}

	return entities.NewAccount(userID, accountName, accType, money, ctx)
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
