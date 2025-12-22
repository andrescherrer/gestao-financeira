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
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(registerUserUseCase *usecases.RegisterUserUseCase) *AuthHandler {
	return &AuthHandler{
		registerUserUseCase: registerUserUseCase,
	}
}

// Register handles user registration requests.
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterUserInput true "User registration data"
// @Success 201 {object} dtos.RegisterUserOutput
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "Conflict - user already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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
