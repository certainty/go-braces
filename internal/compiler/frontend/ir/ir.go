package ir

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/isa"
)

type IRInstruction interface{}

type IRConstant struct {
	Value interface{}
}

func (c IRConstant) String() string {
	return fmt.Sprintf("const %v", c.Value)
}

var _ IRInstruction = (*IRConstant)(nil)

type IRLabel struct {
	Name string
}

type IRGlobalRef struct{}
type IRSet struct{}
type IRClosure struct{}

type IRCall struct {
	Operator IRLabel
	Operands []IRInstruction
}

type IRTailCall struct{}
type IRBranch struct{}
type IRRet struct{}
type IRBlock struct {
	Label        string
	Instructions []IRInstruction
}

func NewConstant(value isa.Value) IRConstant {
	return IRConstant{Value: value}
}

func NewBlock(label string) *IRBlock {
	return &IRBlock{
		Label:        label,
		Instructions: make([]IRInstruction, 0),
	}
}

func (b *IRBlock) String() string {
	return fmt.Sprintf("block %s %v", b.Label, b.Instructions)
}

func (b *IRBlock) AddInstruction(instruction IRInstruction) {
	b.Instructions = append(b.Instructions, instruction)
}

type IR struct {
	Blocks []*IRBlock
}

func NewIR() IR {
	return IR{Blocks: make([]*IRBlock, 0)}
}

func (ir *IR) AddBlock(name string) *IRBlock {
	block := NewBlock(name)
	ir.Blocks = append(ir.Blocks, block)
	return block
}

func LowerToIR(theAST *ast.AST) (*IR, error) {
	var ir IR = NewIR()
	var currentBlock *IRBlock = ir.AddBlock("entry")

	log.Printf("lowering %v", theAST)

	for _, expression := range theAST.Nodes {
		switch exp := expression.(type) {
		case ast.LiteralExpression:
			currentBlock.AddInstruction(NewConstant(exp.Value))
		case ast.UnaryExpression:
			// desugar to a call
		case ast.BinaryExpression:
			// desugar to a call
		default:
			return nil, fmt.Errorf("unhandled expression type %T", expression)
		}
	}

	log.Printf("lowered %v", ir)

	return &ir, nil
}
