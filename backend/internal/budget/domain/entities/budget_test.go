package entities

import (
	"testing"
	"time"

	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	sharedvalueobjects "gestao-financeira/backend/internal/shared/domain/valueobjects"
)

func TestNewBudget(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency) // R$ 1000.00
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	budget, err := NewBudget(userID, categoryID, amount, period, context)
	if err != nil {
		t.Fatalf("NewBudget() error = %v, want nil", err)
	}

	if budget == nil {
		t.Fatal("NewBudget() returned nil budget")
	}

	if budget.ID().IsEmpty() {
		t.Error("NewBudget() returned budget with empty ID")
	}

	if !budget.UserID().Equals(userID) {
		t.Error("NewBudget() returned budget with wrong user ID")
	}

	if !budget.CategoryID().Equals(categoryID) {
		t.Error("NewBudget() returned budget with wrong category ID")
	}

	if budget.Amount().Amount() != amount.Amount() {
		t.Errorf("NewBudget() amount = %v, want %v", budget.Amount().Amount(), amount.Amount())
	}

	if !budget.Period().Equals(period) {
		t.Error("NewBudget() returned budget with wrong period")
	}

	if !budget.Context().Equals(context) {
		t.Error("NewBudget() returned budget with wrong context")
	}

	if !budget.IsActive() {
		t.Error("NewBudget() returned inactive budget")
	}

	if budget.CreatedAt().IsZero() {
		t.Error("NewBudget() returned budget with zero created_at")
	}

	if budget.UpdatedAt().IsZero() {
		t.Error("NewBudget() returned budget with zero updated_at")
	}

	// Check domain events
	events := budget.GetEvents()
	if len(events) == 0 {
		t.Error("NewBudget() should have domain events")
	}
}

func TestNewBudget_InvalidInput(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	tests := []struct {
		name      string
		userID    identityvalueobjects.UserID
		categoryID categoryvalueobjects.CategoryID
		amount    sharedvalueobjects.Money
		period    valueobjects.BudgetPeriod
		context   sharedvalueobjects.AccountContext
		wantError bool
	}{
		{
			name:      "empty user ID",
			userID:    identityvalueobjects.UserID{},
			categoryID: categoryID,
			amount:    amount,
			period:    period,
			context:   context,
			wantError: true,
		},
		{
			name:      "empty category ID",
			userID:    userID,
			categoryID: categoryvalueobjects.CategoryID{},
			amount:    amount,
			period:    period,
			context:   context,
			wantError: true,
		},
		{
			name:      "zero amount",
			userID:    userID,
			categoryID: categoryID,
			amount:    sharedvalueobjects.Money{},
			period:    period,
			context:   context,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBudget(tt.userID, tt.categoryID, tt.amount, tt.period, tt.context)
			if tt.wantError {
				if err == nil {
					t.Errorf("NewBudget() error = nil, want error")
				}
				if got != nil {
					t.Errorf("NewBudget() = %v, want nil", got)
				}
			} else {
				if err != nil {
					t.Errorf("NewBudget() error = %v, want nil", err)
				}
				if got == nil {
					t.Error("NewBudget() = nil, want budget")
				}
			}
		})
	}
}

func TestBudgetFromPersistence(t *testing.T) {
	budgetID := valueobjects.GenerateBudgetID()
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	budget, err := BudgetFromPersistence(budgetID, userID, categoryID, amount, period, context, createdAt, updatedAt, true)
	if err != nil {
		t.Fatalf("BudgetFromPersistence() error = %v, want nil", err)
	}

	if budget == nil {
		t.Fatal("BudgetFromPersistence() returned nil budget")
	}

	if !budget.ID().Equals(budgetID) {
		t.Error("BudgetFromPersistence() returned budget with wrong ID")
	}

	if !budget.UserID().Equals(userID) {
		t.Error("BudgetFromPersistence() returned budget with wrong user ID")
	}

	if budget.IsActive() != true {
		t.Error("BudgetFromPersistence() returned budget with wrong isActive")
	}

	// BudgetFromPersistence should not have domain events
	events := budget.GetEvents()
	if len(events) != 0 {
		t.Error("BudgetFromPersistence() should not have domain events")
	}
}

