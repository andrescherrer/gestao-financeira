package usecases

import (
	"testing"

	"gestao-financeira/backend/internal/budget/application/dtos"
	"gestao-financeira/backend/internal/budget/domain/entities"
	"gestao-financeira/backend/internal/budget/domain/valueobjects"
	categoryvalueobjects "gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// mockBudgetRepository is a mock implementation of BudgetRepository for testing.
type mockBudgetRepository struct {
	budgets map[string]*entities.Budget
}

func newMockBudgetRepository() *mockBudgetRepository {
	return &mockBudgetRepository{
		budgets: make(map[string]*entities.Budget),
	}
}

func (m *mockBudgetRepository) FindByID(id valueobjects.BudgetID) (*entities.Budget, error) {
	budget, exists := m.budgets[id.Value()]
	if !exists {
		return nil, nil
	}
	return budget, nil
}

func (m *mockBudgetRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Budget, error) {
	var result []*entities.Budget
	for _, budget := range m.budgets {
		if budget.UserID().Equals(userID) {
			result = append(result, budget)
		}
	}
	return result, nil
}

func (m *mockBudgetRepository) FindByCategoryAndPeriod(categoryID categoryvalueobjects.CategoryID, period valueobjects.BudgetPeriod) (*entities.Budget, error) {
	for _, budget := range m.budgets {
		if budget.CategoryID().Equals(categoryID) && budget.Period().Equals(period) {
			return budget, nil
		}
	}
	return nil, nil
}

func (m *mockBudgetRepository) Save(budget *entities.Budget) error {
	m.budgets[budget.ID().Value()] = budget
	return nil
}

func (m *mockBudgetRepository) Delete(id valueobjects.BudgetID) error {
	delete(m.budgets, id.Value())
	return nil
}

func (m *mockBudgetRepository) Exists(id valueobjects.BudgetID) (bool, error) {
	_, exists := m.budgets[id.Value()]
	return exists, nil
}

func (m *mockBudgetRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, budget := range m.budgets {
		if budget.UserID().Equals(userID) {
			count++
		}
	}
	return count, nil
}

func (m *mockBudgetRepository) FindByPeriod(userID identityvalueobjects.UserID, period valueobjects.BudgetPeriod) ([]*entities.Budget, error) {
	var result []*entities.Budget
	for _, budget := range m.budgets {
		if budget.UserID().Equals(userID) && budget.Period().Equals(period) {
			result = append(result, budget)
		}
	}
	return result, nil
}

func TestCreateBudgetUseCase_Execute(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockBudgetRepository()
	useCase := NewCreateBudgetUseCase(repository, eventBus)

	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()

	input := dtos.CreateBudgetInput{
		UserID:     userID.Value(),
		CategoryID: categoryID.Value(),
		Amount:     1000.00,
		Currency:   "BRL",
		PeriodType: "MONTHLY",
		Year:       2025,
		Month:      intPtr(12),
		Context:    "PERSONAL",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() returned nil output")
	}

	if output.BudgetID == "" {
		t.Error("Execute() returned output with empty budget ID")
	}

	if output.UserID != userID.Value() {
		t.Errorf("Execute() userID = %v, want %v", output.UserID, userID.Value())
	}

	if output.CategoryID != categoryID.Value() {
		t.Errorf("Execute() categoryID = %v, want %v", output.CategoryID, categoryID.Value())
	}

	if output.Amount != 1000.00 {
		t.Errorf("Execute() amount = %v, want 1000.00", output.Amount)
	}

	if output.Currency != "BRL" {
		t.Errorf("Execute() currency = %v, want BRL", output.Currency)
	}

	if output.PeriodType != "MONTHLY" {
		t.Errorf("Execute() periodType = %v, want MONTHLY", output.PeriodType)
	}

	if output.Year != 2025 {
		t.Errorf("Execute() year = %v, want 2025", output.Year)
	}

	if output.Month == nil || *output.Month != 12 {
		t.Errorf("Execute() month = %v, want 12", output.Month)
	}

	if output.Context != "PERSONAL" {
		t.Errorf("Execute() context = %v, want PERSONAL", output.Context)
	}

	if !output.IsActive {
		t.Error("Execute() returned inactive budget")
	}

	// Verify budget was saved
	budgetID, _ := valueobjects.NewBudgetID(output.BudgetID)
	savedBudget, err := repository.FindByID(budgetID)
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if savedBudget == nil {
		t.Error("Budget was not saved to repository")
	}
}

func TestCreateBudgetUseCase_Execute_DuplicateBudget(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockBudgetRepository()
	useCase := NewCreateBudgetUseCase(repository, eventBus)

	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()

	input := dtos.CreateBudgetInput{
		UserID:     userID.Value(),
		CategoryID: categoryID.Value(),
		Amount:     1000.00,
		Currency:   "BRL",
		PeriodType: "MONTHLY",
		Year:       2025,
		Month:      intPtr(12),
		Context:    "PERSONAL",
	}

	// Create first budget
	_, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("First Execute() error = %v, want nil", err)
	}

	// Try to create duplicate budget (same category and period)
	_, err = useCase.Execute(input)
	if err == nil {
		t.Error("Execute() with duplicate budget should return error")
	}
}

func TestCreateBudgetUseCase_Execute_InvalidInput(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockBudgetRepository()
	useCase := NewCreateBudgetUseCase(repository, eventBus)

	userID := identityvalueobjects.GenerateUserID()
	categoryID := categoryvalueobjects.GenerateCategoryID()

	tests := []struct {
		name      string
		input     dtos.CreateBudgetInput
		wantError bool
	}{
		{
			name: "invalid user ID",
			input: dtos.CreateBudgetInput{
				UserID:     "invalid-uuid",
				CategoryID: categoryID.Value(),
				Amount:     1000.00,
				Currency:   "BRL",
				PeriodType: "MONTHLY",
				Year:       2025,
				Month:      intPtr(12),
				Context:    "PERSONAL",
			},
			wantError: true,
		},
		{
			name: "invalid category ID",
			input: dtos.CreateBudgetInput{
				UserID:     userID.Value(),
				CategoryID: "invalid-uuid",
				Amount:     1000.00,
				Currency:   "BRL",
				PeriodType: "MONTHLY",
				Year:       2025,
				Month:      intPtr(12),
				Context:    "PERSONAL",
			},
			wantError: true,
		},
		{
			name: "monthly period without month",
			input: dtos.CreateBudgetInput{
				UserID:     userID.Value(),
				CategoryID: categoryID.Value(),
				Amount:     1000.00,
				Currency:   "BRL",
				PeriodType: "MONTHLY",
				Year:       2025,
				Month:      nil,
				Context:    "PERSONAL",
			},
			wantError: true,
		},
		{
			name: "invalid currency",
			input: dtos.CreateBudgetInput{
				UserID:     userID.Value(),
				CategoryID: categoryID.Value(),
				Amount:     1000.00,
				Currency:   "INVALID",
				PeriodType: "MONTHLY",
				Year:       2025,
				Month:      intPtr(12),
				Context:    "PERSONAL",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase.Execute(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("Execute() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("Execute() error = %v, want nil", err)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
