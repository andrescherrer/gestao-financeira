package usecases

import (
	"testing"

	accountentities "gestao-financeira/backend/internal/account/domain/entities"
	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/investment/application/dtos"
	"gestao-financeira/backend/internal/investment/domain/entities"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// mockInvestmentRepository is a mock implementation of InvestmentRepository for testing.
type mockInvestmentRepository struct {
	investments map[string]*entities.Investment
	saveErr     error
}

func newMockInvestmentRepository() *mockInvestmentRepository {
	return &mockInvestmentRepository{
		investments: make(map[string]*entities.Investment),
	}
}

func (m *mockInvestmentRepository) FindByID(id investmentvalueobjects.InvestmentID) (*entities.Investment, error) {
	investment, exists := m.investments[id.Value()]
	if !exists {
		return nil, nil
	}
	return investment, nil
}

func (m *mockInvestmentRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Investment, error) {
	var result []*entities.Investment
	for _, investment := range m.investments {
		if investment.UserID().Value() == userID.Value() {
			result = append(result, investment)
		}
	}
	return result, nil
}

func (m *mockInvestmentRepository) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Investment, error) {
	var result []*entities.Investment
	for _, investment := range m.investments {
		if investment.AccountID().Value() == accountID.Value() {
			result = append(result, investment)
		}
	}
	return result, nil
}

func (m *mockInvestmentRepository) FindByType(userID identityvalueobjects.UserID, investmentType investmentvalueobjects.InvestmentType) ([]*entities.Investment, error) {
	var result []*entities.Investment
	for _, investment := range m.investments {
		if investment.UserID().Value() == userID.Value() && investment.InvestmentType().Value() == investmentType.Value() {
			result = append(result, investment)
		}
	}
	return result, nil
}

func (m *mockInvestmentRepository) Save(investment *entities.Investment) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.investments[investment.ID().Value()] = investment
	return nil
}

func (m *mockInvestmentRepository) Delete(id investmentvalueobjects.InvestmentID) error {
	delete(m.investments, id.Value())
	return nil
}

func (m *mockInvestmentRepository) Exists(id investmentvalueobjects.InvestmentID) (bool, error) {
	_, exists := m.investments[id.Value()]
	return exists, nil
}

func (m *mockInvestmentRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, investment := range m.investments {
		if investment.UserID().Value() == userID.Value() {
			count++
		}
	}
	return count, nil
}

