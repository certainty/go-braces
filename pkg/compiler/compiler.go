package compiler

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/pkg/compiler/backend/codegen"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/parser"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/types"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/compiler/middleend/optimization"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	"github.com/certainty/go-braces/pkg/shared/isa"
)

type Compiler struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewCompiler(options CompilerOptions) *Compiler {
	return &Compiler{
		instrumentation: options.instrumentation,
	}
}

func (c Compiler) CompileString(code string, label string) (*isa.AssemblyModule, error) {
	input := lexer.NewStringInput(label, code)
	return c.CompileModule(input)
}

func (c Compiler) CompileModule(input *lexer.Input) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterCompilerModule(input.Origin, string(*input.Buffer))

	ast, err := c.parse(input)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}
	log.Printf("AST: %s", ast.ASTString())

	tpeUniverse, err := c.typeCheck(ast)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	// middleend
	ir, err := c.lowerToIR(ast, tpeUniverse, input.Origin)
	if err != nil {
		return nil, fmt.Errorf("IRError: %w", err)
	}
	log.Printf("IR %v", ir)

	optimizedIr, err := c.optimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	log.Printf("Optimized IR %v", optimizedIr)

	// backend
	assemblyModule, err := c.generateCode(optimizedIr)
	if err != nil {
		return nil, fmt.Errorf("CodeGenError: %w", err)
	}
	log.Printf("AssemblyModule %v", assemblyModule)

	c.instrumentation.LeaveCompilerModule(*assemblyModule)
	return assemblyModule, nil
}

func (c Compiler) parse(input *lexer.Input) (*ast.Source, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	parser := parser.NewParser(c.instrumentation)
	return parser.Parse(input)
}

func (c Compiler) typeCheck(theAST *ast.Source) (types.TypeUniverse, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	typechecker := types.NewChecker(c.instrumentation)
	typeUniverse, err := typechecker.Check(theAST)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	return typeUniverse, nil
}

func (c Compiler) lowerToIR(theAST *ast.Source, tpeUniverse types.TypeUniverse, origin token.Origin) (*ir.Module, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseLowerToIR)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseLowerToIR)

	return ir.LowerToIR(origin, theAST, tpeUniverse)
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

func (c Compiler) generateCode(ir *ir.Module) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	codegen := codegen.NewCodegenerator(c.instrumentation)
	return codegen.GenerateModule(ir)
}
