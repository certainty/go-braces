package ast

import (
	"fmt"
	"strings"
)

type PrintingOptions struct {
	IncludeIDs           bool // include node IDs in the output
	IncludeLocations     bool // include position information in the output
	TokensAsStrings      bool // don't create a node for tokens but inline their textual representation
	LiteralAsStrings     bool // don't create a node for literals but inline their textual representation
	IdentifiersAsStrings bool // don't create a node for identifiers but inline their textual representation
	PrettyPrint          bool // pretty print the output
}

func PrintCanonical() PrintingOptions {
	return PrintingOptions{
		IncludeIDs:           false,
		IncludeLocations:     false,
		TokensAsStrings:      true,
		LiteralAsStrings:     true,
		IdentifiersAsStrings: true,
		PrettyPrint:          false,
	}
}

func PrintTruthfully() PrintingOptions {
	return PrintingOptions{
		IncludeIDs:           true,
		IncludeLocations:     true,
		TokensAsStrings:      false,
		LiteralAsStrings:     false,
		IdentifiersAsStrings: false,
		PrettyPrint:          true,
	}
}

type ASTPrinter struct {
	options PrintingOptions
	output  strings.Builder
}

func Print(node Node, options PrintingOptions) string {
	printer := &ASTPrinter{
		options: options,
		output:  strings.Builder{},
	}

	return printer.print(node)
}

func (p *ASTPrinter) print(node Node) string {
	Walk(p, node)

	return p.output.String()
}

func (p *ASTPrinter) Enter(node Node) bool {
	switch n := node.(type) {
	case *Source:
		p.output.WriteString("(source ")
		p.printCommonNodeProps(n)
	case *BadDecl:
		p.output.WriteString("(bad-decl ")
		p.printCommonNodeProps(n)
	case *TypeSpec:
		p.output.WriteString("(type-spec ")
		p.printCommonNodeProps(n)
	case *Identifier:
		p.printIdentifier(n)
	case *ProcDecl:
		p.output.WriteString("(proc-decl ")
		p.printCommonNodeProps(n)
	case *Field:
		p.output.WriteString("(field ")
		p.printCommonNodeProps(n)
	case *BadStmt:
		p.output.WriteString("(bad-stmt ")
		p.printCommonNodeProps(n)
	case *ExprStmt:
		p.output.WriteString("(expr-stmt ")
		p.printCommonNodeProps(n)
	case *BadExpr:
		p.output.WriteString("(bad-expr ")
		p.printCommonNodeProps(n)
	case *BasicLitExpr:
		p.printLiteral(n)
	case *BlockExpr:
		p.output.WriteString("(block-expr ")
		p.printCommonNodeProps(n)
	case *ParenExpr:
		p.output.WriteString("(paren-expr ")
		p.printCommonNodeProps(n)
	case *UnaryExpr:
		p.output.WriteString("(unary-expr ")
		p.printCommonNodeProps(n)
	case *BinaryExpr:
		p.output.WriteString("(binary-expr ")
		p.printCommonNodeProps(n)
	}

	return true
}

func (p *ASTPrinter) Leave(node Node) {
	p.output.WriteString(") ")
}

func (p *ASTPrinter) printCommonNodeProps(node Node) {
	p.printId(node.ID())
	p.printLocation(node)
}

func (p *ASTPrinter) printIdentifier(node *Identifier) {
	if p.options.IdentifiersAsStrings {
		p.output.WriteString(node.Name)
		p.output.WriteRune(' ')
	} else {
		p.output.WriteString("(id ")
		p.printCommonNodeProps(node)
		p.output.WriteString(fmt.Sprintf("\"%s\"", node.Name))
		p.output.WriteString(") ")
	}
}

func (p *ASTPrinter) printLiteral(node *BasicLitExpr) {
	if p.options.LiteralAsStrings {
		p.output.WriteString(node.Token.Type.String())
		p.output.WriteRune(' ')
	} else {
		p.output.WriteString("(lit ")
		p.printCommonNodeProps(node)
		p.output.WriteString(" (tok")
		p.output.WriteString(node.Token.Type.String())
		p.output.WriteString(fmt.Sprintf("\"%s\"", string(node.Token.Text)))
		p.output.WriteString(")) ")
	}
}

func (p *ASTPrinter) printId(id NodeId) {
	if p.options.IncludeIDs {
		p.output.WriteString(fmt.Sprintf("(id %d) ", id))
	}
}

func (p *ASTPrinter) printLocation(node Node) {
	if p.options.IncludeLocations {
		p.output.WriteString(fmt.Sprintf("(loc %s) ", node.Location()))
	}
}

// type ASTWriter struct{}

// func NewASTWriter() *ASTWriter {
// 	return &ASTWriter{}
// }

// func (w *ASTWriter) WriteNode(node Node) string {
// 	switch n := node.(type) {
// 	case Identifier:
// 		return fmt.Sprintf("(id %s)#%d", n.Name, n.ID())
// 	case BadExpr:
// 		return fmt.Sprintf("(badExpr)#%d", n.ID())
// 	case BasicLitExpr:
// 		return fmt.Sprintf("(lit %v)#%d", n.Token, n.ID())
// 	case ParenExpr:
// 		return fmt.Sprintf("(paren %s)#%d", w.WriteNode(n.Expr), n.ID())
// 	case UnaryExpr:
// 		return fmt.Sprintf("(unExp %s %s)#%d", string(n.Op.Text), w.WriteNode(n.Expr), n.ID())
// 	case BinaryExpr:
// 		return fmt.Sprintf("(binExp %s %s %s)#%d", n.Op, w.WriteNode(n.Left), w.WriteNode(n.Right), n.ID())
// 	case BlockExpr:
// 		stms := []string{}
// 		for _, stm := range n.Statements {
// 			stms = append(stms, w.WriteNode(stm))
// 		}
// 		return fmt.Sprintf("(block %s)#%d", strings.Join(stms, "\n"), n.ID())
// 	case BadStmt:
// 		return fmt.Sprintf("(badStmt)#%d", n.ID())
// 	case ExprStmt:
// 		return w.WriteNode(n.Expr)
// 	case BadDecl:
// 		return fmt.Sprintf("(badDecl)#%d", n.ID())
// 	case TypeSpec:
// 		return fmt.Sprintf("(typeSpec %s)#%d", w.WriteNode(n.Name), n.ID())
// 	case ProcDecl:
// 		return fmt.Sprintf("(defn %s () %s)#%d", n.Name.Name, w.WriteNode(n.Body), n.ID())
// 	default:
// 		return ""
// 	}
// }
