package ir

import (
	"fmt"
	"log"

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

func LowerToIR(coreAst *parser.CoreAST) (*IR, error) {
	var ir IR = NewIR()
	var currentBlock *IRBlock = ir.AddBlock("entry")

	log.Printf("lowering %v", coreAst)

	for _, expression := range coreAst.Expressions {
		switch exp := expression.(type) {
		case parser.LiteralExpression:
			value, err := isa.ValueFromDatum(exp.Datum)
			if err != nil {
				return nil, fmt.Errorf("could not convert literal to value: %w", err)
			}
			currentBlock.AddInstruction(NewConstant(value))
		default:
			return nil, fmt.Errorf("unhandled expression type %T", expression)
		}
	}

	log.Printf("lowered %v", ir)

	return &ir, nil
}
