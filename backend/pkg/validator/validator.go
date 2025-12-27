package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"gestao-financeira/backend/pkg/errors"
)

var (
	// validate is the singleton validator instance
	validate *validator.Validate
)

// Init initializes the validator with custom validations
func Init() {
	validate = validator.New()

	// Use JSON tag names instead of struct field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	registerCustomValidations()
}

// Validate validates a struct using the validate tags
func Validate(s interface{}) error {
	if validate == nil {
		Init()
	}

	if err := validate.Struct(s); err != nil {
		validationErrors := make(map[string]interface{})
		fieldErrors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				field := fe.Field()
				// Use JSON tag name if available
				if jsonTag := getJSONTagName(reflect.TypeOf(s), field); jsonTag != "" {
					field = jsonTag
				}

				// Build error message
				var message string
				switch fe.Tag() {
				case "required":
					message = fmt.Sprintf("%s is required", field)
				case "email":
					message = fmt.Sprintf("%s must be a valid email address", field)
				case "min":
					message = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
				case "max":
					message = fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
				case "uuid":
					message = fmt.Sprintf("%s must be a valid UUID", field)
				case "oneof":
					message = fmt.Sprintf("%s must be one of: %s", field, fe.Param())
				case "gte":
					message = fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
				case "gt":
					message = fmt.Sprintf("%s must be greater than %s", field, fe.Param())
				case "lte":
					message = fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
				case "lt":
					message = fmt.Sprintf("%s must be less than %s", field, fe.Param())
				default:
					message = fmt.Sprintf("%s is invalid", field)
				}

				fieldErrors[field] = message
			}

			validationErrors["fields"] = fieldErrors
			validationErrors["count"] = len(fieldErrors)

			return errors.NewValidationError("Validation failed", validationErrors)
		}

		return errors.NewValidationError("Validation failed", nil)
	}

	return nil
}

// registerCustomValidations registers custom validation functions
func registerCustomValidations() {
	// Example: Custom validation for date format (YYYY-MM-DD)
	validate.RegisterValidation("date_iso8601", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		// Simple validation for YYYY-MM-DD format
		if len(dateStr) != 10 {
			return false
		}
		parts := strings.Split(dateStr, "-")
		if len(parts) != 3 {
			return false
		}
		// Check if all parts are numeric
		for _, part := range parts {
			if !isNumeric(part) {
				return false
			}
		}
		return true
	})

	log.Info().Msg("Custom validations registered")
}

// getJSONTagName extracts the JSON tag name from a struct field
func getJSONTagName(t reflect.Type, fieldName string) string {
	field, found := t.FieldByName(fieldName)
	if !found {
		return ""
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}

	// Get the first part before comma
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

// isNumeric checks if a string is numeric
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