func (m *mockInvestmentRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	investmentType string,
	offset, limit int,
) ([]*entities.Investment, int64, error) {
	var result []*entities.Investment
	for _, investment := range m.investments {
		if investment.UserID().Value() == userID.Value() {
			if context != "" && investment.Context().Value() != context {
				continue
			}
			if investmentType != "" && investment.InvestmentType().Value() != investmentType {
				continue
			}
			result = append(result, investment)
		}
	}
	total := int64(len(result))
	// Simple pagination
	start := offset
	end := offset + limit
	if start > len(result) {
		return []*entities.Investment{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

// mockAccountRepository is a mock implementation of AccountRepository for testing.
type mockAccountRepository struct {
	accounts map[string]*accountentities.Account
}

func newMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{
		accounts: make(map[string]*accountentities.Account),
	}
}

func (m *mockAccountRepository) FindByID(id accountvalueobjects.AccountID) (*accountentities.Account, error) {
	account, exists := m.accounts[id.Value()]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (m *mockAccountRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*accountentities.Account, error) {
	var result []*accountentities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepository) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*accountentities.Account, error) {
	var result []*accountentities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() && account.Context().Value() == context.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepository) Save(account *accountentities.Account) error {
	m.accounts[account.ID().Value()] = account
	return nil
}

func (m *mockAccountRepository) Delete(id accountvalueobjects.AccountID) error {
	delete(m.accounts, id.Value())
	return nil
}

func (m *mockAccountRepository) Exists(id accountvalueobjects.AccountID) (bool, error) {
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

func (m *mockAccountRepository) FindByUserIDWithPagination(
	userID identityvalueobjects.UserID,
	context string,
	offset, limit int,
) ([]*accountentities.Account, int64, error) {
	var result []*accountentities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			if context != "" && account.Context().Value() != context {
				continue
			}
			result = append(result, account)
		}
	}
	total := int64(len(result))
	start := offset
	end := offset + limit
	if start > len(result) {
		return []*accountentities.Account{}, total, nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func TestCreateInvestmentUseCase_Execute(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	investmentRepo := newMockInvestmentRepository()
	accountRepo := newMockAccountRepository()

	// Create a test account
	userID := identityvalueobjects.GenerateUserID()
	accountName, _ := accountvalueobjects.NewAccountName("Conta Corrente")
	accountType := accountvalueobjects.BankType()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(0, currency)
	context := sharedvalueobjects.PersonalContext()
	account, _ := accountentities.NewAccount(userID, accountName, accountType, balance, context)
	accountRepo.Save(account)

	useCase := NewCreateInvestmentUseCase(
		investmentRepo,
		accountRepo,
		eventBus,
	)

	ticker := "PETR4"
	tests := []struct {
		name      string
		input     dtos.CreateInvestmentInput
		wantError bool
	}{
		{
			name: "valid investment with quantity",
			input: dtos.CreateInvestmentInput{
				UserID:         userID.Value(),
				AccountID:      account.ID().Value(),
				Type:           "STOCK",
				Name:           "Petrobras",
				Ticker:         &ticker,
				PurchaseDate:   "2024-01-15",
				PurchaseAmount: 1000.0,
				Currency:       "BRL",
				Quantity:       floatPtr(100.0),
				Context:        "PERSONAL",
			},
			wantError: false,
		},
		{
			name: "valid investment without quantity (CDB)",
			input: dtos.CreateInvestmentInput{
				UserID:         userID.Value(),
				AccountID:      account.ID().Value(),
				Type:           "CDB",
				Name:           "CDB Banco XYZ",
				PurchaseDate:   "2024-01-15",
				PurchaseAmount: 1000.0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			wantError: false,
		},
		{
			name: "invalid user ID",
			input: dtos.CreateInvestmentInput{
				UserID:         "invalid-uuid",
				AccountID:      account.ID().Value(),
				Type:           "STOCK",
				Name:           "Petrobras",
				PurchaseDate:   "2024-01-15",
				PurchaseAmount: 1000.0,
				Currency:       "BRL",
				Quantity:       floatPtr(100.0),
				Context:        "PERSONAL",
			},
			wantError: true,
		},
		{
			name: "account not found",
			input: dtos.CreateInvestmentInput{
				UserID:         userID.Value(),
				AccountID:      accountvalueobjects.GenerateAccountID().Value(),
				Type:           "STOCK",
				Name:           "Petrobras",
				PurchaseDate:   "2024-01-15",
				PurchaseAmount: 1000.0,
				Currency:       "BRL",
				Quantity:       floatPtr(100.0),
				Context:        "PERSONAL",
			},
			wantError: true,
		},
		{
			name: "stock without quantity",
			input: dtos.CreateInvestmentInput{
				UserID:         userID.Value(),
				AccountID:      account.ID().Value(),
				Type:           "STOCK",
				Name:           "Petrobras",
				PurchaseDate:   "2024-01-15",
				PurchaseAmount: 1000.0,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := useCase.Execute(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("Execute() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if output == nil {
					t.Error("Execute() returned nil output for valid input")
					return
				}

				if output.InvestmentID == "" {
					t.Error("Execute() returned empty investment ID")
				}

				// Verify investment was saved
				investmentID, _ := investmentvalueobjects.NewInvestmentID(output.InvestmentID)
				saved, _ := investmentRepo.FindByID(investmentID)
				if saved == nil {
					t.Error("Execute() failed to save investment")
				}
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
