package lexer

import (
	"fmt"
	"unicode"

	"github.com/certainty/go-braces/internal/compiler/location"
)

type UnknownTokenError struct {
	location location.Location
}

func (m UnknownTokenError) Error() string {
	return fmt.Sprintf("Unknown token at %v", m.location)
}

type UnterminatedLiteralError struct {
	location location.Location
}

func (m UnterminatedLiteralError) Error() string {
	return fmt.Sprintf("Unterminated literal at %v", m.location)
}

type Scanner struct {
	origin location.Origin
	buffer *[]rune
	start  uint64
	cursor uint64
	line   uint64
}

func New(buffer *[]rune, origin location.Origin) *Scanner {
	return &Scanner{
		origin: origin,
		buffer: buffer,
		start:  0,
		cursor: 0,
		line:   1,
	}
}

func NewFromString(input string, origin location.Origin) *Scanner {
	runes := []rune(input)
	return New(&runes, origin)
}

func (s *Scanner) NextToken() (Token, error) {
	s.skipWhitespace()

	s.start = s.cursor
	if s.isEof() {
		return s.makeToken(TOKEN_EOF), nil
	}

	next := s.advance()

	if unicode.IsDigit(next) {
		return s.scanNumber()
	}

	switch next {
	case '{':
		return s.makeToken(TOKEN_LBRACE), nil
	case '}':
		return s.makeToken(TOKEN_RBRACE), nil
	case '[':
		return s.makeToken(TOKEN_LBRACKET), nil
	case ']':
		return s.makeToken(TOKEN_RBRACKET), nil
	case '(':
		return s.makeToken(TOKEN_LPAREN), nil
	case ')':
		return s.makeToken(TOKEN_RPAREN), nil
	case ',':
		return s.makeToken(TOKEN_COMMA), nil
	case '+':
		return s.makeToken(TOKEN_PLUS), nil
	case '-':
		return s.makeToken(TOKEN_MINUS), nil
	case '*':
		return s.makeToken(TOKEN_STAR), nil
	case '/':
		return s.makeToken(TOKEN_SLASH), nil
	case '?':
		return s.makeToken(TOKEN_QUESTION_MARK), nil

		// multi char tokens
	case ':':
		if s.match(':') {
			return s.makeToken(TOKEN_COLON_COLON), nil
		} else {
			return s.makeToken(TOKEN_COLON), nil
		}
	case '=':
		if s.match('=') {
			return s.makeToken(TOKEN_EQUAL_EQUAL), nil
		} else {
			return s.makeToken(TOKEN_EQUAL), nil
		}
	case '!':
		if s.match('=') {
			return s.makeToken(TOKEN_BANG_EQUAL), nil
		} else {
			return s.makeToken(TOKEN_BANG), nil
		}
	case '>':
		if s.match('=') {
			return s.makeToken(TOKEN_GT_EQUAL), nil
		} else {
			return s.makeToken(TOKEN_GT), nil
		}
	case '<':
		if s.match('=') {
			return s.makeToken(TOKEN_LT_EQUAL), nil
		} else {
			return s.makeToken(TOKEN_LT), nil
		}
	case '&':
		if s.match('&') {
			return s.makeToken(TOKEN_AMPERSAND_AMPERSAND), nil
		} else {
			return s.makeToken(TOKEN_AMPERSAND), nil
		}
	case '|':
		if s.match('|') {
			return s.makeToken(TOKEN_PIPE_PIPE), nil
		} else {
			return s.makeToken(TOKEN_PIPE), nil
		}
	// literals
	case '"':
		return s.scanString()
	case '\'':
		return s.scanChar()
	}
	return s.unknownTokenError()
}

// TODO: add support for escaped quaracters
func (s *Scanner) scanString() (Token, error) {
	for !s.isEof() && s.peek(0) != '"' {
		if s.peek(0) == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isEof() {
		return s.unterminatedLiteralError()
	}
	// consume closing quote
	s.advance()
	return s.makeToken(TOKEN_STRING), nil
}

// TODO: add support for unicode escapes \uXXXX
func (s *Scanner) scanChar() (Token, error) {
	for !s.isEof() && s.peek(0) != '\'' {
		if s.peek(0) == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isEof() {
		return s.unterminatedLiteralError()
	}
	// consume closing single quote
	s.advance()
	return s.makeToken(TOKEN_CHARACTER), nil
}

// TODO: add support for scientific notation and literals for binary, octal and hex notation
// See also scheme's number literals
func (s *Scanner) scanNumber() (Token, error) {
	for !s.isEof() && unicode.IsDigit(s.peek(0)) {
		s.advance()
	}
	if s.peek(0) == '.' && unicode.IsDigit(s.peek(1)) {
		// consume the '.'
		s.advance()
		for !s.isEof() && unicode.IsDigit(s.peek(0)) {
			s.advance()
		}
	}
	return s.makeToken(TOKEN_NUMBER), nil
}

func (s *Scanner) isEof() bool {
	return s.cursor >= uint64(len(*s.buffer))
}

func (s *Scanner) advance() rune {
	s.cursor++
	return (*s.buffer)[s.cursor-1]
}

// one rune look ahead, returns the next character without advancing the cursor
func (s *Scanner) match(expected rune) bool {
	if s.isEof() {
		return false
	}
	if (*s.buffer)[s.cursor] != expected {
		return false
	}
	s.cursor++
	return true
}

func (s *Scanner) skipWhitespace() {
	for !s.isEof() {
		switch s.peek(0) {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.advance()
			s.line++
		case '/':
			if s.peek(1) == '/' {
				for !s.isEof() && s.peek(0) != '\n' {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *Scanner) peek(offset uint64) rune {
	if s.isEof() {
		return 0
	}
	return (*s.buffer)[s.cursor+offset]
}

func (s *Scanner) location() location.Location {
	return location.Location{Origin: &s.origin, Line: s.line, StartOffset: s.start, EndOffset: s.cursor}
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	loc := s.location()
	text := (*s.buffer)[s.start:s.cursor]
	return MakeToken(tokenType, text, loc)
}

func (s *Scanner) unknownTokenError() (Token, error) {
	return Token{}, UnknownTokenError{location: s.location()}
}

func (s *Scanner) unterminatedLiteralError() (Token, error) {
	return Token{}, UnterminatedLiteralError{location: s.location()}
}
