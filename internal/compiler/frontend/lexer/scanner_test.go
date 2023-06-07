package lexer_test

import (
	"fmt"
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	testCases := []struct {
		input        string
		expectedType lexer.TokenType
		expectedText string
	}{
		{
			input:        "",
			expectedType: lexer.TOKEN_EOF,
			expectedText: "",
		},
		{
			input:        "!",
			expectedType: lexer.TOKEN_BANG,
			expectedText: "!",
		},
		{
			input:        "*",
			expectedType: lexer.TOKEN_STAR,
			expectedText: "*",
		},
		{
			input:        "*",
			expectedType: lexer.TOKEN_STAR,
			expectedText: "*",
		},
		{
			input:        "{",
			expectedType: lexer.TOKEN_LBRACE,
			expectedText: "{",
		},
		{

			input:        "}",
			expectedType: lexer.TOKEN_RBRACE,
			expectedText: "}",
		},
		{

			input:        ":",
			expectedType: lexer.TOKEN_COLON,
			expectedText: ":",
		},
		{

			input:        "::",
			expectedType: lexer.TOKEN_COLON_COLON,
			expectedText: "::",
		},
		{

			input:        "==",
			expectedType: lexer.TOKEN_EQUAL_EQUAL,
			expectedText: "==",
		},
		{

			input:        ">=",
			expectedType: lexer.TOKEN_GT_EQUAL,
			expectedText: ">=",
		},

		// whitespaces
		{

			input:        "   >=",
			expectedType: lexer.TOKEN_GT_EQUAL,
			expectedText: ">=",
		},
		{

			input:        "//this is a comment\n >=",
			expectedType: lexer.TOKEN_GT_EQUAL,
			expectedText: ">=",
		},

		// numbers
		{

			input:        "35",
			expectedType: lexer.TOKEN_NUMBER,
			expectedText: "35",
		},
		{

			input:        "35.34",
			expectedType: lexer.TOKEN_NUMBER,
			expectedText: "35.34",
		},
		// strings
		{

			input:        "\"hello world\"",
			expectedType: lexer.TOKEN_STRING,
			expectedText: "\"hello world\"",
		},

		// chars
		{
			input:        "'a'",
			expectedType: lexer.TOKEN_CHARACTER,
			expectedText: "'a'",
		},
		// keywords
		{
			input:        "defer",
			expectedType: lexer.TOKEN_DEFER,
			expectedText: "defer",
		},
		{
			input:        "true",
			expectedType: lexer.TOKEN_TRUE,
			expectedText: "true",
		},
		{
			input:        "false",
			expectedType: lexer.TOKEN_FALSE,
			expectedText: "false",
		},
		{
			input:        "otherwise",
			expectedType: lexer.TOKEN_OTHERWISE,
			expectedText: "otherwise",
		},
		{
			input:        "for",
			expectedType: lexer.TOKEN_FOR,
			expectedText: "for",
		},
		{
			input:        "return",
			expectedType: lexer.TOKEN_RETURN,
			expectedText: "return",
		},
		{
			input:        "proc",
			expectedType: lexer.TOKEN_PROC,
			expectedText: "proc",
		},
		{
			input:        "if",
			expectedType: lexer.TOKEN_IF,
			expectedText: "if",
		},
		{
			input:        "else",
			expectedType: lexer.TOKEN_ELSE,
			expectedText: "else",
		},
		{
			input:        "var",
			expectedType: lexer.TOKEN_VAR,
			expectedText: "var",
		},
		{
			input:        "let",
			expectedType: lexer.TOKEN_LET,
			expectedText: "let",
		},
		{
			input:        "export",
			expectedType: lexer.TOKEN_EXPORT,
			expectedText: "export",
		},
		{
			input:        "import",
			expectedType: lexer.TOKEN_IMPORT,
			expectedText: "import",
		},
		{
			input:        "someIdentifier",
			expectedType: lexer.TOKEN_IDENTIFIER,
			expectedText: "someIdentifier",
		},
		{
			input:        "someIdent3_3_3_ifier",
			expectedType: lexer.TOKEN_IDENTIFIER,
			expectedText: "someIdent3_3_3_ifier",
		},
		{
			input:        "foo'",
			expectedType: lexer.TOKEN_IDENTIFIER,
			expectedText: "foo'",
		},
		{
			input:        "_",
			expectedType: lexer.TOKEN_IDENTIFIER,
			expectedText: "_",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Scanning %s", tc.input), func(t *testing.T) {
			s := lexer.NewFromString(tc.input, location.NewStringOrigin("test"))
			token, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, token.Type)
			assert.Equal(t, tc.expectedText, string(token.Text))
		})
	}
}
