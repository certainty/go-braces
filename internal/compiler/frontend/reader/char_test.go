package reader_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/isa"
	"github.com/stretchr/testify/assert"
)

func TestNamedChar(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected rune
	}{
		{
			name:     "Read #\\space",
			input:    "#\\space",
			expected: ' ',
		},
		{
			name:     "Read #\\newline",
			input:    "#\\newline",
			expected: '\n',
		},
		{
			name:     "Read #\\return",
			input:    "#\\return",
			expected: '\r',
		},
		{
			name:     "Read #\\tab",
			input:    "#\\tab",
			expected: '\t',
		},
		{
			name:     "Read #\\alarm",
			input:    "#\\alarm",
			expected: '\u0007',
		},
		{
			name:     "Read #\\backspace",
			input:    "#\\backspace",
			expected: '\u0008',
		},
		{
			name:     "Read #\\delete",
			input:    "#\\delete",
			expected: '\u007f',
		},
		{
			name:     "Read #\\null",
			input:    "#\\null",
			expected: '\u0000',
		},
		{
			name:     "Read #\\escape",
			input:    "#\\escape",
			expected: '\u001b',
		},
	}

	for _, rt := range testCases {
		t.Run(rt.name, func(t *testing.T) {
			result, err := ReadString(rt.input)
			assert.NoError(t, err)
			assert.Equal(t, rt.expected, result.Data[0].(isa.DatumChar).Value)
		})
	}

}

func TestHexChar(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected rune
	}{
		{
			name:     "Read #\\xa",
			input:    "#\\xa",
			expected: '\n',
		},
		{
			name:     "Read #\\xababab",
			input:    "#\\xababab",
			expected: 0xababab,
		},
	}
	for _, rt := range testCases {
		t.Run(rt.name, func(t *testing.T) {
			result, err := ReadString(rt.input)
			assert.NoError(t, err)
			assert.NotNil(t, result.Data[0])
			assert.Equal(t, rt.expected, result.Data[0].(isa.DatumChar).Value)
		})
	}
}

func TestCharLiteral(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected rune
	}{
		{
			name:     "Read #\\h",
			input:    "#\\h",
			expected: 'h',
		},
		{
			name:     "Read #\\☆",
			input:    "#\\☆",
			expected: '☆',
		},
	}
	for _, rt := range testCases {
		t.Run(rt.name, func(t *testing.T) {
			result, err := ReadString(rt.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, result.Data)
			assert.Equal(t, rt.expected, result.Data[0].(isa.DatumChar).Value)
		})
	}

}
