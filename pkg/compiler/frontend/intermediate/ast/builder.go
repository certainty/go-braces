package ast

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	hl "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
)

type Builder struct {
	nodeIds astutils.NodeIdManager
}

type BlockBuilder struct {
	Builder
	expr *BlockExpr
}

func NewBuilder() *Builder {
	return &Builder{
		nodeIds: astutils.NewNodeIdManager("ir-"),
	}
}

func (b *Builder) ProcDecl(name Label, tpe types.Procedure, hlDecl hl.ProcDecl) ProcDecl {
	return ProcDecl{
		id:   b.nodeIds.Next(),
		Type: tpe,
		Name: name,
	}
}

func (b *Builder) AtomicLit(tpe types.Type, hlExpr astutils.NodeId) AtomicLitExpr {
	return AtomicLitExpr{
		id:           b.nodeIds.Next(),
		tpe:          tpe,
		hlExprNodeId: hlExpr,
	}
}

func (b *Builder) BinaryExpr(tpe types.Type, op token.Token, left Expression, right Expression, hlExpr astutils.NodeId) BinaryExpr {
	return BinaryExpr{
		id:           b.nodeIds.Next(),
		tpe:          tpe,
		Left:         left,
		Right:        right,
		hlExprNodeId: hlExpr,
	}
}

func (b *Builder) ExprStatement(expr Expression) Statement {
	return ExprStatement{
		Expr: expr,
	}
}

func (b *Builder) BlockBuilder(blockLabel string) *BlockBuilder {
	return &BlockBuilder{
		Builder: *b,
		expr:    &BlockExpr{Label: Label(blockLabel), Statements: make([]Statement, 0)},
	}
}

func (b *BlockBuilder) AddStatement(statement Statement) {
	b.expr.Statements = append(b.expr.Statements, statement)
}

func (b *BlockBuilder) Close() BlockExpr {
	// todo add terminating return if it's missing
	return *b.expr
}
