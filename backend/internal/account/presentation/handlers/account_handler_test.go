package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gestao-financeira/backend/internal/account/application/dtos"
	"gestao-financeira/backend/internal/account/application/usecases"
	"gestao-financeira/backend/internal/account/domain/entities"
	"gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// mockAccountRepositoryForHandler is a mock implementation of AccountRepository for handler testing.
type mockAccountRepositoryForHandler struct {
	accounts map[string]*entities.Account
	saveErr  error
}

func newMockAccountRepositoryForHandler() *mockAccountRepositoryForHandler {
	return &mockAccountRepositoryForHandler{
		accounts: make(map[string]*entities.Account),
	}
}

func (m *mockAccountRepositoryForHandler) FindByID(id valueobjects.AccountID) (*entities.Account, error) {
	account, exists := m.accounts[id.Value()]
	if !exists {
		return nil, nil
	}
	return account, nil
}

func (m *mockAccountRepositoryForHandler) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Account, error) {
	var result []*entities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepositoryForHandler) FindByUserIDAndContext(userID identityvalueobjects.UserID, context sharedvalueobjects.AccountContext) ([]*entities.Account, error) {
	var result []*entities.Account
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() && account.Context().Value() == context.Value() {
			result = append(result, account)
		}
	}
	return result, nil
}

func (m *mockAccountRepositoryForHandler) Save(account *entities.Account) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.accounts[account.ID().Value()] = account
	return nil
}

func (m *mockAccountRepositoryForHandler) Delete(id valueobjects.AccountID) error {
	delete(m.accounts, id.Value())
	return nil
}

func (m *mockAccountRepositoryForHandler) Exists(id valueobjects.AccountID) (bool, error) {
	_, exists := m.accounts[id.Value()]
	return exists, nil
}

func (m *mockAccountRepositoryForHandler) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, account := range m.accounts {
		if account.UserID().Value() == userID.Value() {
			count++
		}
	}
	return count, nil
}

func TestAccountHandler_Create(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	app := fiber.New()
	mockRepo := newMockAccountRepositoryForHandler()
	eventBus := eventbus.NewEventBus()
	createUseCase := usecases.NewCreateAccountUseCase(mockRepo, eventBus)
	listUseCase := usecases.NewListAccountsUseCase(mockRepo)
	getUseCase := usecases.NewGetAccountUseCase(mockRepo)
	handler := NewAccountHandler(createUseCase, listUseCase, getUseCase)

	app.Post("/accounts", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Create(c)
	})

	tests := []struct {
		name           string
		body           dtos.CreateAccountInput
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful account creation",
			body: dtos.CreateAccountInput{
				Name:           "Conta Corrente",
				Type:           "BANK",
				InitialBalance: 1000.00,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "invalid account name",
			body: dtos.CreateAccountInput{
				Name:           "AB", // Too short
				Type:           "BANK",
				InitialBalance: 1000.00,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid account data",
		},
		{
			name: "invalid account type",
			body: dtos.CreateAccountInput{
				Name:           "Conta Corrente",
				Type:           "INVALID",
				InitialBalance: 1000.00,
				Currency:       "BRL",
				Context:        "PERSONAL",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid account data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.accounts = make(map[string]*entities.Account)

			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(bodyJSON))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				if errorMsg, ok := result["error"].(string); ok {
					assert.Contains(t, errorMsg, tt.expectedError)
				}
			} else if tt.expectedStatus == fiber.StatusCreated {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "Account created successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

func TestAccountHandler_List(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	app := fiber.New()
	mockRepo := newMockAccountRepositoryForHandler()
	listUseCase := usecases.NewListAccountsUseCase(mockRepo)
	createUseCase := usecases.NewCreateAccountUseCase(mockRepo, eventbus.NewEventBus())
	getUseCase := usecases.NewGetAccountUseCase(mockRepo)
	handler := NewAccountHandler(createUseCase, listUseCase, getUseCase)

	app.Get("/accounts", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.List(c)
	})

	// Create test accounts
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountName1, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType1, _ := valueobjects.NewAccountType("BANK")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency) // 1000.00 in cents
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")
	account1, _ := entities.NewAccount(userIDVO, accountName1, accountType1, balance, context)
	mockRepo.accounts[account1.ID().Value()] = account1

	accountName2, _ := valueobjects.NewAccountName("Conta Poupança")
	account2, _ := entities.NewAccount(userIDVO, accountName2, accountType1, balance, context)
	mockRepo.accounts[account2.ID().Value()] = account2

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "list all accounts",
			queryParams:    "",
			expectedStatus: fiber.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "list with context filter",
			queryParams:    "?context=PERSONAL",
			expectedStatus: fiber.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/accounts"+tt.queryParams, nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "Accounts retrieved successfully", result["message"])
			data := result["data"].(map[string]interface{})
			accounts := data["accounts"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(accounts))
		})
	}
}

func TestAccountHandler_Get(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"

	app := fiber.New()
	mockRepo := newMockAccountRepositoryForHandler()
	getUseCase := usecases.NewGetAccountUseCase(mockRepo)
	createUseCase := usecases.NewCreateAccountUseCase(mockRepo, eventbus.NewEventBus())
	listUseCase := usecases.NewListAccountsUseCase(mockRepo)
	handler := NewAccountHandler(createUseCase, listUseCase, getUseCase)

	app.Get("/accounts/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Get(c)
	})

	// Create test account
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountName, _ := valueobjects.NewAccountName("Conta Corrente")
	accountType, _ := valueobjects.NewAccountType("BANK")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	balance, _ := sharedvalueobjects.NewMoney(100000, currency)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")
	account, _ := entities.NewAccount(userIDVO, accountName, accountType, balance, context)
	mockRepo.accounts[account.ID().Value()] = account

	tests := []struct {
		name           string
		accountID     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful get account",
			accountID:      account.ID().Value(),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "account not found",
			accountID:      "00000000-0000-0000-0000-000000000000",
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/accounts/"+tt.accountID, nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedError != "" {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				if errorMsg, ok := result["error"].(string); ok {
					assert.Contains(t, errorMsg, tt.expectedError)
				}
			} else if tt.expectedStatus == fiber.StatusOK {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "Account retrieved successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}
