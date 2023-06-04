package reader

import (
	"github.com/certainty/go-braces/internal/isa"
	"strconv"
)

func (p *Parser) parseChar() isa.Datum {
	startPos := p.scanner.Position()
	nextToken, err := p.scanner.PeekN(2)
	if err != nil {
		return nil
	}

	if nextToken != "#\\" {
		return nil
	}
	p.scanner.Skip()
	p.scanner.Skip()

	if char := p.parseHexLiteral(); char != -1 {
		return isa.NewDatumChar(char, p.locationFromPositions(startPos, p.scanner.Position()))
	} else if char = p.parseNamedCharLiteral(); char != -1 {
		return isa.NewDatumChar(char, p.locationFromPositions(startPos, p.scanner.Position()))
	} else if char = p.parseCharLiteral(); char != -1 {
		return isa.NewDatumChar(char, p.locationFromPositions(startPos, p.scanner.Position()))
	}

	return nil
}

// parse a sequence of up to 3 bytes hex encoded
func (p *Parser) parseHexLiteral() rune {
	p.scanner.SavePosition()
	defer p.scanner.ReleaseSavePoint()

	if !p.scanner.Attempt("x") {
		return -1
	}

	digit := p.parseHexDigit()
	if digit == -1 {
		p.scanner.RestorePosition()
		return -1
	}
	return rune(digit)
}

func (p *Parser) parseHexDigit() int64 {
	hexString := ""
	for i := 0; i < 6; i++ {
		k, err := p.scanner.Peek()
		if err != nil || !isHexDigit(k) {
			break
		}
		hexString += string(k)
		p.scanner.Skip()
	}

	if hexString == "" {
		return -1
	}

	n, _ := strconv.ParseInt(hexString, 16, 64)
	return n
}

func (p *Parser) parseNamedCharLiteral() rune {
	if p.scanner.Attempt("space") {
		return ' '
	}
	if p.scanner.Attempt("newline") {
		return '\n'
	}
	if p.scanner.Attempt("return") {
		return '\r'
	}
	if p.scanner.Attempt("tab") {
		return '\t'
	}
	if p.scanner.Attempt("alarm") {
		return 7
	}
	if p.scanner.Attempt("backspace") {
		return 8
	}
	if p.scanner.Attempt("delete") {
		return 127
	}
	if p.scanner.Attempt("null") {
		return 0
	}
	if p.scanner.Attempt("escape") {
		return 27
	}
	return -1
}

func (p *Parser) parseCharLiteral() rune {
	c, err := p.scanner.Next()
	if err != nil {
		return -1
	}
	return c
}

func isHexDigit(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
