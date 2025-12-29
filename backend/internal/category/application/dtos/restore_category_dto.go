package dtos

// RestoreCategoryInput represents the input for restoring a soft-deleted category.
type RestoreCategoryInput struct {
	CategoryID string `json:"category_id" validate:"required,uuid"`
}

// RestoreCategoryOutput represents the output after category restoration.
type RestoreCategoryOutput struct {
	Message    string `json:"message"`
	CategoryID string `json:"category_id"`
}
