package pagination

import (
	"fmt"
	"strconv"
)

// PaginationParams represents pagination parameters from request.
type PaginationParams struct {
	Page  int // 1-based page number
	Limit int // Items per page
}

// PaginationResult represents pagination metadata in response.
type PaginationResult struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ParsePaginationParams parses pagination parameters from query string.
// Defaults: page=1, limit=10
// Max limit: 100
func ParsePaginationParams(pageStr, limitStr string) PaginationParams {
	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			// Enforce maximum limit
			if limit > 100 {
				limit = 100
			}
		}
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// CalculateOffset calculates the database offset from page and limit.
func (p PaginationParams) CalculateOffset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

// BuildPaginationResult builds pagination metadata from params and total count.
func BuildPaginationResult(params PaginationParams, total int64) PaginationResult {
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationResult{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// Validate validates pagination parameters.
func (p PaginationParams) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if p.Limit < 1 {
		return fmt.Errorf("limit must be greater than 0")
	}
	if p.Limit > 100 {
		return fmt.Errorf("limit cannot exceed 100")
	}
	return nil
}
