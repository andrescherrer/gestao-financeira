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
// @Description Creates a new user account with the provided email, password, first name, and last name. The password is hashed using bcrypt before storage.
//
// **Validações**:
// - Email deve ser válido e único no sistema
// - Senha deve ter no mínimo 8 caracteres
// - Nome e sobrenome são obrigatórios
//
// **Segurança**:
// - Senha é hasheada com bcrypt antes de ser armazenada
// - Email é convertido para lowercase antes de ser armazenado
//
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RegisterUserInput true "User registration data" example({"email":"user@example.com","password":"SecurePass123","first_name":"John","last_name":"Doe"})
// @Success 201 {object} map[string]interface{} "User registered successfully" example({"message":"Usuário registrado com sucesso","data":{"user_id":"550e8400-e29b-41d4-a716-446655440000","email":"user@example.com","first_name":"John","last_name":"Doe","full_name":"John Doe"}})
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input data (e.g., invalid email format, password too short, missing fields)" example({"error":"Invalid email format","error_type":"VALIDATION_ERROR","code":400,"details":{"field":"email","message":"email must be a valid email address"}})
// @Failure 409 {object} map[string]interface{} "Conflict - User with this email already exists" example({"error":"Já existe um usuário com este email","error_type":"CONFLICT","code":409})
// @Failure 422 {object} map[string]interface{} "Unprocessable entity - domain validation failed" example({"error":"Password does not meet security requirements","error_type":"DOMAIN_ERROR","code":422})
// @Failure 500 {object} map[string]interface{} "Internal server error - An unexpected error occurred" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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
// Login handles user login requests.
// @Summary Login user
// @Description Authenticates a user with email and password, returns a JWT token for API access. The token should be included in subsequent requests in the Authorization header.
//
// **Autenticação**:
// - Email e senha são validados
// - Senha é verificada usando bcrypt
// - Token JWT é gerado com informações do usuário
// - Token expira após o tempo configurado (padrão: 24h)
//
// **Uso do Token**:
// - Inclua o token no header: `Authorization: Bearer <token>`
// - Token é válido para todos os endpoints protegidos
// - Token expirado retorna 401 Unauthorized
//
// **Segurança**:
// - Tentativas de login falhadas não revelam se o email existe
// - Senha nunca é retornada na resposta
// - Contas inativas não podem fazer login
//
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginInput true "User login credentials" example({"email":"user@example.com","password":"SecurePass123"})
// @Success 200 {object} map[string]interface{} "Login successful" example({"message":"Login realizado com sucesso","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","user":{"user_id":"550e8400-e29b-41d4-a716-446655440000","email":"user@example.com","first_name":"John","last_name":"Doe","full_name":"John Doe"}})
// @Success 200 {object} dtos.LoginOutput "JWT token and user data"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid input data" example({"error":"Invalid email format","error_type":"VALIDATION_ERROR","code":400})
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid email or password" example({"error":"Email ou senha inválidos","error_type":"UNAUTHORIZED","code":401})
// @Failure 403 {object} map[string]interface{} "Forbidden - user account is inactive" example({"error":"Conta de usuário está inativa","error_type":"FORBIDDEN","code":403})
// @Failure 500 {object} map[string]interface{} "Internal server error" example({"error":"An unexpected error occurred","error_type":"INTERNAL_ERROR","code":500,"request_id":"req-123"})
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
