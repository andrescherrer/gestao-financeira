package usecases

import (
	"fmt"
	"time"

	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

// mockAccountRepository is a mock implementation of AccountRepository for testing.
type mockAccountRepository struct {
	accounts map[string]*entities.Account
	saveErr  error
	findErr  error
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{
		accounts: make(map[string]*entities.Account),
	}
}

func (m *mockAccountRepository) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
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
func (m *mockAccountRepository) FindByUserIDWithPagination(userID identityvalueobjects.UserID, context string, offset, limit int) ([]*entities.Account, int64, error) {
	all, _ := m.FindByUserID(userID)
	var filtered []*entities.Account
	for _, acc := range all {
		if context == "" || acc.Context().Value() == context {
			filtered = append(filtered, acc)
		}
	}
	total := int64(len(filtered))
	start := offset
	end := offset + limit
	if start > len(filtered) {
		return []*entities.Account{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

// Helper function to create a test account with a specific ID
func createTestAccountWithID(userID identityvalueobjects.UserID, accountID valueobjects.AccountID, initialBalance sharedvalueobjects.Money) (*entities.Account, error) {
	accountName, err := valueobjects.NewAccountName("Test Account")
	if err != nil {
		return nil, fmt.Errorf("failed to create account name: %w", err)
	}

	accountType := valueobjects.BankType()
	context := sharedvalueobjects.PersonalContext()

	// Use AccountFromPersistence to create account with specific ID
	account, err := entities.AccountFromPersistence(
		accountID,
		userID,
		accountName,
		accountType,
		initialBalance,
		context,
		time.Now(),
		time.Now(),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// Helper function to create a test account (with generated ID)
func createTestAccount(userID identityvalueobjects.UserID, initialBalance sharedvalueobjects.Money) (*entities.Account, error) {
	accountID := valueobjects.GenerateAccountID()
	return createTestAccountWithID(userID, accountID, initialBalance)
}
