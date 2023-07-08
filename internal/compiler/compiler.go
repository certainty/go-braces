package compiler

import (
	"fmt"
	"log"

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

	ast, err := c.parse(input)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}
	log.Printf("AST: %s", ast.ASTring())

	err = c.typeCheck(ast)
	if err != nil {
		return nil, fmt.Errorf("TypeError: %w", err)
	}

	// middleend
	ir, err := c.lowerToIR(ast)
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

func (c Compiler) parse(input *input.Input) (*ast.AST, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	parser := parser.NewParser(c.instrumentation)
	return parser.Parse(input)
}

func (c Compiler) typeCheck(theAST *ast.AST) error {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	typechecker := typechecker.NewTypeChecker(c.instrumentation)
	if err := typechecker.Check(theAST); err != nil {
		return fmt.Errorf("TypeError: %w", err)
	}

	return nil
}

func (c Compiler) lowerToIR(theAST *ast.AST) (*ir.IR, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseLowerToIR)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseLowerToIR)

	return ir.LowerToIR(theAST)
}

func (c Compiler) optimize(ir *ir.IR) (*ir.IR, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseOptimize)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseOptimize)

	optimizer := optimization.NewOptimizer(c.instrumentation)
	optimized, err := optimizer.Optimize(ir)
	if err != nil {
		return nil, fmt.Errorf("OptimizerError: %w", err)
	}
	return optimized, nil
}

func (c Compiler) generateCode(ir *ir.IR) (*isa.AssemblyModule, error) {
	c.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseCodegen)
	defer c.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseCodegen)

	codegen := codegen.NewCodegenerator(c.instrumentation)
	return codegen.GenerateModule(ir)
}
