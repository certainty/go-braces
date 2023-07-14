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
		return fmt.Sprintf("(unExp %s %s)#%d", unaryOpToString(n.Operator), w.WriteNode(n.Operand), n.ID())
	case BinaryExpression:
		return fmt.Sprintf("(binExp %s %s %s)#%d", binOpToString(n.Operator), w.WriteNode(n.Left), w.WriteNode(n.Right), n.ID())
	case LiteralExpression:
		return fmt.Sprintf("(lit %v)#%d", n.Value, n.ID())
	case Identifier:
		return fmt.Sprintf("(id %s)#%d", n.Label, n.ID())
	case ArgumentDecl:
		return fmt.Sprintf("(arg %s %s)#%d", n.Name.Label, n.TypeDecl().Name.Label, n.ID())
	case CallableDecl:
		var args []string
		for _, arg := range n.Arguments {
			args = append(args, w.WriteNode(arg))
		}
		body := ""
		for _, node := range n.Body.Code {
			body += w.WriteNode(node)
		}
		return fmt.Sprintf("(defn %s (%s) %s)#%d", n.Name.Label, strings.Join(args, " "), body, n.ID())
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
