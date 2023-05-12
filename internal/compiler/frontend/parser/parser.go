package parser

import (
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/expander"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	expander        *expander.Expander
	coreParser      *CoreParser
}

func NewParser(instrumentation compiler_introspection.Instrumentation) *Parser {
	return &Parser{
		instrumentation: instrumentation,
		expander:        expander.NewExpander(instrumentation),
		coreParser:      NewCoreParser(instrumentation),
	}
}

func (p *Parser) Parse(data *reader.DatumAST) (*CoreAST, error) {
	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	p.expander = expander.NewExpander(p.instrumentation)
	p.coreParser = NewCoreParser(p.instrumentation)
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

func (p *Parser) doParse(data isa.Datum) (SchemeExpression, error) {
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
