package lexer

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"github.com/certainty/go-braces/internal/compiler/input"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type (
	Scanner struct {
		// the source code producer
		*input.Input

		// temporary start of a given scan attempt
		start uint64
		// the position until we have scanned
		cursor uint64
		// line of start
		line token.Line
		// column of start
		column token.Column

		// encountered errors
		Errors []error
	}
)

func New(input *input.Input) *Scanner {
	scanner := &Scanner{}
	scanner.Reset(input)

	return scanner
}

func (s *Scanner) Reset(input *input.Input) {
	s.Input = input
	s.start = 0
	s.cursor = 0
	s.line = 1
	s.column = 1
	s.Errors = []error{}
}

// Returns the next token from the input.
// If an error is encoutered a token of type token.ILLEGAL is returned..
// In this case the caller can inspect the scanner's Errors slice for more information.
func (s *Scanner) NextToken() token.Token {
	s.skipWhitespace()

	s.start = s.cursor
	if s.isEof() {
		return s.makeToken(token.EOF)
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
	// token recognizable with one byte lookahead
	case '{':
		return s.makeToken(token.LBRACE)
	case '}':
		return s.makeToken(token.RBRACE)
	case '[':
		return s.makeToken(token.LBRACK)
	case ']':
		return s.makeToken(token.RBRACK)
	case '(':
		return s.makeToken(token.LPAREN)
	case ')':
		return s.makeToken(token.RBRACE)
	case '+':
		return s.makeToken(token.ADD)
	case '/':
		return s.makeToken(token.DIV)
	case '%':
		return s.makeToken(token.REM)
	case '^':
		return s.makeToken(token.XOR)
	case ',':
		return s.makeToken(token.COMMA)

	case '#':
		if s.match('\\') {
			return s.scanChar()
		}
	case ':':
		if s.match(':') {
			return s.makeToken(token.DBLCOLON)
		} else {
			return s.makeToken(token.COLON)
		}
	case '=':
		if s.match('=') {
			return s.makeToken(token.EQ)
		}
	case '!':
		if s.match('=') {
			return s.makeToken(token.NEQ)
		} else {
			return s.makeToken(token.NOT)
		}
	case '-':
		if s.match('>') {
			return s.makeToken(token.ARROW)
		}
		return s.makeToken(token.SUB)
	case '*':
		if s.match('*') {
			return s.makeToken(token.POW)
		}
		return s.makeToken(token.MUL)

	case '>':
		if s.match('=') {
			return s.makeToken(token.GTE)
		} else if s.match('>') {
			return s.makeToken(token.SHR)
		} else {
			return s.makeToken(token.GT)
		}
	case '<':
		if s.match('=') {
			return s.makeToken(token.LTE)
		} else if s.match('<') {
			return s.makeToken(token.SHL)
		} else {
			return s.makeToken(token.LT)
		}
	case '&':
		if s.match('&') {
			return s.makeToken(token.LAND)
		} else if s.match('^') {
			return s.makeToken(token.AND_NOT)
		} else {
			return s.makeToken(token.AND)
		}
	case '|':
		if s.match('|') {
			return s.makeToken(token.LOR)
		} else if s.match('>') {
			return s.makeToken(token.PIPE)
		} else {
			return s.makeToken(token.OR)
		}

	// literals
	case '"':
		return s.scanString()
	}
	return s.illegalToken("Unexpected character")
}

func (s *Scanner) illegalToken(message string) token.Token {
	s.Errors = append(s.Errors, fmt.Errorf("%s '%c' at %v", message, s.peek(), s.location()))
	return s.makeToken(token.ILLEGAL)
}

////////////////////////////////////////////////////////////////////
// Strings
////////////////////////////////////////////////////////////////////

