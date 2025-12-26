package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	"gestao-financeira/backend/internal/transaction/application/dtos"
	"gestao-financeira/backend/internal/transaction/application/usecases"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
)

// Helper function to create a test transaction (from usecase tests)
func createTestTransactionForHandler(
	userID identityvalueobjects.UserID,
	accountID accountvalueobjects.AccountID,
	transactionType string,
	amount float64,
	currency string,
	description string,
	date time.Time,
) (*entities.Transaction, error) {
	typeVO, _ := transactionvalueobjects.NewTransactionType(transactionType)
	currencyVO, _ := sharedvalueobjects.NewCurrency(currency)
	amountCents := int64(amount * 100)
	money, _ := sharedvalueobjects.NewMoney(amountCents, currencyVO)
	descriptionVO, _ := transactionvalueobjects.NewTransactionDescription(description)

	return entities.NewTransaction(userID, accountID, typeVO, money, descriptionVO, date)
}

// mockTransactionRepositoryForHandler is a mock implementation of TransactionRepository for handler testing.
type mockTransactionRepositoryForHandler struct {
	transactions map[string]*entities.Transaction
	saveErr      error
}

func newMockTransactionRepositoryForHandler() *mockTransactionRepositoryForHandler {
	return &mockTransactionRepositoryForHandler{
		transactions: make(map[string]*entities.Transaction),
	}
}

func (m *mockTransactionRepositoryForHandler) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	transaction, exists := m.transactions[id.Value()]
	if !exists {
		return nil, nil
	}
	return transaction, nil
}

func (m *mockTransactionRepositoryForHandler) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepositoryForHandler) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.AccountID().Value() == accountID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepositoryForHandler) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() && transaction.AccountID().Value() == accountID.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepositoryForHandler) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() && transaction.TransactionType().Value() == transactionType.Value() {
			result = append(result, transaction)
		}
	}
	return result, nil
}

func (m *mockTransactionRepositoryForHandler) Save(transaction *entities.Transaction) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.transactions[transaction.ID().Value()] = transaction
	return nil
}

func (m *mockTransactionRepositoryForHandler) Delete(id transactionvalueobjects.TransactionID) error {
	delete(m.transactions, id.Value())
	return nil
}

func (m *mockTransactionRepositoryForHandler) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	_, exists := m.transactions[id.Value()]
	return exists, nil
}

func (m *mockTransactionRepositoryForHandler) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, transaction := range m.transactions {
		if transaction.UserID().Value() == userID.Value() {
			count++
		}
	}
	return count, nil
}

func (m *mockTransactionRepositoryForHandler) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	count := int64(0)
	for _, transaction := range m.transactions {
		if transaction.AccountID().Value() == accountID.Value() {
			count++
		}
	}
	return count, nil
}

func TestTransactionHandler_Create(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := accountvalueobjects.GenerateAccountID().Value()
	date := time.Now().Format("2006-01-02")

	app := fiber.New()
	mockRepo := newMockTransactionRepositoryForHandler()
	eventBus := eventbus.NewEventBus()
	createUseCase := usecases.NewCreateTransactionUseCase(mockRepo, eventBus)
	listUseCase := usecases.NewListTransactionsUseCase(mockRepo)
	getUseCase := usecases.NewGetTransactionUseCase(mockRepo)
	updateUseCase := usecases.NewUpdateTransactionUseCase(mockRepo, eventBus)
	deleteUseCase := usecases.NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
	handler := NewTransactionHandler(createUseCase, listUseCase, getUseCase, updateUseCase, deleteUseCase)

	app.Post("/transactions", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Create(c)
	})

	tests := []struct {
		name           string
		body           dtos.CreateTransactionInput
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful transaction creation",
			body: dtos.CreateTransactionInput{
				AccountID:   accountID,
				Type:        "INCOME",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Salário",
				Date:        date,
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "invalid transaction type",
			body: dtos.CreateTransactionInput{
				AccountID:   accountID,
				Type:        "INVALID",
				Amount:      100.50,
				Currency:    "BRL",
				Description: "Test",
				Date:        date,
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Invalid transaction data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.transactions = make(map[string]*entities.Transaction)

			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(bodyJSON))
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
				assert.Equal(t, "Transaction created successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

func TestTransactionHandler_List(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := accountvalueobjects.GenerateAccountID().Value()

	app := fiber.New()
	mockRepo := newMockTransactionRepositoryForHandler()
	listUseCase := usecases.NewListTransactionsUseCase(mockRepo)
	createUseCase := usecases.NewCreateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	getUseCase := usecases.NewGetTransactionUseCase(mockRepo)
	updateUseCase := usecases.NewUpdateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	deleteUseCase := usecases.NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
	handler := NewTransactionHandler(createUseCase, listUseCase, getUseCase, updateUseCase, deleteUseCase)

	app.Get("/transactions", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.List(c)
	})

	// Create test transactions
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountIDVO, _ := accountvalueobjects.NewAccountID(accountID)
	date := time.Now()
	transaction1, _ := createTestTransactionForHandler(userIDVO, accountIDVO, "INCOME", 100.0, "BRL", "Test 1", date)
	mockRepo.transactions[transaction1.ID().Value()] = transaction1

	transaction2, _ := createTestTransactionForHandler(userIDVO, accountIDVO, "EXPENSE", 50.0, "BRL", "Test 2", date)
	mockRepo.transactions[transaction2.ID().Value()] = transaction2

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "list all transactions",
			queryParams:    "",
			expectedStatus: fiber.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "list with account filter",
			queryParams:    "?account_id=" + accountID,
			expectedStatus: fiber.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/transactions"+tt.queryParams, nil)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "Transactions retrieved successfully", result["message"])
			data := result["data"].(map[string]interface{})
			transactions := data["transactions"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(transactions))
		})
	}
}

