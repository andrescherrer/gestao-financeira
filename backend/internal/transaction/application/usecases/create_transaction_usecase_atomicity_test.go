package usecases

import (
	"errors"
	"testing"
	"time"

	accountentities "gestao-financeira/backend/internal/account/domain/entities"
	accountrepositories "gestao-financeira/backend/internal/account/domain/repositories"
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
)

// mockUnitOfWorkWithErrors is a mock UnitOfWork that can simulate errors
type mockUnitOfWorkWithErrors struct {
	*mockUnitOfWork
	beginErr   error
	commitErr  error
	accountErr error
}

func newMockUnitOfWorkWithErrors(
	transactionRepo *mockTransactionRepository,
	accountRepo *mockAccountRepository,
) *mockUnitOfWorkWithErrors {
	return &mockUnitOfWorkWithErrors{
		mockUnitOfWork: newMockUnitOfWork(transactionRepo, accountRepo),
	}
}

func (m *mockUnitOfWorkWithErrors) Begin() error {
	if m.beginErr != nil {
		return m.beginErr
	}
	return m.mockUnitOfWork.Begin()
}

func (m *mockUnitOfWorkWithErrors) Commit() error {
	if m.commitErr != nil {
		return m.commitErr
	}
	return m.mockUnitOfWork.Commit()
}

func (m *mockUnitOfWorkWithErrors) AccountRepository() accountrepositories.AccountRepository {
	if m.accountErr != nil {
		return &mockAccountRepositoryWithError{err: m.accountErr}
	}
	return m.mockUnitOfWork.AccountRepository()
}

type mockAccountRepositoryWithError struct {
	err error
}

func (m *mockAccountRepositoryWithError) FindByID(id accountvalueobjects.AccountID) (*accountentities.Account, error) {
	return nil, m.err
}

func (m *mockAccountRepositoryWithError) FindByUserID(userID identityvalueobjects.UserID) ([]*accountentities.Account, error) {
	return nil, m.err
}

func (m *mockAccountRepositoryWithError) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*accountentities.Account, error) {
	return nil, m.err
}

func (m *mockAccountRepositoryWithError) Save(account *accountentities.Account) error {
	return m.err
}

func (m *mockAccountRepositoryWithError) Delete(id accountvalueobjects.AccountID) error {
	return m.err
}

func (m *mockAccountRepositoryWithError) Exists(id accountvalueobjects.AccountID) (bool, error) {
	return false, m.err
}

func (m *mockAccountRepositoryWithError) Count(userID identityvalueobjects.UserID) (int64, error) {
	return 0, m.err
}
func (m *mockAccountRepositoryWithError) FindByUserIDWithPagination(userID identityvalueobjects.UserID, context string, offset, limit int) ([]*accountentities.Account, int64, error) {
	return nil, 0, m.err
}

func TestCreateTransactionUseCase_Atomicity(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	eventBus := eventbus.NewEventBus()
	date := time.Now().Format("2006-01-02")

	tests := []struct {
		name           string
		input          dtos.CreateTransactionInput
		setupMock      func(*mockTransactionRepository, *mockAccountRepository, *mockUnitOfWorkWithErrors)
		wantError      bool
		errorMsg       string
		validateAtomic func(*testing.T, *mockTransactionRepository, *mockAccountRepository)
	}{
		{
			name: "atomicity: error returned if account not found",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Test transaction",
				Date:        date,
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWorkWithErrors) {
				// Don't create account - simulate account not found
			},
			wantError: true,
			errorMsg:  "account not found",
			validateAtomic: func(t *testing.T, txRepo *mockTransactionRepository, accRepo *mockAccountRepository) {
				// Note: Atomicity is tested in integration tests with real database
				// Mock repositories don't support true transaction rollback
			},
		},
		{
			name: "atomicity: account not updated if transaction save fails",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Test transaction",
				Date:        date,
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWorkWithErrors) {
				// Create account
				currency, _ := sharedvalueobjects.NewCurrency("BRL")
				initialBalance, _ := sharedvalueobjects.NewMoney(0, currency)
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)

				// Simulate transaction save failure
				txRepo.saveErr = errors.New("failed to save transaction")
			},
			wantError: true,
			errorMsg:  "failed to save transaction",
			validateAtomic: func(t *testing.T, txRepo *mockTransactionRepository, accRepo *mockAccountRepository) {
				// Verify account balance was NOT updated (atomicity)
				account, _ := accRepo.FindByID(accountID)
				if account != nil {
					expectedBalance, _ := sharedvalueobjects.NewMoney(0, sharedvalueobjects.MustCurrency("BRL"))
					if !account.Balance().Equals(expectedBalance) {
						t.Error("Account balance should not be updated if transaction save fails")
					}
				}
			},
		},
		{
			name: "atomicity: both saved on success",
			input: dtos.CreateTransactionInput{
				UserID:      userID.Value(),
				AccountID:   accountID.Value(),
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Test transaction",
				Date:        date,
			},
			setupMock: func(txRepo *mockTransactionRepository, accRepo *mockAccountRepository, uow *mockUnitOfWorkWithErrors) {
				// Create account
				currency, _ := sharedvalueobjects.NewCurrency("BRL")
				initialBalance, _ := sharedvalueobjects.NewMoney(0, currency)
				account, _ := createTestAccountWithID(userID, accountID, initialBalance)
				_ = accRepo.Save(account)
			},
			wantError: false,
			validateAtomic: func(t *testing.T, txRepo *mockTransactionRepository, accRepo *mockAccountRepository) {
				// Verify transaction was saved
				if len(txRepo.transactions) == 0 {
					t.Error("Transaction should be saved on success")
				}

				// Verify account balance was updated
				account, _ := accRepo.FindByID(accountID)
				if account != nil {
					expectedBalance, _ := sharedvalueobjects.NewMoney(10050, sharedvalueobjects.MustCurrency("BRL")) // 100.50 in cents
					if !account.Balance().Equals(expectedBalance) {
						t.Errorf("Account balance should be updated. Got %v, want %v", account.Balance(), expectedBalance)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransactionRepo := newMockTransactionRepository()
			mockAccountRepo := newMockAccountRepository()
			mockUOW := newMockUnitOfWorkWithErrors(mockTransactionRepo, mockAccountRepo)
			tt.setupMock(mockTransactionRepo, mockAccountRepo, mockUOW)

			useCase := NewCreateTransactionUseCase(mockUOW, eventBus)
			output, err := useCase.Execute(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if output == nil {
					t.Errorf("expected output but got nil")
					return
				}
			}

			if tt.validateAtomic != nil {
				tt.validateAtomic(t, mockTransactionRepo, mockAccountRepo)
			}
		})
	}
}
