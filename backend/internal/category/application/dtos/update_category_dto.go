package dtos

// UpdateCategoryInput represents the input for updating a category.
type UpdateCategoryInput struct {
	CategoryID  string
	Name        *string
	Description *string
}

// UpdateCategoryOutput represents the output after updating a category.
type UpdateCategoryOutput struct {
	CategoryID  string `json:"category_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	UpdatedAt   string `json:"updated_at"`
}
