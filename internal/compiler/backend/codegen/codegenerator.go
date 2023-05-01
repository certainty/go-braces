package codegen

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type Codegenerator struct {
	introspectionAPI introspection.API
}

func NewCodegenerator(introspectionAPI introspection.API) *Codegenerator {
	return &Codegenerator{introspectionAPI: introspectionAPI}
}

func (c *Codegenerator) GenerateModule(intermediate *ir.IR) (*isa.AssemblyModule, error) {
	return nil, nil
}
