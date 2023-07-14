package ast

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/lexer"
	"github.com/certainty/go-braces/internal/compiler/location"
)

type Builder struct {
	nodes   []Node
	nodeIds NodeId
}

func NewBuilder() *Builder {
	return &Builder{
		nodes:   make([]Node, 0),
		nodeIds: 0,
	}
}

func (a *Builder) Result() *AST {
	return &AST{
		Nodes: a.nodes,
	}
}

func (a *Builder) nextID() NodeId {
	a.nodeIds++
	return a.nodeIds
}

func (a *Builder) AddNode(node Node) {
	a.nodes = append(a.nodes, node)
}

func (a *Builder) NewIdentifier(id string, location location.Location) Identifier {
	return Identifier{
		id:       a.nextID(),
		Label:    id,
		location: location,
	}
}

func (a *Builder) NewTypeDecl(name Identifier, location location.Location) TypeDecl {
	return TypeDecl{
		id:       a.nextID(),
		Name:     name,
		location: location,
	}
}

func (a *Builder) NewPackageDecl(name Identifier, location location.Location) PackageDecl {
	return PackageDecl{
		id:       a.nextID(),
		Name:     name,
		location: location,
	}
}

func (a *Builder) NewFunctionDecl(
	t TypeDecl,
	name Identifier,
	arguments []ArgumentDecl,
	body Block,
	location location.Location,
) CallableDecl {
	return CallableDecl{
		id:          a.nextID(),
		TpeDecl:     t,
		IsProcedure: false,
		Name:        name,
		Arguments:   arguments,
		Body:        body,
		location:    location,
	}
}

func (a *Builder) NewProcedureDecl(
	t TypeDecl,
	name Identifier,
	arguments []ArgumentDecl,
	body Block,
	location location.Location,
) CallableDecl {
	return CallableDecl{
		id:          a.nextID(),
		TpeDecl:     t,
		IsProcedure: true,
		Name:        name,
		Arguments:   arguments,
		Body:        body,
		location:    location,
	}
}

func (a *Builder) NewArgumentDecl(name Identifier, t TypeDecl, location location.Location) ArgumentDecl {
	return ArgumentDecl{
		id:       a.nextID(),
		Name:     name,
		TpeDecl:  t,
		location: location,
	}
}

func (a *Builder) NewBlock(body []Node, location location.Location) Block {
	return Block{
		id:       a.nextID(),
		Code:     body,
		location: location,
	}
}

func (a *Builder) NewBinOp(location location.Location, operator BinaryOperator, left Expression, right Expression) BinaryExpression {
	return BinaryExpression{
		id:       a.nextID(),
		location: location,
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (a *Builder) NewUnaryOp(location location.Location, operator UnaryOperator, operand Expression) UnaryExpression {
	return UnaryExpression{
		location: location,
		Operator: operator,
		Operand:  operand,
	}
}

func (a *Builder) NewLiteralExpression(token lexer.Token, location location.Location) LiteralExpression {
	return LiteralExpression{
		id:       a.nextID(),
		Token:    token,
		Value:    token.Value,
		location: location,
	}
}
