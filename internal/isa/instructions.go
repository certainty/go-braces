package isa

import "fmt"

type OpCode uint8

const (
	OP_LOAD  = iota // load from memory
	OP_LOADA        // load address of label
	OP_LOADI        // load immediate value
	OP_STORE        // store to memory
	OP_ADD
	OP_ADDI
	OP_SUB
	OP_SUBI
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

type ImmediateValue uint8

func NewInstruction(code OpCode, operands ...Operand) Instruction {
	return Instruction{
		Opcode:   code,
		Operands: operands,
	}
}

func (i Instruction) String() string {
	return fmt.Sprintf("%d %v", i.Opcode, i.Operands)
}

func InstAdd(left, right, target Register) Instruction {
	return NewInstruction(OP_ADD, left, right, target)
}

func InstAddI(left Register, right ImmediateValue, target Register) Instruction {
	return NewInstruction(OP_ADDI, left, right, target)
}

func InstSub(left, right, target Register) Instruction {
	return NewInstruction(OP_SUB, left, right, target)
}

func InstSubI(left Register, right ImmediateValue, target Register) Instruction {
	return NewInstruction(OP_SUBI, left, right, target)
}

func InstHalt(returnValueRegister Register) Instruction {
	return NewInstruction(OP_HALT, returnValueRegister)
}

func InstRet(returnValueRegister Register) Instruction {
	return NewInstruction(OP_RET, returnValueRegister)
}

func InstLoad(address ConstantAddress, target Register) Instruction {
	return NewInstruction(OP_LOAD, address, target)
}
