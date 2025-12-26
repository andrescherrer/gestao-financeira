package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
)

// DeleteCategoryUseCase handles category deletion.
type DeleteCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
}

// NewDeleteCategoryUseCase creates a new DeleteCategoryUseCase instance.
func NewDeleteCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		categoryRepository: categoryRepository,
	}
}

// Execute performs the category deletion.
func (uc *DeleteCategoryUseCase) Execute(input dtos.DeleteCategoryInput) (*dtos.DeleteCategoryOutput, error) {
	// Create category ID value object
	categoryID, err := valueobjects.NewCategoryID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Check if category exists
	category, err := uc.categoryRepository.FindByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	// Delete category (soft delete)
	if err := uc.categoryRepository.Delete(categoryID); err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	output := &dtos.DeleteCategoryOutput{
		Message:    "Category deleted successfully",
		CategoryID: categoryID.Value(),
	}

	return output, nil
}
