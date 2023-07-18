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
			expected: "(binary-expr (tok  ADD \"+\") (basic-lit-expr (tok  FIXNUM \"1\")) (basic-lit-expr (tok  FIXNUM \"2\")))",
		},
		// {
		// 	name:     "simple subtraction",
		// 	input:    "3 - 2",
		// 	expected: "(- 3 2)",
		// },
		// {
		// 	name:     "simple multiplication",
		// 	input:    "2 * 3",
		// 	expected: "(* 2 3)",
		// },
		// {
		// 	name:     "simple division",
		// 	input:    "6 / 2",
		// 	expected: "(/ 6 2)",
		// },
		// {
		// 	name:     "exponentiation",
		// 	input:    "2 ** 3",
		// 	expected: "(** 2 3)",
		// },
		// {
		// 	name:     "unary plus",
		// 	input:    "+2",
		// 	expected: "(+ 2)",
		// },
		// {
		// 	name:     "unary minus",
		// 	input:    "-2",
		// 	expected: "(- 2)",
		// },
		// {
		// 	name:     "parentheses",
		// 	input:    "(1 + 2) * 3",
		// 	expected: "(* (+ 1 2) 3)",
		// },
		// {
		// 	name:     "precedence",
		// 	input:    "1 + 2 * 3 ** 4 ",
		// 	expected: "(+ 1 (* 2 (** 3 4)))",
		// },
		// {
		// 	name:     "more precedence",
		// 	input:    "3 ** 4 * 3 + 4",
		// 	expected: "(+ (* (** 3 4) 3) 4)",
		// },
		// {
		// 	name:     "even more precedence",
		// 	input:    "3 ** 4 * 3 - 4",
		// 	expected: "(- (* (** 3 4) 3) 4)",
		// },
		// {
		// 	name:     "right associativity",
		// 	input:    "1 ** 4 ** 2",
		// 	expected: "(** 1 (** 4 2))",
		// },
		// {
		// 	name:     "grouping",
		// 	input:    "(2 ** 3) ** 2",
		// 	expected: "(** (** 2 3) 2)",
		// },

		// {
		// 	name:     "left associativity",
		// 	input:    "2 + 3 + 2",
		// 	expected: "(+ (+ 2 3) 2)",
		// },
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

// // let's test errors
// func TestParser_Parse_Errors(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		expected parser.ParseErrorId
// 	}{
// 		{
// 			name:     "unexpected EOF",
// 			input:    "1 + ",
// 			expected: parser.ParseErrorIdUnexpectedEOF,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			p := parser.NewParser(compiler_introspection.NewNullInstrumentation())
// 			input := lexer.NewStringInput("test", test.input)
// 			_, err := p.Parse(input)
// 			assert.Error(t, err)

// 			if err != nil {
// 				allErrors := err.(parser.ParseErrors)
// 				parseError := allErrors.Errors[0]
// 				assert.Equal(t, test.expected, parseError.Id)
// 			}
// 		})
// 	}

// }
