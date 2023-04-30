package codegen

import (
	"github.com/certainty/go-braces/internal/compiler/middleend/ir"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa/assembly"
)

type Codegenerator struct {
	introspectionAPI introspection.API
}

func NewCodegenerator(introspectionAPI introspection.API) *Codegenerator {
	return &Codegenerator{introspectionAPI: introspectionAPI}
}

func (c *Codegenerator) GenerateModule(ssa *ir.SSA) (*assembly.AssemblyModule, error) {
	return nil, nil
}
