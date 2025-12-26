package usecases

import (
	"testing"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// mockCategoryRepository is a mock implementation of CategoryRepository for testing.
type mockCategoryRepository struct {
	categories map[string]*entities.Category
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[string]*entities.Category),
	}
}

func (m *mockCategoryRepository) FindByID(id valueobjects.CategoryID) (*entities.Category, error) {
	category, exists := m.categories[id.Value()]
	if !exists {
		return nil, nil
	}
	return category, nil
}

func (m *mockCategoryRepository) FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Category, error) {
	var result []*entities.Category
	for _, category := range m.categories {
		if category.UserID().Equals(userID) {
			result = append(result, category)
		}
	}
	return result, nil
}

func (m *mockCategoryRepository) FindByUserIDAndActive(userID identityvalueobjects.UserID, isActive bool) ([]*entities.Category, error) {
	var result []*entities.Category
	for _, category := range m.categories {
		if category.UserID().Equals(userID) && category.IsActive() == isActive {
			result = append(result, category)
		}
	}
	return result, nil
}

func (m *mockCategoryRepository) Save(category *entities.Category) error {
	m.categories[category.ID().Value()] = category
	return nil
}

func (m *mockCategoryRepository) Delete(id valueobjects.CategoryID) error {
	delete(m.categories, id.Value())
	return nil
}

func (m *mockCategoryRepository) Exists(id valueobjects.CategoryID) (bool, error) {
	_, exists := m.categories[id.Value()]
	return exists, nil
}

func (m *mockCategoryRepository) Count(userID identityvalueobjects.UserID) (int64, error) {
	count := int64(0)
	for _, category := range m.categories {
		if category.UserID().Equals(userID) {
			count++
		}
	}
	return count, nil
}

func TestCreateCategoryUseCase_Execute(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockCategoryRepository()
	useCase := NewCreateCategoryUseCase(repository, eventBus)

	userID := identityvalueobjects.GenerateUserID()

	input := dtos.CreateCategoryInput{
		UserID:      userID.Value(),
		Name:        "Alimentação",
		Description: "Gastos com alimentação",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() returned nil output")
	}

	if output.CategoryID == "" {
		t.Error("Execute() returned output with empty category ID")
	}

	if output.Name != input.Name {
		t.Errorf("Execute() output.Name = %v, want %v", output.Name, input.Name)
	}

	if output.Description != input.Description {
		t.Errorf("Execute() output.Description = %v, want %v", output.Description, input.Description)
	}

	if !output.IsActive {
		t.Error("Execute() returned output with inactive category")
	}

	// Verify category was saved
	savedCategory, err := repository.FindByID(valueobjects.MustCategoryID(output.CategoryID))
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if savedCategory == nil {
		t.Fatal("Category was not saved to repository")
	}
}

func TestCreateCategoryUseCase_Execute_InvalidInput(t *testing.T) {
	eventBus := eventbus.NewEventBus()
	repository := newMockCategoryRepository()
	useCase := NewCreateCategoryUseCase(repository, eventBus)

	tests := []struct {
		name      string
		input     dtos.CreateCategoryInput
		wantError bool
	}{
		{
			name: "invalid user ID",
			input: dtos.CreateCategoryInput{
				UserID:      "invalid-uuid",
				Name:        "Alimentação",
				Description: "Description",
			},
			wantError: true,
		},
		{
			name: "invalid category name",
			input: dtos.CreateCategoryInput{
				UserID:      identityvalueobjects.GenerateUserID().Value(),
				Name:        "A", // Too short
				Description: "Description",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase.Execute(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("Execute() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
