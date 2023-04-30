package reader

import (
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

	for {
		datum := p.parseDatum()
		if datum == nil && !p.scanner.IsEof() {
			p.error("expected datum")
			return nil
		}

		if datum != nil {
			data = append(data, datum)
		}

		if p.scanner.IsEof() {
			return data
		}
	}
}

func (p *Parser) parseDatum() Datum {
	return p.parseBoolean()
}

func (p *Parser) parseBoolean() Datum {
	matched := false
	var value bool

	prevPos := p.scanner.Position().Offset

	if p.scanner.Attempt("#true") || p.scanner.Attempt("#t") {
		matched = true
		value = true
	}

	if p.scanner.Attempt("#false") || p.scanner.Attempt("#f") {
		matched = true
		value = false
	}

	if matched {
		pos := p.scanner.Position()
		return NewDatumBool(value, p.makeLocation(pos.Line, prevPos, pos.Offset))
	}

	return nil
}

func (p *Parser) makeLocation(line, start, end uint64) location.Location {
	return location.Location{
		Origin:      &p.input.Origin,
		Line:        line,
		StartOffset: start,
		EndOffset:   end,
	}
}
