package ast

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
)

type Builder struct {
	nodeIds astutils.NodeIdManager
}

func NewBuilder() *Builder {
	return &Builder{
		nodeIds: astutils.NewNodeIdManager(),
	}
}

func (a *Builder) NewBadExpr(location token.Location) BadExpr {
	return BadExpr{
		id:       a.nodeIds.Next(),
		location: location,
	}
}

func (a *Builder) NewBasicLitExpr(location token.Location, token token.Token) BasicLitExpr {
	return BasicLitExpr{
		id:    a.nodeIds.Next(),
		Token: token,
	}
}

func (a *Builder) NewIdentifier(location token.Location, name string) Identifier {
	return Identifier{
		location: location,
		id:       a.nodeIds.Next(),
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
		id:       a.nodeIds.Next(),
		location: location,
		Type:     tpe,
		Name:     name,
		Body:     body,
	}
}

func (a *Builder) NewTypeSpec(location token.Location, name Identifier) TypeSpec {
	return TypeSpec{
		id:       a.nodeIds.Next(),
		location: location,
		Name:     name,
	}
}

func (a *Builder) NewBlockExpr(location token.Location, statements []Statement) BlockExpr {
	return BlockExpr{
		id:         a.nodeIds.Next(),
		location:   location,
		Statements: statements,
	}
}

func (a *Builder) NewUnaryExpr(location token.Location, op token.Token, expr Expression) UnaryExpr {
	return UnaryExpr{
		id:   a.nodeIds.Next(),
		Op:   op,
		Expr: expr,
	}
}

func (a *Builder) NewParenExpr(location token.Location, expr Expression) ParenExpr {
	return ParenExpr{
		id:   a.nodeIds.Next(),
		Expr: expr,
	}
}

func (a *Builder) NewBinaryExpr(op token.Token, left, right Expression) BinaryExpr {
	return BinaryExpr{
		id:    a.nodeIds.Next(),
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

func (a *Builder) NewBadDecl(location token.Location) BadDecl {
	return BadDecl{
		id:       a.nodeIds.Next(),
		location: location,
	}
}
