package usecases

import (
	"errors"
	"fmt"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/domain/repositories"
	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/identity/infrastructure/services"
)

// LoginUseCase handles user authentication.
type LoginUseCase struct {
	userRepository repositories.UserRepository
	jwtService     *services.JWTService
}

// NewLoginUseCase creates a new LoginUseCase instance.
func NewLoginUseCase(
	userRepository repositories.UserRepository,
	jwtService *services.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

// Execute performs user authentication.
// It validates credentials, checks if user exists and is active,
// verifies the password, and generates a JWT token.
func (uc *LoginUseCase) Execute(input dtos.LoginInput) (*dtos.LoginOutput, error) {
	// Create email value object
	email, err := valueobjects.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Find user by email
	user, err := uc.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if !user.VerifyPassword(input.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user.ID().Value(), user.Email().Value())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build output
	output := &dtos.LoginOutput{
		Token:     token,
		UserID:    user.ID().Value(),
		Email:     user.Email().Value(),
		FirstName: user.Name().FirstName(),
		LastName:  user.Name().LastName(),
		FullName:  user.Name().FullName(),
		ExpiresIn: 24 * 60 * 60, // 24 hours in seconds
	}

	return output, nil
}
