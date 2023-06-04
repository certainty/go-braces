package reader

import "github.com/certainty/go-braces/internal/isa"

func (p *Parser) parseBoolean() isa.Datum {
	p.scanner.SavePosition()

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
		p.scanner.ReleaseSavePoint()
		pos := p.scanner.Position()
		return isa.NewDatumBool(value, p.makeLocation(pos.Line, prevPos, pos.Offset))
	}

	p.scanner.RestorePosition()
	return nil
}
