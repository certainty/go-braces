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

type Register = uint64
type ConstantAddress = uint64
type InstructionAddress = uint64
type ImmediateValue = uint8

type Operand = int64

// On the instruction level all type safety of operands is lost.
// Operands are just i64 values and the opcode denotes how they are interpreted
// This makes the VM faster, but also more dangerous.
type Instruction struct {
	Opcode OpCode
	// either register, address or immediate value
	Operands []Operand
}

func NewInstruction(code OpCode, operands ...Operand) Instruction {
	return Instruction{
		Opcode: code,
		// operands are always target-register, op1, op2
		Operands: operands,
	}
}

func (i Instruction) String() string {
	return fmt.Sprintf("%d %v", i.Opcode, i.Operands)
}

func InstAdd(target, left, right Register) Instruction {
	return NewInstruction(OP_ADD, Operand(target), Operand(left), Operand(right))
}

func InstMul(target, left, right Register) Instruction {
	return NewInstruction(OP_MUL, Operand(target), Operand(left), Operand(right))
}

func InstAddI(target Register, left Register, right ImmediateValue) Instruction {
	return NewInstruction(OP_ADDI, Operand(target), Operand(left), Operand(right))
}

func InstSub(target, left, right Register) Instruction {
	return NewInstruction(OP_SUB, Operand(target), Operand(left), Operand(right))
}

func InstSubI(target Register, left Operand, right Operand) Instruction {
	return NewInstruction(OP_SUBI, Operand(target), Operand(left), Operand(right))
}

func InstHalt(reg Register) Instruction {
	return NewInstruction(OP_HALT, Operand(reg))
}

func InstRet(reg Register) Instruction {
	return NewInstruction(OP_RET, Operand(reg))
}

func InstLoad(target Register, address ConstantAddress) Instruction {
	return NewInstruction(OP_LOAD, Operand(target), Operand(address))
}

func InstLoadI(target Register, value ImmediateValue) Instruction {
	return NewInstruction(OP_LOADI, Operand(target), Operand(value))
}
