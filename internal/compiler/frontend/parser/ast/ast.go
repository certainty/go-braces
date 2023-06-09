package ast

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
)

type AST struct {
	Nodes []Node
}

type Node interface {
	Location() location.Location
}

type LValue interface{ Node }
type RValue interface{ Node }

type UnaryMinusExpression struct {
	location location.Location
	Operand  Expression
}

func NewUnaryMinusExpression(location location.Location, operand Expression) UnaryMinusExpression {
	return UnaryMinusExpression{
		location: location,
		Operand:  operand,
	}
}

func (u UnaryMinusExpression) Location() location.Location {
	return u.location
}

type UnaryNotExpression struct {
	location location.Location
	Operand  Expression
}

func NewUnaryNotExpression(location location.Location, operand Expression) UnaryNotExpression {
	return UnaryNotExpression{
		location: location,
		Operand:  operand,
	}
}

func (n UnaryNotExpression) Location() location.Location {
	return n.location
}

func New() *AST {
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
