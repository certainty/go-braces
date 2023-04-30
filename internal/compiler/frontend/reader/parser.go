package reader

import (
	"fmt"
	"io"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
)

type ReadError struct {
	Msg string
	pos Position
}

func (e ReadError) Error() string {
	return fmt.Sprintf("%s at %d:%d", e.Msg, e.pos.Line, e.pos.Col)
}

type Parser struct {
	introspectionAPI introspection.API
	scanner          *Scanner
	errors           []ReadError
}

func NewParser(introspectionAPI introspection.API) *Parser {
	return &Parser{introspectionAPI: introspectionAPI}
}

func (p *Parser) Parse(input location.Input) (*DatumAST, []ReadError) {
	p.scanner = NewScanner(input.Reader())
	p.errors = []ReadError{}

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
		p.scanner.SkipIrrelevant()
		_, err := p.scanner.Peek()

		if err == io.EOF {
			break
		}

		if err != nil {
			p.error("Error reading input")
			break
		}

		datum := p.parseDatum()
		data = append(data, datum)
	}
	return data
}

func (p *Parser) parseDatum() Datum {
	return nil
}
