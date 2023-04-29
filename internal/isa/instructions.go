package isa

type Opcode uint8

const (
	OP_PUSH_TRUE Opcode = iota
	OP_RET
)

type Instruction struct {
	Opcode   Opcode
	Operands []interface{}
}
