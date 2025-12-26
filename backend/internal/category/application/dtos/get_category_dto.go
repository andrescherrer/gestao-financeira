package dtos

// GetCategoryOutput represents the output for getting a category.
type GetCategoryOutput struct {
	CategoryID  string
	UserID      string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   string
	UpdatedAt   string
}
