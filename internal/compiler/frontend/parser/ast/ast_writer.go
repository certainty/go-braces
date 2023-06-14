package ast

import (
	"fmt"
	"strings"
)

type ASTWriter struct{}

func NewASTWriter() *ASTWriter {
	return &ASTWriter{}
}

func (w *ASTWriter) WriteNode(node Node) string {
	switch n := node.(type) {
	case UnaryExpression:
		return fmt.Sprintf("(%s %s)", unaryOpToString(n.Operator), w.WriteNode(n.Operand))
	case BinaryExpression:
		return fmt.Sprintf("(%s %s %s)", binOpToString(n.Operator), w.WriteNode(n.Left), w.WriteNode(n.Right))
	case LiteralExpression:
		return fmt.Sprintf("%v", n.Value)
	default:
		return ""
	}
}

func unaryOpToString(op UnaryOperator) string {
	switch op {
	case UnaryOpNot:
		return "not"
	case UnaryOpNeg:
		return "-"
	case UnaryOpPos:
		return "+"
	default:
		return ""
	}
}

func binOpToString(op BinaryOperator) string {
	switch op {
	case BinOpAdd:
		return "+"
	case BinOpSub:
		return "-"
	case BinOpMul:
		return "*"
	case BinOpDiv:
		return "/"
	case BinOpPow:
		return "**"
	case BinOpMod:
		return "%"
	case BinOpAnd:
		return "&&"
	case BinOpOr:
		return "||"
	default:
		return ""
	}
}

// write canonical
func (w *ASTWriter) Write(ast *AST) string {
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
