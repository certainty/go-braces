package ast

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
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

func (a *Builder) Result() []Node {
	return a.nodes
}

func (a *Builder) AddNode(node Node) {
	a.nodes = append(a.nodes, node)
}

func (a *Builder) NewBadExpr(location token.Location) BadExpr {
	return BadExpr{
		id:       a.nextID(),
		location: location,
	}
}

func (a *Builder) NewBasicLitExpr(location token.Location, token token.Token) BasicLitExpr {
	return BasicLitExpr{
		id:    a.nextID(),
		Token: token,
	}
}

func (a *Builder) NewIdentifier(location token.Location, name string) Identifier {
	return Identifier{
		location: location,
		id:       a.nextID(),
		Name:     name,
	}
}

func (a *Builder) NewField(name Identifier, tpe *TypeSpec) Field {
	return Field{
		Name: name,
		Type: tpe,
	}
}

func (a *Builder) NewProcDecl(location token.Location, name Identifier, args []Field, result *TypeSpec, body BlockExpr) ProcDecl {
	tpe := ProcType{
		Params: args,
		Result: result,
	}

	return ProcDecl{
		id:       a.nextID(),
		location: location,
		Type:     tpe,
		Name:     name,
		Body:     body,
	}
}

func (a *Builder) NewTypeSpec(location token.Location, name Identifier) TypeSpec {
	return TypeSpec{
		id:       a.nextID(),
		location: location,
		Name:     name,
	}
}

func (a *Builder) NewBlockExpr(location token.Location, statements []Statement) BlockExpr {
	return BlockExpr{
		id:         a.nextID(),
		location:   location,
		Statements: statements,
	}
}

func (a *Builder) NewUnaryExpr(location token.Location, op token.Token, expr Expression) UnaryExpr {
	return UnaryExpr{
		id:   a.nextID(),
		Op:   op,
		Expr: expr,
	}
}

func (a *Builder) NewParenExpr(location token.Location, expr Expression) ParenExpr {
	return ParenExpr{
		id:   a.nextID(),
		Expr: expr,
	}
}

func (a *Builder) NewBinaryExpr(op token.Token, left, right Expression) BinaryExpr {
	return BinaryExpr{
		id:    a.nextID(),
		Op:    op,
		Left:  left,
		Right: right,
	}
}

func (a *Builder) NewExprStatement(expr Expression) ExprStmt {
	return ExprStmt{
		Expr: expr,
	}
}

func (a *Builder) nextID() NodeId {
	a.nodeIds++
	return a.nodeIds
}
