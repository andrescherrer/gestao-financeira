package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/internal/identity/application/dtos"
	"gestao-financeira/backend/internal/identity/application/usecases"
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
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var input dtos.RegisterUserInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate input (basic validation - can be enhanced with a validator)
	if input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	if input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	if input.FirstName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First name is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	if input.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Last name is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Execute use case
	output, err := h.registerUserUseCase.Execute(input)
	if err != nil {
		return h.handleUseCaseError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    output,
	})
}

// handleUseCaseError handles errors from use cases and returns appropriate HTTP responses.
func (h *AuthHandler) handleUseCaseError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "already exists") {
		log.Warn().Err(err).Msg("User registration failed: email already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
			"code":  fiber.StatusConflict,
		})
	}

	if strings.Contains(errMsg, "invalid email") {
		log.Warn().Err(err).Msg("User registration failed: invalid email")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "invalid password") {
		log.Warn().Err(err).Msg("User registration failed: invalid password")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters long",
			"code":  fiber.StatusBadRequest,
		})
	}

	if strings.Contains(errMsg, "invalid name") {
		log.Warn().Err(err).Msg("User registration failed: invalid name")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid name format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Generic error handling
	log.Error().Err(err).Msg("User registration failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
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
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var input dtos.LoginInput
	if err := c.BodyParser(&input); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Validate input
	if input.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	if input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Execute use case
	output, err := h.loginUseCase.Execute(input)
	if err != nil {
		return h.handleLoginError(c, err)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    output,
	})
}

// handleLoginError handles errors from login use case and returns appropriate HTTP responses.
func (h *AuthHandler) handleLoginError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Check for specific error types
	if strings.Contains(errMsg, "invalid email or password") {
		log.Warn().Err(err).Msg("Login failed: invalid credentials")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
			"code":  fiber.StatusUnauthorized,
		})
	}

	if strings.Contains(errMsg, "inactive") {
		log.Warn().Err(err).Msg("Login failed: user account is inactive")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User account is inactive",
			"code":  fiber.StatusForbidden,
		})
	}

	if strings.Contains(errMsg, "invalid email") {
		log.Warn().Err(err).Msg("Login failed: invalid email format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
			"code":  fiber.StatusBadRequest,
		})
	}

	// Generic error handling
	log.Error().Err(err).Msg("Login failed: unexpected error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An unexpected error occurred",
		"code":  fiber.StatusInternalServerError,
	})
}
