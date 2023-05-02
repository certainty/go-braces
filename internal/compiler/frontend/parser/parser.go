package parser

import (
	"log"

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
	}
}

func (p *Parser) Parse(data *reader.DatumAST) (*CoreAST, error) {
	p.expander = expander.NewExpander(p.introspectionAPI)
	p.coreParser = NewCoreParser(p.introspectionAPI)
	coreAst := NewCoreAST()

	for _, datum := range data.Data {
		expr, err := p.doParse(datum)

		if err != nil {
			// track errors, try to recover and go on
		} else {
			log.Printf("Adding expression %v", expr)
			coreAst.AddExpression(expr)
		}
	}

	return coreAst, nil
}

func (p *Parser) doParse(data reader.Datum) (SchemeExpression, error) {
	expanded, err := p.expander.Expand(data)
	if err != nil {
		return nil, err
	}

	// expansions might result in nothing, which is ok
	if expanded == nil {
		return nil, nil
	}

	return p.coreParser.Parse(expanded)
}
