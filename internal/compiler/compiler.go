package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/introspection/introspection_events"
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
	introspectionAPI introspection.API
}

func NewCompiler(options CompilerOptions) *Compiler {
	return &Compiler{
		introspectionAPI: options.introspectionAPI,
	}
}

func (c Compiler) CompileString(code string) (*isa.AssemblyModule, error) {
	c.introspectionAPI.SendEvent(introspection_events.EventBeginCompileString{Input: code})
	input := input.NewStringInput("ADHOC", code)
	return c.CompileModule(input)
}

func (c Compiler) CompileModule(input *input.Input) (*isa.AssemblyModule, error) {
	c.introspectionAPI.SendEvent(introspection_events.EventBeginCompileModule{})
	c.introspectionAPI.SingleStepBarrier(&c)

	reader := reader.NewReader(c.introspectionAPI)
	parser := parser.NewParser(c.introspectionAPI)
	coreCompiler := NewCoreCompiler(c.introspectionAPI)

	datum, err := reader.Read(input)
	if err != nil {
		return nil, fmt.Errorf("ReadError: %w", err)
	}

	coreAst, err := parser.Parse(datum)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}

	assemblyModule, err := coreCompiler.CompileModule(coreAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	c.introspectionAPI.SendEvent(introspection_events.EventEndCompileModule{})
	return assemblyModule, nil
}
