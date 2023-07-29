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
		nodeIds: astutils.NewNodeIdManager(),
	}
}

func (b *Builder) AtomicLitExpr(expr ir.AtomicLitExpr) AtomicLitExpr {
	return AtomicLitExpr{
		id:     b.nodeIds.Next(),
		IrExpr: expr,
	}
}

func (b *Builder) BinaryExpr(expr ir.BinaryExpr, leftVar Variable, rightVar Variable) BinaryExpr {
	return BinaryExpr{
		id:     b.nodeIds.Next(),
		IrExpr: expr,
		Left:   leftVar,
		Right:  rightVar,
	}
}

type BasicBlockBuilder struct {
	nodeIds *astutils.NodeIdManager
	block   *BasicBlock
}

func (b *Builder) BlockBuilder(name ir.Label) BasicBlockBuilder {
	return BasicBlockBuilder{
		nodeIds: &b.nodeIds,
		block: &BasicBlock{
			id:         b.nodeIds.Next(),
			label:      name,
			Statements: make([]Statement, 0),
		},
	}
}

func (b *BasicBlockBuilder) AddAssingment(expr Expression) Variable {
	variable := b.Variable("t")
	assignment := SetStmt{
		id:       b.nodeIds.Next(),
		Variable: variable,
	}

	b.block.Statements = append(b.block.Statements, assignment)
	return variable
}

func (b *BasicBlockBuilder) Variable(prefix ir.Label) Variable {
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
