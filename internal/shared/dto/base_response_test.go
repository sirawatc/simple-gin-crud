package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseResponse(t *testing.T) {
	tests := []struct {
		name     string
		code     Code
		data     interface{}
		expected *BaseResponse
	}{
		{
			name: "success response with data",
			code: Success,
			data: map[string]string{"key": "value"},
			expected: &BaseResponse{
				Code:    Success,
				Message: CodeMessage[Success],
				Data:    map[string]string{"key": "value"},
			},
		},
		{
			name: "success response with string data",
			code: Success,
			data: "simple string data",
			expected: &BaseResponse{
				Code:    Success,
				Message: CodeMessage[Success],
				Data:    "simple string data",
			},
		},
		{
			name: "success response with slice data",
			code: Success,
			data: []string{"item1", "item2", "item3"},
			expected: &BaseResponse{
				Code:    Success,
				Message: CodeMessage[Success],
				Data:    []string{"item1", "item2", "item3"},
			},
		},
		{
			name: "success response with nil data",
			code: Success,
			data: nil,
			expected: &BaseResponse{
				Code:    Success,
				Message: CodeMessage[Success],
				Data:    nil,
			},
		},
		{
			name: "bad request response without data",
			code: BadRequest,
			data: nil,
			expected: &BaseResponse{
				Code:    BadRequest,
				Message: CodeMessage[BadRequest],
				Data:    nil,
			},
		},
		{
			name: "not found response without data",
			code: NotFound,
			data: nil,
			expected: &BaseResponse{
				Code:    NotFound,
				Message: CodeMessage[NotFound],
				Data:    nil,
			},
		},
		{
			name: "unprocessable entity response without data",
			code: UnprocessableEntity,
			data: nil,
			expected: &BaseResponse{
				Code:    UnprocessableEntity,
				Message: CodeMessage[UnprocessableEntity],
				Data:    nil,
			},
		},
		{
			name: "internal error response without data",
			code: InternalError,
			data: nil,
			expected: &BaseResponse{
				Code:    InternalError,
				Message: CodeMessage[InternalError],
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildBaseResponse(tt.code, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}
