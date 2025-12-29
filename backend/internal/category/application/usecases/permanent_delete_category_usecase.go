package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/category/application/dtos"
	"gestao-financeira/backend/internal/category/domain/repositories"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
)

// PermanentDeleteCategoryUseCase handles permanent category deletion (hard delete).
// This should only be used by administrators.
type PermanentDeleteCategoryUseCase struct {
	categoryRepository repositories.CategoryRepository
}

// NewPermanentDeleteCategoryUseCase creates a new PermanentDeleteCategoryUseCase instance.
func NewPermanentDeleteCategoryUseCase(
	categoryRepository repositories.CategoryRepository,
) *PermanentDeleteCategoryUseCase {
	return &PermanentDeleteCategoryUseCase{
		categoryRepository: categoryRepository,
	}
}

// Execute performs the permanent category deletion.
func (uc *PermanentDeleteCategoryUseCase) Execute(input dtos.PermanentDeleteCategoryInput) (*dtos.PermanentDeleteCategoryOutput, error) {
	// Create category ID value object
	categoryID, err := valueobjects.NewCategoryID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Try to permanently delete
	repo, ok := uc.categoryRepository.(interface {
		PermanentDelete(valueobjects.CategoryID) error
	})
	if !ok {
		return nil, errors.New("repository does not support permanent delete operation")
	}

	if err := repo.PermanentDelete(categoryID); err != nil {
		return nil, fmt.Errorf("failed to permanently delete category: %w", err)
	}

	output := &dtos.PermanentDeleteCategoryOutput{
		Message:    "Category permanently deleted successfully",
		CategoryID: categoryID.Value(),
	}

	return output, nil
}
