package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/backend/codegen"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/middleend/transformer"
	"github.com/certainty/go-braces/internal/compiler/middleend/typechecker"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa/assembly"
)

// The core compile is essentially the middle and backend combined.
// It takes core form AST as input and compiles it down to byte code.

type CoreCompiler struct {
	introspectionAPI introspection.API

	// semantic analysis
	typechecker *typechecker.TypeChecker

	// transformation and optimization
	cpsTransformer *transformer.CPSTransformer

	// code generation
	codegen *codegen.Codegenerator
}

func NewCoreCompiler(introspectionAPI introspection.API) *CoreCompiler {
	return &CoreCompiler{
		introspectionAPI: introspectionAPI,
		typechecker:      typechecker.NewTypeChecker(introspectionAPI),
		cpsTransformer:   transformer.NewCpsTransformer(introspectionAPI),
		codegen:          codegen.NewCodegenerator(introspectionAPI),
	}
}

func (c *CoreCompiler) CompileModule(coreAst *parser.CoreAST) (*assembly.AssemblyModule, error) {
	if err := c.typechecker.Check(coreAst); err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	cpsAst, err := c.cpsTransformer.Transform(coreAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	assemblyModule, err := c.codegen.GenerateModule(cpsAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	return assemblyModule, nil

}
