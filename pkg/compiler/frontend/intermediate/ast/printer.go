package ast

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/pkg/compiler/frontend/astutils"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
)

type PrintingOptions struct {
	IncludeIDs          bool // include node IDs in the output
	TokensAsStrings     bool // don't create a node for tokens but inline their textual representation
	IncludeReferenceIDs bool // include reference node IDs from the highlevel AST
	LiteralAsStrings    bool // don't create a node for literals but inline their textual representation
	LabelsAsStrings     bool // don't create a node for identifiers but inline their textual representation
	PrettyPrint         bool // pretty print the output
}

func PrintCanonical() PrintingOptions {
	return PrintingOptions{
		IncludeIDs:          false,
		TokensAsStrings:     true,
		IncludeReferenceIDs: false,
		LiteralAsStrings:    true,
		LabelsAsStrings:     true,
		PrettyPrint:         false,
	}
}

func PrintTruthfully() PrintingOptions {
	return PrintingOptions{
		IncludeIDs:          true,
		IncludeReferenceIDs: false,
		TokensAsStrings:     false,
		LiteralAsStrings:    false,
		LabelsAsStrings:     false,
		PrettyPrint:         true,
	}
}

type Printer struct {
	output  strings.Builder
	options PrintingOptions
}

func Print(node Node, options PrintingOptions) string {
	printer := &Printer{
		output: strings.Builder{},
	}

	Walk(printer, node)

	return strings.Trim(printer.output.String(), " ")
}

func (p *Printer) Enter(node Node) bool {
	switch n := node.(type) {
	case *Module:
		p.output.WriteString(fmt.Sprintf(" (module %s", n.Name.Value))
	case *Label:
		p.printLabel(*n)
	case *BlockExpr:
		p.output.WriteString(" (block-expr")
		p.printCommonNodeProps(n)
	case *BinaryExpr:
		p.output.WriteString(" (binary-expr")
		p.printCommonNodeProps(n)
		p.printToken(n.Op)
	case *AtomicLitExpr:
		p.printLiteral(*n)
	case *ExprStatement:
		p.output.WriteString(" (expr-stmt")
		p.printCommonNodeProps(n)
	}
	return true
}

func (p *Printer) Leave(node Node) {
	switch node.(type) {
	case Label:
		return
	}
	p.output.WriteString(")")
}

func (p *Printer) printLiteral(node AtomicLitExpr) {
	p.output.WriteRune(' ')

	if p.options.LiteralAsStrings {
		p.output.WriteString(string(node.Value.Text))
	} else {
		p.output.WriteString("(lit ")
		p.printCommonNodeProps(node)
		p.printToken(node.Value)
		p.output.WriteString(") ")
	}
}

func (p *Printer) printLabel(node Label) {
	p.output.WriteRune(' ')
	if p.options.LabelsAsStrings {
		p.output.WriteString(node.Value)
	} else {
		p.output.WriteString("(id ")
		p.printCommonNodeProps(node)
		p.output.WriteString(fmt.Sprintf("\"%s\"", node.Value))
		p.output.WriteString(")")
	}
}

func (p *Printer) printCommonNodeProps(node Node) {
	p.printId(node.ID())
	p.printReferenceNodeIds(node)
}

func (p *Printer) printId(id astutils.NodeId) {
	if p.options.IncludeIDs {
		p.output.WriteString(fmt.Sprintf(" (id %s)", id))
	}
}

func (p *Printer) printToken(tok token.Token) {
	p.output.WriteRune(' ')
	if p.options.TokensAsStrings {
		p.output.WriteString(string(tok.Text))
	} else {
		p.output.WriteString(fmt.Sprintf("(tok %s \"%s\") ", tok.Type.String(), string(tok.Text)))
	}
}

func (p *Printer) printReferenceNodeIds(node Node) {
	if p.options.IncludeReferenceIDs && len(node.HighlevelNodeIds()) > 0 {
		p.output.WriteString(" (ref-ids ")

		for _, id := range node.HighlevelNodeIds() {
			p.output.WriteString(fmt.Sprintf("%s ", id))
		}

		p.output.WriteString(")")
	}
}
