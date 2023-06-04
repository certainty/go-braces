package parser

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type CoreParser struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewCoreParser(instrumentation compiler_introspection.Instrumentation) *CoreParser {
	return &CoreParser{instrumentation: instrumentation}
}

func (p *CoreParser) Parse(data isa.Datum) (SchemeExpression, error) {
	return p.parseLiteral(data)
}

func (p *CoreParser) parseLiteral(data isa.Datum) (SchemeExpression, error) {
	switch datum := data.(type) {
	case isa.DatumBool:
		return LiteralExpression{Datum: datum}, nil
	case isa.DatumChar:
		return LiteralExpression{Datum: datum}, nil
	default:
		return nil, nil
	}
}
