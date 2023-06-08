package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Parser struct {
	instrumentation compiler_introspection.Instrumentation
	scanner         *lexer.Scanner
}

func NewParser(instrumentation compiler_introspection.Instrumentation) *Parser {
	return &Parser{
		instrumentation: instrumentation,
		scanner:         nil,
	}
}

func (p *Parser) Parse(input *input.Input) (*AST, error) {
	p.scanner = lexer.New(input.Buffer, input.Origin)

	p.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseParse)
	defer p.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseParse)

	return p.parseInput()
}

func (p *Parser) parseInput() (*AST, error) {
	return nil, fmt.Errorf("not implemented")
}
