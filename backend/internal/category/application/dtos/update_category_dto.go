package dtos

// UpdateCategoryInput represents the input for updating a category.
type UpdateCategoryInput struct {
	CategoryID  string
	Name        *string
	Description *string
}

// UpdateCategoryOutput represents the output after updating a category.
type UpdateCategoryOutput struct {
	CategoryID  string
	UserID      string
	Name        string
	Description string
	IsActive    bool
	UpdatedAt   string
}