func TestTransactionHandler_Get(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := accountvalueobjects.GenerateAccountID().Value()

	app := fiber.New()
	mockRepo := newMockTransactionRepositoryForHandler()
	getUseCase := usecases.NewGetTransactionUseCase(mockRepo)
	createUseCase := usecases.NewCreateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	listUseCase := usecases.NewListTransactionsUseCase(mockRepo)
	updateUseCase := usecases.NewUpdateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	deleteUseCase := usecases.NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
	handler := NewTransactionHandler(createUseCase, listUseCase, getUseCase, updateUseCase, deleteUseCase)

	app.Get("/transactions/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Get(c)
	})

	// Create test transaction
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountIDVO, _ := accountvalueobjects.NewAccountID(accountID)
	date := time.Now()
	transaction, _ := createTestTransactionForHandler(userIDVO, accountIDVO, "INCOME", 100.0, "BRL", "Test", date)
	mockRepo.transactions[transaction.ID().Value()] = transaction

	tests := []struct {
		name           string
		transactionID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful get transaction",
			transactionID:  transaction.ID().Value(),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "transaction not found",
			transactionID:  "00000000-0000-0000-0000-000000000000",
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Transaction not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/transactions/"+tt.transactionID, nil)

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
				assert.Equal(t, "Transaction retrieved successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

func TestTransactionHandler_Update(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := accountvalueobjects.GenerateAccountID().Value()

	app := fiber.New()
	mockRepo := newMockTransactionRepositoryForHandler()
	updateUseCase := usecases.NewUpdateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	createUseCase := usecases.NewCreateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	listUseCase := usecases.NewListTransactionsUseCase(mockRepo)
	getUseCase := usecases.NewGetTransactionUseCase(mockRepo)
	deleteUseCase := usecases.NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
	handler := NewTransactionHandler(createUseCase, listUseCase, getUseCase, updateUseCase, deleteUseCase)

	app.Put("/transactions/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Update(c)
	})

	// Create test transaction
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountIDVO, _ := accountvalueobjects.NewAccountID(accountID)
	date := time.Now()
	transaction, _ := createTestTransactionForHandler(userIDVO, accountIDVO, "INCOME", 100.0, "BRL", "Test", date)
	mockRepo.transactions[transaction.ID().Value()] = transaction

	tests := []struct {
		name           string
		transactionID  string
		body           dtos.UpdateTransactionInput
		expectedStatus int
		expectedError  string
	}{
		{
			name:          "successful transaction update",
			transactionID: transaction.ID().Value(),
			body: dtos.UpdateTransactionInput{
				Type: stringPtr("EXPENSE"),
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:          "transaction not found",
			transactionID: "00000000-0000-0000-0000-000000000000",
			body: dtos.UpdateTransactionInput{
				Type: stringPtr("EXPENSE"),
			},
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Transaction not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bodyJSON, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/transactions/"+tt.transactionID, bytes.NewBuffer(bodyJSON))
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
			} else if tt.expectedStatus == fiber.StatusOK {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "Transaction updated successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

func TestTransactionHandler_Delete(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := accountvalueobjects.GenerateAccountID().Value()

	app := fiber.New()
	mockRepo := newMockTransactionRepositoryForHandler()
	deleteUseCase := usecases.NewDeleteTransactionUseCase(mockRepo, eventbus.NewEventBus())
	createUseCase := usecases.NewCreateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	listUseCase := usecases.NewListTransactionsUseCase(mockRepo)
	getUseCase := usecases.NewGetTransactionUseCase(mockRepo)
	updateUseCase := usecases.NewUpdateTransactionUseCase(mockRepo, eventbus.NewEventBus())
	handler := NewTransactionHandler(createUseCase, listUseCase, getUseCase, updateUseCase, deleteUseCase)

	app.Delete("/transactions/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", userID)
		return handler.Delete(c)
	})

	// Create test transaction
	userIDVO, _ := identityvalueobjects.NewUserID(userID)
	accountIDVO, _ := accountvalueobjects.NewAccountID(accountID)
	date := time.Now()
	transaction, _ := createTestTransactionForHandler(userIDVO, accountIDVO, "INCOME", 100.0, "BRL", "Test", date)
	mockRepo.transactions[transaction.ID().Value()] = transaction

	tests := []struct {
		name           string
		transactionID  string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "successful transaction deletion",
			transactionID:  transaction.ID().Value(),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "transaction not found",
			transactionID:  "00000000-0000-0000-0000-000000000000",
			expectedStatus: fiber.StatusNotFound,
			expectedError:  "Transaction not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest("DELETE", "/transactions/"+tt.transactionID, nil)

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
				assert.Equal(t, "Transaction deleted successfully", result["message"])
				assert.NotNil(t, result["data"])
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
