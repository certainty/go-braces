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
