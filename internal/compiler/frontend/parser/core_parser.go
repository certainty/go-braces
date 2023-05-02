package parser

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
)

type CoreParser struct {
	introspectionAPI introspection.API
}

func NewCoreParser(introspectionAPI introspection.API) *CoreParser {
	return &CoreParser{introspectionAPI: introspectionAPI}
}

func (p *CoreParser) Parse(data reader.Datum) (SchemeExpression, error) {
	return p.parseLiteral(data)
}

func (p *CoreParser) parseLiteral(data reader.Datum) (SchemeExpression, error) {
	switch datum := data.(type) {
	case reader.DatumBool:
		{
			return LiteralExpression{Datum: datum}, nil
		}
	default:
		return nil, nil
	}
}
