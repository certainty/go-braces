package lexer_test

import (
	"fmt"
	"testing"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/input"
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
			input:         "+",
			expectedType:  lexer.TOKEN_PLUS,
			expectedText:  "+",
			expectedValue: nil,
		},
		{
			input:         "-",
			expectedType:  lexer.TOKEN_MINUS,
			expectedText:  "-",
			expectedValue: nil,
		},
		{
			input:         "/",
			expectedType:  lexer.TOKEN_SLASH,
			expectedText:  "/",
			expectedValue: nil,
		},
		{

			input:         "%",
			expectedType:  lexer.TOKEN_MOD,
			expectedText:  "%",
			expectedValue: nil,
		},
		{
			input:         "<",
			expectedType:  lexer.TOKEN_LT,
			expectedText:  "<",
			expectedValue: nil,
		},
		{
			input:         ">",
			expectedType:  lexer.TOKEN_GT,
			expectedText:  ">",
			expectedValue: nil,
		},
		{
			input:         "&",
			expectedType:  lexer.TOKEN_AMPERSAND,
			expectedText:  "&",
			expectedValue: nil,
		},
		{
			input:         "|",
			expectedType:  lexer.TOKEN_PIPE,
			expectedText:  "|",
			expectedValue: nil,
		},

		{
			input:         ",",
			expectedType:  lexer.TOKEN_COMMA,
			expectedText:  ",",
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

			input:         "^",
			expectedType:  lexer.TOKEN_CARET,
			expectedText:  "^",
			expectedValue: nil,
		},
		{
			input:         "**",
			expectedType:  lexer.TOKEN_POWER,
			expectedText:  "**",
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
		{

			input:         "<=",
			expectedType:  lexer.TOKEN_LT_EQUAL,
			expectedText:  "<=",
			expectedValue: nil,
		},
		{

			input:         "||",
			expectedType:  lexer.TOKEN_PIPE_PIPE,
			expectedText:  "||",
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
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "35",
			expectedValue: int64(35),
		},
		{

			input:         "35.34",
			expectedType:  lexer.TOKEN_FLOAT,
			expectedText:  "35.34",
			expectedValue: 35.34,
		},
		{

			input:         "#b0",
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "#b0",
			expectedValue: int64(0),
		},
		{

			input:         "#d1344",
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "#d1344",
			expectedValue: int64(1344),
		},
		{

			input:         "#b01011",
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "#b01011",
			expectedValue: int64(11),
		},
		{

			input:         "#x2f2f",
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "#x2f2f",
			expectedValue: int64(0x2f2f),
		},

		{

			input:         "#o777",
			expectedType:  lexer.TOKEN_INTEGER,
			expectedText:  "#o777",
			expectedValue: int64(511),
		},

		// strings
		{

			input:         "\"hello world\"",
			expectedType:  lexer.TOKEN_STRING,
			expectedText:  "\"hello world\"",
			expectedValue: "hello world",
		},

		{

			input:         "\"hello \\\" world\"",
			expectedType:  lexer.TOKEN_STRING,
			expectedText:  "\"hello \\\" world\"",
			expectedValue: "hello \" world",
		},

		// chars
		{
			input:         "#\\a",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\a",
			expectedValue: lexer.CodePoint{'a'},
		},
		{
			input:         "#\\٭",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\٭",
			expectedValue: lexer.CodePoint{'٭'},
		},
		{
			input:         "#\\u1324",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\u1324",
			expectedValue: lexer.CodePoint{rune(1324)},
		},
		{
			input:         "#\\u0024",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\u0024",
			expectedValue: lexer.CodePoint{rune(24)},
		},
		{
			input:         "#\\x00002f",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\x00002f",
			expectedValue: lexer.CodePoint{rune(0x2f)},
		},
		{

			input:         "#\\newline",
			expectedType:  lexer.TOKEN_CHARACTER,
			expectedText:  "#\\newline",
			expectedValue: lexer.CodePoint{'\n'},
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

			input:         "from",
			expectedType:  lexer.TOKEN_FROM,
			expectedText:  "from",
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
			input:         "break",
			expectedType:  lexer.TOKEN_BREAK,
			expectedText:  "break",
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
			input:         "api",
			expectedType:  lexer.TOKEN_API,
			expectedText:  "api",
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
			input:         "setter",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "setter",
			expectedValue: nil,
		},
		{
			input:         "set'",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "set'",
			expectedValue: nil,
		},
		{
			input:         "relse",
			expectedType:  lexer.TOKEN_IDENTIFIER,
			expectedText:  "relse",
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
			inp := input.NewStringInput("test", tc.input)
			s := lexer.New(inp)
			token, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, token.Type)
			assert.Equal(t, tc.expectedText, string(token.Text))
			assert.Equal(t, tc.expectedValue, token.Value)
		})
	}
}

func TestScannerMultipleTokens(t *testing.T) {
	testCases := []struct {
		input    string
		expected []lexer.TokenType
	}{
		{
			input:    "",
			expected: []lexer.TokenType{lexer.TOKEN_EOF},
		},
		{
			input:    "3 + 4",
			expected: []lexer.TokenType{lexer.TOKEN_INTEGER, lexer.TOKEN_PLUS, lexer.TOKEN_INTEGER, lexer.TOKEN_EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Scanning %s", tc.input), func(t *testing.T) {
			inp := input.NewStringInput("test", tc.input)
			s := lexer.New(inp)
			tokens := []lexer.TokenType{}
			for {
				token, err := s.NextToken()
				assert.NoError(t, err)
				tokens = append(tokens, token.Type)
				if token.Type == lexer.TOKEN_EOF {
					break
				}
			}
			assert.Equal(t, tc.expected, tokens)
		})
	}

}
