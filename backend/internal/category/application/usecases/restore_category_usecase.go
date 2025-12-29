package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
)

// RestoreCategoryUseCase handles category restoration from soft delete.
type RestoreCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
}

// NewRestoreCategoryUseCase creates a new RestoreCategoryUseCase instance.
func NewRestoreCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
) *RestoreCategoryUseCase {
	return &RestoreCategoryUseCase{
		categoryRepository: categoryRepository,
	}
}

// Execute performs the category restoration.
func (uc *RestoreCategoryUseCase) Execute(input dtos.RestoreCategoryInput) (*dtos.RestoreCategoryOutput, error) {
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

	// If category exists and is not deleted, return error
	if category != nil {
		return nil, errors.New("category is not deleted")
	}

	// Try to restore
	repo, ok := uc.categoryRepository.(interface {
		Restore(valueobjects.CategoryID) error
	})
	if !ok {
		return nil, errors.New("repository does not support restore operation")
	}

	if err := repo.Restore(categoryID); err != nil {
		return nil, fmt.Errorf("failed to restore category: %w", err)
	}

	output := &dtos.RestoreCategoryOutput{
		Message:    "Category restored successfully",
		CategoryID: categoryID.Value(),
	}

	return output, nil
}
