package repositories

import (
	"gestao-financeira/backend/internal/category/domain/entities"
	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// CategoryRepository defines the interface for category persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type CategoryRepository interface {
	// FindByID finds a category by its ID.
	// Returns nil if the category is not found.
	FindByID(id valueobjects.CategoryID) (*entities.Category, error)

	// FindByUserID finds all categories for a given user.
	// Returns an empty slice if no categories are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Category, error)

	// FindByUserIDAndActive finds all active categories for a given user.
	// Returns an empty slice if no categories are found.
	FindByUserIDAndActive(userID identityvalueobjects.UserID, isActive bool) ([]*entities.Category, error)

	// Save saves or updates a category.
	// If the category already exists (by ID), it updates it.
	// If the category doesn't exist, it creates a new one.
	Save(category *entities.Category) error

	// Delete deletes a category by its ID.
	Delete(id valueobjects.CategoryID) error

	// Exists checks if a category with the given ID already exists.
	// Returns true if the category exists, false otherwise.
	Exists(id valueobjects.CategoryID) (bool, error)

	// Count returns the total number of categories for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)

	// FindByUserIDAndSlug finds a category by user ID and slug.
	// Returns nil if the category is not found.
	FindByUserIDAndSlug(userID identityvalueobjects.UserID, slug valueobjects.CategorySlug) (*entities.Category, error)
}
