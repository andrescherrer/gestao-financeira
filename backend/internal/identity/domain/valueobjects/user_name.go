package valueobjects

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// UserName represents a user's name value object.
// It can represent a full name or be split into first and last name.
type UserName struct {
	firstName string
	lastName  string
	fullName  string
}

const (
	// MinNameLength is the minimum length for a name part.
	MinNameLength = 2
	// MaxNameLength is the maximum length for a name part.
	MaxNameLength = 100
	// MaxFullNameLength is the maximum length for the full name.
	MaxFullNameLength = 200
)

// NewUserName creates a new UserName from first and last name.
func NewUserName(firstName, lastName string) (UserName, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if err := validateNamePart(firstName, "first name"); err != nil {
		return UserName{}, err
	}

	if err := validateNamePart(lastName, "last name"); err != nil {
		return UserName{}, err
	}

	fullName := fmt.Sprintf("%s %s", firstName, lastName)

	return UserName{
		firstName: firstName,
		lastName:  lastName,
		fullName:  fullName,
	}, nil
}

// NewUserNameFromFullName creates a new UserName from a full name string.
// It attempts to split the name into first and last name.
func NewUserNameFromFullName(fullName string) (UserName, error) {
	fullName = strings.TrimSpace(fullName)

	if fullName == "" {
		return UserName{}, errors.New("full name cannot be empty")
	}

	if len(fullName) > MaxFullNameLength {
		return UserName{}, fmt.Errorf("full name is too long (max %d characters)", MaxFullNameLength)
	}

	// Split by spaces
	parts := strings.Fields(fullName)

	if len(parts) == 0 {
		return UserName{}, errors.New("full name must contain at least one word")
	}

	// If only one part, use it as first name and empty last name
	if len(parts) == 1 {
		firstName := parts[0]
		if err := validateNamePart(firstName, "name"); err != nil {
			return UserName{}, err
		}
		return UserName{
			firstName: firstName,
			lastName:  "",
			fullName:  firstName,
		}, nil
	}

	// Multiple parts: first part is first name, rest is last name
	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")

	if err := validateNamePart(firstName, "first name"); err != nil {
		return UserName{}, err
	}

	if err := validateNamePart(lastName, "last name"); err != nil {
		return UserName{}, err
	}

	return UserName{
		firstName: firstName,
		lastName:  lastName,
		fullName:  fullName,
	}, nil
}

// MustUserName creates a new UserName and panics if invalid.
// Use this only when you are certain the name is valid (e.g., in tests).
func MustUserName(firstName, lastName string) UserName {
	un, err := NewUserName(firstName, lastName)
	if err != nil {
		panic(err)
	}
	return un
}

// FirstName returns the first name.
func (un UserName) FirstName() string {
	return un.firstName
}

// LastName returns the last name.
func (un UserName) LastName() string {
	return un.lastName
}

// FullName returns the full name.
func (un UserName) FullName() string {
	return un.fullName
}

// String returns the full name (implements fmt.Stringer).
func (un UserName) String() string {
	return un.fullName
}

// Initials returns the initials of the name (e.g., "JD" for "John Doe").
func (un UserName) Initials() string {
	initials := ""
	if len(un.firstName) > 0 {
		initials += string(unicode.ToUpper(rune(un.firstName[0])))
	}
	if len(un.lastName) > 0 {
		initials += string(unicode.ToUpper(rune(un.lastName[0])))
	}
	return initials
}

// Equals checks if two UserName values are equal.
func (un UserName) Equals(other UserName) bool {
	return un.firstName == other.firstName && un.lastName == other.lastName
}

// IsEmpty checks if the name is empty.
func (un UserName) IsEmpty() bool {
	return un.firstName == "" && un.lastName == ""
}

// HasLastName checks if the user has a last name.
func (un UserName) HasLastName() bool {
	return un.lastName != ""
}

// validateNamePart validates a single name part (first or last name).
func validateNamePart(name, partName string) error {
	if name == "" {
		return fmt.Errorf("%s cannot be empty", partName)
	}

	if len(name) < MinNameLength {
		return fmt.Errorf("%s must be at least %d characters long", partName, MinNameLength)
	}

	if len(name) > MaxNameLength {
		return fmt.Errorf("%s is too long (max %d characters)", partName, MaxNameLength)
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	for _, char := range name {
		if !unicode.IsLetter(char) && char != ' ' && char != '-' && char != '\'' {
			return fmt.Errorf("%s contains invalid characters", partName)
		}
	}

	// Check that it doesn't start or end with space, hyphen, or apostrophe
	if strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") {
		return fmt.Errorf("%s cannot start or end with a space", partName)
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("%s cannot start or end with a hyphen", partName)
	}

	if strings.HasPrefix(name, "'") || strings.HasSuffix(name, "'") {
		return fmt.Errorf("%s cannot start or end with an apostrophe", partName)
	}

	return nil
}
