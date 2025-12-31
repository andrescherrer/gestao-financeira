package validator

import (
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// FieldError represents a validation field error
type FieldError struct {
	Field string
	Tag   string
	Param string
}

func (e FieldError) Error() string {
	return e.Tag
}

// SecurityValidationConfig holds configuration for security validations
type SecurityValidationConfig struct {
	MaxStringLength    int
	MaxDescriptionLength int
	MaxNameLength      int
	AllowHTML          bool
	SanitizeInput      bool
}

// DefaultSecurityConfig returns default security validation configuration
func DefaultSecurityConfig() SecurityValidationConfig {
	return SecurityValidationConfig{
		MaxStringLength:      1000,
		MaxDescriptionLength: 5000,
		MaxNameLength:        255,
		AllowHTML:            false,
		SanitizeInput:        true,
	}
}

// SanitizeString removes potentially dangerous characters from a string
func SanitizeString(input string, config SecurityValidationConfig) string {
	if !config.SanitizeInput {
		return input
	}

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}
	input = result.String()

	// HTML escape if HTML is not allowed
	if !config.AllowHTML {
		input = html.EscapeString(input)
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// ValidateStringLength validates that a string is within acceptable length limits
func ValidateStringLength(s string, maxLength int) error {
	if len(s) > maxLength {
		return FieldError{
			Field: "string",
			Tag:   "max",
			Param: string(rune(maxLength)),
		}
	}
	return nil
}

// ValidateNoSQLInjection checks for common SQL injection patterns
func ValidateNoSQLInjection(input string) error {
	// Common SQL injection patterns
	patterns := []string{
		"(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)",
		"(?i)(--|;|/\\*|\\*/|xp_|sp_)",
		"(?i)(or|and)\\s+\\d+\\s*=\\s*\\d+",
		"(?i)(or|and)\\s+['\"].*['\"]\\s*=\\s*['\"].*['\"]",
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			continue
		}
		if matched {
			return FieldError{
				Field: "input",
				Tag:   "no_sql_injection",
				Param: "",
			}
		}
	}

	return nil
}

// ValidateNoXSS checks for common XSS patterns
func ValidateNoXSS(input string) error {
	// Common XSS patterns
	patterns := []string{
		"(?i)<script",
		"(?i)javascript:",
		"(?i)onerror\\s*=",
		"(?i)onclick\\s*=",
		"(?i)onload\\s*=",
		"(?i)onmouseover\\s*=",
		"(?i)<iframe",
		"(?i)<object",
		"(?i)<embed",
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			continue
		}
		if matched {
			return FieldError{
				Field: "input",
				Tag:   "no_xss",
				Param: "",
			}
		}
	}

	return nil
}

// ValidateNoPathTraversal checks for path traversal patterns
func ValidateNoPathTraversal(input string) error {
	// Path traversal patterns
	patterns := []string{
		"\\.\\./",
		"\\.\\.\\\\",
		"/etc/passwd",
		"\\\\windows\\\\system32",
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			continue
		}
		if matched {
			return FieldError{
				Field: "input",
				Tag:   "no_path_traversal",
				Param: "",
			}
		}
	}

	return nil
}

// ValidateUTF8 validates that a string is valid UTF-8
func ValidateUTF8(input string) error {
	if !utf8.ValidString(input) {
		return FieldError{
			Field: "input",
			Tag:   "utf8",
			Param: "",
		}
	}
	return nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return FieldError{
			Field: "password",
			Tag:   "min",
			Param: "8",
		}
	}

	if len(password) > 128 {
		return FieldError{
			Field: "password",
			Tag:   "max",
			Param: "128",
		}
	}

	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	// Require at least 3 of 4 character types
	typesCount := 0
	if hasUpper {
		typesCount++
	}
	if hasLower {
		typesCount++
	}
	if hasDigit {
		typesCount++
	}
	if hasSpecial {
		typesCount++
	}

	if typesCount < 3 {
		return FieldError{
			Field: "password",
			Tag:   "password_strength",
			Param: "",
		}
	}

	return nil
}

// RegisterSecurityValidations registers security-related custom validations
func RegisterSecurityValidations(v *validator.Validate) {
	// Register no_sql_injection validation
	v.RegisterValidation("no_sql_injection", func(fl validator.FieldLevel) bool {
		err := ValidateNoSQLInjection(fl.Field().String())
		return err == nil
	})

	// Register no_xss validation
	v.RegisterValidation("no_xss", func(fl validator.FieldLevel) bool {
		err := ValidateNoXSS(fl.Field().String())
		return err == nil
	})

	// Register no_path_traversal validation
	v.RegisterValidation("no_path_traversal", func(fl validator.FieldLevel) bool {
		err := ValidateNoPathTraversal(fl.Field().String())
		return err == nil
	})

	// Register utf8 validation
	v.RegisterValidation("utf8", func(fl validator.FieldLevel) bool {
		err := ValidateUTF8(fl.Field().String())
		return err == nil
	})

	// Register password_strength validation
	v.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		err := ValidatePasswordStrength(fl.Field().String())
		return err == nil
	})
}

