package bytecode

type OpCode uint8

const (
	OpPush OpCode = iota
	OpPop
	OpNop
	OpConst
	OpSet
	OpTrue
	OpFalse
	OpNil
	OpUnspecified
)

type Instruction struct {
	OpCode OpCode
	Args   []interface{}
}

func NewInstruction(opCode OpCode, args ...interface{}) Instruction {
	return Instruction{OpCode: opCode, Args: args}
}

func Set() Instruction {
	return NewInstruction(OpSet)
}

func True() Instruction {
	return NewInstruction(OpTrue)
}

func False() Instruction {
	return NewInstruction(OpFalse)
}

func Unspecified() Instruction {
	return NewInstruction(OpUnspecified)
}

func Nil() Instruction {
	return NewInstruction(OpNil)
}

func Const(address ConstAddress) Instruction {
	return NewInstruction(OpConst, address)
}
