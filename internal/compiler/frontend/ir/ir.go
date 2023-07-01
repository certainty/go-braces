package ir

import (
	"fmt"
	"github.com/certainty/go-braces/internal/isa"
)

// linearized intermediate representation

type IRInstruction interface{}

type IRConstant struct {
	Value isa.Value
}

func (c IRConstant) String() string {
	return fmt.Sprintf("const %v", c.Value)
}

type IRLabel struct{}
type IRRet struct{}
type IRBlock struct {
	Label        string
	Instructions []IRInstruction
}

func NewIRConstant(value interface{}) IRConstant {
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

func LowerToIR(coreAst *CoreAST) (*IR, error) {
	var ir IR = NewIR()
	var currentBlock *IRBlock = ir.AddBlock("entry")

	for _, node := range coreAst.Nodes {
		switch node := node.(type) {
		case CoreConstant:
			currentBlock.AddInstruction(NewIRConstant(node.Value))
		case LogicalConnective:
		case Apply:
		default:
			return nil, fmt.Errorf("unhandled expression type %T", node)
		}
	}

	return &ir, nil
}
