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
	case Identifier:
		return fmt.Sprintf("(id %s)#%d", n.Name, n.ID())
	case BadExpr:
		return fmt.Sprintf("(badExpr)#%d", n.ID())
	case BasicLitExpr:
		return fmt.Sprintf("(lit %v)#%d", n.Token, n.ID())
	case ParenExpr:
		return fmt.Sprintf("(paren %s)#%d", w.WriteNode(n.Expr), n.ID())
	case UnaryExpr:
		return fmt.Sprintf("(unExp %s %s)#%d", string(n.Op.Text), w.WriteNode(n.Expr), n.ID())
	case BinaryExpr:
		return fmt.Sprintf("(binExp %s %s %s)#%d", n.Op, w.WriteNode(n.Left), w.WriteNode(n.Right), n.ID())
	case BlockExpr:
		stms := []string{}
		for _, stm := range n.Statements {
			stms = append(stms, w.WriteNode(stm))
		}
		return fmt.Sprintf("(block %s)#%d", strings.Join(stms, "\n"), n.ID())
	case BadStmt:
		return fmt.Sprintf("(badStmt)#%d", n.ID())
	case ExprStmt:
		return w.WriteNode(n.Expr)
	case BadDecl:
		return fmt.Sprintf("(badDecl)#%d", n.ID())
	case TypeSpec:
		return fmt.Sprintf("(typeSpec %s)#%d", w.WriteNode(n.Name), n.ID())
	case ProcDecl:
		return fmt.Sprintf("(defn %s () %s)#%d", n.Name.Name, w.WriteNode(n.Body), n.ID())
	default:
		return ""
	}
}
