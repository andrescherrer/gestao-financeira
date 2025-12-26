package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// UpdateCategoryUseCase handles category updates.
type UpdateCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
	eventBus           *eventbus.EventBus
}

// NewUpdateCategoryUseCase creates a new UpdateCategoryUseCase instance.
func NewUpdateCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
	eventBus *eventbus.EventBus,
) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		categoryRepository: categoryRepository,
		eventBus:           eventBus,
	}
}

// Execute performs the category update.
func (uc *UpdateCategoryUseCase) Execute(input dtos.UpdateCategoryInput) (*dtos.UpdateCategoryOutput, error) {
	// Create category ID value object
	categoryID, err := valueobjects.NewCategoryID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Find category by ID
	category, err := uc.categoryRepository.FindByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	// Update name if provided
	if input.Name != nil {
		categoryName, err := valueobjects.NewCategoryName(*input.Name)
		if err != nil {
			return nil, fmt.Errorf("invalid category name: %w", err)
		}
		if err := category.UpdateName(categoryName); err != nil {
			return nil, fmt.Errorf("failed to update category name: %w", err)
		}
	}

	// Update description if provided
	if input.Description != nil {
		if err := category.UpdateDescription(*input.Description); err != nil {
			return nil, fmt.Errorf("failed to update category description: %w", err)
		}
	}

	// Check if at least one field was provided for update
	if input.Name == nil && input.Description == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Save category to repository
	if err := uc.categoryRepository.Save(category); err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	// Publish domain events
	domainEvents := category.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			_ = err // Ignore for now, but should be logged
		}
	}
	category.ClearEvents()

	// Build output
	output := &dtos.UpdateCategoryOutput{
		CategoryID:  category.ID().Value(),
		UserID:      category.UserID().Value(),
		Name:        category.Name().Value(),
		Description: category.Description(),
		IsActive:    category.IsActive(),
		UpdatedAt:   category.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
