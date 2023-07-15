package lexer_test

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanner(t *testing.T) {
	testCases := []struct {
		input         string
		expectedType  token.Type
		expectedText  string
		expectedValue interface{}
	}{
		{
			input:         "",
			expectedType:  token.EOF,
			expectedText:  "",
			expectedValue: nil,
		},
		{
			input:         "!",
			expectedType:  token.NOT,
			expectedText:  "!",
			expectedValue: nil,
		},
		{
			input:         "*",
			expectedType:  token.MUL,
			expectedText:  "*",
			expectedValue: nil,
		},
		{
			input:         "+",
			expectedType:  token.ADD,
			expectedText:  "+",
			expectedValue: nil,
		},
		{
			input:         "-",
			expectedType:  token.SUB,
			expectedText:  "-",
			expectedValue: nil,
		},
		{
			input:         "/",
			expectedType:  token.DIV,
			expectedText:  "/",
			expectedValue: nil,
		},
		{

			input:         "%",
			expectedType:  token.REM,
			expectedText:  "%",
			expectedValue: nil,
		},
		{
			input:         "<",
			expectedType:  token.LT,
			expectedText:  "<",
			expectedValue: nil,
		},
		{
			input:         ">",
			expectedType:  token.GT,
			expectedText:  ">",
			expectedValue: nil,
		},
		{
			input:         "&",
			expectedType:  token.AND,
			expectedText:  "&",
			expectedValue: nil,
		},
		{
			input:         "|",
			expectedType:  token.OR,
			expectedText:  "|",
			expectedValue: nil,
		},
		{
			input:         ",",
			expectedType:  token.COMMA,
			expectedText:  ",",
			expectedValue: nil,
		},
		{
			input:         ";",
			expectedType:  token.SEMICOLON,
			expectedText:  ";",
			expectedValue: nil,
		},
		{
			input:         ".",
			expectedType:  token.DOT,
			expectedText:  ".",
			expectedValue: nil,
		},
		{
			input:         "{",
			expectedType:  token.LBRACE,
			expectedText:  "{",
			expectedValue: nil,
		},
		{

			input:         "}",
			expectedType:  token.RBRACE,
			expectedText:  "}",
			expectedValue: nil,
		},
		{
			input:         ":",
			expectedType:  token.COLON,
			expectedText:  ":",
			expectedValue: nil,
		},
		{
			input:         "::",
			expectedType:  token.DBLCOLON,
			expectedText:  "::",
			expectedValue: nil,
		},
		{
			input:         "^",
			expectedType:  token.XOR,
			expectedText:  "^",
			expectedValue: nil,
		},
		{
			input:         "**",
			expectedType:  token.POW,
			expectedText:  "**",
			expectedValue: nil,
		},
		{

			input:         "==",
			expectedType:  token.EQ,
			expectedText:  "==",
			expectedValue: nil,
		},
		{
			input:         ">=",
			expectedType:  token.GTE,
			expectedText:  ">=",
			expectedValue: nil,
		},
		{
			input:         "<=",
			expectedType:  token.LTE,
			expectedText:  "<=",
			expectedValue: nil,
		},
		{
			input:         ">>",
			expectedType:  token.SHR,
			expectedText:  ">>",
			expectedValue: nil,
		},
		{
			input:         "<<",
			expectedType:  token.SHL,
			expectedText:  "<<",
			expectedValue: nil,
		},
		{
			input:         "||",
			expectedType:  token.LOR,
			expectedText:  "||",
			expectedValue: nil,
		},
		{
			input:         "&&",
			expectedType:  token.LAND,
			expectedText:  "&&",
			expectedValue: nil,
		},
		{
			input:         "&^",
			expectedType:  token.AND_NOT,
			expectedText:  "&^",
			expectedValue: nil,
		},
		{
			input:         "->",
			expectedType:  token.ARROW,
			expectedText:  "->",
			expectedValue: nil,
		},
		{
			input:         "|>",
			expectedType:  token.PIPE,
			expectedText:  "|>",
			expectedValue: nil,
		},

		// whitespaces
		{
			input:         "   >=",
			expectedType:  token.GTE,
			expectedText:  ">=",
			expectedValue: nil,
		},
		{
			input:         "//this is a comment\n >=",
			expectedType:  token.GTE,
			expectedText:  ">=",
			expectedValue: nil,
		},

		// numbers
		{
			input:         "35",
			expectedType:  token.FIXNUM,
			expectedText:  "35",
			expectedValue: int(35),
		},
		{
			input:         "35.34",
			expectedType:  token.FLONUM,
			expectedText:  "35.34",
			expectedValue: 35.34,
		},
		{
			input:         "#b0",
			expectedType:  token.FIXNUM,
			expectedText:  "#b0",
			expectedValue: int(0),
		},
		{
			input:         "#d1344",
			expectedType:  token.FIXNUM,
			expectedText:  "#d1344",
			expectedValue: int(1344),
		},
		{
			input:         "#b01011",
			expectedType:  token.FIXNUM,
			expectedText:  "#b01011",
			expectedValue: int(11),
		},
		{
			input:         "#x2f2f",
			expectedType:  token.FIXNUM,
			expectedText:  "#x2f2f",
			expectedValue: int(0x2f2f),
		},

		{
			input:         "#o777",
			expectedType:  token.FIXNUM,
			expectedText:  "#o777",
			expectedValue: int(511),
		},

		// strings
		{
			input:         "\"hello world\"",
			expectedType:  token.STRING,
			expectedText:  "\"hello world\"",
			expectedValue: "hello world",
		},
		{
			input:         "\"hello \\\" world\"",
			expectedType:  token.STRING,
			expectedText:  "\"hello \\\" world\"",
			expectedValue: "hello \" world",
		},

		// chars
		{
			input:         "#\\a",
			expectedType:  token.CHAR,
			expectedText:  "#\\a",
			expectedValue: rune('a'),
		},
		{
			input:         "#\\٭",
			expectedType:  token.CHAR,
			expectedText:  "#\\٭",
			expectedValue: rune('٭'),
		},
		{
			input:         "#\\u1324",
			expectedType:  token.CHAR,
			expectedText:  "#\\u1324",
			expectedValue: rune(1324),
		},
		{
			input:         "#\\u0024",
			expectedType:  token.CHAR,
			expectedText:  "#\\u0024",
			expectedValue: rune(24),
		},
		{
			input:         "#\\x00002f",
			expectedType:  token.CHAR,
			expectedText:  "#\\x00002f",
			expectedValue: rune(0x2f),
		},
		{
			input:         "#\\newline",
			expectedType:  token.CHAR,
			expectedText:  "#\\newline",
			expectedValue: rune('\n'),
		},

		// keywords
		{
			input:         "ensure",
			expectedType:  token.ENSURE,
			expectedText:  "ensure",
			expectedValue: nil,
		},
		{
			input:         "true",
			expectedType:  token.TRUE,
			expectedText:  "true",
			expectedValue: true,
		},
		{
			input:         "false",
			expectedType:  token.FALSE,
			expectedText:  "false",
			expectedValue: false,
		},
		{
			input:         "for",
			expectedType:  token.FOR,
			expectedText:  "for",
			expectedValue: nil,
		},
		{
			input:         "return",
			expectedType:  token.RETURN,
			expectedText:  "return",
			expectedValue: nil,
		},
		{
			input:         "proc",
			expectedType:  token.PROC,
			expectedText:  "proc",
			expectedValue: nil,
		},
		{
			input:         "if",
			expectedType:  token.IF,
			expectedText:  "if",
			expectedValue: nil,
		},
		{
			input:         "else",
			expectedType:  token.ELSE,
			expectedText:  "else",
			expectedValue: nil,
		},
		{
			input:         "break",
			expectedType:  token.BREAK,
			expectedText:  "break",
			expectedValue: nil,
		},
		{
			input:         "let",
			expectedType:  token.LET,
			expectedText:  "let",
			expectedValue: nil,
		},
		{
			input:         "set",
			expectedType:  token.SET,
			expectedText:  "set",
			expectedValue: nil,
		},
		{
			input:         "api",
			expectedType:  token.API,
			expectedText:  "api",
			expectedValue: nil,
		},
		{
			input:         "import",
			expectedType:  token.IMPORT,
			expectedText:  "import",
			expectedValue: nil,
		},
		{
			input:         "someIdentifier",
			expectedType:  token.IDENTIFIER,
			expectedText:  "someIdentifier",
			expectedValue: nil,
		},
		{
			input:         "someIdent3_3_3_ifier",
			expectedType:  token.IDENTIFIER,
			expectedText:  "someIdent3_3_3_ifier",
			expectedValue: nil,
		},
		{
			input:         "setter",
			expectedType:  token.IDENTIFIER,
			expectedText:  "setter",
			expectedValue: nil,
		},
		{
			input:         "set'",
			expectedType:  token.IDENTIFIER,
			expectedText:  "set'",
			expectedValue: nil,
		},
		{
			input:         "relse",
			expectedType:  token.IDENTIFIER,
			expectedText:  "relse",
			expectedValue: nil,
		},

		{
			input:         "foo'",
			expectedType:  token.IDENTIFIER,
			expectedText:  "foo'",
			expectedValue: nil,
		},
		{
			input:         "_",
			expectedType:  token.IDENTIFIER,
			expectedText:  "_",
			expectedValue: nil,
		},
		{
			input:         "...",
			expectedType:  token.ELIPSES,
			expectedText:  "...",
			expectedValue: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Scanning %s", tc.input), func(t *testing.T) {
			inp := input.NewStringInput("test", tc.input)
			s := lexer.New(inp)
			scannedToken := s.NextToken()

			assert.NotEqual(t, scannedToken.IsIllegal(), true)
			assert.Equal(t, tc.expectedType, scannedToken.Type)
			assert.Equal(t, tc.expectedText, string(scannedToken.Text))
			assert.Equal(t, tc.expectedValue, scannedToken.LitValue)
		})
	}
}

func TestScannerMultipleTokens(t *testing.T) {
	testCases := []struct {
		input    string
		expected []token.Type
	}{
		{
			input:    "",
			expected: []token.Type{token.EOF},
		},
		{
			input:    "3 + 4",
			expected: []token.Type{token.FIXNUM, token.ADD, token.FIXNUM, token.EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Scanning %s", tc.input), func(t *testing.T) {
			inp := input.NewStringInput("test", tc.input)
			s := lexer.New(inp)
			tokens := []token.Type{}
			for {
				scannedToken := s.NextToken()
				assert.NotEqual(t, scannedToken.IsIllegal(), true)
				tokens = append(tokens, scannedToken.Type)

				if scannedToken.Type == token.EOF {
					break
				}
			}
			assert.Equal(t, tc.expected, tokens)
		})
	}

}
