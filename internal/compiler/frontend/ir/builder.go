package ir

import "github.com/certainty/go-braces/internal/compiler/location"

func CreateModule(name Label, source location.Origin) *Module {
	return &Module{
		Name:      name,
		Source:    source,
		Functions: make([]*Function, 0),
	}
}

func CreateFunction(tpe Type, name Label) *Function {
	return &Function{
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
}

func NewBlockBuilder(label Label, registers *RegisterAllocator) *BlockBuilder {
	return &BlockBuilder{
		Block:             CreateBasicBlock(label),
		RegisterAllocator: registers,
	}
}

func (b *BlockBuilder) OpLit(tpe Type, operand Operand) {
	instruction := CreateAssignmentInstruction(b.RegisterAllocator.Next(""), tpe, operand)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *BlockBuilder) OpAdd(tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(b.RegisterAllocator.Next(""), Add, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *BlockBuilder) OpSub(tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(b.RegisterAllocator.Next(""), Sub, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *BlockBuilder) OpRet(tpe Type, register Register) {
	instruction := CreateReturnInstruction(tpe, register)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}
