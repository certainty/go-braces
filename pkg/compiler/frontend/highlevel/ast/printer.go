package ast

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
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

	return strings.Trim(p.output.String(), " ")
}

func (p *ASTPrinter) Enter(node Node) bool {
	switch n := node.(type) {
	case *Source:
		p.output.WriteString(" (source")
		p.printCommonNodeProps(n)
	case *BadDecl:
		p.output.WriteString(" (bad-decl")
		p.printCommonNodeProps(n)
	case *TypeSpec:
		p.output.WriteString(" (type-spec")
		p.printCommonNodeProps(n)
	case *Identifier:
		p.printIdentifier(n)
	case *ProcDecl:
		p.output.WriteString(" (proc-decl")
		p.printCommonNodeProps(n)
	case *Field:
		p.output.WriteString(" (field")
		p.printCommonNodeProps(n)
	case *BadStmt:
		p.output.WriteString(" (bad-stmt")
		p.printCommonNodeProps(n)
	case *ExprStmt:
		p.output.WriteString(" (expr-stmt")
		p.printCommonNodeProps(n)
	case *BadExpr:
		p.output.WriteString(" (bad-expr")
		p.printCommonNodeProps(n)
	case *BasicLitExpr:
		p.printLiteral(n)
	case *BlockExpr:
		p.output.WriteString(" (block-expr")
		p.printCommonNodeProps(n)
	case *ParenExpr:
		p.output.WriteString(" (paren-expr")
		p.printCommonNodeProps(n)
	case *UnaryExpr:
		p.output.WriteString(" (unary-expr")
		p.printCommonNodeProps(n)
		p.printToken(n.Op)
	case *BinaryExpr:
		p.output.WriteString(" (binary-expr")
		p.printCommonNodeProps(n)
		p.printToken(n.Op)
	}

	return true
}

func (p *ASTPrinter) Leave(node Node) {
	switch node.(type) {
	case BasicLitExpr, Identifier:
		if p.options.LiteralAsStrings {
			return
		}
	}
	p.output.WriteString(")")
}

func (p *ASTPrinter) printCommonNodeProps(node Node) {
	p.printId(node.ID())
	p.printLocation(node)
}

func (p *ASTPrinter) printIdentifier(node *Identifier) {
	p.output.WriteRune(' ')
	if p.options.IdentifiersAsStrings {
		p.output.WriteString(node.Name)
	} else {
		p.output.WriteString("(id ")
		p.printCommonNodeProps(node)
		p.output.WriteString(fmt.Sprintf("\"%s\"", node.Name))
		p.output.WriteString(")")
	}
}

func (p *ASTPrinter) printLiteral(node *BasicLitExpr) {
	p.output.WriteRune(' ')
	if p.options.LiteralAsStrings {
		p.output.WriteString(string(node.Token.Text))
	} else {
		p.output.WriteString("(lit ")
		p.printCommonNodeProps(node)
		p.printToken(node.Token)
		p.output.WriteString(") ")
	}
}

func (p *ASTPrinter) printToken(tok token.Token) {
	p.output.WriteRune(' ')
	if p.options.TokensAsStrings {
		p.output.WriteString(string(tok.Text))
	} else {
		p.output.WriteString(fmt.Sprintf("(tok %s \"%s\") ", tok.Type.String(), string(tok.Text)))
	}
}

func (p *ASTPrinter) printId(id astutils.NodeId) {
	if p.options.IncludeIDs {
		p.output.WriteString(fmt.Sprintf(" (id %s)", id))
	}
}

func (p *ASTPrinter) printLocation(node Node) {
	if p.options.IncludeLocations {
		p.output.WriteString(fmt.Sprintf(" (loc %s)", node.Location()))
	}
}
