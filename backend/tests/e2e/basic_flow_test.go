package e2e

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountusecases "gestao-financeira/backend/internal/account/application/usecases"
	accountpersistence "gestao-financeira/backend/internal/account/infrastructure/persistence"
	accounthandlers "gestao-financeira/backend/internal/account/presentation/handlers"
	accountroutes "gestao-financeira/backend/internal/account/presentation/routes"
	identityusecases "gestao-financeira/backend/internal/identity/application/usecases"
	identitypersistence "gestao-financeira/backend/internal/identity/infrastructure/persistence"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
	identityhandlers "gestao-financeira/backend/internal/identity/presentation/handlers"
	identityroutes "gestao-financeira/backend/internal/identity/presentation/routes"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
	transactionusecases "gestao-financeira/backend/internal/transaction/application/usecases"
	transactionpersistence "gestao-financeira/backend/internal/transaction/infrastructure/persistence"
	transactionhandlers "gestao-financeira/backend/internal/transaction/presentation/handlers"
	transactionroutes "gestao-financeira/backend/internal/transaction/presentation/routes"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestApp creates a test Fiber app with all routes configured
func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB) {
	// Create in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate all models
	err = db.AutoMigrate(
		&identitypersistence.UserModel{},
		&accountpersistence.AccountModel{},
		&transactionpersistence.TransactionModel{},
	)
	require.NoError(t, err)

	// Initialize Event Bus
	eventBus := eventbus.NewEventBus()

	// Initialize JWT Service
	jwtService := services.NewJWTService()

	// Initialize repositories
	userRepo := identitypersistence.NewGormUserRepository(db)
	accountRepo := accountpersistence.NewGormAccountRepository(db)
	transactionRepo := transactionpersistence.NewGormTransactionRepository(db)

	// Initialize use cases
	registerUserUseCase := identityusecases.NewRegisterUserUseCase(userRepo, eventBus)
	loginUseCase := identityusecases.NewLoginUseCase(userRepo, jwtService)
	createAccountUseCase := accountusecases.NewCreateAccountUseCase(accountRepo, eventBus)
	createTransactionUseCase := transactionusecases.NewCreateTransactionUseCase(
		transactionRepo,
		eventBus,
	)

	// Initialize handlers
	authHandler := identityhandlers.NewAuthHandler(registerUserUseCase, loginUseCase)
	accountHandler := accounthandlers.NewAccountHandler(
		createAccountUseCase,
		nil, // listAccountsUseCase not needed for this test
		nil, // getAccountUseCase not needed for this test
	)
	transactionHandler := transactionhandlers.NewTransactionHandler(
		createTransactionUseCase,
		nil, // listTransactionsUseCase not needed for this test
		nil, // getTransactionUseCase not needed for this test
		nil, // updateTransactionUseCase not needed for this test
		nil, // deleteTransactionUseCase not needed for this test
	)

	// Create Fiber app
	app := fiber.New()

	// Setup routes
	api := app.Group("/api/v1")
	identityroutes.SetupAuthRoutes(api, authHandler)
	accountroutes.SetupAccountRoutes(api, accountHandler, jwtService)
	transactionroutes.SetupTransactionRoutes(api, transactionHandler, jwtService)

	return app, db
}

