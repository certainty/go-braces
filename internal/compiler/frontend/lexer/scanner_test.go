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
