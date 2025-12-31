package entities

import (
	"testing"
	"time"

	accountvalueobjects "gestao-financeira/backend/internal/account/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	investmentvalueobjects "gestao-financeira/backend/internal/investment/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

func TestNewInvestment(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investmentType := investmentvalueobjects.StockType()
	name, _ := investmentvalueobjects.NewInvestmentName("Petrobras", stringPtr("PETR4"))
	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	purchaseAmount, _ := sharedvalueobjects.NewMoneyFromFloat(1000.0, sharedvalueobjects.MustCurrency("BRL"))
	quantity := floatPtr(100.0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	tests := []struct {
		name           string
		userID         identityvalueobjects.UserID
		accountID      accountvalueobjects.AccountID
		investmentType investmentvalueobjects.InvestmentType
		invName        investmentvalueobjects.InvestmentName
		purchaseDate   time.Time
		purchaseAmount sharedvalueobjects.Money
		quantity       *float64
		context        sharedvalueobjects.AccountContext
		wantError      bool
	}{
		{
			"valid investment with quantity",
			userID, accountID, investmentType, name, purchaseDate, purchaseAmount, quantity, context,
			false,
		},
		{
			"valid investment without quantity (CDB)",
			userID, accountID, investmentvalueobjects.CDBType(), name, purchaseDate, purchaseAmount, nil, context,
			false,
		},
		{
			"empty user ID",
			identityvalueobjects.UserID{}, accountID, investmentType, name, purchaseDate, purchaseAmount, quantity, context,
			true,
		},
		{
			"empty account ID",
			userID, accountvalueobjects.AccountID{}, investmentType, name, purchaseDate, purchaseAmount, quantity, context,
			true,
		},
		{
			"empty name",
			userID, accountID, investmentType, investmentvalueobjects.InvestmentName{}, purchaseDate, purchaseAmount, quantity, context,
			true,
		},
		{
			"zero purchase date",
			userID, accountID, investmentType, name, time.Time{}, purchaseAmount, quantity, context,
			true,
		},
		{
			"future purchase date",
			userID, accountID, investmentType, name, time.Now().Add(24 * time.Hour), purchaseAmount, quantity, context,
			true,
		},
		{
			"zero purchase amount",
			userID, accountID, investmentType, name, purchaseDate, sharedvalueobjects.Zero(sharedvalueobjects.MustCurrency("BRL")), quantity, context,
			true,
		},
		{
			"negative purchase amount",
			userID, accountID, investmentType, name, purchaseDate, mustMoney(-100.0, "BRL"), quantity, context,
			true,
		},
		{
			"stock without quantity",
			userID, accountID, investmentvalueobjects.StockType(), name, purchaseDate, purchaseAmount, nil, context,
			true,
		},
		{
			"stock with zero quantity",
			userID, accountID, investmentvalueobjects.StockType(), name, purchaseDate, purchaseAmount, floatPtr(0.0), context,
			true,
		},
		{
			"stock with negative quantity",
			userID, accountID, investmentvalueobjects.StockType(), name, purchaseDate, purchaseAmount, floatPtr(-10.0), context,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investment, err := NewInvestment(
				tt.userID,
				tt.accountID,
				tt.investmentType,
				tt.invName,
				tt.purchaseDate,
				tt.purchaseAmount,
				tt.quantity,
				tt.context,
			)
			if (err != nil) != tt.wantError {
				t.Errorf("NewInvestment() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if investment == nil {
					t.Error("NewInvestment() returned nil for valid input")
				}
				if investment.ID().IsEmpty() {
					t.Error("NewInvestment() returned investment with empty ID")
				}
				if len(investment.GetEvents()) == 0 {
					t.Error("NewInvestment() should emit domain events")
				}
			}
		})
	}
}

func TestInvestment_UpdateCurrentValue(t *testing.T) {
	// Helper function tested separately

	tests := []struct {
		name      string
		value     sharedvalueobjects.Money
		wantError bool
	}{
		{
			"valid positive value",
			mustMoney(1000.0, "BRL"),
			false,
		},
		{
			"valid higher value",
			mustMoney(1500.0, "BRL"),
			false,
		},
		{
			"valid lower value",
			mustMoney(800.0, "BRL"),
			false,
		},
		{
			"negative value",
			mustMoney(-100.0, "BRL"),
			true,
		},
		{
			"different currency",
			mustMoney(1000.0, "USD"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investment := createTestInvestment(t)
			err := investment.UpdateCurrentValue(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("Investment.UpdateCurrentValue() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if !investment.CurrentValue().Equals(tt.value) {
					t.Errorf("Investment.UpdateCurrentValue() currentValue = %v, want %v", investment.CurrentValue(), tt.value)
				}
				// Should emit events
				events := investment.GetEvents()
				if len(events) == 0 {
					t.Error("Investment.UpdateCurrentValue() should emit domain events")
				}
			}
		})
	}
}

func TestInvestment_CalculateReturn(t *testing.T) {
	investment := createTestInvestment(t)

	// Update to higher value
	newValue := mustMoney(1200.0, "BRL")
	err := investment.UpdateCurrentValue(newValue)
	if err != nil {
		t.Fatalf("Failed to update current value: %v", err)
	}

	returnObj := investment.CalculateReturn()

	if !returnObj.IsPositive() {
		t.Error("Investment.CalculateReturn() should return positive for gain")
	}

	if returnObj.Percentage() != 20.0 {
		t.Errorf("Investment.CalculateReturn() percentage = %v, want 20.0", returnObj.Percentage())
	}

	// Update to lower value
	newValue2 := mustMoney(800.0, "BRL")
	err = investment.UpdateCurrentValue(newValue2)
	if err != nil {
		t.Fatalf("Failed to update current value: %v", err)
	}

	returnObj2 := investment.CalculateReturn()

	if !returnObj2.IsNegative() {
		t.Error("Investment.CalculateReturn() should return negative for loss")
	}
}

func TestInvestment_AddQuantity(t *testing.T) {
	// Helper function tested separately

	tests := []struct {
		name      string
		quantity  float64
		wantError bool
	}{
		{"valid positive quantity", 50.0, false},
		{"zero quantity", 0.0, true},
		{"negative quantity", -10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investment := createTestInvestment(t)
			err := investment.AddQuantity(tt.quantity)
			if (err != nil) != tt.wantError {
				t.Errorf("Investment.AddQuantity() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}

	// Test adding to CDB (should fail)
	cdbInvestment := createTestCDBInvestment(t)
	err := cdbInvestment.AddQuantity(10.0)
	if err == nil {
		t.Error("Investment.AddQuantity() should fail for investment type that doesn't require quantity")
	}
}

func TestInvestment_RemoveQuantity(t *testing.T) {
	tests := []struct {
		name      string
		quantity  float64
		wantError bool
	}{
		{"valid quantity", 50.0, false},
		{"remove all", 100.0, false},
		{"zero quantity", 0.0, true},
		{"negative quantity", -10.0, true},
		{"insufficient quantity", 200.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			investment := createTestInvestment(t)
			err := investment.RemoveQuantity(tt.quantity)
			if (err != nil) != tt.wantError {
				t.Errorf("Investment.RemoveQuantity() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}

	// Test removing from CDB (should fail)
	cdbInvestment := createTestCDBInvestment(t)
	err := cdbInvestment.RemoveQuantity(10.0)
	if err == nil {
		t.Error("Investment.RemoveQuantity() should fail for investment type that doesn't require quantity")
	}
}

func TestInvestmentFromPersistence(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investmentID := investmentvalueobjects.GenerateInvestmentID()
	investmentType := investmentvalueobjects.StockType()
	name, _ := investmentvalueobjects.NewInvestmentName("Petrobras", stringPtr("PETR4"))
	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	purchaseAmount := mustMoney(1000.0, "BRL")
	currentValue := mustMoney(1200.0, "BRL")
	quantity := floatPtr(100.0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")
	createdAt := time.Now()
	updatedAt := time.Now()

	investment, err := InvestmentFromPersistence(
		investmentID,
		userID,
		accountID,
		investmentType,
		name,
		purchaseDate,
		purchaseAmount,
		currentValue,
		quantity,
		context,
		createdAt,
		updatedAt,
	)

	if err != nil {
		t.Fatalf("InvestmentFromPersistence() error = %v", err)
	}

	if investment.ID() != investmentID {
		t.Error("InvestmentFromPersistence() ID mismatch")
	}

	if len(investment.GetEvents()) != 0 {
		t.Error("InvestmentFromPersistence() should not emit domain events")
	}
}

// Helper functions

func createTestInvestment(t *testing.T) *Investment {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investmentType := investmentvalueobjects.StockType()
	name, _ := investmentvalueobjects.NewInvestmentName("Petrobras", stringPtr("PETR4"))
	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	purchaseAmount := mustMoney(1000.0, "BRL")
	quantity := floatPtr(100.0)
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	investment, err := NewInvestment(
		userID,
		accountID,
		investmentType,
		name,
		purchaseDate,
		purchaseAmount,
		quantity,
		context,
	)
	if err != nil {
		t.Fatalf("Failed to create test investment: %v", err)
	}

	return investment
}

func createTestCDBInvestment(t *testing.T) *Investment {
	userID := identityvalueobjects.GenerateUserID()
	accountID := accountvalueobjects.GenerateAccountID()
	investmentType := investmentvalueobjects.CDBType()
	name, _ := investmentvalueobjects.NewInvestmentName("CDB Banco XYZ", nil)
	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	purchaseAmount := mustMoney(1000.0, "BRL")
	context := sharedvalueobjects.MustAccountContext("PERSONAL")

	investment, err := NewInvestment(
		userID,
		accountID,
		investmentType,
		name,
		purchaseDate,
		purchaseAmount,
		nil,
		context,
	)
	if err != nil {
		t.Fatalf("Failed to create test CDB investment: %v", err)
	}

	return investment
}

func mustMoney(amount float64, currencyCode string) sharedvalueobjects.Money {
	money, err := sharedvalueobjects.NewMoneyFromFloatString(amount, currencyCode)
	if err != nil {
		panic(err)
	}
	return money
}

func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
