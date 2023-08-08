package ssa

import "github.com/certainty/go-braces/pkg/compiler/frontend/astutils"

type SSAWalker struct{}

func NewSSAWalker() *SSAWalker {
	return &SSAWalker{}
}

func (w *SSAWalker) Walk(visitor astutils.Visitor[Node], node Node) {
	cont := visitor.Enter(node)
	if !cont {
		return
	}
	defer visitor.Leave(node)

	switch node := node.(type) {
	case BasicBlock:
	case BinaryExpr:
	case AtomicLitExpr:
	case ProcDecl:
	case VariableExpr:
	case Phi:
	case SetStmt:
		w.Walk(visitor, node.Value)
	case ReturnStmt:
		w.Walk(visitor, node.Value)
	case ExprStatement:
		w.Walk(visitor, node.Expr)
	}
}
