package ssa

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
)

type (
	Node interface {
		astutils.NodeIded
	}

	Expression interface {
		Node
		ssaExprNode()
	}

	Declaration interface {
		ssaDeclNode()
	}

	Statement interface {
		ssaStmtNode()
	}

	ProcDecl struct {
		id     astutils.NodeId
		Name   ir.Label
		irDecl ir.ProcDecl
		Blocks []*BasicBlock
	}

	BasicBlock struct {
		id           astutils.NodeId
		label        ir.Label
		Statements   []Statement
		Predecessors []*BasicBlock
		Successors   []*BasicBlock
	}

	BinaryExpr struct {
		id       astutils.NodeId
		IrExprId astutils.NodeId
		Op       token.Token
		Left     Expression
		Right    Expression
	}

	AtomicLitExpr struct {
		id       astutils.NodeId
		Value    token.Token
		IrExprId astutils.NodeId
	}

	VariableExpr struct {
		id       astutils.NodeId
		Variable Variable
	}

	Phi struct {
		id          astutils.NodeId
		Assignments map[*BasicBlock]Expression
	}

	SetStmt struct {
		id       astutils.NodeId
		Variable Variable
		Value    Expression
	}

	ReturnStmt struct {
		id    astutils.NodeId
		Value Expression
	}

	ExprStatement struct {
		id   astutils.NodeId
		Expr Expression
	}

	Variable struct {
		Prefix  string
		Version astutils.Version
	}

	Module struct {
		Name         ir.Label
		Declarations []Declaration
	}
)

func (BinaryExpr) ssaExprNode()    {}
func (AtomicLitExpr) ssaExprNode() {}
func (Phi) ssaExprNode()           {}
func (VariableExpr) ssaExprNode()  {}

func (e AtomicLitExpr) ID() astutils.NodeId { return e.id }
func (e BinaryExpr) ID() astutils.NodeId    { return e.id }
func (e Phi) ID() astutils.NodeId           { return e.id }
func (e VariableExpr) ID() astutils.NodeId  { return e.id }

func (ProcDecl) ssaDeclNode()      {}
func (BasicBlock) ssaStmtNode()    {}
func (SetStmt) ssaStmtNode()       {}
func (ReturnStmt) ssaStmtNode()    {}
func (ExprStatement) ssaStmtNode() {}

func (p ProcDecl) ID() astutils.NodeId   { return p.id }
func (r ReturnStmt) ID() astutils.NodeId { return r.id }
func (s SetStmt) ID() astutils.NodeId    { return s.id }
func (b BasicBlock) ID() astutils.NodeId { return b.id }
func (e ExprStatement) ID() astutils.NodeId {
	return e.id
}

func (b *BasicBlock) AddPredecessor(block *BasicBlock) {
	if b.Predecessors == nil {
		b.Predecessors = make([]*BasicBlock, 0)
	}
	b.Predecessors = append(b.Predecessors, block)
}

func (b *BasicBlock) AddSuccessor(block *BasicBlock) {
	if b.Successors == nil {
		b.Successors = make([]*BasicBlock, 0)
	}
	b.Successors = append(b.Successors, block)
}
