package dtos

// PermanentDeleteCategoryInput represents the input for permanently deleting a category.
type PermanentDeleteCategoryInput struct {
	CategoryID string `json:"category_id" validate:"required,uuid"`
}

// PermanentDeleteCategoryOutput represents the output after permanent category deletion.
type PermanentDeleteCategoryOutput struct {
	Message    string `json:"message"`
	CategoryID string `json:"category_id"`
}
