package isa

import "fmt"

// load store architecture

type OpCode uint8

const (
	OP_TRUE OpCode = iota
	OP_FALSE
	OP_CONST
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_NEG
	OP_MOD
	OP_REM
	OP_AND
	OP_OR
	OP_XOR
	OP_NOT
	OP_MOV
	OP_SHL
	OP_SHR
	OP_BEQ
	OP_BNE
	OP_PRNT
	OP_NEW_TABLE
	OP_GET
	OP_SET
	OP_HALT
	OP_RET
)

type Operand interface{}

type Register uint32

type ConstantAddress uint64

type InstructionAddress uint64

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

// loads the constant into the given register
func InstConst(address ConstantAddress, target Register) Instruction {
	return NewInstruction(OP_CONST, address, target)
}
