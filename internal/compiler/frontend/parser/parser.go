package parser

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/expander"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
)

type Parser struct {
	introspectionAPI introspection.API
	expander         *expander.Expander
	coreParser       *CoreParser
}

func NewParser(introspectionAPI introspection.API) *Parser {
	return &Parser{
		introspectionAPI: introspectionAPI,
		expander:         expander.NewExpander(introspectionAPI),
		coreParser:       NewCoreParser(introspectionAPI),
	}
}

func (p *Parser) Parse(data *reader.DatumAST) (*CoreAST, error) {
	return nil, nil
}
