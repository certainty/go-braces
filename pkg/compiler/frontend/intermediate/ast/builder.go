package ir

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
)

func CreateModule(name Label, source token.Origin) *Module {
	return &Module{
		Name:       name,
		Source:     source,
		Procedures: make([]Procedure, 0),
	}
}

func CreateProcedure(tpe Type, name Label) Procedure {
	return Procedure{
		tpe:    tpe,
		Name:   name,
		Args:   make([]Argument, 0),
		Blocks: make([]*BasicBlock, 0),
	}
}

func CreateBasicBlock(label Label) *BasicBlock {
	return &BasicBlock{
		Label:        label,
		Instructions: make([]Instruction, 0),
	}
}

func CreateSimpleInstruction(register Register, operation Operation, tpe Type, operands ...Operand) SimpleInstruction {
	return SimpleInstruction{
		tpe:       tpe,
		Register:  register,
		Operation: operation,
		Operands:  operands,
	}
}

func CreateAssignmentInstruction(register Register, tpe Type, operand Operand) AssignmentInstruction {
	return AssignmentInstruction{
		Register: register,
		tpe:      tpe,
		Operand:  operand,
	}
}

func CreateReturnInstruction(tpe Type, register Register) ReturnInstruction {
	return ReturnInstruction{
		tpe:      tpe,
		Register: register,
	}
}

type BlockBuilder struct {
	Block             *BasicBlock
	RegisterAllocator *RegisterAllocator
	ModuleBuilder     *IrBuilder
}

func NewBlockBuilder(label Label, registers *RegisterAllocator, builder *IrBuilder) *BlockBuilder {
	return &BlockBuilder{
		Block:             CreateBasicBlock(label),
		RegisterAllocator: registers,
		ModuleBuilder:     builder,
	}
}

func (b *BlockBuilder) IsEmpty() bool {
	return len(b.Block.Instructions) == 0
}

func (b *BlockBuilder) LastInstruction() Instruction {
	if b.IsEmpty() {
		return nil
	}
	return b.Block.Instructions[len(b.Block.Instructions)-1]
}

func (b *BlockBuilder) OpLit(tpe Type, operand Operand) Register {
	target := b.RegisterAllocator.Next("")
	instruction := CreateAssignmentInstruction(target, tpe, operand)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
	return target
}

func (b *BlockBuilder) OpAdd(tpe Type, lhs Operand, rhs Operand) Register {
	target := b.RegisterAllocator.Next("")
	instruction := CreateSimpleInstruction(target, Add, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
	return target
}

func (b *BlockBuilder) OpMul(tpe Type, lhs Operand, rhs Operand) Register {
	target := b.RegisterAllocator.Next("")
	instruction := CreateSimpleInstruction(target, Mul, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
	return target
}

func (b *BlockBuilder) OpSub(tpe Type, lhs Operand, rhs Operand) Register {
	target := b.RegisterAllocator.Next("")
	instruction := CreateSimpleInstruction(target, Sub, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
	return target
}

func (b *BlockBuilder) OpRet(tpe Type, register Register) {
	instruction := CreateReturnInstruction(tpe, register)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}
