package ast

import (
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

func New() *AST {
	return &AST{
		Nodes: []Node{},
	}
}

func (ast *AST) ASTring() string {
	return NewASTWriter().Write(ast)
}

func (ast *AST) AddExpression(expression Expression) {
	ast.Nodes = append(ast.Nodes, expression)
}
