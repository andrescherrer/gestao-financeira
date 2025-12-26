package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// CreateCategoryUseCase handles category creation.
type CreateCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
	eventBus           *eventbus.EventBus
}

// NewCreateCategoryUseCase creates a new CreateCategoryUseCase instance.
func NewCreateCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
	eventBus *eventbus.EventBus,
) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepository: categoryRepository,
		eventBus:           eventBus,
	}
}

// Execute performs the category creation.
// It validates the input, creates value objects, creates a new category entity,
// saves it to the repository, and publishes domain events.
func (uc *CreateCategoryUseCase) Execute(input dtos.CreateCategoryInput) (*dtos.CreateCategoryOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create category name value object
	categoryName, err := valueobjects.NewCategoryName(input.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid category name: %w", err)
	}

	// Create category entity
	category, err := entities.NewCategory(userID, categoryName, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Save category to repository
	if err := uc.categoryRepository.Save(category); err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	// Publish domain events
	domainEvents := category.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the category creation
			_ = err // Ignore for now, but should be logged
		}
	}
	category.ClearEvents()

	// Build output
	output := &dtos.CreateCategoryOutput{
		CategoryID:  category.ID().Value(),
		UserID:      category.UserID().Value(),
		Name:        category.Name().Value(),
		Slug:        category.Slug().Value(),
		Description: category.Description(),
		IsActive:    category.IsActive(),
		CreatedAt:   category.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
