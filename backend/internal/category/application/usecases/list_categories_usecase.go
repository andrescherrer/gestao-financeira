package usecases

import (
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/repositories"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/pkg/pagination"
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
// Supports pagination when page and limit parameters are provided.
func (uc *ListCategoriesUseCase) Execute(input dtos.ListCategoriesInput) (*dtos.ListCategoriesOutput, error) {
	// Create user ID value object
	userID, err := identityvalueobjects.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse pagination parameters
	paginationParams := pagination.ParsePaginationParams(input.Page, input.Limit)
	usePagination := input.Page != "" || input.Limit != ""

	var categories []*entities.Category
	var total int64

	// Check if we should use pagination
	if usePagination {
		// Use paginated query
		categories, total, err = uc.categoryRepository.FindByUserIDWithPagination(
			userID,
			input.IsActive,
			paginationParams.CalculateOffset(),
			paginationParams.Limit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list categories: %w", err)
		}
	} else {
		// Use non-paginated query (backward compatibility)
		if input.IsActive != nil {
			categories, err = uc.categoryRepository.FindByUserIDAndActive(userID, *input.IsActive)
		} else {
			categories, err = uc.categoryRepository.FindByUserID(userID)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list categories: %w", err)
		}

		// Count total categories
		total, err = uc.categoryRepository.Count(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to count categories: %w", err)
		}
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
		Count:      int64(len(categoryOutputs)),
	}

	// Add pagination metadata if pagination was used
	if usePagination {
		paginationResult := pagination.BuildPaginationResult(paginationParams, total)
		output.Pagination = &paginationResult
		output.Count = total // Use total from query for accurate count
	}

	return output, nil
}
