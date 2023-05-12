package isa

import "fmt"

type OpCode uint8

const (
	OP_TRUE OpCode = iota
	OP_FALSE
	OP_HALT
	OP_RET
)

type Operand interface{}

type Register uint32

type Instruction struct {
	Opcode   OpCode
	Operands []Operand
}

func NewInstruction(code OpCode, operands ...Operand) Instruction {
	return Instruction{
		Opcode:   code,
		Operands: operands,
	}
}

func (i Instruction) String() string {
	return fmt.Sprintf("%d %v", i.Opcode, i.Operands)
}

func InstTrue(register Register) Instruction {
	return NewInstruction(OP_TRUE, register)
}

func InstFalse(register Register) Instruction {
	return NewInstruction(OP_FALSE, register)
}

func InstHalt(returnValueRegister Register) Instruction {
	return NewInstruction(OP_HALT, returnValueRegister)
}
