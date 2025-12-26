package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
)

// GetCategoryUseCase handles getting a single category.
type GetCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
}

// NewGetCategoryUseCase creates a new GetCategoryUseCase instance.
func NewGetCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
) *GetCategoryUseCase {
	return &GetCategoryUseCase{
		categoryRepository: categoryRepository,
	}
}

// Execute performs the category retrieval.
func (uc *GetCategoryUseCase) Execute(categoryID string) (*dtos.GetCategoryOutput, error) {
	// Create category ID value object
	id, err := valueobjects.NewCategoryID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Find category by ID
	category, err := uc.categoryRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	// Build output
	output := &dtos.GetCategoryOutput{
		CategoryID:  category.ID().Value(),
		UserID:      category.UserID().Value(),
		Name:        category.Name().Value(),
		Description: category.Description(),
		IsActive:    category.IsActive(),
		CreatedAt:   category.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}
