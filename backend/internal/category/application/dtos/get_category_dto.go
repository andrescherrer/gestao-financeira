package dtos

// GetCategoryOutput represents the output for getting a category.
type GetCategoryOutput struct {
	CategoryID  string `json:"category_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
