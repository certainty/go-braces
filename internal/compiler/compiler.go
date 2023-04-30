package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa/assembly"
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
	reader           *reader.Reader
	parser           *parser.Parser
	coreCompiler     *CoreCompiler
}

func NewCompiler(options CompilerOptions) *Compiler {
	return &Compiler{
		introspectionAPI: options.introspectionAPI,
		reader:           reader.NewReader(options.introspectionAPI),
		parser:           parser.NewParser(options.introspectionAPI),
		coreCompiler:     NewCoreCompiler(options.introspectionAPI),
	}
}

func (c *Compiler) CompileString(code string) (*assembly.AssemblyModule, error) {
	input := location.NewStringInput(code, "ADHOC")
	return c.CompileModule(input)
}

func (c *Compiler) CompileModule(input location.Input) (*assembly.AssemblyModule, error) {
	c.introspectionAPI.SendEvent(introspection.EventStartCompileModule())

	datum, err := c.reader.Read(input)
	if err != nil {
		return nil, fmt.Errorf("ReadError: %w", err)
	}

	coreAst, err := c.parser.Parse(datum)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}

	assemblyModule, err := c.coreCompiler.CompileModule(coreAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	c.introspectionAPI.SendEvent(introspection.EventEndCompileModule())
	return assemblyModule, nil
}
