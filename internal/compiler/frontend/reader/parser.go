package reader

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
)

type Parser struct {
	introspectionAPI introspection.API
	scanner          *Scanner
	errors           []ReadError
	input            *input.Input
}

func NewParser(introspectionAPI introspection.API) *Parser {
	return &Parser{introspectionAPI: introspectionAPI}
}

func (p *Parser) Parse(input *input.Input) (*DatumAST, []ReadError) {
	p.scanner = NewScanner(input.Buffer)
	p.errors = []ReadError{}
	p.input = input

	data := p.parseAll()
	if len(p.errors) > 0 {
		return nil, p.errors
	} else {
		return &DatumAST{Data: data}, nil
	}
}

func (p *Parser) error(msg string) {
	pos := p.scanner.Position()
	p.errors = append(p.errors, ReadError{Msg: msg, pos: pos})
}

func (p *Parser) recover() {
	// simple recovery strategy for now
	for {
		ch, err := p.scanner.Peek()
		if err != nil || ch == '(' || ch == '\n' {
			break
		}
		p.scanner.Next()
	}
}

func (p *Parser) parseAll() []Datum {
	data := []Datum{}
	attempts := 0

	for {
		attempts++
		if attempts > 10 {
			return nil
		}

		fmt.Printf("Parse all: %d\n", p.scanner.pos)

		datum := p.parseDatum()
		if datum == nil {
			p.error("Failed to parse datum")
			//p.recover()
		} else {
			data = append(data, datum)
		}

		println("scanner.IsEof():", p.scanner.IsEof())
		println("scanner pos", p.scanner.pos)

		if p.scanner.IsEof() {
			break
		}
	}
	return data
}

func (p *Parser) parseDatum() Datum {
	return p.parseBoolean()
}

func (p *Parser) parseBoolean() Datum {
	matched, err := p.scanner.Attempt("#t")

	if err != nil {
		return nil
	}

	if matched {
		pos := p.scanner.Position()
		return NewDatumBool(true, p.makeLocation(pos.Line, pos.Offset-2, pos.Offset))
	} else {
		return nil
	}
}

func (p *Parser) makeLocation(line, start, end uint64) location.Location {
	return location.Location{
		Origin:      &p.input.Origin,
		Line:        line,
		StartOffset: start,
		EndOffset:   end,
	}
}
