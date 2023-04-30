package reader_test

import (
	"strings"
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/location"
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
			result, err := reader.Read(location.NewStringInput("TESTS", tt.input))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSkipIrrelevantTokens(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected rune
	}{
		{
			name:     "Skip intraline whitespace",
			input:    "   \t\t  1",
			expected: '1',
		},
		{
			name:     "Skip line endings",
			input:    "\n\n\r\n1",
			expected: '1',
		},
		{
			name:     "Skip comments",
			input:    "; comment line 1\n; comment line 2\n1",
			expected: '1',
		},
		{
			name:     "Skip nested comments",
			input:    "#| comment |# 1",
			expected: '1',
		},
		{
			name:     "Skip mixed irrelevant tokens",
			input:    "  ;comment\n\n\t#| comment |#\n \t 1",
			expected: '1',
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scanner := reader.NewScanner(strings.NewReader(tc.input))
			err := scanner.SkipIrrelevant()

			assert.NoError(t, err)

			tok, err := scanner.Next()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tok)
		})
	}
}
