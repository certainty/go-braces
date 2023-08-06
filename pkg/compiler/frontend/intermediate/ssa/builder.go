package ssa

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
)

type Builder struct {
	nodeIds  astutils.NodeIdManager
	versions astutils.VersionManager
}

func NewBuilder() Builder {
	return Builder{
		nodeIds:  astutils.NewNodeIdManager("ssa-"),
		versions: astutils.NewVersionManager(),
	}
}

func (b *Builder) AtomicLitExpr(expr ir.AtomicLitExpr) AtomicLitExpr {
	return AtomicLitExpr{
		id:       b.nodeIds.Next(),
		Value:    expr.Value,
		IrExprId: expr.ID(),
	}
}

func (b *Builder) BinaryExpr(expr ir.BinaryExpr, leftVar Expression, rightVar Expression) BinaryExpr {
	return BinaryExpr{
		id:       b.nodeIds.Next(),
		IrExprId: expr.ID(),
		Op:       expr.Op,
		Left:     leftVar,
		Right:    rightVar,
	}
}

func (b *Builder) VariableExpr(variable Variable) VariableExpr {
	return VariableExpr{
		id:       b.nodeIds.Next(),
		Variable: variable,
	}
}

type BasicBlockBuilder struct {
	nodeIds  *astutils.NodeIdManager
	versions *astutils.VersionManager
	block    *BasicBlock
}

func (b *Builder) BlockBuilder(name ir.Label) *BasicBlockBuilder {
	return &BasicBlockBuilder{
		nodeIds:  &b.nodeIds,
		versions: &b.versions,
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

func (b *BasicBlockBuilder) AddExpr(expr Expression) {
	b.block.Statements = append(b.block.Statements, ExprStatement{
		id:   b.nodeIds.Next(),
		Expr: expr,
	})
}

func (b *BasicBlockBuilder) Variable(prefix string) Variable {
	return Variable{
		Prefix:  prefix,
		Version: b.versions.Next(),
	}
}

func (b *BasicBlockBuilder) AddStatement(statement Statement) {
	b.block.Statements = append(b.block.Statements, statement)
}

func (b *BasicBlockBuilder) Close() *BasicBlock {
	return b.block
}
