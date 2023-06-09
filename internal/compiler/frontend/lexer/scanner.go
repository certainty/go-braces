package lexer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/certainty/go-braces/internal/compiler/input"
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

type InvalidStringLiteralError struct {
	location location.Location
}

func (m InvalidStringLiteralError) Error() string {
	return fmt.Sprintf("Invalid string literal at %v", m.location)
}

type InvalidNumberLiteralError struct {
	location location.Location
}

func (m InvalidNumberLiteralError) Error() string {
	return fmt.Sprintf("Invalid number literal at %v", m.location)
}

type Scanner struct {
	*input.Input
	start  uint64
	cursor uint64
	line   uint64
}

func New(input *input.Input) *Scanner {
	return &Scanner{Input: input, start: 0, cursor: 0, line: 1}
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

	//numbers
	if unicode.IsDigit(next) {
		return s.scanNumber()
	}

	// different base numbers
	if next == '#' && s.peekN(1) != '\\' {
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
	case '^':
		return s.makeToken(TOKEN_CARET), nil
	case '%':
		return s.makeToken(TOKEN_MOD), nil
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

////////////////////////////////////////////////////////////////////
// Strings
////////////////////////////////////////////////////////////////////

func (s *Scanner) scanString() (Token, error) {
	value := ""
	for {
		if s.isEof() {
			return s.unterminatedLiteralError()
		} else if s.match('\n') {
			value += "\n"
			s.line++
		} else if s.match('\\') {
			switch s.peek() {
			case '"':
				value += "\""
			case 'n':
				value += "\n"
			case 't':
				value += "\t"
			case 'b':
				value += "\b"
			case '0':
				value += "\000"
			case 'r':
				value += "\r"
			case '\\':
				value += "\\"
			default:
				return s.invalidStringLiteral()
			}
			s.advance()
		} else if s.match('"') {
			return s.makeTokenWithValue(TOKEN_STRING, value), nil
		} else {
			value += string(s.advance())
		}
	}
}

// //////////////////////////////////////////////////////////////////
// Chars
// //////////////////////////////////////////////////////////////////

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
			return s.makeTokenWithValue(TOKEN_CHARACTER, CodePoint{char}), nil
		}
	}

	next := s.advance()
	if next == 'u' && unicode.IsDigit(s.peek()) {
		return s.scanCharUnicodeEscape()
	} else if next == 'x' && isHexDigit(s.peek()) {
		return s.scanCharHexEscape()
	} else if unicode.IsPrint(next) {
		return s.makeTokenWithValue(TOKEN_CHARACTER, CodePoint{next}), nil
	} else {
		return s.invalidCharacterLiteral()
	}
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
	return s.makeTokenWithValue(TOKEN_CHARACTER, CodePoint{rune(value)}), nil
}

func (s *Scanner) scanCharHexEscape() (Token, error) {
	// exactly 3 bytes hex encoded, so six chars
	for i := 0; i < 6; i++ {
		if s.isEof() {
			return s.invalidCharacterLiteral()
		} else if !isHexDigit(s.peek()) {
			return s.invalidCharacterLiteral()
		}
		s.advance()
	}

	text := string(s.tokenText()[3:])
	value, err := strconv.ParseInt(text, 16, 32)
	if err != nil {
		return s.invalidCharacterLiteral()
	}
	return s.makeTokenWithValue(TOKEN_CHARACTER, CodePoint{rune(value)}), nil
}

// //////////////////////////////////////////////////////////////////
// Numbers
// //////////////////////////////////////////////////////////////////
func (s *Scanner) scanNumber() (Token, error) {
	if s.match('#') {
		return s.scanIntWithBase()
	} else if s.matchString("+nan.0") || s.matchString("-nan.0") {
		return s.makeTokenWithValue(TOKEN_FLOAT, math.NaN()), nil
	} else if s.matchString("+inf.0") {
		return s.makeTokenWithValue(TOKEN_FLOAT, math.Inf(1)), nil
	} else if s.matchString("-inf.0") {
		return s.makeTokenWithValue(TOKEN_FLOAT, math.Inf(-1)), nil
	} else {
		return s.scanFloatOrInt()
	}
}

