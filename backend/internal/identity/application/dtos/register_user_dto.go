package dtos

// RegisterUserInput represents the input data for user registration.
// @Description User registration input with email, password, first name and last name
type RegisterUserInput struct {
	// Email address of the user (must be a valid email format, maximum 255 characters)
	Email string `json:"email" example:"user@example.com" validate:"required,email,max=255,no_sql_injection,no_xss,utf8"`
	// Password for the user account (minimum 8 characters, must meet strength requirements)
	Password string `json:"password" example:"SecurePass123" validate:"required,min=8,max=128,password_strength,no_sql_injection,no_xss"`
	// First name of the user (minimum 2 characters, maximum 100)
	FirstName string `json:"first_name" example:"John" validate:"required,min=2,max=100,no_sql_injection,no_xss,utf8"`
	// Last name of the user (minimum 2 characters, maximum 100)
	LastName string `json:"last_name" example:"Doe" validate:"required,min=2,max=100,no_sql_injection,no_xss,utf8"`
}

// RegisterUserOutput represents the output data after user registration.
type RegisterUserOutput struct {
	// Unique identifier for the user
	UserID string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Email address of the registered user
	Email string `json:"email" example:"user@example.com"`
	// First name of the user
	FirstName string `json:"first_name" example:"John"`
	// Last name of the user
	LastName string `json:"last_name" example:"Doe"`
	// Full name of the user (first name + last name)
	FullName string `json:"full_name" example:"John Doe"`
}
