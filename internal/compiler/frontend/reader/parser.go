package reader

import (
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	scanner         *Scanner
	errors          []ReadError
	input           *input.Input
}

func NewParser(instrumentation compiler_introspection.Instrumentation) *Parser {
	return &Parser{instrumentation: instrumentation}
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

func (p *Parser) recover() error {
	// simple recovery strategy for now
	for {
		ch, err := p.scanner.Peek()
		if err != nil || ch == '(' || ch == '\n' {
			break
		}
		if err := p.scanner.Skip(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseAll() []isa.Datum {
	data := []isa.Datum{}

	for {
		if err := p.scanner.SkipIrrelevant(); err != nil {
			p.error(err.Error())
			return nil
		}

		datum := p.parseDatum()
		if datum == nil && !p.scanner.IsEof() {
			p.error("expected datum")
			if err := p.recover(); err != nil {
				return nil
			}
		}

		if datum != nil {
			p.instrumentation.Breakpoint(compiler_introspection.BPReaderAccepted, p)
			data = append(data, datum)
		}

		if p.scanner.IsEof() {
			return data
		}
	}
}

func (p *Parser) parseDatum() isa.Datum {
	p.instrumentation.Breakpoint(compiler_introspection.BPReaderParseDatum, p)

	// the contract for the following parsers is that they only consume input if they match.
	//If not they have to take care of restoring the scanner's state.

	datum := p.parseBoolean()
	if datum != nil {
		return datum
	}

	datum = p.parseChar()
	if datum != nil {
		return datum
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

func (p *Parser) locationFromPositions(start Position, end Position) location.Location {
	return p.makeLocation(start.Line, start.Offset, end.Offset)
}
