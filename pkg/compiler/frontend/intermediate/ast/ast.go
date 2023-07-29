package ast

// The IR is a simplified, typed, version of the highlevel language.
// It reduces it to couple core constructs:
// * procedures (used for both functions and procedures of the high level language)
// * blocks (named sequences of instructions that can be addressed)
// * binary expressions (all expressions are binary expressions in this language)
// * literals (for highlevel values, labels and structs, no more support for sets, maps, etc.)
// * assignments (set a settable place to a new value)
// * return (return a value from a procedure or block)
// * call (call a procedure)
// * branch (branch to a block depending on a condition)

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	hlast "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
)

type (
	Node interface {
		astutils.NodeIded
		HighlevelNodes() []hlast.Node // the AST nodes that were used to create this node
	}

	Expression interface {
		Node
		exprNode()
	}

	Declaration interface {
		Node
		declNode()
	}

	Statement interface {
		Node
		stmtNode()
	}
)

type (
	BlockExpr struct {
		id         astutils.NodeId
		Label      Label
		Statements []Statement
	}

	BinaryExpr struct {
		id     astutils.NodeId
		Op     token.Token
		tpe    types.Type
		Left   Expression
		Right  Expression
		hlExpr hlast.BinaryExpr
	}

	AtomicLitExpr struct {
		id     astutils.NodeId
		tpe    types.Type
		HlExpr hlast.BasicLitExpr
	}

	Label string
)

func (BlockExpr) exprNode()     {}
func (BinaryExpr) exprNode()    {}
func (AtomicLitExpr) exprNode() {}

func (e BlockExpr) ID() astutils.NodeId     { return e.id }
func (e BinaryExpr) ID() astutils.NodeId    { return e.id }
func (e AtomicLitExpr) ID() astutils.NodeId { return e.id }

func (e BlockExpr) HighlevelNodes() []hlast.Node {
	hnodes := make([]hlast.Node, len(e.Statements))

	for _, stmt := range e.Statements {
		hnodes = append(hnodes, stmt.HighlevelNodes()...)
	}
	return hnodes
}
func (e BinaryExpr) HighlevelNodes() []hlast.Node    { return []hlast.Node{e.hlExpr} }
func (e AtomicLitExpr) HighlevelNodes() []hlast.Node { return []hlast.Node{e.HlExpr} }

type (
	ExprStatement struct {
		Expr Expression
	}
)

func (ExprStatement) stmtNode()             {}
func (e ExprStatement) ID() astutils.NodeId { return e.Expr.ID() }
func (e ExprStatement) HighlevelNodes() []hlast.Node {
	return e.Expr.HighlevelNodes()
}

type (
	ProcDecl struct {
		id     astutils.NodeId
		hlDecl hlast.ProcDecl
		Blocks []BlockExpr
		Type   types.Procedure
		Name   Label
	}
)

func (ProcDecl) declNode()             {}
func (d ProcDecl) ID() astutils.NodeId { return d.id }
func (d ProcDecl) HighlevelNodes() []hlast.Node {
	return []hlast.Node{d.hlDecl}
}

type Module struct {
	Name         Label
	Declarations []Declaration
}

func (b BlockExpr) IsEmpty() bool {
	return len(b.Statements) == 0
}

func (b BlockExpr) LastStatement() Statement {
	if b.IsEmpty() {
		// TODO: synthesize a return statement
		return nil
	}

	return b.Statements[len(b.Statements)-1]
}
