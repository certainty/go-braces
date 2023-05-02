package reader_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/stretchr/testify/assert"
)

func ReadString(s string) (*reader.DatumAST, error) {
	r := reader.NewReader(introspection.NullAPI())
	return r.Read(input.NewStringInput("TESTS", s))
}

func TestReadBoolean(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Read #t",
			input:    "#t",
			expected: true,
		},

		{
			name:     "Read #true",
			input:    "#t",
			expected: true,
		},
		{
			name:     "Read #f",
			input:    "#f",
			expected: false,
		},
		{
			name:     "Read #false",
			input:    "#false",
			expected: false,
		},
	}

	for _, rt := range testCases {
		t.Run(rt.name, func(t *testing.T) {
			result, err := ReadString(rt.input)
			assert.NoError(t, err)
			assert.Equal(t, rt.expected, result.Data[0].(reader.DatumBool).Value)
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
		// {
		// 	name:     "Skip comments",
		// 	input:    "; comment line 1\n; comment line 2\n1",
		// 	expected: '1',
		// },
		// {
		// 	name:     "Skip nested comments",
		// 	input:    "#| comment |# 1",
		// 	expected: '1',
		// },
		{
			name:     "Skip mixed irrelevant tokens",
			input:    "  ;comment\n\n\t#| comment |#\n \t 1",
			expected: '1',
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := []rune(tc.input)
			scanner := reader.NewScanner(&input)
			err := scanner.SkipIrrelevant()

			assert.NoError(t, err)
			tok, err := scanner.Next()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tok)
		})
	}
}
