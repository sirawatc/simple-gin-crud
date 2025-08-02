package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Username string `validate:"required,min=3,max=20"`
	Password string `validate:"required,min=8"`
	Age      int    `validate:"gte=18"`
	Email    string `validate:"required,email"`
	Website  string `validate:"url"`
}

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
}

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Username: "John Doe",
				Password: "password123",
				Age:      25,
				Email:    "john@example.com",
				Website:  "https://example.com",
			},
			expected: nil,
		},
		{
			name: "invalid struct",
			input: TestStruct{
				Username: "",
				Password: "123",
				Age:      15,
				Email:    "invalid-email",
				Website:  "not-a-url",
			},
			expected: []string{
				"Username is a required field",
				"Password must be at least 8 characters in length",
				"Age must be 18 or greater",
				"Email must be a valid email address",
				"Website must be a valid URL",
			},
		},
		{
			name:  "empty struct",
			input: TestStruct{},
			expected: []string{
				"Username is a required field",
				"Password is a required field",
				"Age must be 18 or greater",
				"Email is a required field",
				"Website must be a valid URL",
			},
		},
		{
			name:  "nil interface",
			input: nil,
			expected: []string{
				"validator: (nil)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := v.Validate(test.input)
			assert.Equal(t, test.expected, errors)
		})
	}
}

func TestValidator_TranslateErrors(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name     string
		input    validator.ValidationErrors
		expected []string
	}{
		{
			name: "no error found",
			input: func() validator.ValidationErrors {
				validStruct := TestStruct{
					Username: "John Doe",
					Password: "password123",
					Age:      25,
					Email:    "john@example.com",
					Website:  "https://example.com",
				}
				err := v.validate.Struct(validStruct)
				if err != nil {
					return err.(validator.ValidationErrors)
				}
				return nil
			}(),
			expected: []string{},
		},
		{
			name:     "nil error",
			input:    nil,
			expected: []string{},
		},
		{
			name: "validation errors",
			input: func() validator.ValidationErrors {
				invalidStruct := TestStruct{
					Username: "",
					Password: "123",
					Age:      15,
					Email:    "invalid-email",
					Website:  "not-a-url",
				}
				err := v.validate.Struct(invalidStruct)
				if err != nil {
					return err.(validator.ValidationErrors)
				}
				return nil
			}(),
			expected: []string{
				"Username is a required field",
				"Password must be at least 8 characters in length",
				"Age must be 18 or greater",
				"Email must be a valid email address",
				"Website must be a valid URL",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := v.TranslateErrors(test.input)
			assert.Equal(t, test.expected, errors)
		})
	}
}
