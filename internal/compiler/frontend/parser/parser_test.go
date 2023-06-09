package parser_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple addition",
			input:    "1 + 2",
			expected: "(+ 1 2)",
		},
		{
			name:     "simple subtraction",
			input:    "3 - 2",
			expected: "(- 3 2)",
		},
		{
			name:     "simple multiplication",
			input:    "2 * 3",
			expected: "(* 2 3)",
		},
		{
			name:     "simple division",
			input:    "6 / 2",
			expected: "(/ 6 2)",
		},
		{
			name:     "exponentiation",
			input:    "2 ^ 3",
			expected: "(^ 2 3)",
		},
		{
			name:     "unary plus",
			input:    "+2",
			expected: "(+ 2)",
		},
		{
			name:     "unary minus",
			input:    "-2",
			expected: "(- 2)",
		},
		{
			name:     "parentheses",
			input:    "(1 + 2) * 3",
			expected: "(* (+ 1 2) 3)",
		},
		{
			name:     "precedence",
			input:    "1 + 2 * 3 ^ 4 ",
			expected: "(+ 1 (* 2 (^ 3 4)))",
		},
		{
			name:     "right associativity",
			input:    "2 ^ 3 ^ 2",
			expected: "(^ 2 (^ 3 2))",
		},
		{
			name:     "left associativity",
			input:    "2 + 3 + 2",
			expected: "(+ (+ 2 3) 2)",
		},
	}

	astWriter := ast.NewASTWriter()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(compiler_introspection.NewNullInstrumentation())
			input := input.NewStringInput("test", test.input)
			ast, err := parser.Parse(input)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(ast.Nodes)) // single expression
			if err != nil {
				result := astWriter.WriteNode(ast.Nodes[0])
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

// let's test errors
func TestParser_Parse_Errors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parser.ParseErrorId
	}{
		{
			name:     "unexpected EOF",
			input:    "1 + ",
			expected: parser.ParseErrorIdUnexpectedEOF,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(compiler_introspection.NewNullInstrumentation())
			input := input.NewStringInput("test", test.input)
			_, err := p.Parse(input)
			assert.Error(t, err)

			if err != nil {
				allErrors := err.(parser.ParseErrors)
				parseError := allErrors.Errors[0]
				assert.Equal(t, test.expected, parseError.Id)
			}
		})
	}

}
