package repositories

import (
	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// UserRepository defines the interface for user persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type UserRepository interface {
	// FindByID finds a user by its ID.
	// Returns nil if the user is not found.
	FindByID(id valueobjects.UserID) (*entities.User, error)

	// FindByEmail finds a user by its email address.
	// Returns nil if the user is not found.
	FindByEmail(email valueobjects.Email) (*entities.User, error)

	// Save saves or updates a user.
	// If the user already exists (by ID), it updates it.
	// If the user doesn't exist, it creates a new one.
	Save(user *entities.User) error

	// Delete deletes a user by its ID.
	Delete(id valueobjects.UserID) error

	// Exists checks if a user with the given email already exists.
	// Returns true if the user exists, false otherwise.
	Exists(email valueobjects.Email) (bool, error)

	// Count returns the total number of users.
	Count() (int64, error)
}
