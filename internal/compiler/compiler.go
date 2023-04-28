package compiler

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
)

// The compile follows a traditional compile design of frontend, middleend and backend
// Since scheme has rather rich meta syntactical capabilities with its macro system
// we separate the core compiler, which deals with core forms, afer they have been transformed
// parsed and expanded.
//
// This struct is the main interface to the compiler and houses the compiler frontend (syntactic analysisee and macro expansion)
// as well as the core compiler which deals with the rest.
type Compiler struct {
	introspectionAPI *introspection.API
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

func (c *Compiler) JitCompile(code string) (*CompilationUnit, error) {
	return c.compileUnit(code, "JIT")
}

func (c *Compiler) compileUnit(code string, name string) (*CompilationUnit, error) {
	datum, err := c.reader.Read(code)
	if err != nil {
		return nil, fmt.Errorf("ReadError: %w", err)
	}

	coreAst, err := c.parser.Parse(datum)
	if err != nil {
		return nil, fmt.Errorf("ParseError: %w", err)
	}

	compilationUnit, err := c.coreCompiler.Compile(coreAst)
	if err != nil {
		return nil, fmt.Errorf("CompilerBug: %w", err)
	}

	return compilationUnit, nil
}
