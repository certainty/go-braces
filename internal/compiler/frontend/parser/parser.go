package parser

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
)

type Parser struct {
	introspectionAPI introspection.API
}

func NewParser(introspectionAPI introspection.API) *Parser {
	return &Parser{
		introspectionAPI: introspectionAPI,
	}
}

func (p *Parser) Parse(data *reader.DatumAST) (*CoreAST, error) {
	// expander := expander.NewExpander(p.introspectionAPI)
	// coreParser := NewCoreParser(p.introspectionAPI)

	return nil, nil
}
