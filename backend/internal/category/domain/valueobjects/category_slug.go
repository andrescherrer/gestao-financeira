package valueobjects

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CategorySlug represents a category slug value object.
type CategorySlug struct {
	value string
}

const (
	// MinCategorySlugLength is the minimum length for a category slug.
	MinCategorySlugLength = 2
	// MaxCategorySlugLength is the maximum length for a category slug.
	MaxCategorySlugLength = 100
)

// NewCategorySlug creates a new CategorySlug value object.
func NewCategorySlug(slug string) (CategorySlug, error) {
	slug = strings.TrimSpace(slug)

	if slug == "" {
		return CategorySlug{}, errors.New("category slug cannot be empty")
	}

	if len(slug) < MinCategorySlugLength {
		return CategorySlug{}, errors.New("category slug is too short")
	}

	if len(slug) > MaxCategorySlugLength {
		return CategorySlug{}, errors.New("category slug is too long")
	}

	// Validate slug format (lowercase letters, numbers, hyphens)
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return CategorySlug{}, errors.New("category slug contains invalid characters")
		}
	}

	return CategorySlug{value: slug}, nil
}

// GenerateSlugFromName generates a slug from a category name.
// It removes accents, converts to lowercase, replaces spaces with hyphens,
// and removes invalid characters.
func GenerateSlugFromName(name string) CategorySlug {
	// Remove accents
	slug := removeAccents(strings.ToLower(strings.TrimSpace(name)))

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove invalid characters (keep only letters, numbers, and hyphens)
	var builder strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			builder.WriteRune(char)
		}
	}
	slug = builder.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// If empty after processing, use default
	if slug == "" {
		slug = "categoria"
	}

	// Ensure minimum length
	if len(slug) < MinCategorySlugLength {
		slug = slug + "-cat"
	}

	// Truncate if too long
	if len(slug) > MaxCategorySlugLength {
		slug = slug[:MaxCategorySlugLength]
		slug = strings.TrimSuffix(slug, "-")
	}

	return CategorySlug{value: slug}
}

// removeAccents removes accents from a string.
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// MustCategorySlug creates a new CategorySlug and panics if invalid.
func MustCategorySlug(slug string) CategorySlug {
	cs, err := NewCategorySlug(slug)
	if err != nil {
		panic(err)
	}
	return cs
}

// Value returns the category slug as a string.
func (cs CategorySlug) Value() string {
	return cs.value
}

// String returns the category slug as a string (implements fmt.Stringer).
func (cs CategorySlug) String() string {
	return cs.value
}

// Equals checks if two CategorySlug values are equal.
func (cs CategorySlug) Equals(other CategorySlug) bool {
	return cs.value == other.value
}

// IsEmpty checks if the category slug is empty.
func (cs CategorySlug) IsEmpty() bool {
	return cs.value == ""
}
