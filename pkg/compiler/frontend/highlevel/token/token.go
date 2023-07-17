// Package token provides functionality to represent token and locations for the braces language.
// A token represents a parsed piece of the source code, while a location represents a specific point or range in the source code.
// The package includes the ability to differentiate between various types of tokens including literals, operators, and keywords.
// It provides tools for creating tokens, checking their type, and handling illegal or unrecognized tokens.
package token

import (
	"fmt"
	"strings"
)

type Type uint8

// A token is essentially a tuple of the token type, the text that.
type Token struct {
	Type Type

	// The text that was used to create the token.
	Text []rune
	// Some tokens may have literal values associated with them, which have been created from the token text during lexing.
	// This can only be the case for literals, like strings, numbers, etc.
	LitValue interface{}

	// The location of the token in the source code.
	Location Location
}

const (
	ILLEGAL Type = iota
	EOF

	literal_begin // stolen from the go parser to implement token type checks efficiently
	IDENTIFIER
	FIXNUM
	FLONUM
	BYTE
	BOOLEAN
	UNIT
	CHAR
	STRING
	literal_end

	operator_begin
	// arithmetic
	ADD // +
	SUB // -
	MUL // *
	POW // **
	DIV // /
	REM // %

	// boolean
	LAND // &&
	LOR  // ||

	// bitwise
	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	EQ  // ==
	LT  // <
	LTE // <=
	GT  // >
	GTE // >=
	NOT // !
	NEQ // !=

	ELLIPSES // ...

	LPAREN    // (
	RPAREN    // )
	LBRACK    // [
	RBRACK    // ]
	LBRACE    // {
	RBRACE    // }
	COMMA     // ,
	DOT       // .
	COLON     // :
	DBLCOLON  // ::
	SEMICOLON // ;
	ARROW     // ->
	PIPE      // |>

	operator_end

	keyword_begin

	PACKAGE
	IMPORT
	API

	DATA
	ALIAS

	FUN
	PROC

	IF
	COND
	FOR
	BREAK
	RETURN
	ENSURE

	LET
	SET

	// we might be able to turn that later into constants in the universe (outermost scope) of the language
	TRUE
	FALSE

	keyword_end
)

var tokenStrings = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	IDENTIFIER: "IDENTIFIER",
	FIXNUM:     "FIXNUM",
	FLONUM:     "FLONUM",
	BYTE:       "BYTE",
	BOOLEAN:    "BOOLEAN",
	UNIT:       "UNIT",
	CHAR:       "CHAR",
	STRING:     "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	REM: "%",
	POW: "**",

	LAND: "&&",
	LOR:  "||",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	EQ:  "==",
	LT:  "<",
	LTE: "<=",
	GT:  ">",
	GTE: ">=",
	NOT: "!",
	NEQ: "!=",

	ELLIPSES: "...",

	LPAREN:    "(",
	RPAREN:    ")",
	LBRACK:    "[",
	RBRACK:    "]",
	LBRACE:    "{",
	RBRACE:    "}",
	COMMA:     ",",
	DOT:       ".",
	COLON:     ":",
	DBLCOLON:  "::",
	SEMICOLON: ";",
	ARROW:     "->",
	PIPE:      "|>",

	PACKAGE: "package",
	IMPORT:  "import",
	API:     "api",

	DATA:  "data",
	ALIAS: "alias",

	FUN:  "fun",
	PROC: "proc",

	IF:     "if",
	COND:   "cond",
	FOR:    "for",
	BREAK:  "break",
	RETURN: "return",
	ENSURE: "ensure",

	LET: "let",
	SET: "set",

	TRUE:  "true",
	FALSE: "false",
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type, keyword_end-(keyword_begin+1))
	for i := keyword_begin + 1; i < keyword_end; i++ {
		keywords[tokenStrings[i]] = i
	}
}

// Create a token, optionally providing a value
// If more than one value is given, only the first one is taken and the rest is ignored.
func New(location Location, tokenType Type, text []rune, value ...interface{}) Token {
	if len(value) == 1 {
		return Token{Type: tokenType, Text: text, LitValue: value[0], Location: location}
	} else {
		return Token{Type: tokenType, Text: text, Location: location}
	}
}

// Create an illegal token, which can be used during scanning to collect more than one error.
func Illegal(location Location, text []rune) Token {
	return New(location, ILLEGAL, text)
}

// Create a keyword token from the given string.
// It will return an ILLEGAL token if the string is not a known keyword.
func ByKeyword(location Location, keyword string) Token {
	for i := keyword_begin + 1; i < keyword_end; i++ {
		if tokenStrings[i] == keyword {
			return New(location, i, []rune(keyword))
		}
	}

	return Illegal(location, []rune(keyword))
}

func (t Token) IsLiteral() bool {
	return literal_begin < t.Type && t.Type < literal_end
}

func (t Token) IsOperator() bool {
	return operator_begin < t.Type && t.Type < operator_end
}

func (t Token) IsKeyword() bool {
	return keyword_begin < t.Type && t.Type < keyword_end
}

func (t Token) IsIdentifier() bool {
	return t.Type == IDENTIFIER
}

func (t Token) IsEOF() bool {
	return t.Type == EOF
}

func (t Token) IsIllegal() bool {
	return t.Type == ILLEGAL
}

func (t Type) String() string {
	return strings.ToUpper(tokenStrings[t])
}

func (t Token) String() string {
	return fmt.Sprintf("(%s, %s)", t.Type.String(), string(t.Text))
}

func (t Token) Sexp() string {
	return fmt.Sprintf("(tok %s \"%s\" %s)", t.Type.Sexp(), string(t.Text), t.Location.Sexp())
}

var tokenSexp = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	IDENTIFIER: "IDENTIFIER",
	FIXNUM:     "FIXNUM",
	FLONUM:     "FLONUM",
	BYTE:       "BYTE",
	BOOLEAN:    "BOOLEAN",
	UNIT:       "UNIT",
	CHAR:       "CHAR",
	STRING:     "STRING",

	ADD: "ADD",
	SUB: "SUB",
	MUL: "MUL",
	DIV: "DIV",
	REM: "REM",
	POW: "POW",

	LAND: "LAND",
	LOR:  "LOR",

	AND:     "AND",
	OR:      "OR",
	XOR:     "XOR",
	SHL:     "SHL",
	SHR:     "SHR",
	AND_NOT: "AND_NOT",

	EQ:  "EQ",
	LT:  "LT",
	LTE: "LTE",
	GT:  "GT",
	GTE: "GTE",
	NOT: "NOT",
	NEQ: "NEQ",

	ELLIPSES: "ELLIPSES",

	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACK:    "LBRACK",
	RBRACK:    "RBRACK",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	COMMA:     "COMMMA",
	DOT:       "DOT",
	COLON:     "COLON",
	DBLCOLON:  "DBLCOLON",
	SEMICOLON: "SEMICOLON",
	ARROW:     "ARROW",
	PIPE:      "PIPE",

	PACKAGE: "PACKAGE",
	IMPORT:  "IMPORT",
	API:     "API",

	DATA:  "DATA",
	ALIAS: "ALIAS",

	FUN:  "FUN",
	PROC: "PROC",

	IF:     "IF",
	COND:   "COND",
	FOR:    "FOR",
	BREAK:  "BREAK",
	RETURN: "RETURN",
	ENSURE: "ENSURE",

	LET: "LET",
	SET: "SET",

	TRUE:  "TRUE",
	FALSE: "TRUE",
}

func (t Type) Sexp() string {
	return tokenSexp[t]
}
