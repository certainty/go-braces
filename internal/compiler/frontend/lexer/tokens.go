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

type CodePoint struct {
	Char rune
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
	TOKEN_POWER
	TOKEN_COLON
	TOKEN_CARET
	TOKEN_MOD

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

	// keywords
	TOKEN_PACKAGE
	TOKEN_IMPORT
	TOKEN_API

	TOKEN_DATA
	TOKEN_ALIAS

	TOKEN_FUN
	TOKEN_PROC

	TOKEN_IF
	TOKEN_ELSE
	TOKEN_LET
	TOKEN_SET
	TOKEN_MATCH
	TOKEN_FOR
	TOKEN_BREAK
	TOKEN_FROM

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
	switch t {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_PLUS:
		return "PLUS"
	case TOKEN_MINUS:
		return "MINUS"
	case TOKEN_STAR:
		return "STAR"
	case TOKEN_POWER:
		return "POWER"
	case TOKEN_SLASH:
		return "SLASH"
	case TOKEN_COLON:
		return "COLON"
	case TOKEN_CARET:
		return "CARET"
	case TOKEN_MOD:
		return "MOD"
	case TOKEN_EQUAL:
		return "EQUAL"
	case TOKEN_EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case TOKEN_BANG:
		return "BANG"
	case TOKEN_BANG_EQUAL:
		return "BANG_EQUAL"
	case TOKEN_GT:
		return "GT"
	case TOKEN_GT_EQUAL:
		return "GT_EQUAL"
	case TOKEN_LT:
		return "LT"
	case TOKEN_LT_EQUAL:
		return "LT_EQUAL"
	case TOKEN_AMPERSAND:
		return "AMPERSAND"
	case TOKEN_AMPERSAND_AMPERSAND:
		return "AMPERSAND_AMPERSAND"
	case TOKEN_PIPE:
		return "PIPE"
	case TOKEN_PIPE_PIPE:
		return "PIPE_PIPE"
	case TOKEN_PIPE_GT:
		return "PIPE_GT"
	case TOKEN_COLON_COLON:
		return "COLON_COLON"
	case TOKEN_ARROW:
		return "ARROW"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_INTEGER:
		return "INTEGER"
	case TOKEN_FLOAT:
		return "FLOAT"
	case TOKEN_BOOLEAN:
		return "BOOLEAN"
	case TOKEN_CHARACTER:
		return "CHARACTER"
	case TOKEN_FUN:
		return "FUN"
	case TOKEN_PROC:
		return "PROC"
	case TOKEN_PACKAGE:
		return "PACKAGE"
	case TOKEN_IMPORT:
		return "IMPORT"
	case TOKEN_API:
		return "API"
	case TOKEN_DATA:
		return "DATA"
	case TOKEN_ALIAS:
		return "ALIAS"
	case TOKEN_IF:
		return "IF"
	case TOKEN_ELSE:
		return "ELSE"
	case TOKEN_LET:
		return "LET"
	case TOKEN_SET:
		return "SET"
	case TOKEN_MATCH:
		return "MATCH"
	case TOKEN_FOR:
		return "FOR"
	case TOKEN_FROM:
		return "FROM"
	case TOKEN_BREAK:
		return "BREAK"
	case TOKEN_RETURN:
		return "RETURN"
	case TOKEN_DEFER:
		return "DEFER"
	case TOKEN_ELLIPSIS:
		return "ELLIPSIS"
	case TOKEN_TRUE:
		return "TRUE"
	case TOKEN_FALSE:
		return "FALSE"
	default:
		return "UNKNOWN"
	}
}
