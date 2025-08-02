package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*http.Request)
		expected string
	}{
		{
			name: "request ID provided",
			setup: func(r *http.Request) {
				r.Header.Set(RequestIDHeader, "existing-req-id")
			},
			expected: "existing-req-id",
		},
		{
			name: "empty request ID",
			setup: func(r *http.Request) {
				r.Header.Set(RequestIDHeader, "")
			},
			expected: "",
		},
		{
			name:     "request ID not provided",
			setup:    nil,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(RequestIDMiddleware())

			router.GET("/test", func(c *gin.Context) {
				requestID := GetRequestID(c.Request.Context())

				if tc.expected != "" {
					assert.Equal(t, tc.expected, requestID)
				} else {
					err := uuid.Validate(requestID)
					assert.NoError(t, err)
				}

				c.JSON(http.StatusOK, nil)
			})

			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			if tc.setup != nil {
				tc.setup(req)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(context.Context) context.Context
		expected string
	}{
		{
			name: "request ID provided",
			setup: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, requestIDKey{}, "test-request-id")
			},
			expected: "test-request-id",
		},
		{
			name: "empty request ID",
			setup: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, requestIDKey{}, "")
			},
			expected: "",
		},
		{
			name: "invalid type",
			setup: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, requestIDKey{}, 123)
			},
			expected: "",
		},
		{
			name:     "without request ID in context",
			setup:    nil,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			if tc.setup != nil {
				ctx = tc.setup(ctx)
			}

			result := GetRequestID(ctx)
			assert.Equal(t, tc.expected, result)
		})
	}
}
