package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler/backend/codegen"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/parser"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/types"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/compiler/middleend/optimization"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	"github.com/certainty/go-braces/pkg/shared/isa"
	log "github.com/sirupsen/logrus"
)

type Compiler struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewCompiler(options CompilerOptions) *Compiler {
	return &Compiler{
		instrumentation: options.instrumentation,
	}
}

func (c *Compiler) CompileFile(path string) (*isa.AssemblyModule, error) {
	input, err := lexer.NewFileInput(path)
	if err != nil {
		return nil, err
	}
	return c.CompileModule(input)
}

func (c *Compiler) CompileString(code string, label string) (*isa.AssemblyModule, error) {
	input := lexer.NewStringInput(label, code)
	return c.CompileModule(input)
}

func (c Compiler) CompileModule(input *lexer.Input) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterCompilerModule(input.Origin, string(*input.Buffer))

	theAST, err := c.parse(input)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}
	log.Debugf("AST: %v", theAST)
	log.Debugf("AST: %s", ast.Print(theAST, ast.PrintTruthfully()))

	tpeUniverse, err := c.typeCheck(theAST)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	// middleend
	irModule, err := c.lowerToIR(theAST, tpeUniverse, input.Origin)
	if err != nil {
		return nil, fmt.Errorf("IRError: %w", err)
	}
	log.Debugf("IR %v", irModule)
	log.Debugf("IR %s", ir.Print(irModule, ir.PrintTruthfully()))

	ssa, err := c.optimize(irModule)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	log.Debugf("Optimized IR %v", ssa)
	log.Debugf("Optmized IR %v", ir.Print(ssa, ir.PrintTruthfully().ForSSA()))

	// backend
	assemblyModule, err := c.generateCode(ssa)
	if err != nil {
		return nil, fmt.Errorf("CodeGenError: %w", err)
	}
	log.Debugf("AssemblyModule %v", assemblyModule)

	c.instrumentation.LeaveCompilerModule(*assemblyModule)
	return assemblyModule, nil
}

func (c Compiler) parse(input *lexer.Input) (*ast.Source, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	parser := parser.NewParser(c.instrumentation)
	return parser.Parse(input)
}

func (c Compiler) typeCheck(theAST *ast.Source) (*types.TypeUniverse, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	typechecker := types.NewChecker(c.instrumentation)
	typeUniverse, err := typechecker.Check(theAST)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}
	return typeUniverse, nil
}

func (c Compiler) lowerToIR(theAST *ast.Source, tpeUniverse *types.TypeUniverse, origin token.Origin) (*ir.Module, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseLowerToIR)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseLowerToIR)

	return intermediate.Lower(origin, theAST, tpeUniverse)
}

func (c Compiler) optimize(ir *ir.Module) (*ir.Module, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseOptimize)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseOptimize)

	optimizer := optimization.NewOptimizer(c.instrumentation)
	optimized, err := optimizer.Optimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	return optimized, nil
}

func (c Compiler) generateCode(ssa *ir.Module) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	codegen := codegen.NewCodegenerator(c.instrumentation)
	return codegen.GenerateModule(ssa)
}
