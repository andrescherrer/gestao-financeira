package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/domain/entities"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/infrastructure/eventbus"
)

// RegisterUserUseCase handles user registration.
type RegisterUserUseCase struct {
	userRepository repositories.UserRepository
	eventBus       *eventbus.EventBus
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase instance.
func NewRegisterUserUseCase(
	userRepository repositories.UserRepository,
	eventBus *eventbus.EventBus,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository: userRepository,
		eventBus:       eventBus,
	}
}

// Execute performs the user registration.
// It validates the input, checks if the email already exists, creates a new user,
// saves it to the repository, and publishes domain events.
func (uc *RegisterUserUseCase) Execute(input dtos.RegisterUserInput) (*dtos.RegisterUserOutput, error) {
	// Create email value object
	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Check if user already exists
	exists, err := uc.userRepository.Exists(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Create password hash value object
	passwordHash, err := valueobjects.NewPasswordHashFromPlain(input.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// Create user name value object
	name, err := valueobjects.NewUserName(input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	// Create user entity
	user, err := entities.NewUser(email, passwordHash, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save user to repository
	if err := uc.userRepository.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Publish domain events
	domainEvents := user.GetEvents()
	for _, event := range domainEvents {
		if err := uc.eventBus.Publish(event); err != nil {
			// Log error but don't fail the registration
			// In production, you might want to handle this differently
			// (e.g., store events in an outbox pattern)
			_ = err // Ignore for now, but should be logged
		}
	}
	user.ClearEvents()

	// Build output
	output := &dtos.RegisterUserOutput{
		UserID:    user.ID().Value(),
		Email:     user.Email().Value(),
		FirstName: user.Name().FirstName(),
		LastName:  user.Name().LastName(),
		FullName:  user.Name().FullName(),
	}

	return output, nil
}
