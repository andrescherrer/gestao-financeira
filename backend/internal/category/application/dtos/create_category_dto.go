package dtos

// CreateCategoryInput represents the input for creating a new category.
type CreateCategoryInput struct {
	UserID      string
	Name        string `json:"name" validate:"required,min=2,max=100,no_sql_injection,no_xss,utf8"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000,no_sql_injection,no_xss,utf8"`
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
