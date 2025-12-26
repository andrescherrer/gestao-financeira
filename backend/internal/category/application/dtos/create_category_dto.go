package dtos

// CreateCategoryInput represents the input for creating a new category.
type CreateCategoryInput struct {
	UserID      string
	Name        string
	Description string
}

// CreateCategoryOutput represents the output after creating a category.
type CreateCategoryOutput struct {
	CategoryID  string
	UserID      string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   string
}
