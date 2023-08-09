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
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
)

type (
	Node interface {
		astutils.NodeIded
		HighlevelNodeIds() []astutils.NodeId
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
		id           astutils.NodeId
		Op           token.Token
		Type         types.Type
		Left         Expression
		Right        Expression
		hlExprNodeId astutils.NodeId
	}

	AtomicLitExpr struct {
		id           astutils.NodeId
		tpe          types.Type
		Value        token.Token
		hlExprNodeId astutils.NodeId
	}

	Label struct {
		id     astutils.NodeId
		Origin *astutils.NodeId // if it stems from an identifier in the source code
		Value  string
	}

	Variable struct {
		id      astutils.NodeId
		Version astutils.Version
		Origin  *astutils.NodeId // if it stems from an identifier in the source code
		Name    string
	}
)

func (BlockExpr) exprNode()     {}
func (BinaryExpr) exprNode()    {}
func (AtomicLitExpr) exprNode() {}
func (Label) exprNode()         {}
func (Variable) exprNode()      {}

func (e BlockExpr) ID() astutils.NodeId     { return e.id }
func (e BinaryExpr) ID() astutils.NodeId    { return e.id }
func (e AtomicLitExpr) ID() astutils.NodeId { return e.id }
func (e Label) ID() astutils.NodeId         { return e.id }
func (e Variable) ID() astutils.NodeId      { return e.id }

func (e BlockExpr) HighlevelNodeIds() []astutils.NodeId {
	hnodes := make([]astutils.NodeId, len(e.Statements))

	for _, stmt := range e.Statements {
		hnodes = append(hnodes, stmt.HighlevelNodeIds()...)
	}
	return hnodes
}
func (e BinaryExpr) HighlevelNodeIds() []astutils.NodeId    { return []astutils.NodeId{e.hlExprNodeId} }
func (e AtomicLitExpr) HighlevelNodeIds() []astutils.NodeId { return []astutils.NodeId{e.hlExprNodeId} }
func (e Label) HighlevelNodeIds() []astutils.NodeId {
	if e.Origin == nil {
		return []astutils.NodeId{}
	}
	return []astutils.NodeId{*e.Origin}
}

func (e Variable) HighlevelNodeIds() []astutils.NodeId {
	if e.Origin == nil {
		return []astutils.NodeId{}
	}
	return []astutils.NodeId{*e.Origin}
}

type (
	ExprStatement struct {
		Expr Expression
	}

	ReturnStmt struct {
		id    astutils.NodeId
		Value Expression
	}

	AssignStmt struct {
		id       astutils.NodeId
		Variable *Variable
		Expr     Expression
	}

	// SSA specific nodes
	Phi struct {
		id       astutils.NodeId
		Variable *Variable
		Choices  []PhiChoice
	}

	PhiChoice struct {
		Predecessor *SSABlock
		Value       Expression
	}

	SSABlock struct {
		*BlockExpr
		Predecessors []*SSABlock
		Successors   []*SSABlock
	}
)

func (ExprStatement) stmtNode() {}
func (ReturnStmt) stmtNode()    {}
func (AssignStmt) stmtNode()    {}
func (Phi) stmtNode()           {}

func (e ExprStatement) ID() astutils.NodeId { return e.Expr.ID() }
func (e ReturnStmt) ID() astutils.NodeId    { return e.id }
func (e AssignStmt) ID() astutils.NodeId    { return e.id }
func (e Phi) ID() astutils.NodeId           { return e.id }

func (e ExprStatement) HighlevelNodeIds() []astutils.NodeId {
	return e.Expr.HighlevelNodeIds()
}

func (e ReturnStmt) HighlevelNodeIds() []astutils.NodeId {
	return e.Value.HighlevelNodeIds()
}

func (e AssignStmt) HighlevelNodeIds() []astutils.NodeId {
	if e.Variable.Origin == nil {
		return []astutils.NodeId{}
	}
	return []astutils.NodeId{*e.Variable.Origin}
}

func (e Phi) HighlevelNodeIds() []astutils.NodeId {
	return []astutils.NodeId{}
}

type (
	ProcDecl struct {
		id        astutils.NodeId
		hlDeclID  astutils.NodeId
		Blocks    []*BlockExpr
		SSABlocks []*SSABlock // will be set after SSA conversion
		Type      types.Procedure
		Name      Label
	}
)

func (ProcDecl) declNode()             {}
func (d ProcDecl) ID() astutils.NodeId { return d.id }
func (d ProcDecl) HighlevelNodeIds() []astutils.NodeId {
	return []astutils.NodeId{d.hlDeclID}
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
		return nil
	}

	return b.Statements[len(b.Statements)-1]
}

func (Module) declNode()             {}
func (m Module) ID() astutils.NodeId { return m.Name.ID() }
func (m Module) HighlevelNodeIds() []astutils.NodeId {
	hnodes := make([]astutils.NodeId, len(m.Declarations))

	for _, decl := range m.Declarations {
		hnodes = append(hnodes, decl.HighlevelNodeIds()...)
	}

	return hnodes
}
