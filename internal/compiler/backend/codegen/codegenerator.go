package codegen

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa/assembly"
)

type Codegenerator struct {
	introspectionAPI introspection.API
}

func NewCodegenerator(introspectionAPI introspection.API) *Codegenerator {
	return &Codegenerator{introspectionAPI: introspectionAPI}
}

func (c *Codegenerator) GenerateModule(cpsAst *parser.CoreAST) (*assembly.AssemblyModule, error) {
	return nil, nil
}
