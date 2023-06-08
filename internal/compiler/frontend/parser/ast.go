package parser

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type AST struct {
	Nodes []Node
}
type Node interface {
	Location() location.Location
}

type Expression interface {
	Node
}

type LiteralExpression struct {
	Value    isa.Value
	location location.Location
}

func NewLiteralExpression(value isa.Value, location location.Location) LiteralExpression {
	return LiteralExpression{
		Value:    value,
		location: location,
	}
}

func (l LiteralExpression) String() string {
	return fmt.Sprintf("Lit{ %s }[%d:%d]", l.Value, l.Location().Line, l.Location().StartOffset)
}

func (l LiteralExpression) Location() location.Location {
	return l.location
}

func NewAST() *AST {
	return &AST{
		Nodes: []Node{},
	}
}

func (ast *AST) String() string {
	return fmt.Sprintf("AST %v ", ast.Nodes)
}

func (ast *AST) AddExpression(expression Expression) {
	ast.Nodes = append(ast.Nodes, expression)
}
