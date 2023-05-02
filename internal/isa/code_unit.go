package isa

type CodeUnit struct {
	Constants    []Value
	Instructions []Instruction
}

func EmptyCodeUnit() CodeUnit {
	return CodeUnit{
		Constants:    []Value{},
		Instructions: []Instruction{},
	}
}
