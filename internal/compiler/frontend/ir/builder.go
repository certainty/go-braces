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
		Blocks: make([]*BasicBlock, 0),
	}
}

func CreateBasicBlock(label Label) *BasicBlock {
	return &BasicBlock{
		Label:        label,
		Instructions: make([]Instruction, 0),
	}
}

func CreateSimpleInstruction(register Register, operation Operation, tpe Type, operands ...Operand) *SimpleInstruction {
	return &SimpleInstruction{
		tpe:       tpe,
		Register:  register,
		Operation: operation,
		Operand:   operands,
	}
}

func CreateReturnInstruction(tpe Type, register Register) *ReturnInstruction {
	return &ReturnInstruction{
		tpe:      tpe,
		Register: register,
	}
}

type Builder struct {
	// the module
	Module *Module

	// the current block
	Block *BasicBlock
}

func NewBuilder() *Builder {
	return &Builder{}
}

func SetInsertPoint(builder *Builder, block *BasicBlock) {
	builder.Block = block
}

func (b *Builder) OpAdd(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Add, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpSub(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Sub, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpMul(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Mul, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpDiv(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Div, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpOr(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Or, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpAnd(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, And, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpXor(register Register, tpe Type, lhs Operand, rhs Operand) {
	instruction := CreateSimpleInstruction(register, Xor, tpe, lhs, rhs)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpNeg(register Register, tpe Type, operand Operand) {
	instruction := CreateSimpleInstruction(register, Neg, tpe, operand)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}

func (b *Builder) OpRet(tpe Type, register Register) {
	instruction := CreateReturnInstruction(tpe, register)
	b.Block.Instructions = append(b.Block.Instructions, instruction)
}
