package dtos

// ListCategoriesInput represents the input for listing categories.
type ListCategoriesInput struct {
	UserID   string
	IsActive *bool // Optional filter for active status
}

// ListCategoriesOutput represents the output for listing categories.
type ListCategoriesOutput struct {
	Categories []GetCategoryOutput
	Count      int64
}
