package ast

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/types"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type AST struct {
	Nodes []Node
}

type Node interface {
	Location() location.Location
}

type Declaration interface{ Node }
type Expression interface{ Node }
type Statement interface{ Node }

type Identifier struct {
	ID       string
	location location.Location
}

var _ Node = (*Identifier)(nil)

func (id Identifier) Location() location.Location {
	return id.location
}

var _ Node = (*Identifier)(nil)
var _ Declaration = (*Identifier)(nil)

type ArgumentDecl struct {
	Name     Identifier
	Type     types.Type
	location location.Location
}

var _ Node = (*ArgumentDecl)(nil)
var _ Declaration = (*ArgumentDecl)(nil)

func (def ArgumentDecl) Location() location.Location {
	return def.location
}

type FunctionDecl struct {
	Type      types.Type
	Name      Identifier
	Arguments []ArgumentDecl
	Body      []Node
	location  location.Location
}

var _ Node = (*FunctionDecl)(nil)
var _ Declaration = (*FunctionDecl)(nil)

func (d FunctionDecl) Location() location.Location {
	return d.location
}

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
