package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/repositories"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// ListCategoriesUseCase handles listing categories.
type ListCategoriesUseCase struct {
	categoryRepository repositories.CategoryRepository
}

// NewListCategoriesUseCase creates a new ListCategoriesUseCase instance.
func NewListCategoriesUseCase(
	categoryRepository repositories.CategoryRepository,
) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		categoryRepository: categoryRepository,
	}
}

// Execute performs the category listing.
func (uc *ListCategoriesUseCase) Execute(input dtos.ListCategoriesInput) (*dtos.ListCategoriesOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var categories []*entities.Category

	// Filter by active status if provided
	if input.IsActive != nil {
		categories, err = uc.categoryRepository.FindByUserIDAndActive(userID, *input.IsActive)
	} else {
		categories, err = uc.categoryRepository.FindByUserID(userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Count total categories
	count, err := uc.categoryRepository.Count(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	// Convert to output DTOs
	categoryOutputs := make([]dtos.GetCategoryOutput, 0, len(categories))
	for _, category := range categories {
		categoryOutputs = append(categoryOutputs, dtos.GetCategoryOutput{
			CategoryID:  category.ID().Value(),
			UserID:      category.UserID().Value(),
			Name:        category.Name().Value(),
			Slug:        category.Slug().Value(),
			Description: category.Description(),
			IsActive:    category.IsActive(),
			CreatedAt:   category.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   category.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	output := &dtos.ListCategoriesOutput{
		Categories: categoryOutputs,
		Count:      count,
	}

	return output, nil
}
