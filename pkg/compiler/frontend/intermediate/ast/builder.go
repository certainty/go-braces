package ast

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	hl "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
)

type Builder struct {
	nodeIds   astutils.NodeIdManager
	variables astutils.VersionManager
}

type BlockBuilder struct {
	*Builder
	expr *BasicBlock
}

func NewBuilder() *Builder {
	return &Builder{
		nodeIds: astutils.NewNodeIdManager("ir-"),
	}
}

func (b *Builder) ProcDecl(name Label, tpe types.Procedure, hlDecl hl.ProcDecl) *ProcDecl {
	return &ProcDecl{
		id:   b.nodeIds.Next(),
		Type: tpe,
		Name: name,
	}
}

func (b *Builder) AtomicLit(tpe types.Type, value token.Token, hlExpr astutils.NodeId) *AtomicLitExpr {
	return &AtomicLitExpr{
		id:           b.nodeIds.Next(),
		tpe:          tpe,
		Value:        value,
		hlExprNodeId: hlExpr,
	}
}

func (b *Builder) BinaryExpr(tpe types.Type, op token.Token, left Expression, right Expression, hlExpr astutils.NodeId) *BinaryExpr {
	return &BinaryExpr{
		id:           b.nodeIds.Next(),
		Type:         tpe,
		Op:           op,
		Left:         left,
		Right:        right,
		hlExprNodeId: hlExpr,
	}
}

func (b *Builder) ExprStatement(expr Expression) Statement {
	return &ExprStatement{
		Expr: expr,
	}
}

func (b *Builder) BlockBuilder(blockLabel Label, origin *astutils.NodeId) *BlockBuilder {
	return &BlockBuilder{
		Builder: b,
		expr:    &BasicBlock{Label: blockLabel, Statements: make([]Statement, 0)},
	}
}

func (b *Builder) Label(name string, origin *astutils.NodeId) Label {
	return Label{
		id:     b.nodeIds.Next(),
		Origin: origin,
		Value:  name,
	}
}

func (b *BlockBuilder) AddStatement(statement Statement) {
	b.expr.Statements = append(b.expr.Statements, statement)
}

func (b *BlockBuilder) ReplaceLastStatement(statement Statement) {
	b.expr.Statements[len(b.expr.Statements)-1] = statement
}

func (b *BlockBuilder) AddAssignment(variable *Variable, value Expression) {
	b.AddStatement(&AssignStmt{
		id:       b.nodeIds.Next(),
		Variable: variable,
		Expr:     value,
	})
}

func (b *Builder) ReturnStmt(expr Expression) Statement {
	return &ReturnStmt{
		id:    b.nodeIds.Next(),
		Value: expr,
	}
}

func (b *Builder) Variable(name string) *Variable {
	return &Variable{
		id:      b.nodeIds.Next(),
		Version: b.variables.Next(),
	}
}

func (b *BlockBuilder) Close() *BasicBlock {
	return b.expr
}
