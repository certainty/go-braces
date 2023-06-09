package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/backend/codegen"
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/compiler/frontend/typechecker"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/compiler/middleend/optimization"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

// The compile follows a traditional compile design of frontend, middleend and backend
// Since scheme has rather rich meta syntactical capabilities with its macro system
// we separate the core compiler, which deals with core forms, after they have been transformed
// parsed and expanded, from the rest.
//
// This struct is the main interface to the compiler and houses the compiler frontend
// (syntactic analysises and macro expansion) as well as the core compiler which deals with the rest.
type Compiler struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewCompiler(options CompilerOptions) *Compiler {
	return &Compiler{
		instrumentation: options.instrumentation,
	}
}

func (c Compiler) CompileString(code string, label string) (*isa.AssemblyModule, error) {
	input := input.NewStringInput(label, code)
	return c.CompileModule(input)
}

func (c Compiler) CompileModule(input *input.Input) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterCompilerModule(input.Origin, string(*input.Buffer))

	coreAst, err := c.phaseParse(input)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}
	coreAst, err = c.phaseTypeCheck(coreAst)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}
	ir, err := c.phaseLowerToIR(coreAst)
	if err != nil {
		return nil, fmt.Errorf("IRError: %w", err)
	}
	optimizedIr, err := c.phaseOptimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	assemblyModule, err := c.phaseCodeGen(optimizedIr)
	if err != nil {
		return nil, fmt.Errorf("CodeGenError: %w", err)
	}

	c.instrumentation.LeaveCompilerModule(*assemblyModule)
	return assemblyModule, nil
}

func (c Compiler) phaseParse(input *input.Input) (*ast.AST, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	parser := parser.NewParser(c.instrumentation)
	return parser.Parse(input)
}

func (c Compiler) phaseTypeCheck(theAST *ast.AST) (*ast.AST, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	typechecker := typechecker.NewTypeChecker(c.instrumentation)
	if err := typechecker.Check(theAST); err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	return theAST, nil
}

func (c Compiler) phaseLowerToIR(theAST *ast.AST) (*ir.IR, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseLowerToIR)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseLowerToIR)

	return ir.LowerToIR(theAST)
}

func (c Compiler) phaseOptimize(ir *ir.IR) (*ir.IR, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseOptimize)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseOptimize)

	optimizer := optimization.NewOptimizer(c.instrumentation)
	optimized, err := optimizer.Optimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	return optimized, nil
}

func (c Compiler) phaseCodeGen(ir *ir.IR) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	codegen := codegen.NewCodegenerator(c.instrumentation)
	return codegen.GenerateModule(ir)
}
