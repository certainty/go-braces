package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/backend/codegen"
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/typechecker"
	"github.com/certainty/go-braces/internal/compiler/middleend/optimization"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

// The CoreCompiler deals with core scheme forms
// It has a small frontend for type analysis and also bundles
// the optimizer and the code generator.

// The core compiler is invoked after and expression has been
// macro-expanded during the parsing process of the compiler.

type CoreCompiler struct {
	instrumentation compiler_introspection.Instrumentation

	// semantic analysis
	typechecker *typechecker.TypeChecker

	// transformation and optimization
	optimizer *optimization.Optimizer

	// code generation
	codegen *codegen.Codegenerator
}

func NewCoreCompiler(instrumentation compiler_introspection.Instrumentation) *CoreCompiler {
	return &CoreCompiler{
		instrumentation: instrumentation,
		typechecker:     typechecker.NewTypeChecker(instrumentation),
		optimizer:       optimization.NewOptimizer(instrumentation),
		codegen:         codegen.NewCodegenerator(instrumentation),
	}
}

func (c *CoreCompiler) CompileModule(coreAst *parser.CoreAST) (*isa.AssemblyModule, error) {
	if err := c.typechecker.Check(coreAst); err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	ir, err := ir.LowerToIR(coreAst)
	if err != nil {
		return nil, fmt.Errorf("IRError: %w", err)
	}

	optimized, err := c.optimizer.Optimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}

	assemblyModule, err := c.codegen.GenerateModule(optimized)
	if err != nil {
		return nil, fmt.Errorf("CodegeneratorError: %w", err)
	}

	return assemblyModule, nil
}