func (s *Scanner) scanIntWithBase() (Token, error) {
	base := uint(10)
	baseSign := s.peek()

	switch baseSign {
	case 'b':
		base = 2
	case 'o':
		base = 8
	case 'x':
		base = 16
	case 'd':
		base = 10
	default:
		return s.invalidNumberLiteral()
	}

	s.advance()
	s.scanDigits(base)
	text := string(s.tokenText())[2:]
	value, err := strconv.ParseInt(text, int(base), 64)
	if err != nil {
		return s.invalidNumberLiteral()
	}
	return s.makeTokenWithValue(TOKEN_INTEGER, value), nil
}

func (s *Scanner) scanFloatOrInt() (Token, error) {
	sign := s.peek()

	if sign == '+' || sign == '-' {
		s.advance()
	}

	s.scanDigits(10)

	if s.peek() == '.' && unicode.IsDigit(s.peekN(1)) {
		s.advance()
		s.scanDigits(10)

		text := string(s.tokenText())
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return s.invalidNumberLiteral()
		}
		return s.makeTokenWithValue(TOKEN_FLOAT, value), nil
	} else {
		text := string(s.tokenText())
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return s.invalidNumberLiteral()
		}
		return s.makeTokenWithValue(TOKEN_INTEGER, value), nil
	}
}

func (s *Scanner) scanDigits(base uint) {
	for !s.isEof() && isDigit(s.peek(), base) {
		s.advance()
	}
}

func isHexDigit(c rune) bool {
	return unicode.IsDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isBinaryDigit(c rune) bool {
	return c == '0' || c == '1'
}

func isOctalDigit(c rune) bool {
	return c >= '0' && c <= '7'
}

func isDigit(c rune, base uint) bool {
	switch base {
	case 2:
		return isBinaryDigit(c)
	case 8:
		return isOctalDigit(c)
	case 10:
		return unicode.IsDigit(c)
	case 16:
		return isHexDigit(c)
	default:
		return false
	}
}

// //////////////////////////////////////////////////////////////////
// Identifiers
// //////////////////////////////////////////////////////////////////
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

	for kw, token := range keywords {
		if string(s.tokenText()) == kw {
			if token == TOKEN_TRUE {
				return s.makeTokenWithValue(TOKEN_TRUE, true), nil
			} else if token == TOKEN_FALSE {
				return s.makeTokenWithValue(TOKEN_FALSE, false), nil
			} else {
				return s.makeToken(token), nil
			}
		}
	}

	return s.makeToken(TOKEN_IDENTIFIER), nil
}

// //////////////////////////////////////////////////////////////////
// Helpers
// //////////////////////////////////////////////////////////////////
func (s *Scanner) isEof() bool {
	return s.cursor >= uint64(len(*s.Buffer))
}

func (s *Scanner) advance() rune {
	s.cursor++
	return (*s.Buffer)[s.cursor-1]
}

// one rune look ahead, returns the next character without advancing the cursor
func (s *Scanner) match(expected rune) bool {
	if s.isEof() {
		return false
	}
	if (*s.Buffer)[s.cursor] != expected {
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
	} else if end > uint64(len(*s.Buffer)) {
		return false
	} else if string((*s.Buffer)[start:end]) != expected {
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

	if s.isEof() || nextCursor >= uint64(len(*s.Buffer)) {
		return 0
	}
	return (*s.Buffer)[nextCursor]
}

func (s *Scanner) location() location.Location {
	return location.Location{Origin: &s.Origin, Line: s.line, StartOffset: s.start, EndOffset: s.cursor}
}

func (s *Scanner) makeToken(tokenType TokenType) Token {
	return MakeToken(tokenType, s.tokenText(), s.location())
}

func (s *Scanner) makeTokenWithValue(tokenType TokenType, value interface{}) Token {
	return MakeTokenWithValue(tokenType, s.tokenText(), value, s.location())
}

func (s *Scanner) tokenText() []rune {
	return (*s.Buffer)[s.start:s.cursor]
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

func (s *Scanner) invalidStringLiteral() (Token, error) {
	return Token{}, InvalidStringLiteralError{location: s.location()}
}

func (s *Scanner) invalidNumberLiteral() (Token, error) {
	return Token{}, InvalidNumberLiteralError{location: s.location()}
}
