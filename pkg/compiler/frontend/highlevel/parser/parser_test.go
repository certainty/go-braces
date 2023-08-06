package parser_test

import (
	"testing"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/parser"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
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
			expected: "(binary-expr + 1 2)",
		},
		{
			name:     "simple subtraction",
			input:    "3 - 2",
			expected: "(binary-expr - 3 2)",
		},
		{
			name:     "simple multiplication",
			input:    "2 * 3",
			expected: "(binary-expr * 2 3)",
		},
		{
			name:     "simple division",
			input:    "6 / 2",
			expected: "(binary-expr / 6 2)",
		},
		{
			name:     "exponentiation",
			input:    "2 ** 3",
			expected: "(binary-expr ** 2 3)",
		},
		{
			name:     "unary plus",
			input:    "+2",
			expected: "(unary-expr + 2)",
		},
		{
			name:     "unary minus",
			input:    "-2",
			expected: "(unary-expr - 2)",
		},
		{
			name:     "parentheses",
			input:    "(1 + 2) * 3",
			expected: "(binary-expr * (binary-expr + 1 2) 3)",
		},
		{
			name:     "mixed expressions",
			input:    "3 ** 4 * 3 + (-4)",
			expected: "(binary-expr + (binary-expr * (binary-expr ** 3 4) 3) (unary-expr - 4))",
		},

		{
			name:     "precedence",
			input:    "1 + 2 * 3 ** 4 ",
			expected: "(binary-expr + 1 (binary-expr * 2 (binary-expr ** 3 4)))",
		},
		{
			name:     "more precedence",
			input:    "3 ** 4 * 3 + 4",
			expected: "(binary-expr + (binary-expr * (binary-expr ** 3 4) 3) 4)",
		},
		{
			name:     "even more precedence",
			input:    "3 ** 4 * 3 - 4",
			expected: "(binary-expr - (binary-expr * (binary-expr ** 3 4) 3) 4)",
		},
		{
			name:     "right associativity",
			input:    "1 ** 4 ** 2",
			expected: "(binary-expr ** 1 (binary-expr ** 4 2))",
		},
		{
			name:     "grouping",
			input:    "(2 ** 3) ** 2",
			expected: "(binary-expr ** (binary-expr ** 2 3) 2)",
		},

		{
			name:     "left associativity",
			input:    "2 + 3 + 2",
			expected: "(binary-expr + (binary-expr + 2 3) 2)",
		},
	}

	printOptions := ast.PrintCanonical()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(compiler_introspection.NewNullInstrumentation())
			input := lexer.NewStringInput("test", test.input)
			expr := parser.ParseExpression(input)
			assert.Empty(t, parser.Errors())
			result := ast.Print(expr, printOptions)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseProcedure(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple procedure",
			input:    "proc main() { 3+3 }",
			expected: "(source (proc-decl main (block-expr (expr-stmt (binary-expr + 3 3)))))",
		},
	}

	printOptions := ast.PrintCanonical()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := parser.NewParser(compiler_introspection.NewNullInstrumentation())
			input := lexer.NewStringInput("test", test.input)
			decl, _ := parser.Parse(input)
			assert.Empty(t, parser.Errors())
			result := ast.Print(*decl, printOptions)
			assert.Equal(t, test.expected, result)
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
			expected: parser.ParseErrorIdUnexpectedToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(compiler_introspection.NewNullInstrumentation())
			input := lexer.NewStringInput("test", test.input)
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
