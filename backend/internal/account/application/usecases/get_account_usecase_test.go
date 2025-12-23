package usecases

import (
	"errors"
	"testing"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// Extended mock repository with error support for GetAccountUseCase
type mockGetAccountRepository struct {
	*mockAccountRepository
	findByIDErr error
}

func newMockGetAccountRepository() *mockGetAccountRepository {
	return &mockGetAccountRepository{
		mockAccountRepository: newMockAccountRepository(),
	}
}

func (m *mockGetAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.mockAccountRepository.FindByID(id)
}

func TestGetAccountUseCase_Execute(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()

	tests := []struct {
		name      string
		input     dtos.GetAccountInput
		setupMock func(*mockGetAccountRepository)
		wantError bool
		errorMsg  string
		validate  func(*testing.T, *dtos.GetAccountOutput)
	}{
		{
			name: "get existing account",
			input: dtos.GetAccountInput{
				AccountID: "",
			},
			setupMock: func(m *mockGetAccountRepository) {
				account, _ := createTestAccount(userID, "Conta Corrente", "BANK", 1000.00, "BRL", "PERSONAL")
				_ = m.Save(account)
				// Set AccountID in input after account is created
			},
			wantError: false,
			validate: func(t *testing.T, output *dtos.GetAccountOutput) {
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
				if output.AccountID == "" {
					t.Errorf("expected account ID to be set")
				}
				if output.UserID != userID.Value() {
					t.Errorf("expected user ID %s, got %s", userID.Value(), output.UserID)
				}
				if output.Name != "Conta Corrente" {
					t.Errorf("expected name 'Conta Corrente', got %s", output.Name)
				}
				if output.Type != "BANK" {
					t.Errorf("expected type 'BANK', got %s", output.Type)
				}
				if output.Balance != 1000.00 {
					t.Errorf("expected balance 1000.00, got %f", output.Balance)
				}
				if output.Currency != "BRL" {
					t.Errorf("expected currency 'BRL', got %s", output.Currency)
				}
				if output.Context != "PERSONAL" {
					t.Errorf("expected context 'PERSONAL', got %s", output.Context)
				}
				if !output.IsActive {
					t.Errorf("expected account to be active")
				}
			},
		},
		{
			name: "account not found",
			input: dtos.GetAccountInput{
				AccountID: valueobjects.GenerateAccountID().Value(),
			},
			setupMock: func(m *mockGetAccountRepository) {
				// No accounts in repository
			},
			wantError: true,
			errorMsg:  "account not found",
		},
		{
			name: "invalid account ID",
			input: dtos.GetAccountInput{
				AccountID: "invalid-uuid",
			},
			setupMock: func(m *mockGetAccountRepository) {
				// No setup needed
			},
			wantError: true,
			errorMsg:  "invalid account ID",
		},
		{
			name: "repository error",
			input: dtos.GetAccountInput{
				AccountID: valueobjects.GenerateAccountID().Value(),
			},
			setupMock: func(m *mockGetAccountRepository) {
				m.findByIDErr = errors.New("database error")
			},
			wantError: true,
			errorMsg:  "failed to find account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockGetAccountRepository()
			tt.setupMock(mockRepo)

			// For the first test case, we need to set the AccountID after creating the account
			if tt.name == "get existing account" {
				// Create account first
				account, _ := createTestAccount(userID, "Conta Corrente", "BANK", 1000.00, "BRL", "PERSONAL")
				_ = mockRepo.Save(account)
				tt.input.AccountID = account.ID().Value()
			}

			useCase := NewGetAccountUseCase(mockRepo)
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

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

