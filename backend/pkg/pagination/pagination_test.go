package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name     string
		pageStr  string
		limitStr string
		expected PaginationParams
	}{
		{
			name:     "default values",
			pageStr:  "",
			limitStr: "",
			expected: PaginationParams{Page: 1, Limit: 10},
		},
		{
			name:     "valid values",
			pageStr:  "2",
			limitStr: "20",
			expected: PaginationParams{Page: 2, Limit: 20},
		},
		{
			name:     "invalid page defaults to 1",
			pageStr:  "invalid",
			limitStr: "10",
			expected: PaginationParams{Page: 1, Limit: 10},
		},
		{
			name:     "invalid limit defaults to 10",
			pageStr:  "1",
			limitStr: "invalid",
			expected: PaginationParams{Page: 1, Limit: 10},
		},
		{
			name:     "limit exceeds maximum",
			pageStr:  "1",
			limitStr: "200",
			expected: PaginationParams{Page: 1, Limit: 100},
		},
		{
			name:     "zero page defaults to 1",
			pageStr:  "0",
			limitStr: "10",
			expected: PaginationParams{Page: 1, Limit: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePaginationParams(tt.pageStr, tt.limitStr)
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.Limit, result.Limit)
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		params   PaginationParams
		expected int
	}{
		{
			name:     "page 1, limit 10",
			params:   PaginationParams{Page: 1, Limit: 10},
			expected: 0,
		},
		{
			name:     "page 2, limit 10",
			params:   PaginationParams{Page: 2, Limit: 10},
			expected: 10,
		},
		{
			name:     "page 3, limit 20",
			params:   PaginationParams{Page: 3, Limit: 20},
			expected: 40,
		},
		{
			name:     "page 0 defaults to offset 0",
			params:   PaginationParams{Page: 0, Limit: 10},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.params.CalculateOffset()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPaginationResult(t *testing.T) {
	tests := []struct {
		name     string
		params   PaginationParams
		total    int64
		expected PaginationResult
	}{
		{
			name:   "first page, has next",
			params: PaginationParams{Page: 1, Limit: 10},
			total:  25,
			expected: PaginationResult{
				Page:       1,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
				HasNext:    true,
				HasPrev:    false,
			},
		},
		{
			name:   "middle page",
			params: PaginationParams{Page: 2, Limit: 10},
			total:  25,
			expected: PaginationResult{
				Page:       2,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
				HasNext:    true,
				HasPrev:    true,
			},
		},
		{
			name:   "last page",
			params: PaginationParams{Page: 3, Limit: 10},
			total:  25,
			expected: PaginationResult{
				Page:       3,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
				HasNext:    false,
				HasPrev:    true,
			},
		},
		{
			name:   "empty result",
			params: PaginationParams{Page: 1, Limit: 10},
			total:  0,
			expected: PaginationResult{
				Page:       1,
				Limit:      10,
				Total:      0,
				TotalPages: 1,
				HasNext:    false,
				HasPrev:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPaginationResult(tt.params, tt.total)
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.Limit, result.Limit)
			assert.Equal(t, tt.expected.Total, result.Total)
			assert.Equal(t, tt.expected.TotalPages, result.TotalPages)
			assert.Equal(t, tt.expected.HasNext, result.HasNext)
			assert.Equal(t, tt.expected.HasPrev, result.HasPrev)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		params  PaginationParams
		wantErr bool
	}{
		{
			name:    "valid params",
			params:  PaginationParams{Page: 1, Limit: 10},
			wantErr: false,
		},
		{
			name:    "invalid page",
			params:  PaginationParams{Page: 0, Limit: 10},
			wantErr: true,
		},
		{
			name:    "invalid limit",
			params:  PaginationParams{Page: 1, Limit: 0},
			wantErr: true,
		},
		{
			name:    "limit exceeds maximum",
			params:  PaginationParams{Page: 1, Limit: 101},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
