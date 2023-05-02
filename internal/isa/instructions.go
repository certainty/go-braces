package isa

type OpCode uint8
type Register uint32

const (
	OP_TRUE OpCode = iota
	OP_FALSE
	OP_HALT
	OP_RET
)

type Instruction struct {
	Opcode   OpCode
	Operands []interface{}
}

func NewInstruction(code OpCode, operands ...interface{}) Instruction {
	return Instruction{
		Opcode:   code,
		Operands: operands,
	}
}

func InstTrue(register Register) Instruction {
	return NewInstruction(OP_TRUE, register)
}

func InstHalt(returnValueRegister Register) Instruction {
	return NewInstruction(OP_HALT, returnValueRegister)
}
