package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/reporting/application/usecases"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
	"gestao-financeira/backend/internal/transaction/domain/entities"
	transactionvalueobjects "gestao-financeira/backend/internal/transaction/domain/valueobjects"
	"time"
)

// mockTransactionRepositoryForReports is a mock implementation for testing reports.
type mockTransactionRepositoryForReports struct {
	transactions []*entities.Transaction
}

func (m *mockTransactionRepositoryForReports) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Transaction, error) {
	var result []*entities.Transaction
	for _, tx := range m.transactions {
		if tx.UserID().Equals(userID) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// Implement other required methods with empty implementations
func (m *mockTransactionRepositoryForReports) FindByID(id transactionvalueobjects.TransactionID) (*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepositoryForReports) FindByAccountID(accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepositoryForReports) FindByUserIDAndAccountID(userID identityvalueobjects.UserID, accountID accountvalueobjects.AccountID) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepositoryForReports) FindByUserIDAndType(userID identityvalueobjects.UserID, transactionType transactionvalueobjects.TransactionType) ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepositoryForReports) Save(transaction *entities.Transaction) error {
	return nil
}
func (m *mockTransactionRepositoryForReports) Delete(id transactionvalueobjects.TransactionID) error {
	return nil
}
func (m *mockTransactionRepositoryForReports) Exists(id transactionvalueobjects.TransactionID) (bool, error) {
	return false, nil
}
func (m *mockTransactionRepositoryForReports) Count(userID identityvalueobjects.UserID) (int64, error) {
	return 0, nil
}
func (m *mockTransactionRepositoryForReports) CountByAccountID(accountID accountvalueobjects.AccountID) (int64, error) {
	return 0, nil
}
func (m *mockTransactionRepositoryForReports) FindActiveRecurringTransactions() ([]*entities.Transaction, error) {
	return nil, nil
}
func (m *mockTransactionRepositoryForReports) FindByParentIDAndDate(parentID transactionvalueobjects.TransactionID, date time.Time) (*entities.Transaction, error) {
	return nil, nil
}

func TestReportHandler_GetMonthlyReport(t *testing.T) {
	// Create test data
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	accountID, _ := accountvalueobjects.NewAccountID("123e4567-e89b-12d3-a456-426614174001")
	currency, _ := sharedvalueobjects.NewCurrency("BRL")

	incomeAmount, _ := sharedvalueobjects.NewMoney(100000, currency)
	expenseAmount, _ := sharedvalueobjects.NewMoney(50000, currency)

	incomeTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("INCOME"),
		incomeAmount,
		transactionvalueobjects.MustTransactionDescription("Salary"),
		time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	)

	expenseTx, _ := entities.NewTransaction(
		userID,
		accountID,
		transactionvalueobjects.MustTransactionType("EXPENSE"),
		expenseAmount,
		transactionvalueobjects.MustTransactionDescription("Groceries"),
		time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
	)

	mockRepo := &mockTransactionRepositoryForReports{
		transactions: []*entities.Transaction{incomeTx, expenseTx},
	}

	// Create use cases
	monthlyUseCase := usecases.NewMonthlyReportUseCase(mockRepo, nil) // nil cache for tests
	annualUseCase := usecases.NewAnnualReportUseCase(mockRepo)
	categoryUseCase := usecases.NewCategoryReportUseCase(mockRepo)
	incomeVsExpenseUseCase := usecases.NewIncomeVsExpenseUseCase(mockRepo)

	// Create handler
	handler := NewReportHandler(monthlyUseCase, annualUseCase, categoryUseCase, incomeVsExpenseUseCase)

	// Create Fiber app
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set user ID in context (simulating auth middleware)
		c.Locals("userID", userID.Value())
		c.Locals("request_id", "test-request-id")
		return c.Next()
	})

	// Setup route
	app.Get("/reports/monthly", handler.GetMonthlyReport)

	// Create request
	req := httptest.NewRequest("GET", "/reports/monthly?year=2025&month=1&currency=BRL", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Validate response structure
	data, ok := result["data"].(map[string]interface{})
	require.True(t, ok, "response should have 'data' field")

	assert.Equal(t, float64(2025), data["year"])
	assert.Equal(t, float64(1), data["month"])
	assert.Equal(t, "BRL", data["currency"])
	assert.Equal(t, 1000.00, data["total_income"])
	assert.Equal(t, 500.00, data["total_expense"])
}

func TestReportHandler_GetMonthlyReport_Unauthorized(t *testing.T) {
	mockRepo := &mockTransactionRepositoryForReports{transactions: []*entities.Transaction{}}
	monthlyUseCase := usecases.NewMonthlyReportUseCase(mockRepo, nil) // nil cache for tests
	handler := NewReportHandler(monthlyUseCase, nil, nil, nil)

	app := fiber.New()
	app.Get("/reports/monthly", handler.GetMonthlyReport)

	req := httptest.NewRequest("GET", "/reports/monthly?year=2025&month=1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestReportHandler_GetMonthlyReport_MissingParams(t *testing.T) {
	userID, _ := identityvalueobjects.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	mockRepo := &mockTransactionRepositoryForReports{transactions: []*entities.Transaction{}}
	monthlyUseCase := usecases.NewMonthlyReportUseCase(mockRepo, nil) // nil cache for tests
	handler := NewReportHandler(monthlyUseCase, nil, nil, nil)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", userID.Value())
		return c.Next()
	})
	app.Get("/reports/monthly", handler.GetMonthlyReport)

	// Missing year
	req := httptest.NewRequest("GET", "/reports/monthly?month=1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Missing month
	req = httptest.NewRequest("GET", "/reports/monthly?year=2025", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