func TestBudget_UpdateAmount(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	budget, _ := NewBudget(userID, categoryID, amount, period, context)
	originalUpdatedAt := budget.UpdatedAt()

	// Wait a bit to ensure updatedAt changes
	time.Sleep(10 * time.Millisecond)

	newAmount, _ := sharedvalueobjects.NewMoney(200000, currency)
	err := budget.UpdateAmount(newAmount)
	if err != nil {
		t.Fatalf("UpdateAmount() error = %v, want nil", err)
	}

	if budget.Amount().Amount() != newAmount.Amount() {
		t.Errorf("UpdateAmount() amount = %v, want %v", budget.Amount().Amount(), newAmount.Amount())
	}

	if !budget.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdateAmount() should update updatedAt timestamp")
	}

	// Check domain events
	events := budget.GetEvents()
	if len(events) == 0 {
		t.Error("UpdateAmount() should add domain event")
	}
}

func TestBudget_UpdateAmount_Invalid(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	budget, _ := NewBudget(userID, categoryID, amount, period, context)

	// Test negative amount
	negativeAmount, _ := sharedvalueobjects.NewMoney(-10000, currency)
	err := budget.UpdateAmount(negativeAmount)
	if err == nil {
		t.Error("UpdateAmount() with negative amount should return error")
	}

	// Test zero amount
	zeroAmount, _ := sharedvalueobjects.NewMoney(0, currency)
	err = budget.UpdateAmount(zeroAmount)
	if err == nil {
		t.Error("UpdateAmount() with zero amount should return error")
	}

	// Test different currency
	usdCurrency, _ := sharedvalueobjects.NewCurrency("USD")
	differentCurrencyAmount, _ := sharedvalueobjects.NewMoney(100000, usdCurrency)
	err = budget.UpdateAmount(differentCurrencyAmount)
	if err == nil {
		t.Error("UpdateAmount() with different currency should return error")
	}
}

func TestBudget_UpdatePeriod(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	budget, _ := NewBudget(userID, categoryID, amount, period, context)
	originalUpdatedAt := budget.UpdatedAt()

	// Wait a bit to ensure updatedAt changes
	time.Sleep(10 * time.Millisecond)

	newPeriod, _ := valueobjects.NewMonthlyBudgetPeriod(2026, 1)
	err := budget.UpdatePeriod(newPeriod)
	if err != nil {
		t.Fatalf("UpdatePeriod() error = %v, want nil", err)
	}

	if !budget.Period().Equals(newPeriod) {
		t.Error("UpdatePeriod() period not updated correctly")
	}

	if !budget.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatePeriod() should update updatedAt timestamp")
	}
}

func TestBudget_ActivateDeactivate(t *testing.T) {
	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()
	currency, _ := sharedvalueobjects.NewCurrency("BRL")
	amount, _ := sharedvalueobjects.NewMoney(100000, currency)
	period, _ := valueobjects.NewMonthlyBudgetPeriod(2025, 12)
	context, _ := sharedvalueobjects.NewAccountContext("PERSONAL")

	budget, _ := NewBudget(userID, categoryID, amount, period, context)

	if !budget.IsActive() {
		t.Error("New budget should be active")
	}

	// Deactivate
	err := budget.Deactivate()
	if err != nil {
		t.Fatalf("Deactivate() error = %v, want nil", err)
	}

	if budget.IsActive() {
		t.Error("Budget should be inactive after Deactivate()")
	}

	// Try to deactivate again
	err = budget.Deactivate()
	if err == nil {
		t.Error("Deactivate() on inactive budget should return error")
	}

	// Activate
	err = budget.Activate()
	if err != nil {
		t.Fatalf("Activate() error = %v, want nil", err)
	}

	if !budget.IsActive() {
		t.Error("Budget should be active after Activate()")
	}

	// Try to activate again
	err = budget.Activate()
	if err == nil {
		t.Error("Activate() on active budget should return error")
	}
}

