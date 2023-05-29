package compiler

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/input"
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
	reader := reader.NewReader(c.instrumentation)
	parser := parser.NewParser(c.instrumentation)
	coreCompiler := NewCoreCompiler(c.instrumentation)

	c.instrumentation.Breakpoint("compiler::before_read", &c)
	datum, err := reader.Read(input)
	if err != nil {
		return nil, fmt.Errorf("ReadError: %w", err)
	}

	c.instrumentation.Breakpoint("compiler::before_parse", &c)
	coreAst, err := parser.Parse(datum)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}

	c.instrumentation.Breakpoint("compiler::before_core_compile", &c)
	assemblyModule, err := coreCompiler.CompileModule(coreAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	c.instrumentation.LeaveCompilerModule(*assemblyModule)
	return assemblyModule, nil
}
