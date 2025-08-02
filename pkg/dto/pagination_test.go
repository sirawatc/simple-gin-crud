package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginationRequest(t *testing.T) {
	tests := []struct {
		name        string
		page        string
		pageSize    string
		expected    *PaginationRequest
		expectError bool
	}{
		{
			name:     "valid parameters",
			page:     "2",
			pageSize: "15",
			expected: &PaginationRequest{
				Page:     2,
				PageSize: 15,
			},
			expectError: false,
		},
		{
			name:     "empty strings should use defaults",
			page:     "",
			pageSize: "",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: false,
		},
		{
			name:     "invalid page should return error",
			page:     "invalid",
			pageSize: "10",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
		{
			name:     "invalid pageSize should return error",
			page:     "1",
			pageSize: "invalid",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
		{
			name:     "zero page should return error",
			page:     "0",
			pageSize: "10",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
		{
			name:     "negative page should return error",
			page:     "-1",
			pageSize: "10",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
		{
			name:     "zero pageSize should return error",
			page:     "1",
			pageSize: "0",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
		{
			name:     "negative pageSize should return error",
			page:     "1",
			pageSize: "-1",
			expected: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errors := NewPaginationRequest(tt.page, tt.pageSize)

			if tt.expectError {
				assert.NotEmpty(t, errors)
			} else {
				assert.Empty(t, errors)
			}

			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.PageSize, result.PageSize)
		})
	}
}

func TestPaginationRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		request  *PaginationRequest
		expected int
	}{
		{
			name: "first page",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expected: 0,
		},
		{
			name: "second page",
			request: &PaginationRequest{
				Page:     2,
				PageSize: 10,
			},
			expected: 10,
		},
		{
			name: "third page with custom page size",
			request: &PaginationRequest{
				Page:     3,
				PageSize: 25,
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetOffset()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaginationRequest_GetLimit(t *testing.T) {
	tests := []struct {
		name     string
		request  *PaginationRequest
		expected int
	}{
		{
			name: "standard page size",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			expected: 10,
		},
		{
			name: "custom page size",
			request: &PaginationRequest{
				Page:     2,
				PageSize: 25,
			},
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.GetLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewPaginationResponse(t *testing.T) {
	tests := []struct {
		name       string
		request    *PaginationRequest
		totalItems int64
		expected   *PaginationResponse
	}{
		{
			name: "first page with items",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			totalItems: 25,
			expected: &PaginationResponse{
				Page:       1,
				PageSize:   10,
				TotalPages: 3,
				TotalItems: 25,
			},
		},
		{
			name: "middle page",
			request: &PaginationRequest{
				Page:     2,
				PageSize: 10,
			},
			totalItems: 25,
			expected: &PaginationResponse{
				Page:       2,
				PageSize:   10,
				TotalPages: 3,
				TotalItems: 25,
			},
		},
		{
			name: "last page",
			request: &PaginationRequest{
				Page:     3,
				PageSize: 10,
			},
			totalItems: 25,
			expected: &PaginationResponse{
				Page:       3,
				PageSize:   10,
				TotalPages: 3,
				TotalItems: 25,
			},
		},
		{
			name: "exact page size",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			totalItems: 10,
			expected: &PaginationResponse{
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
				TotalItems: 10,
			},
		},
		{
			name: "no items",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			totalItems: 0,
			expected: &PaginationResponse{
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
				TotalItems: 0,
			},
		},
		{
			name: "items less than page size",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			totalItems: 5,
			expected: &PaginationResponse{
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
				TotalItems: 5,
			},
		},
		{
			name: "items requiring multiple pages",
			request: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
			totalItems: 100,
			expected: &PaginationResponse{
				Page:       1,
				PageSize:   10,
				TotalPages: 10,
				TotalItems: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPaginationResponse(tt.request, tt.totalItems)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewPaginationDataResponse(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	request := &PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
	totalItems := int64(25)

	result := NewPaginationDataResponse(data, request, totalItems)

	assert.Equal(t, data, result.Items)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.PageSize)
	assert.Equal(t, 3, result.Pagination.TotalPages)
	assert.Equal(t, int64(25), result.Pagination.TotalItems)
}
