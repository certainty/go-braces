package ir

import (
	"fmt"
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
	case Call:
		operands := make([]string, len(n.Operands))
		for i, operand := range n.Operands {
			operands[i] = w.WriteNode(operand)
		}
		return fmt.Sprintf("(call %s %s)", w.WriteNode(n.Operator), strings.Join(operands, " "))
	case Junction:
		left := w.WriteNode(n.Left)
		right := w.WriteNode(n.Right)
		return fmt.Sprintf("(junct %s %s %s)", w.writeJunctor(n.Junctor), left, right)
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
	case PrimitiveOpAdd:
		return "prim#plus"
	case PrimitiveOpSub:
		return "prim#minus"
	case PrimitiveOpMul:
		return "prim#mul"
	case PrimitiveOpDiv:
		return "prim#div"
	case PrimitiveOpPow:
		return "prim#pow"
	case PrimitiveOpMod:
		return "prim#mod"
	case PrimitiveOpAnd:
		return "prim#and"
	case PrimitiveOpOr:
		return "prim#or"
	default:
		return ""
	}
}

func (w *CoreASTWriter) writeJunctor(j JunctionOp) string {
	switch j {
	case JunctionOpAnd:
		return "prim#and"
	case JunctionOpOr:
		return "prim#or"
	default:
		return ""
	}
}
