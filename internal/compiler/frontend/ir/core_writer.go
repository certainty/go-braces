package ir

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"strings"
)

type CoreASTWriter struct{}

func NewCoreASTWriter() *CoreASTWriter {
	return &CoreASTWriter{}
}

func (w *CoreASTWriter) WriteNode(node CoreNode) string {
	switch n := node.(type) {
	case CoreConstant:
		return fmt.Sprintf("%v", n.Value)
	case Apply:
		operands := make([]string, len(n.Operands))
		for i, operand := range n.Operands {
			operands[i] = w.WriteNode(operand)
		}
		return fmt.Sprintf("(apply %s %s)", w.WriteNode(n.Operator), strings.Join(operands, " "))
	case LogicalConnective:
		left := w.WriteNode(n.Left)
		right := w.WriteNode(n.Right)
		return fmt.Sprintf("(logic-apply %s %s %s)", w.writeConnective(n.Operator), left, right)
	case Primitive:
		return w.writePrimitiveOp(n.Op)
	default:
		return ""
	}
}

func (w *CoreASTWriter) Write(ast CoreAST) string {
	var sb strings.Builder
	sb.WriteRune('(')
	for idx, node := range ast.Nodes {
		sb.WriteString(w.WriteNode(node))
		// check if it's the last node
		if idx != len(ast.Nodes)-1 {
			sb.WriteString(" ")
		}
	}
	sb.WriteRune(')')
	return sb.String()
}

func (w *CoreASTWriter) writePrimitiveOp(op PrimitiveOp) string {
	switch op {
	case PrimitiveOp(ast.BinOpAdd):
		return "prim#plus"
	case PrimitiveOp(ast.BinOpSub):
		return "prim#minus"
	case PrimitiveOp(ast.BinOpMul):
		return "prim#mul"
	case PrimitiveOp(ast.BinOpDiv):
		return "prim#div"
	case PrimitiveOp(ast.BinOpPow):
		return "prim#pow"
	case PrimitiveOp(ast.BinOpMod):
		return "prim#mod"
	case PrimitiveOp(ast.BinOpAnd):
		return "prim#and"
	case PrimitiveOp(ast.BinOpOr):
		return "prim#or"
	default:
		return ""
	}
}

func (w *CoreASTWriter) writeConnective(op LogicalOperator) string {
	switch op {
	case LogicalOperator(ast.BinOpAnd):
		return "prim#and"
	case LogicalOperator(ast.BinOpOr):
		return "prim#or"
	default:
		return ""
	}
}
