package persistence

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"

	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM.
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository.
func NewGormUserRepository(db *gorm.DB) repositories.UserRepository {
	return &GormUserRepository{db: db}
}

// FindByID finds a user by its ID.
func (r *GormUserRepository) FindByID(id valueobjects.UserID) (*entities.User, error) {
	var model UserModel
	if err := r.db.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return r.toDomain(&model)
}

// FindByEmail finds a user by its email address.
func (r *GormUserRepository) FindByEmail(email valueobjects.Email) (*entities.User, error) {
	var model UserModel
	if err := r.db.Where("email = ?", email.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.toDomain(&model)
}

// Save saves or updates a user.
func (r *GormUserRepository) Save(user *entities.User) error {
	model := r.toModel(user)

	// Check if user exists
	var existing UserModel
	err := r.db.Where("id = ?", model.ID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new user
		if err := r.db.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	} else {
		// Update existing user
		if err := r.db.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

// Delete deletes a user by its ID.
func (r *GormUserRepository) Delete(id valueobjects.UserID) error {
	if err := r.db.Where("id = ?", id.Value()).Delete(&UserModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Exists checks if a user with the given email already exists.
func (r *GormUserRepository) Exists(email valueobjects.Email) (bool, error) {
	var count int64
	if err := r.db.Model(&UserModel{}).Where("email = ?", email.Value()).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

// Count returns the total number of users.
func (r *GormUserRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// toDomain converts a UserModel to a User entity.
func (r *GormUserRepository) toDomain(model *UserModel) (*entities.User, error) {
	// Convert value objects
	userID, err := valueobjects.NewUserID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	email, err := valueobjects.NewEmail(model.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	passwordHash, err := valueobjects.NewPasswordHashFromHash(model.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("invalid password hash: %w", err)
	}

	name, err := valueobjects.NewUserName(model.FirstName, model.LastName)
	if err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	// Reconstruct user entity from persisted data
	return entities.FromPersistence(
		userID,
		email,
		passwordHash,
		name,
		model.CreatedAt,
		model.UpdatedAt,
		model.IsActive,
	)
}

// toModel converts a User entity to a UserModel.
func (r *GormUserRepository) toModel(user *entities.User) *UserModel {
	name := user.Name()

	return &UserModel{
		ID:           user.ID().Value(),
		Email:        user.Email().Value(),
		PasswordHash: user.PasswordHash().Value(),
		FirstName:    name.FirstName(),
		LastName:     name.LastName(),
		IsActive:     user.IsActive(),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}
}
