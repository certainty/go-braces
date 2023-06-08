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
		input         string
		expectedType  lexer.TokenType
		expectedText  string
		expectedValue interface{}
	}{
		{
			input:         "",
			expectedType:  lexer.TOKEN_EOF,
			expectedText:  "",
			expectedValue: nil,
		},
		{
			input:         "!",
			expectedType:  lexer.TOKEN_BANG,
			expectedText:  "!",
			expectedValue: nil,
		},
		{
			input:         "*",
			expectedType:  lexer.TOKEN_STAR,
			expectedText:  "*",
			expectedValue: nil,
		},
		{
			input:         "*",
			expectedType:  lexer.TOKEN_STAR,
			expectedText:  "*",
			expectedValue: nil,
		},
		{
			input:         "{",
			expectedType:  lexer.TOKEN_LBRACE,
			expectedText:  "{",
			expectedValue: nil,
		},
		{

			input:         "}",
			expectedType:  lexer.TOKEN_RBRACE,
			expectedText:  "}",
			expectedValue: nil,
		},
		{

			input:         ":",
			expectedType:  lexer.TOKEN_COLON,
			expectedText:  ":",
			expectedValue: nil,
		},
		{

			input:         "::",
			expectedType:  lexer.TOKEN_COLON_COLON,
			expectedText:  "::",
			expectedValue: nil,
		},
		{

			input:         "==",
			expectedType:  lexer.TOKEN_EQUAL_EQUAL,
			expectedText:  "==",
			expectedValue: nil,
		},
		{

			input:         ">=",
			expectedType:  lexer.TOKEN_GT_EQUAL,
			expectedText:  ">=",
			expectedValue: nil,
		},

		// whitespaces
		{

			input:         "   >=",
			expectedType:  lexer.TOKEN_GT_EQUAL,
			expectedText:  ">=",
			expectedValue: nil,
		},
		{

			input:         "//this is a comment\n >=",
			expectedType:  lexer.TOKEN_GT_EQUAL,
			expectedText:  ">=",
			expectedValue: nil,
		},

		// numbers
		{

			input:         "35",
			expectedType:  lexer.TOKEN_NUMBER,
			expectedText:  "35",
			expectedValue: int64(35),
		},
		{

			input:         "35.34",
			expectedType:  lexer.TOKEN_NUMBER,
			expectedText:  "35.34",
			expectedValue: 35.34,
		},
		// strings
		{

			input:         "\"hello world\"",
			expectedType:  lexer.TOKEN_STRING,
			expectedText:  "\"hello world\"",
			expectedValue: "hello world",
		},

		// chars
		{
			input:         "#\\a",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\a",
			expectedValue: 'a',
		},
		{
			input:         "#\\٭",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\٭",
			expectedValue: '٭',
		},
		{
			input:         "#\\u1324",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\u1324",
			expectedValue: rune(1324),
		},
		{
			input:         "#\\u0024",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\u0024",
			expectedValue: rune(24),
		},
		{
			input:         "#\\x00002f",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\x00002f",
			expectedValue: rune(0x2f),
		},
		{

			input:         "#\\newline",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\newline",
			expectedValue: '\n',
		},

		// keywords
		{
			input:         "defer",
			expectedType:  lexer.TOKEN_DEFER,
			expectedText:  "defer",
			expectedValue: nil,
		},
		{
			input:         "true",
			expectedType:  lexer.TOKEN_TRUE,
			expectedText:  "true",
			expectedValue: true,
		},
		{
			input:         "false",
			expectedType:  lexer.TOKEN_FALSE,
			expectedText:  "false",
			expectedValue: false,
		},
		{
			input:         "for",
			expectedType:  lexer.TOKEN_FOR,
			expectedText:  "for",
			expectedValue: nil,
		},
		{
			input:         "return",
			expectedType:  lexer.TOKEN_RETURN,
			expectedText:  "return",
			expectedValue: nil,
		},
		{
			input:         "proc",
			expectedType:  lexer.TOKEN_PROC,
			expectedText:  "proc",
			expectedValue: nil,
		},
		{
			input:         "if",
			expectedType:  lexer.TOKEN_IF,
			expectedText:  "if",
			expectedValue: nil,
		},
		{
			input:         "else",
			expectedType:  lexer.TOKEN_ELSE,
			expectedText:  "else",
			expectedValue: nil,
		},
		{
			input:         "var",
			expectedType:  lexer.TOKEN_VAR,
			expectedText:  "var",
			expectedValue: nil,
		},
		{
			input:         "let",
			expectedType:  lexer.TOKEN_LET,
			expectedText:  "let",
			expectedValue: nil,
		},
		{
			input:         "set",
			expectedType:  lexer.TOKEN_SET,
			expectedText:  "set",
			expectedValue: nil,
		},
		{
			input:         "export",
			expectedType:  lexer.TOKEN_EXPORT,
			expectedText:  "export",
			expectedValue: nil,
		},
		{
			input:         "import",
			expectedType:  lexer.TOKEN_IMPORT,
			expectedText:  "import",
			expectedValue: nil,
		},
		{
			input:         "someIdentifier",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "someIdentifier",
			expectedValue: nil,
		},
		{
			input:         "someIdent3_3_3_ifier",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "someIdent3_3_3_ifier",
			expectedValue: nil,
		},
		{
			input:         "foo'",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "foo'",
			expectedValue: nil,
		},
		{
			input:         "_",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "_",
			expectedValue: nil,
		},
		{
			input:         "->",
			expectedType:  lexer.TOKEN_ARROW,
			expectedText:  "->",
			expectedValue: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Scanning %s", tc.input), func(t *testing.T) {
			s := lexer.NewFromString(tc.input, location.NewStringOrigin("test"))
			token, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, token.Type)
			assert.Equal(t, tc.expectedText, string(token.Text))
			assert.Equal(t, tc.expectedValue, token.Value)
		})
	}
}