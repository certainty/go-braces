package lexer

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
)

type TokenType uint8
type Token struct {
	Type     TokenType
	Text     []rune
	Value    interface{} // optionally we can also transport values back
	Location location.Location
}

const (
	TOKEN_EOF TokenType = iota

	// single char tokens
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_COMMA
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_STAR
	TOKEN_SLASH
	TOKEN_QUESTION_MARK
	TOKEN_COLON

	// one or two
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_GT
	TOKEN_GT_EQUAL
	TOKEN_LT
	TOKEN_LT_EQUAL
	TOKEN_AMPERSAND
	TOKEN_AMPERSAND_AMPERSAND
	TOKEN_PIPE
	TOKEN_PIPE_PIPE
	TOKEN_PIPE_GT
	TOKEN_COLON_COLON
	TOKEN_ARROW

	// literal
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_INTEGER
	TOKEN_FLOAT
	TOKEN_BOOLEAN
	TOKEN_CHARACTER
	TOKEN_SYMBOL

	// keywords
	TOKEN_FUN
	TOKEN_PROC
	TOKEN_PACKAGE
	TOKEN_IMPORT
	TOKEN_EXPORT

	TOKEN_DATA
	TOKEN_ALIAS

	TOKEN_VAR
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_LET
	TOKEN_SET
	TOKEN_MATCH
	TOKEN_FOR

	TOKEN_RETURN
	TOKEN_DEFER

	TOKEN_ELLIPSIS

	TOKEN_TRUE
	TOKEN_FALSE
)

func MakeToken(tokenType TokenType, text []rune, location location.Location) Token {
	return Token{Type: tokenType, Text: text, Value: nil, Location: location}
}

func MakeTokenWithValue(tokenType TokenType, text []rune, value interface{}, location location.Location) Token {
	return Token{Type: tokenType, Text: text, Value: value, Location: location}
}

func (t Token) String() string {
	return fmt.Sprintf("(%s, %s)", t.Type.String(), string(t.Text))
}

func (t TokenType) String() string {
	return fmt.Sprintf("%d", t)
}
