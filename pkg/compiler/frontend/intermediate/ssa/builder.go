package ssa

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
)

type Builder struct {
	nodeIds astutils.NodeIdManager
}

func NewBuilder() Builder {
	return Builder{
		nodeIds: astutils.NewNodeIdManager("ssa-"),
	}
}

func (b *Builder) AtomicLitExpr(expr ir.AtomicLitExpr) AtomicLitExpr {
	return AtomicLitExpr{
		id:       b.nodeIds.Next(),
		Value:    expr.Value,
		IrExprId: expr.ID(),
	}
}

func (b *Builder) BinaryExpr(expr ir.BinaryExpr, leftVar Variable, rightVar Variable) BinaryExpr {
	return BinaryExpr{
		id:       b.nodeIds.Next(),
		IrExprId: expr.ID(),
		Op:       expr.Op,
		Left:     leftVar,
		Right:    rightVar,
	}
}

type BasicBlockBuilder struct {
	nodeIds *astutils.NodeIdManager
	block   *BasicBlock
}

func (b *Builder) BlockBuilder(name ir.Label) *BasicBlockBuilder {
	return &BasicBlockBuilder{
		nodeIds: &b.nodeIds,
		block: &BasicBlock{
			id:         b.nodeIds.Next(),
			label:      name,
			Statements: make([]Statement, 0),
		},
	}
}

func (b *BasicBlockBuilder) AddAssignment(expr Expression) Variable {
	variable := b.Variable("tmp")

	assignment := SetStmt{
		id:       b.nodeIds.Next(),
		Variable: variable,
		Value:    expr,
	}

	b.block.Statements = append(b.block.Statements, assignment)
	return variable
}

func (b *BasicBlockBuilder) Variable(prefix string) Variable {
	return Variable{
		Prefix: prefix,
		id:     b.nodeIds.Next(),
	}
}

func (b *BasicBlockBuilder) AddStatement(statement Statement) {
	b.block.Statements = append(b.block.Statements, statement)
}

func (b *BasicBlockBuilder) Close() *BasicBlock {
	return b.block
}
