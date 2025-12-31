package repositories

import (
	"gestao-financeira/backend/internal/goal/domain/entities"
	goalvalueobjects "gestao-financeira/backend/internal/goal/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
)

// GoalRepository defines the interface for goal persistence operations.
// This interface belongs to the domain layer and will be implemented in the infrastructure layer.
type GoalRepository interface {
	// FindByID finds a goal by its ID.
	// Returns nil if the goal is not found.
	FindByID(id goalvalueobjects.GoalID) (*entities.Goal, error)

	// FindByUserID finds all goals for a given user.
	// Returns an empty slice if no goals are found.
	FindByUserID(userID identityvalueobjects.UserID) ([]*entities.Goal, error)

	// FindByStatus finds all goals for a given user filtered by status.
	// Returns an empty slice if no goals are found.
	FindByStatus(userID identityvalueobjects.UserID, status goalvalueobjects.GoalStatus) ([]*entities.Goal, error)

	// Save saves or updates a goal.
	// If the goal already exists (by ID), it updates it.
	// If the goal doesn't exist, it creates a new one.
	Save(goal *entities.Goal) error

	// Delete deletes a goal by its ID.
	Delete(id goalvalueobjects.GoalID) error

	// Exists checks if a goal with the given ID already exists.
	// Returns true if the goal exists, false otherwise.
	Exists(id goalvalueobjects.GoalID) (bool, error)

	// Count returns the total number of goals for a given user.
	Count(userID identityvalueobjects.UserID) (int64, error)

	// FindByUserIDWithPagination finds goals for a given user with pagination.
	// Returns the goals, total count, and any error.
	FindByUserIDWithPagination(
		userID identityvalueobjects.UserID,
		context string,
		status string,
		offset, limit int,
	) ([]*entities.Goal, int64, error)
}
