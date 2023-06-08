package lexer

import (
	"fmt"
	"strconv"
	"strings"
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

type InvalidCharacterLiteralError struct {
	location location.Location
}

func (m InvalidCharacterLiteralError) Error() string {
	return fmt.Sprintf("Invalid character literal at %v", m.location)
}

type InvalidNumberLiteralError struct {
	location location.Location
}

func (m InvalidNumberLiteralError) Error() string {
	return fmt.Sprintf("Invalid number literal at %v", m.location)
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

	next := s.peek()
	if unicode.IsLetter(next) || next == '_' {
		return s.scanIdentifier()
	}
	if unicode.IsDigit(next) || (next == '#' && s.peekN(1) != '\\') {
		return s.scanNumber()
	}

	next = s.advance()
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
		if s.match('>') {
			return s.makeToken(TOKEN_ARROW), nil
		}
		return s.makeToken(TOKEN_MINUS), nil
	case '*':
		return s.makeToken(TOKEN_STAR), nil
	case '/':
		return s.makeToken(TOKEN_SLASH), nil
	case '?':
		return s.makeToken(TOKEN_QUESTION_MARK), nil
	case '#':
		if s.match('\\') {
			return s.scanChar()
		}

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
		} else if s.match('>') {
			return s.makeToken(TOKEN_PIPE_GT), nil
		} else {
			return s.makeToken(TOKEN_PIPE), nil
		}
	// literals
	case '"':
		return s.scanString()
	}
	return s.unknownTokenError()
}

// TODO: add support for escaped quaracters
func (s *Scanner) scanString() (Token, error) {
	for {
		if s.isEof() {
			return s.unterminatedLiteralError()
		} else if s.match('\n') {
			s.line++
		} else if s.match('"') {
			value := strings.Trim(string(s.tokenText()), "\"")
			return s.makeTokenWithValue(TOKEN_STRING, value), nil
		} else {
			s.advance()
		}
	}
}

var (
	namedChars = map[string]rune{
		"newline":   '\n',
		"space":     ' ',
		"tab":       '\t',
		"backspace": '\b',
		"null":      '\000',
		"return":    '\r',
		"escape":    '\033',
		"delete":    '\177',
	}
)

func (s *Scanner) scanChar() (Token, error) {
	if s.isEof() {
		return s.unterminatedLiteralError()
	}

	for name, char := range namedChars {
		if s.matchString(name) {
			return s.makeTokenWithValue(TOKEN_CHARACTER, char), nil
		}
	}

	next := s.advance()
	if next == 'u' && unicode.IsDigit(s.peek()) {
		return s.scanCharUnicodeEscape()
	} else if next == 'x' && isHexCharacter(s.peek()) {
		return s.scanCharHexEscape()
	} else if unicode.IsPrint(next) {
		return s.makeTokenWithValue(TOKEN_CHARACTER, next), nil
	} else {
		return s.invalidCharacterLiteral()
	}
}

func isHexCharacter(c rune) bool {
	return unicode.IsDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (s *Scanner) scanCharUnicodeEscape() (Token, error) {
	// exactly 4 digits (leading zeros are expected)
	for i := 0; i < 4; i++ {
		if s.isEof() {
			return s.invalidCharacterLiteral()
		} else if !unicode.IsDigit(s.peek()) {
			return s.invalidCharacterLiteral()
		}
		s.advance()
	}

	text := strings.Trim(string(s.tokenText()[3:]), "0") // remove leading zeros
	value, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return s.invalidCharacterLiteral()
	}
	return s.makeTokenWithValue(TOKEN_CHARACTER, rune(value)), nil
}

func (s *Scanner) scanCharHexEscape() (Token, error) {
	// exactly 3 bytes hex encoded, so six chars
	for i := 0; i < 6; i++ {
		if s.isEof() {
			return s.invalidCharacterLiteral()
		} else if !isHexCharacter(s.peek()) {
			return s.invalidCharacterLiteral()
		}
		s.advance()
	}

	text := string(s.tokenText()[3:])
	value, err := strconv.ParseInt(text, 16, 32)
	if err != nil {
		return s.invalidCharacterLiteral()
	}
	return s.makeTokenWithValue(TOKEN_CHARACTER, rune(value)), nil
}

