package dtos

import "gestao-financeira/backend/pkg/pagination"

// ListCategoriesInput represents the input for listing categories.
type ListCategoriesInput struct {
	UserID   string
	IsActive *bool  // Optional filter for active status
	Page     string `json:"page,omitempty"`  // Query parameter
	Limit    string `json:"limit,omitempty"` // Query parameter
}

// ListCategoriesOutput represents the output for listing categories.
type ListCategoriesOutput struct {
	Categories []GetCategoryOutput          `json:"categories"`
	Count      int64                        `json:"count"`
	Pagination *pagination.PaginationResult `json:"pagination,omitempty"`
}
