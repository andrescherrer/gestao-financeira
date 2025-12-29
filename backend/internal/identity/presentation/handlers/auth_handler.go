package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/application/usecases"
	apperrors "gestao-financeira/backend/pkg/errors"
	"gestao-financeira/backend/pkg/middleware"
	"gestao-financeira/backend/pkg/validator"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	registerUserUseCase *usecases.RegisterUserUseCase
	loginUseCase        *usecases.LoginUseCase
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(
	registerUserUseCase *usecases.RegisterUserUseCase,
	loginUseCase *usecases.LoginUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUserUseCase: registerUserUseCase,
		loginUseCase:        loginUseCase,
	}
}

// Register handles user registration requests.
// @Summary Register a new user
// @Description Creates a new user account with the provided email, password, first name, and last name.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterUserInput true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully" "Example: {\"message\": \"User registered successfully\", \"data\": {\"user_id\": \"uuid\", \"email\": \"user@example.com\", \"first_name\": \"John\", \"last_name\": \"Doe\", \"full_name\": \"John Doe\"}}"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input data (e.g., invalid email format, password too short, missing fields)"
// @Failure 409 {object} map[string]interface{} "Conflict - User with this email already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error - An unexpected error occurred"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var input dtos.RegisterUserInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Corpo da requisição inválido",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate input using validator
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.registerUserUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuário registrado com sucesso",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *AuthHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeConflict {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("User registration failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("User registration failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}

// Login handles user login requests.
// @Summary Login user
// @Description Authenticates a user with email and password, returns a JWT token for API access
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginInput true "User login credentials (email and password)"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Success 200 {object} dtos.LoginOutput "JWT token and user data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid email or password"
// @Failure 403 {object} map[string]interface{} "Forbidden - user account is inactive"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var input dtos.LoginInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Str("request_id", middleware.GetRequestID(c)).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Corpo da requisição inválido",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate input using validator
	if err := validator.Validate(&input); err != nil {
		// Validation error is already an AppError, just return it
		return err
	}

	// Execute use case
	output, err := h.loginUseCase.Execute(input)
	if err != nil {
		return h.handleLoginError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login realizado com sucesso",
		"data":    output,
	})
}

// handleLoginError handles errors from login use case and returns appropriate HTTP responses.
// Uses AppError for consistent error handling instead of string matching.
func (h *AuthHandler) handleLoginError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Map domain errors to AppError
	appErr := apperrors.MapDomainError(err)

	// Log error with appropriate level
	if appErr.Type == apperrors.ErrorTypeValidation || appErr.Type == apperrors.ErrorTypeUnauthorized || appErr.Type == apperrors.ErrorTypeForbidden {
		log.Warn().Err(err).Str("error_type", string(appErr.Type)).Msg("Login failed")
	} else {
		log.Error().Err(err).Str("error_type", string(appErr.Type)).Msg("Login failed")
	}

	// Return error - the middleware will handle the response formatting
	return appErr
}