func (s *Scanner) scanNumber() (Token, error) {
	for !s.isEof() && unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// dot followed by digit?
	if s.peek() == '.' && unicode.IsDigit(s.peekN(1)) {
		s.advance()

		// all digits after the dot
		for !s.isEof() && unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	text := string(s.tokenText())
	if strings.Contains(text, ".") {
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return s.invalidNumberLiteral()
		}
		return s.makeTokenWithValue(TOKEN_NUMBER, value), nil
	} else {
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return s.invalidNumberLiteral()
		}
		return s.makeTokenWithValue(TOKEN_NUMBER, value), nil
	}

}

var keywords = map[string]TokenType{
	"fun":     TOKEN_FUN,
	"proc":    TOKEN_PROC,
	"package": TOKEN_PACKAGE,
	"import":  TOKEN_IMPORT,
	"export":  TOKEN_EXPORT,
	"data":    TOKEN_DATA,
	"alias":   TOKEN_ALIAS,
	"var":     TOKEN_VAR,
	"if":      TOKEN_IF,
	"else":    TOKEN_ELSE,
	"let":     TOKEN_LET,
	"set":     TOKEN_SET,
	"match":   TOKEN_MATCH,
	"for":     TOKEN_FOR,
	"return":  TOKEN_RETURN,
	"defer":   TOKEN_DEFER,
	"...":     TOKEN_ELLIPSIS,
	"true":    TOKEN_TRUE,
	"false":   TOKEN_FALSE,
}

func (s *Scanner) scanIdentifier() (Token, error) {
	// keywords
	for kw, token := range keywords {
		if s.matchString(kw) {
			if token == TOKEN_TRUE {
				return s.makeTokenWithValue(TOKEN_TRUE, true), nil
			} else if token == TOKEN_FALSE {
				return s.makeTokenWithValue(TOKEN_FALSE, false), nil
			} else {
				return s.makeToken(token), nil
			}
		}
	}

	// identifiers
	for {
		if s.isEof() {
			break
		} else {
			next := s.peek()
			if unicode.IsLetter(next) || unicode.IsDigit(next) || next == '_' || next == '\'' {
				s.advance()
			} else {
				break
			}
		}
	}

	return s.makeToken(TOKEN_IDENTIFIER), nil
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

func (s *Scanner) matchString(expected string) bool {
	start := s.cursor
	end := start + uint64(len(expected))
	if s.isEof() {
		return false
	} else if end > uint64(len(*s.buffer)) {
		return false
	} else if string((*s.buffer)[start:end]) != expected {
		return false
	}
	s.cursor = end
	return true
}

func (s *Scanner) skipWhitespace() {
	for !s.isEof() {
		switch s.peek() {
		case ' ', '\r', '\t':
			s.advance()
		case '\n':
			s.advance()
			s.line++
		case '/':
			if s.peekN(1) == '/' {
				for !s.isEof() && s.peek() != '\n' {
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

func (s Scanner) peek() rune {
	return s.peekN(0)
}

func (s *Scanner) peekN(offset uint64) rune {
	nextCursor := s.cursor + offset

	if s.isEof() || nextCursor >= uint64(len(*s.buffer)) {
		return 0
	}
	return (*s.buffer)[nextCursor]
}

func (s *Scanner) location() location.Location {
	return location.Location{Origin: &s.origin, Line: s.line, StartOffset: s.start, EndOffset: s.cursor}
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	return MakeToken(tokenType, s.tokenText(), s.location())
}

func (s *Scanner) makeTokenWithValue(tokenType TokenType, value interface{}) Token {
	return MakeTokenWithValue(tokenType, s.tokenText(), value, s.location())
}

func (s *Scanner) tokenText() []rune {
	return (*s.buffer)[s.start:s.cursor]
}

func (s *Scanner) unknownTokenError() (Token, error) {
	return Token{}, UnknownTokenError{location: s.location()}
}

func (s *Scanner) unterminatedLiteralError() (Token, error) {
	return Token{}, UnterminatedLiteralError{location: s.location()}
}

func (s *Scanner) invalidCharacterLiteral() (Token, error) {
	return Token{}, InvalidCharacterLiteralError{location: s.location()}
}

func (s *Scanner) invalidNumberLiteral() (Token, error) {
	return Token{}, InvalidNumberLiteralError{location: s.location()}
}