func (s *Scanner) scanString() token.Token {
	value := ""
	for {
		if s.isEof() {
			return s.illegalToken("Unterminated string literal")
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
				return s.illegalToken("invalid escape sequence")
			}
			s.advance()
		} else if s.match('"') {
			return s.makeToken(token.STRING, value)
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

func (s *Scanner) scanChar() token.Token {
	if s.isEof() {
		return s.illegalToken("unterminated character literal")
	}

	// identify named characters
	for name, char := range namedChars {
		if s.matchString(name) {
			return s.makeToken(token.CHAR, char)
		}
	}

	next := s.advance()

	// scan a unicode escape sequence
	// which is \u followed by exactly 4 digits (leading zeros are expected)
	if next == 'u' && unicode.IsDigit(s.peek()) {
		return s.scanCharUnicodeEscape()

		// scan a unicode escape sequence
		// which is \x followed hex digits
	} else if next == 'x' && isHexDigit(s.peek()) {
		return s.scanCharHexEscape()
	} else if unicode.IsPrint(next) {
		return s.makeToken(token.CHAR, next)
	} else {
		return s.illegalToken("invalid character literal")
	}
}

func (s *Scanner) scanCharUnicodeEscape() token.Token {
	for i := 0; i < 4; i++ {
		if s.isEof() {
			return s.illegalToken("unterminated character literal")
		} else if !unicode.IsDigit(s.peek()) {
			return s.illegalToken("invalid unicode escape sequence")
		}
		s.advance()
	}

	text := strings.Trim(string(s.tokenText()[3:]), "0") // remove leading zeros
	value, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return s.illegalToken("invalid unicode escape sequence")
	}
	return s.makeToken(token.CHAR, rune(value))
}

// scan exactly 3 bytes hex encoded, so six chars
func (s *Scanner) scanCharHexEscape() token.Token {
	for i := 0; i < 6; i++ {
		if s.isEof() {
			return s.illegalToken("unexpected EOF")
		} else if !isHexDigit(s.peek()) {
			return s.illegalToken("invalid hex escape sequence")
		}
		s.advance()
	}

	text := string(s.tokenText()[3:])
	value, err := strconv.ParseInt(text, 16, 32)
	if err != nil {
		return s.illegalToken("invalid hex escape sequence")
	}
	return s.makeToken(token.CHAR, rune(value))
}

// //////////////////////////////////////////////////////////////////
// Numbers
// //////////////////////////////////////////////////////////////////
func (s *Scanner) scanNumber() token.Token {
	if s.match('#') {
		return s.scanIntWithBase()
	} else if s.matchString("+nan.0") || s.matchString("-nan.0") {
		return s.makeToken(token.FLONUM, math.NaN())
	} else if s.matchString("+inf.0") {
		return s.makeToken(token.FLONUM, math.Inf(1))
	} else if s.matchString("-inf.0") {
		return s.makeToken(token.FLONUM, math.Inf(-1))
	} else {
		return s.scanFloatOrInt()
	}
}

func (s *Scanner) scanIntWithBase() token.Token {
	var base uint8
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
		return s.illegalToken("invalid base signifier for number literal")
	}

	s.advance()
	s.scanDigits(base)
	text := string(s.tokenText())[2:]
	// we need to know the widt of ints on this platform
	value, err := strconv.ParseInt(text, int(base), 64)
	if err != nil {
		return s.illegalToken("invalid fixnum literal")
	}
	return s.makeToken(token.FIXNUM, int(value))
}

func (s *Scanner) scanFloatOrInt() token.Token {
	if sign := s.peek(); sign == '+' || sign == '-' {
		s.advance()
	}
	s.scanDigits(10)

	if s.peek() == '.' && unicode.IsDigit(s.peekN(1)) {
		s.advance()
		s.scanDigits(10)

		text := string(s.tokenText())
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return s.illegalToken("invalid flonum literal")
		}
		return s.makeToken(token.FLONUM, value)
	} else {
		text := string(s.tokenText())
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return s.illegalToken("invalid fixnum literal")
		}
		return s.makeToken(token.FIXNUM, int(value))
	}
}

func (s *Scanner) scanDigits(base uint8) {
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

func isDigit(c rune, base uint8) bool {
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
func (s *Scanner) scanIdentifier() token.Token {
	for {
		if s.isEof() {
			break
		}

		next := s.peek()
		if unicode.IsLetter(next) || unicode.IsDigit(next) || next == '_' || next == '\'' {
			s.advance()
		} else {
			break
		}
	}

	keyword := token.ByKeyword(s.location(), string(s.tokenText()))

	if keyword.IsKeyword() {
		if keyword.Type == token.TRUE || keyword.Type == token.FALSE {
			keyword.LitValue = keyword.Type == token.TRUE
		}

		return keyword
	}

	return s.makeToken(token.IDENTIFIER)
}

// //////////////////////////////////////////////////////////////////
// Helpers
// //////////////////////////////////////////////////////////////////
func (s *Scanner) isEof() bool {
	return s.cursor >= uint64(len(*s.Buffer))
}

func (s *Scanner) advance() rune {
	s.column++
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
			s.column = 1
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

func (s *Scanner) makeToken(tokenType token.Type, value ...interface{}) token.Token {
	return token.New(s.location(), tokenType, s.tokenText(), value)
}

func (s *Scanner) location() token.Location {
	return token.NewLocation(
		s.Origin,
		token.Line(s.line),
		token.Column(s.column),
		token.From(s.start),
		token.To(s.cursor),
	)
}

func (s *Scanner) tokenText() []rune {
	return (*s.Buffer)[s.start:s.cursor]
}
