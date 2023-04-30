package reader_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/stretchr/testify/assert"
)

func TestReadBoolean(t *testing.T) {
	r := reader.NewReader(introspection.NullAPI())
	result, err := r.Read(input.NewStringInput("TESTS", "#t"))

	assert.NoError(t, err)
	assert.Len(t, result.Data, 1)
	assert.IsType(t, reader.DatumBool{}, result.Data[0])
	assert.Equal(t, true, result.Data[0].(reader.DatumBool).Value)
}

// func TestSkipIrrelevantTokens(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		input    string
// 		expected rune
// 	}{
// 		{
// 			name:     "Skip intraline whitespace",
// 			input:    "   \t\t  1",
// 			expected: '1',
// 		},
// 		{
// 			name:     "Skip line endings",
// 			input:    "\n\n\r\n1",
// 			expected: '1',
// 		},
// 		{
// 			name:     "Skip comments",
// 			input:    "; comment line 1\n; comment line 2\n1",
// 			expected: '1',
// 		},
// 		{
// 			name:     "Skip nested comments",
// 			input:    "#| comment |# 1",
// 			expected: '1',
// 		},
// 		{
// 			name:     "Skip mixed irrelevant tokens",
// 			input:    "  ;comment\n\n\t#| comment |#\n \t 1",
// 			expected: '1',
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			scanner := reader.NewScanner(strings.NewReader(tc.input))
// 			err := scanner.SkipIrrelevant()

// 			assert.NoError(t, err)

// 			tok, err := scanner.Next()
// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.expected, tok)
// 		})
// 	}
// }
