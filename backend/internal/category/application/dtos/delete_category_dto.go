package dtos

// DeleteCategoryInput represents the input for deleting a category.
type DeleteCategoryInput struct {
	CategoryID string
}

// DeleteCategoryOutput represents the output after deleting a category.
type DeleteCategoryOutput struct {
	Message    string
	CategoryID string
}