// TestE2E_BasicFlow tests the complete flow: Register → Login → Create Account → Create Transaction
func TestE2E_BasicFlow(t *testing.T) {
	app, _ := setupTestApp(t)

	// Step 1: Register a new user
	registerPayload := map[string]interface{}{
		"email":      "e2e@example.com",
		"password":   "password123",
		"first_name": "Teste",
		"last_name":  "Usuario",
	}
	registerBody, _ := json.Marshal(registerPayload)

	registerReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp, err := app.Test(registerReq)
	require.NoError(t, err)

	if registerResp.StatusCode != fiber.StatusCreated {
		var errorResult map[string]interface{}
		json.NewDecoder(registerResp.Body).Decode(&errorResult)
		t.Fatalf("Registration failed with status %d: %v", registerResp.StatusCode, errorResult)
	}

	var registerResult map[string]interface{}
	err = json.NewDecoder(registerResp.Body).Decode(&registerResult)
	require.NoError(t, err)

	// user_id might be in "user_id" or "data.user_id" field
	var userID string
	if id, ok := registerResult["user_id"].(string); ok {
		userID = id
	} else if data, ok := registerResult["data"].(map[string]interface{}); ok {
		if id, ok := data["user_id"].(string); ok {
			userID = id
		}
	}
	assert.NotEmpty(t, userID, "user_id should be present in register response")

	// Step 2: Login to get JWT token
	loginPayload := map[string]interface{}{
		"email":    "e2e@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := app.Test(loginReq)
	require.NoError(t, err)

	if loginResp.StatusCode != fiber.StatusOK {
		var errorResult map[string]interface{}
		json.NewDecoder(loginResp.Body).Decode(&errorResult)
		t.Fatalf("Login failed with status %d: %v", loginResp.StatusCode, errorResult)
	}

	var loginResult map[string]interface{}
	err = json.NewDecoder(loginResp.Body).Decode(&loginResult)
	require.NoError(t, err)

	// Token is in "data.token" field
	data, ok := loginResult["data"].(map[string]interface{})
	require.True(t, ok, "data should be present in login response")

	token, ok := data["token"].(string)
	require.True(t, ok, "token should be present in data")
	require.NotEmpty(t, token)

	// Step 3: Create an account (protected route)
	accountPayload := map[string]interface{}{
		"name":            "Conta E2E",
		"type":            "BANK",
		"initial_balance": 1000.00,
		"currency":        "BRL",
		"context":         "PERSONAL",
	}
	accountBody, _ := json.Marshal(accountPayload)

	accountReq := httptest.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(accountBody))
	accountReq.Header.Set("Content-Type", "application/json")
	accountReq.Header.Set("Authorization", "Bearer "+token)
	accountResp, err := app.Test(accountReq)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, accountResp.StatusCode)

	var accountResult map[string]interface{}
	err = json.NewDecoder(accountResp.Body).Decode(&accountResult)
	require.NoError(t, err)

	// Account data might be in "data" field or directly in response
	var accountData map[string]interface{}
	if data, ok := accountResult["data"].(map[string]interface{}); ok {
		accountData = data
	} else {
		accountData = accountResult
	}
	require.NotEmpty(t, accountData)

	accountID, ok := accountData["account_id"].(string)
	if !ok {
		t.Logf("Account response: %+v", accountResult)
		t.Fatal("account_id should be present in account response")
	}
	require.NotEmpty(t, accountID)

	// Step 4: Create a transaction (protected route)
	transactionPayload := map[string]interface{}{
		"account_id":  accountID,
		"type":        "EXPENSE",
		"amount":      100.00,
		"currency":    "BRL",
		"description": "Teste E2E",
		"date":        "2025-12-27",
	}
	transactionBody, _ := json.Marshal(transactionPayload)

	transactionReq := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(transactionBody))
	transactionReq.Header.Set("Content-Type", "application/json")
	transactionReq.Header.Set("Authorization", "Bearer "+token)
	transactionResp, err := app.Test(transactionReq)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, transactionResp.StatusCode)

	var transactionResult map[string]interface{}
	err = json.NewDecoder(transactionResp.Body).Decode(&transactionResult)
	require.NoError(t, err)

	// Transaction data might be in "data" field or directly in response
	var transactionData map[string]interface{}
	if data, ok := transactionResult["data"].(map[string]interface{}); ok {
		transactionData = data
	} else {
		transactionData = transactionResult
	}
	require.NotEmpty(t, transactionData)

	assert.Equal(t, "EXPENSE", transactionData["type"])
	// Amount might be float64 or string, check both
	amount := transactionData["amount"]
	if amountFloat, ok := amount.(float64); ok {
		assert.Equal(t, 100.00, amountFloat)
	} else if amountStr, ok := amount.(string); ok {
		assert.Equal(t, "100.00", amountStr)
	} else {
		t.Errorf("Unexpected amount type: %T", amount)
	}
	assert.Equal(t, "Teste E2E", transactionData["description"])
}

// TestE2E_UnauthorizedAccess tests that protected routes require authentication
func TestE2E_UnauthorizedAccess(t *testing.T) {
	app, _ := setupTestApp(t)

	// Try to access protected route without token
	req := httptest.NewRequest("GET", "/api/v1/accounts", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestE2E_InvalidToken tests that invalid tokens are rejected
func TestE2E_InvalidToken(t *testing.T) {
	app, _ := setupTestApp(t)

	// Try to access protected route with invalid token
	req := httptest.NewRequest("GET", "/api/v1/accounts", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
