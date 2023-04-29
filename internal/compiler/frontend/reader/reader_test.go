package reader_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/stretchr/testify/assert"
)

// Assuming your Read function has the following signature:
// func Read(input string) (interface{}, error)
func TestReadBoolean(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    interface{}
		expectError bool
	}{
		{
			name:        "Read #t",
			input:       "#t",
			expected:    true,
			expectError: false,
		},
		{
			name:        "Read #f",
			input:       "#f",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Read empty string",
			input:       "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Read unsupported literal",
			input:       "123",
			expected:    nil,
			expectError: true,
		},
	}

	reader := reader.NewReader(introspection.NullAPI())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte(tt.input)
			result, err := reader.Read(&input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
