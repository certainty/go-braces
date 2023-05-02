package ir

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/isa"
)

type IRInstruction interface{}

type IRConstant struct {
	Value isa.Value
}

type IRLabel struct{}
type IRGlobalRef struct{}
type IRSet struct{}
type IRClosure struct{}
type IRCall struct{}
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

func (b *IRBlock) AddInstruction(instruction IRInstruction) {
	b.Instructions = append(b.Instructions, instruction)
}

type IR struct {
	Blocks []IRBlock
}

func NewIR() IR {
	return IR{Blocks: make([]IRBlock, 0)}
}

func (ir *IR) AddBlock(name string) *IRBlock {
	block := NewBlock(name)
	ir.Blocks = append(ir.Blocks, *block)
	return block
}

func LowerToIR(coreAst *parser.CoreAST) (*IR, error) {
	var ir IR = IR{Blocks: make([]IRBlock, 0)}
	var currentBlock *IRBlock = ir.AddBlock("entry")

	for _, expression := range coreAst.Expressions {
		switch exp := expression.(type) {
		case parser.LiteralExpression:
			{
				currentBlock.AddInstruction(NewConstant(exp.Datum))
			}
		default:
			return nil, fmt.Errorf("unhandled expression type %T", expression)
		}
	}
	return &ir, nil
}
