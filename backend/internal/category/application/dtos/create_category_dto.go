package dtos

// CreateCategoryInput represents the input for creating a new category.
type CreateCategoryInput struct {
	UserID      string
	Name        string
	Description string
}

// CreateCategoryOutput represents the output after creating a category.
type CreateCategoryOutput struct {
	CategoryID  string `json:"category_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}
