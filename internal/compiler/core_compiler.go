package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/backend/codegen"
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/typechecker"
	"github.com/certainty/go-braces/internal/compiler/middleend/optimization"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
)

// The core compile is essentially the middle and backend combined.
// It takes core form AST as input and compiles it down to byte code.

type CoreCompiler struct {
	introspectionAPI introspection.API

	// semantic analysis
	typechecker *typechecker.TypeChecker

	// transformation and optimization
	optimizer *optimization.Optimizer

	// code generation
	codegen *codegen.Codegenerator
}

func NewCoreCompiler(introspectionAPI introspection.API) *CoreCompiler {
	return &CoreCompiler{
		introspectionAPI: introspectionAPI,
		typechecker:      typechecker.NewTypeChecker(introspectionAPI),
		optimizer:        optimization.NewOptimizer(introspectionAPI),
		codegen:          codegen.NewCodegenerator(introspectionAPI),
	}
}

func (c *CoreCompiler) CompileModule(coreAst *parser.CoreAST) (*isa.AssemblyModule, error) {
	if err := c.typechecker.Check(coreAst); err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}
	ir, err := ir.LowerToIR(coreAst)

	optimized, err := c.optimizer.Optimize(ir)

	assemblyModule, err := c.codegen.GenerateModule(optimized)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	return assemblyModule, nil

}
